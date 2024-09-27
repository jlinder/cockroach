#!/usr/bin/env bash

BASEDIR="y/add-py"

mkdir -p ${BASEDIR}

rg --type=py --files-without-match -e 'Use of this software is governed by the CockroachDB Software License' |
  rg -v 'c-deps' |
  sort > ${BASEDIR}/step1.txt

truncate -s 0 ${BASEDIR}/head.txt

cat ${BASEDIR}/step1.txt | while read -r FILENAME; do
  head -2 "${FILENAME}" >> ${BASEDIR}/head.txt
  echo >> ${BASEDIR}/head.txt
done;

cat ${BASEDIR}/step1.txt | xargs -n 50 rg --files-without-match -e '#! ?/usr/bin/env python' > ${BASEDIR}/files-not-executable.txt
cat ${BASEDIR}/step1.txt | xargs -n 50 rg --l -e '#! ?/usr/bin/env python' > ${BASEDIR}/files-executable.txt

