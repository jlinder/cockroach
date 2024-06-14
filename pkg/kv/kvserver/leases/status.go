// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package leases

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvserverpb"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/redact"
)

// StatusInput is the input to the Status function.
type StatusInput struct {
	// Information about the local store and replica.
	LocalStoreID       roachpb.StoreID
	MaxOffset          time.Duration
	MinProposedTs      hlc.ClockTimestamp
	MinValidObservedTs hlc.ClockTimestamp

	// The current time and the time of the request to evaluate the lease for.
	Now       hlc.ClockTimestamp
	RequestTs hlc.Timestamp

	// The lease to evaluate.
	Lease roachpb.Lease
}

// Status returns a lease status. The lease status is linked to the desire to
// serve a request at a specific timestamp (which may be a future timestamp)
// under the lease, as well as a notion of the current hlc time (now).
//
// # Explanation
//
// A status of ERROR indicates a failure to determine the correct lease status,
// and should not occur under normal operations. The caller's only recourse is
// to give up or to retry.
//
// If the lease is expired according to the now timestamp (and, in the case of
// epoch-based leases, the liveness epoch), a status of EXPIRED is returned.
// Note that this ignores the timestamp of the request, which may well
// technically be eligible to be served under the lease. The key feature of an
// EXPIRED status is that it reflects that a new lease with a start timestamp
// greater than or equal to now can be acquired non-cooperatively.
//
// If the lease is not EXPIRED, the lease's start timestamp is checked against
// the minProposedTimestamp. This timestamp indicates the oldest timestamp that
// a lease can have as its start time and still be used by the node. It is set
// both in cooperative lease transfers and to prevent reuse of leases across
// node restarts (which would result in latching violations). Leases with start
// times preceding this timestamp are assigned a status of PROSCRIBED and can
// not be used. Instead, a new lease should be acquired by callers.
//
// If the lease is not EXPIRED or PROSCRIBED, the request timestamp is taken
// into account. The expiration timestamp is adjusted for clock offset; if the
// request timestamp falls into the so-called "stasis period" at the end of the
// lifetime of the lease, or if the request timestamp is beyond the end of the
// lifetime of the lease, the status is UNUSABLE. Callers typically want to
// react to an UNUSABLE lease status by extending the lease, if they are in a
// position to do so.
//
// Finally, for requests timestamps falling before the stasis period of a lease
// that is not EXPIRED and also not PROSCRIBED, the status is VALID.
//
// # Implementation Note
//
// On the surface, it might seem like we could easily abandon the lease stasis
// concept in favor of consulting a request's uncertainty interval. We would
// then define a request's timestamp as the maximum of its read_timestamp and
// its global_uncertainty_limit, and simply check whether this timestamp falls
// below a lease's expiration. This could allow certain transactional requests
// to operate more closely to a lease's expiration. But not all requests that
// expect linearizability use an uncertainty interval (e.g. non-transactional
// requests), and so the lease stasis period serves as a kind of catch-all
// uncertainty interval for non-transactional and admin requests.
//
// Without that stasis period, the following linearizability violation could
// occur for two non-transactional requests operating on a single register
// during a lease change:
//
//   - a range lease gets committed on the new lease holder (but not the old).
//   - client proposes and commits a write on new lease holder (with a timestamp
//     just greater than the expiration of the old lease).
//   - client tries to read what it wrote, but hits a slow coordinator (which
//     assigns a timestamp covered by the old lease).
//   - the read is served by the old lease holder (which has not processed the
//     change in lease holdership).
//   - the client fails to read their own write.
func Status(ctx context.Context, nl NodeLiveness, i StatusInput) kvserverpb.LeaseStatus {
	lease := i.Lease
	status := kvserverpb.LeaseStatus{
		Lease: lease,
		// NOTE: it would not be correct to accept either only the request time
		// or only the current time in this method, we need both. We need the
		// request time to determine whether the current lease can serve a given
		// request, even if that request has a timestamp in the future of
		// present time. We need the current time to distinguish between an
		// EXPIRED lease and an UNUSABLE lease. Only an EXPIRED lease can change
		// hands through a lease acquisition.
		Now:                       i.Now,
		RequestTime:               i.RequestTs,
		MinValidObservedTimestamp: i.MinValidObservedTs,
	}
	var expiration hlc.Timestamp
	if lease.Type() == roachpb.LeaseExpiration {
		expiration = lease.GetExpiration()
	} else {
		l, ok := nl.GetLiveness(lease.Replica.NodeID)
		status.Liveness = l.Liveness
		if !ok || status.Liveness.Epoch < lease.Epoch {
			// If lease validity can't be determined (e.g. gossip is down
			// and liveness info isn't available for owner), we can neither
			// use the lease nor do we want to attempt to acquire it.
			var msg redact.StringBuilder
			if !ok {
				msg.Printf("can't determine lease status of %s due to node liveness error: %v",
					lease.Replica, liveness.ErrRecordCacheMiss)
			} else {
				msg.Printf("can't determine lease status of %s because node liveness info for n%d is stale. lease: %s, liveness: %s",
					lease.Replica, lease.Replica.NodeID, lease, l.Liveness)
			}
			if leaseStatusLogLimiter.ShouldLog() {
				log.Infof(ctx, "%s", msg)
			}
			status.State = kvserverpb.LeaseState_ERROR
			status.ErrInfo = msg.String()
			return status
		}
		if status.Liveness.Epoch > lease.Epoch {
			status.State = kvserverpb.LeaseState_EXPIRED
			return status
		}
		expiration = status.Liveness.Expiration.ToTimestamp()
	}
	maxOffset := i.MaxOffset
	stasis := expiration.Add(-int64(maxOffset), 0)
	ownedLocally := lease.OwnedBy(i.LocalStoreID)
	// NB: the order of these checks is important, and goes from stronger to
	// weaker reasons why the lease may be considered invalid. For example,
	// EXPIRED or PROSCRIBED must take precedence over UNUSABLE, because some
	// callers consider UNUSABLE as valid. For an example issue that this ordering
	// may cause, see https://github.com/cockroachdb/cockroach/issues/100101.
	if expiration.LessEq(i.Now.ToTimestamp()) {
		status.State = kvserverpb.LeaseState_EXPIRED
	} else if ownedLocally && lease.ProposedTS.Less(i.MinProposedTs) {
		// If the replica owns the lease, additional verify that the lease's
		// proposed timestamp is not earlier than the min proposed timestamp.
		status.State = kvserverpb.LeaseState_PROSCRIBED
	} else if stasis.LessEq(i.RequestTs) {
		status.State = kvserverpb.LeaseState_UNUSABLE
	} else {
		status.State = kvserverpb.LeaseState_VALID
	}
	return status
}

var leaseStatusLogLimiter = func() *log.EveryN {
	e := log.Every(15 * time.Second)
	e.ShouldLog() // waste the first shot
	return &e
}()
