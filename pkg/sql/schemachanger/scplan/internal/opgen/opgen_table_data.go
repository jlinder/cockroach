// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package opgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

func init() {
	opRegistry.register((*scpb.TableData)(nil),
		toPublic(scpb.Status_ABSENT,
			equiv(scpb.Status_DROPPED),
			to(scpb.Status_PUBLIC)),
		toAbsent(
			scpb.Status_PUBLIC,
			to(scpb.Status_DROPPED,
				revertible(false),
				emit(func(this *scpb.TableData, md *opGenContext) *scop.CreateGCJobForTable {
					return &scop.CreateGCJobForTable{
						TableID:             this.TableID,
						DatabaseID:          this.DatabaseID,
						StatementForDropJob: statementForDropJob(this, md),
					}
				}),
			),
			to(scpb.Status_ABSENT),
		),
	)
}
