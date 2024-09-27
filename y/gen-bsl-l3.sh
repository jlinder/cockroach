#!/usr/bin/env bash

mkdir -p y/bsl-line3

rg -n --multiline '// Use of this software is governed by the Business Source License\
// included in the file licenses/BSL.txt.\
//\
// As of the Change Date specified in that file, in accordance with\
// the Business Source License, use of this software will be governed\
// by the Apache License, Version 2.0, included in the file\
// licenses/APL.txt.' |
  rg 'governed by the Business Source License' |
  rg ':3://' |
  sort > y/bsl-line3/step1.txt

sed -e 's/:3:\/\/.*$//' y/bsl-line3/step1.txt | sort | uniq > y/bsl-line3/step2.txt



