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
	"sync/atomic"
	"time"

	ocstats "go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// ClientStatsHandler is a stats.Handler implementation
// that collects stats for a gRPC client. Predefined
// measures and views can be used to access the collected data.
type ClientStatsHandler struct{}

var _ stats.Handler = &ClientStatsHandler{}

// NewClientStatsHandler returns a stats.Handler implementation
// that collects stats for a gRPC client. Predefined
// measures and views can be used to access the collected data.
func NewClientStatsHandler() *ClientStatsHandler {
	return &ClientStatsHandler{}
}

// TODO(jbd): Remove NewClientStatsHandler and NewServerStatsHandler
// given they are not doing anything than returning a zero value pointer.

// TagConn adds connection related data to the given context and returns the
// new context.
func (h *ClientStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// Do nothing. This is here to satisfy the interface "google.golang.org/grpc/stats.Handler"
	return ctx
}

// HandleConn processes the connection events.
func (h *ClientStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	// Do nothing. This is here to satisfy the interface "google.golang.org/grpc/stats.Handler"
}

// TagRPC gets the tag.Map populated by the application code, serializes
// its tags into the GRPC metadata in order to be sent to the server.
func (h *ClientStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	startTime := time.Now()
	if info == nil {
		if grpclog.V(2) {
			grpclog.Infof("clientHandler.TagRPC called with nil info.", info.FullMethodName)
		}
		return ctx
	}

	d := &rpcData{startTime: startTime}
	ts := tag.FromContext(ctx)
	encoded := tag.Encode(ts)
	ctx = stats.SetTags(ctx, encoded)
	ctx, _ = tag.New(ctx,
		tag.Upsert(KeyMethod, methodName(info.FullMethodName)),
	)
	// TODO(acetechnologist): should we be recording this later? What is the
	// point of updating d.reqLen & d.reqCount if we update now?
	ocstats.Record(ctx, RPCClientStartedCount.M(1))

	return context.WithValue(ctx, grpcClientRPCKey, d)
}

// HandleRPC processes the RPC events.
func (h *ClientStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	switch st := s.(type) {
	case *stats.Begin, *stats.OutHeader, *stats.InHeader, *stats.InTrailer, *stats.OutTrailer:
		// do nothing for client
	case *stats.OutPayload:
		h.handleRPCOutPayload(ctx, st)
	case *stats.InPayload:
		h.handleRPCInPayload(ctx, st)
	case *stats.End:
		h.handleRPCEnd(ctx, st)
	default:
		grpclog.Infof("unexpected stats: %T", st)
	}
}

func (h *ClientStatsHandler) handleRPCOutPayload(ctx context.Context, s *stats.OutPayload) {
	d, ok := ctx.Value(grpcClientRPCKey).(*rpcData)
	if !ok {
		if grpclog.V(2) {
			grpclog.Infoln("clientHandler.handleRPCOutPayload failed to retrieve *rpcData from context")
		}
		return
	}

	ocstats.Record(ctx, RPCClientRequestBytes.M(int64(s.Length)))
	atomic.AddInt64(&d.reqCount, 1)
}

func (h *ClientStatsHandler) handleRPCInPayload(ctx context.Context, s *stats.InPayload) {
	d, ok := ctx.Value(grpcClientRPCKey).(*rpcData)
	if !ok {
		if grpclog.V(2) {
			grpclog.Infoln("clientHandler.handleRPCInPayload failed to retrieve *rpcData from context")
		}
		return
	}

	ocstats.Record(ctx, RPCClientResponseBytes.M(int64(s.Length)))
	atomic.AddInt64(&d.respCount, 1)
}

func (h *ClientStatsHandler) handleRPCEnd(ctx context.Context, s *stats.End) {
	d, ok := ctx.Value(grpcClientRPCKey).(*rpcData)
	if !ok {
		if grpclog.V(2) {
			grpclog.Infoln("clientHandler.handleRPCEnd failed to retrieve *rpcData from context")
		}
		return
	}

	elapsedTime := time.Since(d.startTime)
	reqCount := atomic.LoadInt64(&d.reqCount)
	respCount := atomic.LoadInt64(&d.respCount)

	m := []ocstats.Measurement{
		RPCClientRequestCount.M(reqCount),
		RPCClientResponseCount.M(respCount),
		RPCClientFinishedCount.M(1),
		RPCClientRoundTripLatency.M(float64(elapsedTime) / float64(time.Millisecond)),
	}

	if s.Error != nil {
		s, ok := status.FromError(s.Error)
		if ok {
			ctx, _ = tag.New(ctx,
				tag.Upsert(KeyStatus, s.Code().String()),
			)
		}
		m = append(m, RPCClientErrorCount.M(1))
	}

	ocstats.Record(ctx, m...)
}
