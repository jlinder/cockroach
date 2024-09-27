#!/usr/bin/env bash

BASEDIR="y/add-py"

cat ${BASEDIR}/files-not-executable.txt | while read -r FILENAME; do

  DATE="$(git log --graph --date=format:'%Y' --pretty=format:'%ad' "${FILENAME}" |
    tail -1 |
    sed -e 's/^..//')"

  sed -i -e '1s/^/#\
# Use of this software is governed by the CockroachDB Software License\
# included in the \/LICENSE file.\
\
/' "${FILENAME}"

  sed -i -e "1s/^/# Copyright $(echo -n ${DATE}) The Cockroach Authors.\n/" "${FILENAME}"

done


cat ${BASEDIR}/files-executable.txt | while read -r FILENAME; do

  DATE="$(git log --graph --date=format:'%Y' --pretty=format:'%ad' "${FILENAME}" |
    tail -1 |
    sed -e 's/^..//')"

  sed -i -e '2s/^/\
#\
# Use of this software is governed by the CockroachDB Software License\
# included in the \/LICENSE file.\
\
/' "${FILENAME}"

  sed -i -e "3s/^/# Copyright $(echo -n ${DATE}) The Cockroach Authors.\n/" "${FILENAME}"

done


