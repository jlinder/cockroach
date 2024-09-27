#!/usr/bin/env bash

mkdir -p y/bsl-not-line3

rg -n --multiline '// Use of this software is governed by the Business Source License\
// included in the file licenses/BSL.txt.\
//\
// As of the Change Date specified in that file, in accordance with\
// the Business Source License, use of this software will be governed\
// by the Apache License, Version 2.0, included in the file\
// licenses/APL.txt.' |
  rg 'governed by the Business Source License' |
  rg -v ':3://' |
  sort > y/bsl-not-line3/step1.txt

gsed -e 's/:[0-9][0-9]*:\/\/.*$//' y/bsl-not-line3/step1.txt | sort | uniq > y/bsl-not-line3/step2.txt

gsed -e 's/:[0-9][0-9]*:\/\/.*$//' y/bsl-not-line3/step1.txt | sort | uniq -c | sort > y/bsl-not-line3/count-per-file.txt

sed -e 's/\([^:]*\):\([0-9][0-9]*\):\/\/.*$/\2 \1/' y/bsl-not-line3/step1.txt | sort > y/bsl-not-line3/by-line-number.txt



