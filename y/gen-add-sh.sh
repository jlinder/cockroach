#!/usr/bin/env bash

BASEDIR="y/add-sh"

mkdir -p ${BASEDIR}

rg --type=sh --files-without-match -e 'Use of this software is governed by the CockroachDB Software License' |
  rg -v 'c-deps|^y/' |
  sort > ${BASEDIR}/step1.txt

sed -e 's/^.*\(\.[^.]*\)$/\1/' ${BASEDIR}/step1.txt | sort | uniq -c > ${BASEDIR}/file-types.txt

truncate -s 0 ${BASEDIR}/head.txt
cat ${BASEDIR}/step1.txt | while read -r FILENAME; do
  head -1 "${FILENAME}" >> ${BASEDIR}/head.txt
done;

truncate -s 0 ${BASEDIR}/head-2.txt
cat ${BASEDIR}/step1.txt | while read -r FILENAME; do
  head -2 "${FILENAME}" >> ${BASEDIR}/head-2.txt
  echo >> ${BASEDIR}/head-2.txt
done;

sort ${BASEDIR}/head.txt | uniq -c > ${BASEDIR}/first-lines.txt


cat ${BASEDIR}/step1.txt |
  xargs -n 50 rg --files-without-match -e '^(#! ?/bin/(ba)?sh|#!/usr/bin/env (ba|z)sh)' > ${BASEDIR}/files-not-executable.txt

cat ${BASEDIR}/step1.txt |
  xargs -n 50 rg --files-with-matches -e '^(#! ?/bin/(ba)?sh|#!/usr/bin/env (ba|z)sh)' > ${BASEDIR}/files-executable.txt


