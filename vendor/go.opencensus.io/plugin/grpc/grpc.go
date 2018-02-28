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

// Package grpc contains OpenCensus stats and trace
// integrations with gRPC.
package grpc

import (
	"golang.org/x/net/context"

	"go.opencensus.io/plugin/grpc/grpcstats"
	"go.opencensus.io/plugin/grpc/grpctrace"

	"google.golang.org/grpc/stats"
)

// NewClientStatsHandler enables OpenCensus stats and trace
// for gRPC clients. If these features need to be indiviually turned
// on, see grpcstats and grpctrace packages.
func NewClientStatsHandler() stats.Handler {
	return handler{
		grpcstats.NewClientStatsHandler(),
		grpctrace.NewClientStatsHandler(),
	}
}

// NewServerStatsHandler enables OpenCensus stats and trace
// for gRPC servers. If these features need to be indiviually turned
// on, see grpcstats and grpctrace packages.
func NewServerStatsHandler() stats.Handler {
	return handler{
		grpcstats.NewServerStatsHandler(),
		grpctrace.NewServerStatsHandler(),
	}
}

type handler []stats.Handler

func (h handler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	for _, hh := range h {
		hh.HandleConn(ctx, cs)
	}
}

func (h handler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	for _, hh := range h {
		hh.HandleRPC(ctx, rs)
	}
}

func (h handler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	for _, hh := range h {
		ctx = hh.TagConn(ctx, cti)
	}
	return ctx
}

func (h handler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	for _, hh := range h {
		ctx = hh.TagRPC(ctx, rti)
	}
	return ctx
}
