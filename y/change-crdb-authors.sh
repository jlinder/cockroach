#!/usr/bin/env bash

BASEDIR="y/crdb-authors"

cat ${BASEDIR}/step2.txt | while read -r FILENAME; do

  gsed -i -e 's/Cockroach Labs, Inc./The Cockroach Authors./' "${FILENAME}"

done



