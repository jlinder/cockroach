#!/usr/bin/env bash

set -euo pipefail

PR_NUMBER="$(echo "$GITHUB_REF" | awk -F '/' '{print $3}')"

# TODO(jamesl): figure out how to get github creds for querying the github API
bin/bazel run //pkg/testutils/lint-epic-issue-refs:lint-epic-issue-refs "$PR_NUMBER"
