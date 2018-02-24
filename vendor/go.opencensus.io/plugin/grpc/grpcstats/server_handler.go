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
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"

	ocstats "go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// ServerStatsHandler is a stats.Handler implementation
// that collects stats for a gRPC server. Predefined
// measures and views can be used to access the collected data.
type ServerStatsHandler struct{}

var _ stats.Handler = &ServerStatsHandler{}

// NewServerStatsHandler returns a stats.Handler implementation
// that collects stats for a gRPC server. Predefined
// measures and views can be used to access the collected data.
func NewServerStatsHandler() *ServerStatsHandler {
	return &ServerStatsHandler{}
}

// TagConn adds connection related data to the given context and returns the
// new context.
func (h *ServerStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// Do nothing. This is here to satisfy the interface "google.golang.org/grpc/stats.Handler"
	return ctx
}

// HandleConn processes the connection events.
func (h *ServerStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	// Do nothing. This is here to satisfy the interface "google.golang.org/grpc/stats.Handler"
}

// TagRPC gets the metadata from gRPC context, extracts the encoded tags from
// it and creates a new tag.Map and puts them into the returned context.
func (h *ServerStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	startTime := time.Now()
	if info == nil {
		if grpclog.V(2) {
			grpclog.Infof("serverHandler.TagRPC called with nil info.", info.FullMethodName)
		}
		return ctx
	}
	d := &rpcData{startTime: startTime}
	ctx, _ = h.createTags(ctx, info.FullMethodName)
	ocstats.Record(ctx, RPCServerStartedCount.M(1))
	return context.WithValue(ctx, grpcServerRPCKey, d)
}

// HandleRPC processes the RPC events.
func (h *ServerStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	switch st := s.(type) {
	case *stats.Begin, *stats.InHeader, *stats.InTrailer, *stats.OutHeader, *stats.OutTrailer:
		// Do nothing for server
	case *stats.InPayload:
		h.handleRPCInPayload(ctx, st)
	case *stats.OutPayload:
		// For stream it can be called multiple times per RPC.
		h.handleRPCOutPayload(ctx, st)
	case *stats.End:
		h.handleRPCEnd(ctx, st)
	default:
		grpclog.Infof("unexpected stats: %T", st)
	}
}

func (h *ServerStatsHandler) handleRPCInPayload(ctx context.Context, s *stats.InPayload) {
	d, ok := ctx.Value(grpcServerRPCKey).(*rpcData)
	if !ok {
		if grpclog.V(2) {
			grpclog.Infoln("serverHandler.handleRPCInPayload failed to retrieve *rpcData from context")
		}
		return
	}

	ocstats.Record(ctx, RPCServerRequestBytes.M(int64(s.Length)))
	atomic.AddInt64(&d.reqCount, 1)
}

func (h *ServerStatsHandler) handleRPCOutPayload(ctx context.Context, s *stats.OutPayload) {
	d, ok := ctx.Value(grpcServerRPCKey).(*rpcData)
	if !ok {
		if grpclog.V(2) {
			grpclog.Infoln("serverHandler.handleRPCOutPayload failed to retrieve *rpcData from context")
		}
		return
	}

	ocstats.Record(ctx, RPCServerResponseBytes.M(int64(s.Length)))
	atomic.AddInt64(&d.respCount, 1)
}

func (h *ServerStatsHandler) handleRPCEnd(ctx context.Context, s *stats.End) {
	d, ok := ctx.Value(grpcServerRPCKey).(*rpcData)
	if !ok {
		if grpclog.V(2) {
			grpclog.Infoln("serverHandler.handleRPCEnd failed to retrieve *rpcData from context")
		}
		return
	}

	elapsedTime := time.Since(d.startTime)
	reqCount := atomic.LoadInt64(&d.reqCount)
	respCount := atomic.LoadInt64(&d.respCount)

	m := []ocstats.Measurement{
		RPCServerRequestCount.M(reqCount),
		RPCServerResponseCount.M(respCount),
		RPCServerFinishedCount.M(1),
		RPCServerServerElapsedTime.M(float64(elapsedTime) / float64(time.Millisecond)),
	}

	if s.Error != nil {
		s, ok := status.FromError(s.Error)
		if ok {
			ctx, _ = tag.New(ctx,
				tag.Upsert(KeyStatus, s.Code().String()),
			)
		}
		m = append(m, RPCServerErrorCount.M(1))
	}

	ocstats.Record(ctx, m...)
}

// createTags creates a new tag map containing the tags extracted from the
// gRPC metadata.
func (h *ServerStatsHandler) createTags(ctx context.Context, fullinfo string) (context.Context, error) {
	mods := []tag.Mutator{
		tag.Upsert(KeyMethod, methodName(fullinfo)),
	}
	if tagsBin := stats.Tags(ctx); tagsBin != nil {
		old, err := tag.Decode([]byte(tagsBin))
		if err != nil {
			return nil, fmt.Errorf("serverHandler.createTags failed to decode tagsBin %v: %v", tagsBin, err)
		}
		return tag.New(tag.NewContext(ctx, old), mods...)
	}
	return tag.New(ctx, mods...)
}
