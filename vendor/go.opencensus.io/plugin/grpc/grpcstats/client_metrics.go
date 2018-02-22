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
// Client connection needs to use a ClientStatsHandler in order to enable collection.
var (
	// Available client measures
	RPCClientErrorCount       *stats.Int64Measure
	RPCClientRoundTripLatency *stats.Float64Measure
	RPCClientRequestBytes     *stats.Int64Measure
	RPCClientResponseBytes    *stats.Int64Measure
	RPCClientStartedCount     *stats.Int64Measure
	RPCClientFinishedCount    *stats.Int64Measure
	RPCClientRequestCount     *stats.Int64Measure
	RPCClientResponseCount    *stats.Int64Measure

	// Predefined client views
	RPCClientErrorCountView       *view.View
	RPCClientRoundTripLatencyView *view.View
	RPCClientRequestBytesView     *view.View
	RPCClientResponseBytesView    *view.View
	RPCClientRequestCountView     *view.View
	RPCClientResponseCountView    *view.View
)

// TODO(acetechnologist): This is temporary and will need to be replaced by a
// mechanism to load these defaults from a common repository/config shared by
// all supported languages. Likely a serialized protobuf of these defaults.

func defaultClientMeasures() {
	var err error

	// Creating client measures
	if RPCClientErrorCount, err = stats.Int64("grpc.io/client/error_count", "RPC Errors", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/error_count: %v", err)
	}
	if RPCClientRoundTripLatency, err = stats.Float64("grpc.io/client/roundtrip_latency", "RPC roundtrip latency in msecs", unitMillisecond); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/roundtrip_latency: %v", err)
	}
	if RPCClientRequestBytes, err = stats.Int64("grpc.io/client/request_bytes", "Request bytes", unitByte); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/request_bytes: %v", err)
	}
	if RPCClientResponseBytes, err = stats.Int64("grpc.io/client/response_bytes", "Response bytes", unitByte); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/response_bytes: %v", err)
	}
	if RPCClientStartedCount, err = stats.Int64("grpc.io/client/started_count", "Number of client RPCs (streams) started", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/started_count: %v", err)
	}
	if RPCClientFinishedCount, err = stats.Int64("grpc.io/client/finished_count", "Number of client RPCs (streams) finished", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/finished_count: %v", err)
	}
	if RPCClientRequestCount, err = stats.Int64("grpc.io/client/request_count", "Number of client RPC request messages", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/request_count: %v", err)
	}
	if RPCClientResponseCount, err = stats.Int64("grpc.io/client/response_count", "Number of client RPC response messages", unitCount); err != nil {
		log.Fatalf("Cannot create measure grpc.io/client/response_count: %v", err)
	}
}

func defaultClientViews() {
	RPCClientErrorCountView, _ = view.New(
		"grpc.io/client/error_count/cumulative",
		"RPC Errors",
		[]tag.Key{KeyStatus, KeyMethod},
		RPCClientErrorCount,
		aggMean)
	RPCClientRoundTripLatencyView, _ = view.New(
		"grpc.io/client/roundtrip_latency/cumulative",
		"Latency in msecs",
		[]tag.Key{KeyMethod},
		RPCClientRoundTripLatency,
		aggDistMillis)
	RPCClientRequestBytesView, _ = view.New(
		"grpc.io/client/request_bytes/cumulative",
		"Request bytes",
		[]tag.Key{KeyMethod},
		RPCClientRequestBytes,
		aggDistBytes)
	RPCClientResponseBytesView, _ = view.New(
		"grpc.io/client/response_bytes/cumulative",
		"Response bytes",
		[]tag.Key{KeyMethod},
		RPCClientResponseBytes,
		aggDistBytes)
	RPCClientRequestCountView, _ = view.New(
		"grpc.io/client/request_count/cumulative",
		"Count of request messages per client RPC",
		[]tag.Key{KeyMethod},
		RPCClientRequestCount,
		aggDistCounts)
	RPCClientResponseCountView, _ = view.New(
		"grpc.io/client/response_count/cumulative",
		"Count of response messages per client RPC",
		[]tag.Key{KeyMethod},
		RPCClientResponseCount,
		aggDistCounts)

	clientViews = append(clientViews,
		RPCClientErrorCountView,
		RPCClientRoundTripLatencyView,
		RPCClientRequestBytesView,
		RPCClientResponseBytesView,
		RPCClientRequestCountView,
		RPCClientResponseCountView,
	)
	// TODO(jbd): Add roundtrip_latency, uncompressed_request_bytes, uncompressed_response_bytes, request_count, response_count.
}

// initClient registers the default metrics (measures and views)
// for a GRPC client.
func initClient() {
	defaultClientMeasures()
	defaultClientViews()
}

var clientViews []*view.View
