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

// Package grpctrace is a package to assist with tracing incoming and outgoing gRPC requests.
package grpctrace

import (
	"strings"

	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

// ClientStatsHandler is a an implementation of grpc.StatsHandler
// that can be passed to grpc.Dial
// using grpc.WithStatsHandler to enable trace context propagation and
// automatic span creation for outgoing gRPC requests.
type ClientStatsHandler struct{}

var _ stats.Handler = &ClientStatsHandler{}

// NewClientStatsHandler returns a StatsHandler that can be passed to grpc.Dial
// using grpc.WithStatsHandler to enable trace context propagation and
// automatic span creation for outgoing gRPC requests.
func NewClientStatsHandler() *ClientStatsHandler {
	return &ClientStatsHandler{}
}

// TODO(jbd): Remove NewClientStatsHandler and NewServerStatsHandler
// given they are not doing anything than returning a zero value pointer.

// ServerStatsHandler is a an implementation of grpc.StatsHandler
// that can be passed to grpc.NewServer using grpc.StatsHandler
// to enable trace context propagation and automatic span creation
// for incoming gRPC requests..
type ServerStatsHandler struct{}

// NewServerStatsHandler returns a StatsHandler that can be passed to
// grpc.NewServer using grpc.StatsHandler to enable trace context propagation
// and automatic span creation for incoming gRPC requests.
func NewServerStatsHandler() *ServerStatsHandler {
	return &ServerStatsHandler{}
}

var _ stats.Handler = &ServerStatsHandler{}

const traceContextKey = "grpc-trace-bin"

// TagRPC creates a new trace span for the client side of the RPC.
//
// It returns ctx with the new trace span added and a serialization of the
// SpanContext added to the outgoing gRPC metadata.
func (c *ClientStatsHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	name := "Sent" + strings.Replace(rti.FullMethodName, "/", ".", -1)
	ctx, _ = trace.StartSpanWithOptions(ctx, name, trace.StartOptions{RecordEvents: true, RegisterNameForLocalSpanStore: true})
	traceContextBinary := propagation.Binary(trace.FromContext(ctx).SpanContext())
	if len(traceContextBinary) == 0 {
		return ctx
	}
	md := metadata.Pairs(traceContextKey, string(traceContextBinary))
	if oldMD, ok := metadata.FromOutgoingContext(ctx); ok {
		md = metadata.Join(oldMD, md)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// TagRPC creates a new trace span for the server side of the RPC.
//
// It checks the incoming gRPC metadata in ctx for a SpanContext, and if
// it finds one, uses that SpanContext as the parent context of the new span.
//
// It returns ctx, with the new trace span added.
func (s *ServerStatsHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	name := "Recv" + strings.Replace(rti.FullMethodName, "/", ".", -1)
	opt := trace.StartOptions{RecordEvents: true, RegisterNameForLocalSpanStore: true}
	if s := md[traceContextKey]; len(s) > 0 {
		if parent, ok := propagation.FromBinary([]byte(s[0])); ok {
			ctx, _ = trace.StartSpanWithRemoteParent(ctx, name, parent, opt)
			return ctx
		}
	}
	ctx, _ = trace.StartSpanWithOptions(ctx, name, opt)
	return ctx
}

// HandleRPC processes the RPC stats, adding information to the current trace span.
func (c *ClientStatsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	handleRPC(ctx, rs)
}

// HandleRPC processes the RPC stats, adding information to the current trace span.
func (s *ServerStatsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	handleRPC(ctx, rs)
}

func handleRPC(ctx context.Context, rs stats.RPCStats) {
	span := trace.FromContext(ctx)
	// TODO: compressed and uncompressed sizes are not populated in every message.
	switch rs := rs.(type) {
	case *stats.Begin:
		span.SetAttributes(
			trace.BoolAttribute{Key: "Client", Value: rs.Client},
			trace.BoolAttribute{Key: "FailFast", Value: rs.FailFast})
	case *stats.InPayload:
		span.AddMessageReceiveEvent(0 /* TODO: messageID */, int64(rs.Length), int64(rs.WireLength))
	case *stats.InHeader:
		span.AddMessageReceiveEvent(0, int64(rs.WireLength), int64(rs.WireLength))
	case *stats.InTrailer:
		span.AddMessageReceiveEvent(0, int64(rs.WireLength), int64(rs.WireLength))
	case *stats.OutPayload:
		span.AddMessageSendEvent(0, int64(rs.Length), int64(rs.WireLength))
	case *stats.OutHeader:
		span.AddMessageSendEvent(0, 0, 0)
	case *stats.OutTrailer:
		span.AddMessageSendEvent(0, int64(rs.WireLength), int64(rs.WireLength))
	case *stats.End:
		if rs.Error != nil {
			code, desc := grpc.Code(rs.Error), grpc.ErrorDesc(rs.Error)
			span.SetStatus(trace.Status{Code: int32(code), Message: desc})
		}
		span.End()
	}
}

// TagConn is a no-op for this StatsHandler.
func (c *ClientStatsHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	return ctx
}

// TagConn is a no-op for this StatsHandler.
func (s *ServerStatsHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn is a no-op for this StatsHandler.
func (c *ClientStatsHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
}

// HandleConn is a no-op for this StatsHandler.
func (s *ServerStatsHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
}
