#!/usr/bin/env bash


sed -i -e ':a;N;$!ba;s/\/\/ Use of this software is governed by the Business Source License\
\/\/ included in the file licenses\/BSL.txt.\
\/\/\
\/\/ As of the Change Date specified in that file, in accordance with\
\/\/ the Business Source License, use of this software will be governed\
\/\/ by the Apache License, Version 2.0, included in the file\
\/\/ licenses\/APL.txt./\/\/ Use of this software is governed by the CockroachDB Software License\
\/\/ included in the \/LICENSE file./' pkg/ui/workspaces/db-console/scripts/topo.js




