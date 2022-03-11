// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseNoteRequiresRef(t *testing.T) {
	testCases := []struct {
		message  string
		ci       commitInfo
		expected bool
	}{
		{
			message: "nil release note categories",
			ci: commitInfo{},
			expected: false,
		},
		{
			message: "empty release note categories map",
			ci: commitInfo{releaseNoteCategories: map[string]int{}},
			expected: false,
		},
		{
			message: "release note category 'None'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"None": 1}},
			expected: false,
		},
		{
			message: "nil release note category 'bugfix'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"bugfix": 1}},
			expected: false,
		},
		{
			message: "nil release note category 'bug fix'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"bug fix": 1}},
			expected: false,
		},
		{
			message: "nil release note category 'general change'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"general change": 1}},
			expected: true,
		},
		{
			message: "nil release note category is blank",
			ci: commitInfo{releaseNoteCategories: map[string]int{"": 1}},
			expected: true,
		},
		{
			message: "nil release note category 'anything else'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"anything else": 1}},
			expected: true,
		},
		{
			message: "nil release note category 'sql change'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"bug fix": 1, "sql change": 1}},
			expected: true,
		},
		{
			message: "nil release note category 'security'",
			ci: commitInfo{releaseNoteCategories: map[string]int{"security": 1, "bugfix": 1}},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			result := tc.ci.releaseNoteRequiresRef()
			assert.Equal(t, tc.expected, result)
		})
	}
}


func TestExtractFixIssuesIDs(t *testing.T) {
	testCases := []struct {
		message  string
		expected map[string]int
	}{
		{
			message: `workload: Fix folder name generation.

Fixes #75200 #98922
close #75201
closed #592
RESOLVE #5555

Release Notes: None`,
			expected: map[string]int{"#75200": 1, "#98922": 1, "#75201": 1, "#592": 1, "#5555": 1},
		},
		{
			message: `logictestccl: skip flaky TestCCLLogic/fakedist-metadata/partitioning_enum

Informs #75227
Epic CRDB-491

Release note (bug fix): Fixin the bug`,
			expected: map[string]int{},
		},
		{
			message: `lots of variations

fixes #74932; we were reading from the replicas map without...
Closes #74889.
Resolves #74482, #74784  #65117   #79299.
Fix:  #73834
epic: CRDB-9234.
epic CRDB-235, CRDB-40192;   DEVINF-392
Fixed:  #29833 example/repo#941
see also:  #9243
informs: #912,   #4729   #2911  cockroachdb/cockroach#2934

Release note (sql change): Import now checks readability...`,
			expected: map[string]int{"#74932": 1, "#74889": 1, "#74482": 1, "#74784": 1, "#65117": 1, "#79299": 1, "#73834": 1, "#29833": 1, "example/repo#941": 1},
		},
		{
			message: `lots of variations 2

Resolved: #4921
This comes w/ support for Applie Silicon Macs. Closes #72829.
This doesn't completely fix #71901 in that some...
      Fixes #491
Resolves #71614, Resolves #71607
Thereby, fixes #69765
Informs #69765 (point 4).
Fixes #65200. The last remaining 21.1 version (V21_1) can be removed as

Release note (build change): Upgrade to Go 1.17.6
Release note (ops change): Added a metric
Release notes (enterprise change): Client certificates may...`,
			expected: map[string]int{"#4921": 1, "#72829": 1, "#71901": 1, "#491": 1, "#71614": 1, "#71607": 1, "#69765": 1, "#65200": 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			result := extractFixIssueIDs(tc.message)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractInformIssueIDs(t *testing.T) {
	testCases := []struct {
		message  string
		expected map[string]int
	}{
		{
			message: `logictestccl: skip flaky TestCCLLogic/fakedist-metadata/partitioning_enum

Informs #75227
Epic CRDB-491
Fix:  #73834

Release note (bug fix): Fixin the bug`,
			expected: map[string]int{"#75227": 1},
		},
		{
			message: `lots of variations

Fixed:  #29833 example/repo#941
see also:  #9243
informs: #912,   #4729   #2911  cockroachdb/cockroach#2934
Informs #69765 (point 4).
This informs #59293 with these additions:

Release note (sql change): Import now checks readability...`,
			expected: map[string]int{"#9243": 1, "#912": 1, "#4729": 1, "#2911": 1, "cockroachdb/cockroach#2934": 1, "#69765": 1, "#59293": 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			result := extractInformIssueIDs(tc.message)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractEpicIDs(t *testing.T) {
	testCases := []struct {
		message  string
		expected map[string]int
	}{
		{
			message: `logictestccl: skip flaky TestCCLLogic/fakedist-metadata/partitioning_enum

Informs #75227
Epic CRDB-491
Fix:  #73834

Release note (bug fix): Fixin the bug`,
			expected: map[string]int{"CRDB-491": 1},
		},
		{
			message: `lots of variations

epic: CRDB-9234.
epic CRDB-235, CRDB-40192;   DEVINF-392

Release note (sql change): Import now checks readability...`,
			expected: map[string]int{"CRDB-9234": 1, "CRDB-235": 1, "CRDB-40192": 1, "DEVINF-392": 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			result := extractEpicIDs(tc.message)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractReleaseNoteCategories(t *testing.T) {
	testCases := []struct {
		message  string
		expected map[string]int
	}{
		{
			message: `workload: Fix folder name generation.

Fixes #75200 #98922

Release Notes: None`,
			expected: map[string]int{"None": 1},
		},
		{
			message: `lots of variations

Fix:  #73834
epic: CRDB-9234.
see also:  #9243

Release note (sql change): Import now checks readability...`,
			expected: map[string]int{"sql change": 1},
		},
		{
			message: `lots of variations 2

This doesn't completely fix #71901 in that some...

Release note (bug fix): Fixin the bug
Release note (build change): Upgrade to Go 1.17.6
Release note (ops change): Added a metric
Release notes (enterprise change): Client certificates may...`,
			expected: map[string]int{"bug fix": 1, "build change": 1, "ops change": 1, "enterprise change": 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			result := extractReleaseNoteCategories(tc.message)
			assert.Equal(t, tc.expected, result)
		})
	}
}
