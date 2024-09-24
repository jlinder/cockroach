#!/usr/bin/env bash

# Copyright 2022 The Cockroach Authors.
#
# Use of this software is governed by the CockroachDB Software License
# included in the /LICENSE file.

set -euo pipefail

# assignment_casts_gen.sh generates a CSV file of test cases for use by
# TestAssignmentCastsMatchPostgres, based on the files 'literals.txt' and
# 'types.txt'. To use this script, Postgres must be installed locally with the
# PostGIS extension and must already be running.
#
# Usage:
#   ./assignment_casts_gen.sh > assignment_casts.csv

pgversion=$(psql -AXqtc "SELECT substring(version(), 'PostgreSQL (\d+\.\d+)')")

echo "# Testcases for TestAssignmentCastsMatchPostgres."
echo "#"
echo "# Results captured from PostgreSQL ${pgversion}."
echo "#"
echo "# This file was automatically generated by assignment_casts_gen.sh from the"
echo "# contents of 'literals.txt' and 'types.txt'. To skip a testcase please add it"
echo "# to assignment_casts_skip.csv rather than commenting it out here."
echo "literal,type,expect"
while read -r type; do
  psql -Xqc "CREATE TABLE assignment_casts (val $type)"
  while read -r literal; do
    # Quote literal and type in case they contain quotes or commas.
    printf '"%s","%s",' "${literal//\"/\"\"}" "${type//\"/\"\"}"
    psql --csv -Xqtc "INSERT INTO assignment_casts VALUES ($literal) RETURNING quote_nullable(val)" 2>/dev/null || echo 'error'
  done <literals.txt
  psql -Xqc "DROP TABLE IF EXISTS assignment_casts" 2>/dev/null
done <types.txt
