/*
Copyright 2024 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"context"
	"strings"

	vkit "cloud.google.com/go/spanner/apiv1"
	"cloud.google.com/go/spanner/apiv1/spannerpb"
	"cloud.google.com/go/spanner/internal"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	directPathIPV6Prefix = "[2001:4860:8040"
	directPathIPV4Prefix = "34.126"
)

type contextKey string

const metricsTracerKey contextKey = "metricsTracer"

// spannerClient is an interface that defines the methods available from Cloud Spanner API.
type spannerClient interface {
	CallOptions() *vkit.CallOptions
	Close() error
	Connection() *grpc.ClientConn
	CreateSession(context.Context, *spannerpb.CreateSessionRequest, ...gax.CallOption) (*spannerpb.Session, error)
	BatchCreateSessions(context.Context, *spannerpb.BatchCreateSessionsRequest, ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error)
	GetSession(context.Context, *spannerpb.GetSessionRequest, ...gax.CallOption) (*spannerpb.Session, error)
	ListSessions(context.Context, *spannerpb.ListSessionsRequest, ...gax.CallOption) *vkit.SessionIterator
	DeleteSession(context.Context, *spannerpb.DeleteSessionRequest, ...gax.CallOption) error
	ExecuteSql(context.Context, *spannerpb.ExecuteSqlRequest, ...gax.CallOption) (*spannerpb.ResultSet, error)
	ExecuteStreamingSql(context.Context, *spannerpb.ExecuteSqlRequest, ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error)
	ExecuteBatchDml(context.Context, *spannerpb.ExecuteBatchDmlRequest, ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error)
	Read(context.Context, *spannerpb.ReadRequest, ...gax.CallOption) (*spannerpb.ResultSet, error)
	StreamingRead(context.Context, *spannerpb.ReadRequest, ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error)
	BeginTransaction(context.Context, *spannerpb.BeginTransactionRequest, ...gax.CallOption) (*spannerpb.Transaction, error)
	Commit(context.Context, *spannerpb.CommitRequest, ...gax.CallOption) (*spannerpb.CommitResponse, error)
	Rollback(context.Context, *spannerpb.RollbackRequest, ...gax.CallOption) error
	PartitionQuery(context.Context, *spannerpb.PartitionQueryRequest, ...gax.CallOption) (*spannerpb.PartitionResponse, error)
	PartitionRead(context.Context, *spannerpb.PartitionReadRequest, ...gax.CallOption) (*spannerpb.PartitionResponse, error)
	BatchWrite(context.Context, *spannerpb.BatchWriteRequest, ...gax.CallOption) (spannerpb.Spanner_BatchWriteClient, error)
}

// grpcSpannerClient is the gRPC API implementation of the transport-agnostic
// spannerClient interface.
type grpcSpannerClient struct {
	raw                  *vkit.Client
	metricsTracerFactory *builtinMetricsTracerFactory
}

var (
	// Ensure that grpcSpannerClient implements spannerClient.
	_ spannerClient = (*grpcSpannerClient)(nil)
)

// newGRPCSpannerClient initializes a new spannerClient that uses the gRPC
// Spanner API.
func newGRPCSpannerClient(ctx context.Context, sc *sessionClient, opts ...option.ClientOption) (spannerClient, error) {
	raw, err := vkit.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	g := &grpcSpannerClient{raw: raw, metricsTracerFactory: sc.metricsTracerFactory}
	clientInfo := []string{"gccl", internal.Version}
	if sc.userAgent != "" {
		agentWithVersion := strings.SplitN(sc.userAgent, "/", 2)
		if len(agentWithVersion) == 2 {
			clientInfo = append(clientInfo, agentWithVersion[0], agentWithVersion[1])
		}
	}
	raw.SetGoogleClientInfo(clientInfo...)
	if sc.callOptions != nil {
		raw.CallOptions = mergeCallOptions(raw.CallOptions, sc.callOptions)
	}
	return g, nil
}

func (g *grpcSpannerClient) newBuiltinMetricsTracer(ctx context.Context) *builtinMetricsTracer {
	mt := g.metricsTracerFactory.createBuiltinMetricsTracer(ctx)
	return &mt
}

func (g *grpcSpannerClient) CallOptions() *vkit.CallOptions {
	return g.raw.CallOptions
}

func (g *grpcSpannerClient) Close() error {
	return g.raw.Close()
}

func (g *grpcSpannerClient) Connection() *grpc.ClientConn {
	return g.raw.Connection()
}

func (g *grpcSpannerClient) CreateSession(ctx context.Context, req *spannerpb.CreateSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.CreateSession(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) BatchCreateSessions(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest, opts ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.BatchCreateSessions(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) GetSession(ctx context.Context, req *spannerpb.GetSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.GetSession(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) ListSessions(ctx context.Context, req *spannerpb.ListSessionsRequest, opts ...gax.CallOption) *vkit.SessionIterator {
	return g.raw.ListSessions(ctx, req, opts...)
}

func (g *grpcSpannerClient) DeleteSession(ctx context.Context, req *spannerpb.DeleteSessionRequest, opts ...gax.CallOption) error {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	err := g.raw.DeleteSession(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return err
}

func (g *grpcSpannerClient) ExecuteSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.ExecuteSql(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) ExecuteStreamingSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error) {
	return g.raw.ExecuteStreamingSql(peer.NewContext(ctx, &peer.Peer{}), req, opts...)
}

func (g *grpcSpannerClient) ExecuteBatchDml(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest, opts ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.ExecuteBatchDml(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) Read(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.Read(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) StreamingRead(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error) {
	return g.raw.StreamingRead(peer.NewContext(ctx, &peer.Peer{}), req, opts...)
}

func (g *grpcSpannerClient) BeginTransaction(ctx context.Context, req *spannerpb.BeginTransactionRequest, opts ...gax.CallOption) (*spannerpb.Transaction, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.BeginTransaction(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) Commit(ctx context.Context, req *spannerpb.CommitRequest, opts ...gax.CallOption) (*spannerpb.CommitResponse, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.Commit(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) Rollback(ctx context.Context, req *spannerpb.RollbackRequest, opts ...gax.CallOption) error {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	err := g.raw.Rollback(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return err
}

func (g *grpcSpannerClient) PartitionQuery(ctx context.Context, req *spannerpb.PartitionQueryRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.PartitionQuery(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) PartitionRead(ctx context.Context, req *spannerpb.PartitionReadRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	mt := g.newBuiltinMetricsTracer(ctx)
	defer recordOperationCompletion(mt)
	ctx = context.WithValue(ctx, metricsTracerKey, mt)
	resp, err := g.raw.PartitionRead(ctx, req, opts...)
	statusCode, _ := status.FromError(err)
	mt.currOp.setStatus(statusCode.Code().String())
	return resp, err
}

func (g *grpcSpannerClient) BatchWrite(ctx context.Context, req *spannerpb.BatchWriteRequest, opts ...gax.CallOption) (spannerpb.Spanner_BatchWriteClient, error) {
	return g.raw.BatchWrite(peer.NewContext(ctx, &peer.Peer{}), req, opts...)
}
