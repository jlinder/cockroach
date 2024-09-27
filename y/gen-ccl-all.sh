#!/usr/bin/env bash

BASEDIR=y/ccl-all
mkdir -p ${BASEDIR}

rg -n --multiline '// Licensed as a CockroachDB Enterprise file under the Cockroach Community\
// License \(the "License"\); you may not use this file except in compliance with\
// the License. You may obtain a copy of the License at\
//\
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt' |
  rg 'Licensed as a CockroachDB Enterprise file under the Cockroach Community' |
  sort > ${BASEDIR}/step1.txt

sed -e 's/:3:\/\/.*$//' ${BASEDIR}/step1.txt | sort | uniq > ${BASEDIR}/step2.txt



