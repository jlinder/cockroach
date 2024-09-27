#!/usr/bin/env bash

BASEDIR="y/add-make"

mkdir -p ${BASEDIR}

rg --type=make --files-without-match -e 'Use of this software is governed by the CockroachDB Software License' |
  rg -v 'c-deps' |
  sort > ${BASEDIR}/step1.txt

truncate -s 0 ${BASEDIR}/head.txt

cat ${BASEDIR}/step1.txt | while read -r FILENAME; do
  head -2 "${FILENAME}" >> ${BASEDIR}/head.txt
  echo >> ${BASEDIR}/head.txt
done;


#rg -n --multiline '// Licensed as a CockroachDB Enterprise file under the Cockroach Community\
#// License \(the "License"\); you may not use this file except in compliance with\
#// the License. You may obtain a copy of the License at\
#//\
#//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt' |
  #rg 'Licensed as a CockroachDB Enterprise file under the Cockroach Community' |
  #rg ':3://' |
  #sort > y/ccl-line3/step1.txt

#gsed -e 's/:[0-9][0-9]*:\/\/.*$//' y/bsl-not-slash/step1.txt | sort | uniq > y/bsl-not-slash/step2.txt

#gsed -e 's/:[0-9][0-9]*:\/\/.*$//' y/bsl-not-slash/step1.txt | sort | uniq -c | sort > y/bsl-not-slash/count-per-file.txt

#sed -e 's/\([^:]*\):\([0-9][0-9]*\):\/\/.*$/\2 \1/' y/bsl-not-slash/step1.txt | sort > y/bsl-not-slash/by-line-number.txt


