#!/usr/bin/env bash

mkdir -p y/bsl-not-slash

rg -n 'Use of this software is governed by the Business Source License' |
  rg -v '^y/|^search' |
  sort > y/bsl-not-slash/step1.txt

#rg -n --multiline '// Licensed as a CockroachDB Enterprise file under the Cockroach Community\
#// License \(the "License"\); you may not use this file except in compliance with\
#// the License. You may obtain a copy of the License at\
#//\
#//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt' |
  #rg 'Licensed as a CockroachDB Enterprise file under the Cockroach Community' |
  #rg ':3://' |
  #sort > y/ccl-line3/step1.txt

  sed -e 's/:[0-9][0-9]*:\/\/.*$//' -e 's/:[0-9][0-9]*:#.*$//' y/bsl-not-slash/step1.txt | sort | uniq > y/bsl-not-slash/step2.txt

  sed -e 's/:[0-9][0-9]*:\/\/.*$//' -e 's/:[0-9][0-9]*:#.*$//' y/bsl-not-slash/step1.txt | sort | uniq -c | sort > y/bsl-not-slash/count-per-file.txt

  sed -e 's/\([^:]*\):\([0-9][0-9]*\):\/\/.*$/\2 \1/' -e 's/\([^:]*\):\([0-9][0-9]*\):#.*$/\2 \1/' y/bsl-not-slash/step1.txt | sort > y/bsl-not-slash/by-line-number.txt


