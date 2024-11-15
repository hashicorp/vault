package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/admin_query_v1"

	"google.golang.org/grpc"
)

type routingImpl_QueryAdminV1 struct {
	client *RoutingClient
}

var _ admin_query_v1.QueryAdminServiceClient = (*routingImpl_QueryAdminV1)(nil)

func (c *routingImpl_QueryAdminV1) GetAllIndexes(ctx context.Context, in *admin_query_v1.GetAllIndexesRequest, opts ...grpc.CallOption) (*admin_query_v1.GetAllIndexesResponse, error) {
	if in.BucketName != nil {
		return c.client.fetchConnForBucket(*in.BucketName).QueryAdminV1().GetAllIndexes(ctx, in, opts...)
	} else {
		return c.client.fetchConn().QueryAdminV1().GetAllIndexes(ctx, in, opts...)
	}
}

func (c *routingImpl_QueryAdminV1) CreatePrimaryIndex(ctx context.Context, in *admin_query_v1.CreatePrimaryIndexRequest, opts ...grpc.CallOption) (*admin_query_v1.CreatePrimaryIndexResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).QueryAdminV1().CreatePrimaryIndex(ctx, in, opts...)
}

func (c *routingImpl_QueryAdminV1) CreateIndex(ctx context.Context, in *admin_query_v1.CreateIndexRequest, opts ...grpc.CallOption) (*admin_query_v1.CreateIndexResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).QueryAdminV1().CreateIndex(ctx, in, opts...)
}

func (c *routingImpl_QueryAdminV1) DropPrimaryIndex(ctx context.Context, in *admin_query_v1.DropPrimaryIndexRequest, opts ...grpc.CallOption) (*admin_query_v1.DropPrimaryIndexResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).QueryAdminV1().DropPrimaryIndex(ctx, in, opts...)
}

func (c *routingImpl_QueryAdminV1) DropIndex(ctx context.Context, in *admin_query_v1.DropIndexRequest, opts ...grpc.CallOption) (*admin_query_v1.DropIndexResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).QueryAdminV1().DropIndex(ctx, in, opts...)
}

func (c *routingImpl_QueryAdminV1) BuildDeferredIndexes(ctx context.Context, in *admin_query_v1.BuildDeferredIndexesRequest, opts ...grpc.CallOption) (*admin_query_v1.BuildDeferredIndexesResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).QueryAdminV1().BuildDeferredIndexes(ctx, in, opts...)
}

func (c *routingImpl_QueryAdminV1) WaitForIndexOnline(ctx context.Context, in *admin_query_v1.WaitForIndexOnlineRequest, opts ...grpc.CallOption) (*admin_query_v1.WaitForIndexOnlineResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).QueryAdminV1().WaitForIndexOnline(ctx, in, opts...)
}
