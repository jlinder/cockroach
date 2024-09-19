// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

import Long from "long";
import * as protobuf from "protobufjs/minimal";

protobuf.util.Long = Long;
protobuf.configure();
