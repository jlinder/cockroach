#!/usr/bin/env bash

mkdir -p y/bsl-all

rg -n --multiline '// Use of this software is governed by the Business Source License\
// included in the file licenses/BSL.txt.\
//\
// As of the Change Date specified in that file, in accordance with\
// the Business Source License, use of this software will be governed\
// by the Apache License, Version 2.0, included in the file\
// licenses/APL.txt.' |
  rg 'governed by the Business Source License' |
  sort > y/bsl-all/step1.txt

sed -e 's/:[0-9][0-9]*:\/\/.*$//' y/bsl-all/step1.txt | sort | uniq > y/bsl-all/step2.txt

sed -e 's/:[0-9][0-9]*:\/\/.*$//' y/bsl-all/step1.txt | sort | uniq -c | sort > y/bsl-all/count-per-file.txt



