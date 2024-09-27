#!/usr/bin/env bash

BASEDIR="y/crdb-authors"
mkdir -p ${BASEDIR}

rg -n --multiline '// This code has been modified from its original form by Cockroach Labs, Inc.\
// All modifications are Copyright 20\d\d Cockroach Labs, Inc.' |
  rg 'This code has been modified from its original form by Cockroach' |
  sort > ${BASEDIR}/step1.txt

gsed -e 's/:[0-9][0-9]*:\/\/.*$//' ${BASEDIR}/step1.txt | sort | uniq > ${BASEDIR}/step2.txt

sed -e 's/\([^:]*\):\([0-9][0-9]*\):\/\/.*$/\2 \1/' ${BASEDIR}/step1.txt | sort > ${BASEDIR}/by-line-number.txt


