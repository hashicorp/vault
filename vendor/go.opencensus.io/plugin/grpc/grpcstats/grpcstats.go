// Copyright 2017, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package grpcstats provides OpenCensus stats support for gRPC clients and servers.
package grpcstats // import "go.opencensus.io/plugin/grpc/grpcstats"

import (
	"log"
	"strings"
	"time"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type grpcInstrumentationKey string

// rpcData holds the instrumentation RPC data that is needed between the start
// and end of an call. It holds the info that this package needs to keep track
// of between the various GRPC events.
type rpcData struct {
	// startTime represents the time at which TagRPC was invoked at the
	// beginning of an RPC. It is an appoximation of the time when the
	// application code invoked GRPC code.
	startTime           time.Time
	reqCount, respCount int64 // access atomically
}

// The following variables define the default hard-coded auxiliary data used by
// both the default GRPC client and GRPC server metrics.
// These are Go objects instances mirroring the some of the proto definitions
// found at "github.com/google/instrumentation-proto/census.proto".
// A complete description of each can be found there.
// TODO(acetechnologist): This is temporary and will need to be replaced by a
// mechanism to load these defaults from a common repository/config shared by
// all supported languages. Likely a serialized protobuf of these defaults.
var (
	unitByte             = "By"
	unitCount            = "1"
	unitMillisecond      = "ms"
	slidingTimeSubuckets = 6

	rpcBytesBucketBoundaries  = []float64{0, 1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296}
	rpcMillisBucketBoundaries = []float64{0, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000}
	rpcCountBucketBoundaries  = []float64{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

	aggCount      = view.CountAggregation{}
	aggMean       = view.MeanAggregation{}
	aggDistBytes  = view.DistributionAggregation(rpcBytesBucketBoundaries)
	aggDistMillis = view.DistributionAggregation(rpcMillisBucketBoundaries)
	aggDistCounts = view.DistributionAggregation(rpcCountBucketBoundaries)

	KeyMethod tag.Key
	KeyStatus tag.Key
)

func init() {
	var err error
	if KeyMethod, err = tag.NewKey("method"); err != nil {
		log.Fatalf("Cannot create method key: %v", err)
	}
	if KeyStatus, err = tag.NewKey("canonical_status"); err != nil {
		log.Fatalf("Cannot create canonical_status key: %v", err)
	}
	initServer()
	initClient()
}

var (
	grpcServerConnKey = grpcInstrumentationKey("server-conn")
	grpcServerRPCKey  = grpcInstrumentationKey("server-rpc")
	grpcClientRPCKey  = grpcInstrumentationKey("client-rpc")
)

func methodName(fullname string) string {
	return strings.TrimLeft(fullname, "/")
}
