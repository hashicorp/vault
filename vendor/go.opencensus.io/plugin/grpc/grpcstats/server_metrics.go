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

package grpcstats

import (
	"log"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following variables are measures and views made available for gRPC clients.
// Server needs to use a ServerStatsHandler in order to enable collection.
var (
	// Available server measures
	RPCServerErrorCount        *stats.Int64Measure
	RPCServerServerElapsedTime *stats.Float64Measure
	RPCServerRequestBytes      *stats.Int64Measure
	RPCServerResponseBytes     *stats.Int64Measure
	RPCServerStartedCount      *stats.Int64Measure
	RPCServerFinishedCount     *stats.Int64Measure
	RPCServerRequestCount      *stats.Int64Measure
	RPCServerResponseCount     *stats.Int64Measure

	// Predefined server views
	RPCServerErrorCountView        *view.View
	RPCServerServerElapsedTimeView *view.View
	RPCServerRequestBytesView      *view.View
	RPCServerResponseBytesView     *view.View
	RPCServerRequestCountView      *view.View
	RPCServerResponseCountView     *view.View
)

// TODO(acetechnologist): This is temporary and will need to be replaced by a
// mechanism to load these defaults from a common repository/config shared by
// all supported languages. Likely a serialized protobuf of these defaults.

func defaultServerMeasures() {
	var err error

	if RPCServerErrorCount, err = stats.Int64("grpc.io/server/error_count", "RPC Errors", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/error_count: %v", err)
	}
	if RPCServerServerElapsedTime, err = stats.Float64("grpc.io/server/server_elapsed_time", "Server elapsed time in msecs", unitMillisecond); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/server_elapsed_time: %v", err)
	}
	if RPCServerRequestBytes, err = stats.Int64("grpc.io/server/request_bytes", "Request bytes", unitByte); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/request_bytes: %v", err)
	}
	if RPCServerResponseBytes, err = stats.Int64("grpc.io/server/response_bytes", "Response bytes", unitByte); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/response_bytes: %v", err)
	}
	if RPCServerStartedCount, err = stats.Int64("grpc.io/server/started_count", "Number of server RPCs (streams) started", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/started_count: %v", err)
	}
	if RPCServerFinishedCount, err = stats.Int64("grpc.io/server/finished_count", "Number of server RPCs (streams) finished", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/finished_count: %v", err)
	}
	if RPCServerRequestCount, err = stats.Int64("grpc.io/server/request_count", "Number of server RPC request messages", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/request_count: %v", err)
	}
	if RPCServerResponseCount, err = stats.Int64("grpc.io/server/response_count", "Number of server RPC response messages", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/server/response_count: %v", err)
	}
}

func defaultServerViews() {
	RPCServerErrorCountView, _ = view.New(
		"grpc.io/server/error_count/cumulative",
		"RPC Errors",
		[]tag.Key{KeyMethod, KeyStatus},
		RPCServerErrorCount,
		aggCount)
	RPCServerServerElapsedTimeView, _ = view.New(
		"grpc.io/server/server_elapsed_time/cumulative",
		"Server elapsed time in msecs",
		[]tag.Key{KeyMethod},
		RPCServerServerElapsedTime,
		aggDistMillis)
	RPCServerRequestBytesView, _ = view.New(
		"grpc.io/server/request_bytes/cumulative",
		"Request bytes",
		[]tag.Key{KeyMethod},
		RPCServerRequestBytes,
		aggDistBytes)
	RPCServerResponseBytesView, _ = view.New(
		"grpc.io/server/response_bytes/cumulative",
		"Response bytes",
		[]tag.Key{KeyMethod},
		RPCServerResponseBytes,
		aggDistBytes)
	RPCServerRequestCountView, _ = view.New(
		"grpc.io/server/request_count/cumulative",
		"Count of request messages per server RPC",
		[]tag.Key{KeyMethod},
		RPCServerRequestCount,
		aggDistCounts)
	RPCServerResponseCountView, _ = view.New(
		"grpc.io/server/response_count/cumulative",
		"Count of response messages per server RPC",
		[]tag.Key{KeyMethod},
		RPCServerResponseCount,
		aggDistCounts)

	serverViews = append(serverViews,
		RPCServerErrorCountView,
		RPCServerServerElapsedTimeView,
		RPCServerRequestBytesView,
		RPCServerResponseBytesView,
		RPCServerRequestCountView,
		RPCServerResponseCountView)

	// TODO(jbd): Add roundtrip_latency, uncompressed_request_bytes, uncompressed_response_bytes, request_count, response_count.
}

func initServer() {
	defaultServerMeasures()
	defaultServerViews()
}

var serverViews []*view.View
