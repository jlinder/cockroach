#!/usr/bin/env bash

cat y/ccl-line3/step2.txt | while read -r FILENAME; do

  sed -i -e ':a;N;$!ba;s/\/\/ Licensed as a CockroachDB Enterprise file under the Cockroach Community\
\/\/ License (the "License"); you may not use this file except in compliance with\
\/\/ the License. You may obtain a copy of the License at\
\/\/\
\/\/     https:\/\/github.com\/cockroachdb\/cockroach\/blob\/master\/licenses\/CCL.txt/\/\/ Use of this software is governed by the CockroachDB Software License\
\/\/ included in the \/LICENSE file./' "${FILENAME}"

done



