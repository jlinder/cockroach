#!/usr/bin/env bash

mkdir -p y/ccl-not-slash

rg -n 'Licensed as a CockroachDB Enterprise file under the Cockroach Community' |
  rg -v '^y/|^search' |
  sort > y/ccl-not-slash/step1.txt

#rg -n --multiline '// Licensed as a CockroachDB Enterprise file under the Cockroach Community\
#// License \(the "License"\); you may not use this file except in compliance with\
#// the License. You may obtain a copy of the License at\
#//\
#//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt' |
  #rg 'Licensed as a CockroachDB Enterprise file under the Cockroach Community' |
  #rg ':3://' |
  #sort > y/ccl-line3/step1.txt

#sed -e 's/:3:\/\/.*$//' y/ccl-line3/step1.txt | sort | uniq > y/ccl-line3/step2.txt



