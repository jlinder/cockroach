#!/usr/bin/env bash

mkdir -p y/ccl-not-line3

rg -n --multiline '// Licensed as a CockroachDB Enterprise file under the Cockroach Community\
// License \(the "License"\); you may not use this file except in compliance with\
// the License. You may obtain a copy of the License at\
//\
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt' |
  rg 'Licensed as a CockroachDB Enterprise file under the Cockroach Community' |
  rg -v ':3://' |
  sort > y/ccl-not-line3/step1.txt

gsed -e 's/:[0-9][0-9]*:\/\/.*$//' y/ccl-not-line3/step1.txt | sort | uniq > y/ccl-not-line3/step2.txt



