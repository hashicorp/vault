package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/admin_collection_v1"
	"google.golang.org/grpc"
)

type routingImpl_CollectionV1 struct {
	client *RoutingClient
}

var _ admin_collection_v1.CollectionAdminServiceClient = (*routingImpl_CollectionV1)(nil)

func (c *routingImpl_CollectionV1) ListCollections(ctx context.Context, in *admin_collection_v1.ListCollectionsRequest, opts ...grpc.CallOption) (*admin_collection_v1.ListCollectionsResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).CollectionV1().ListCollections(ctx, in, opts...)
}
func (c *routingImpl_CollectionV1) CreateScope(ctx context.Context, in *admin_collection_v1.CreateScopeRequest, opts ...grpc.CallOption) (*admin_collection_v1.CreateScopeResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).CollectionV1().CreateScope(ctx, in, opts...)
}
func (c *routingImpl_CollectionV1) DeleteScope(ctx context.Context, in *admin_collection_v1.DeleteScopeRequest, opts ...grpc.CallOption) (*admin_collection_v1.DeleteScopeResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).CollectionV1().DeleteScope(ctx, in, opts...)
}
func (c *routingImpl_CollectionV1) CreateCollection(ctx context.Context, in *admin_collection_v1.CreateCollectionRequest, opts ...grpc.CallOption) (*admin_collection_v1.CreateCollectionResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).CollectionV1().CreateCollection(ctx, in, opts...)
}
func (c *routingImpl_CollectionV1) DeleteCollection(ctx context.Context, in *admin_collection_v1.DeleteCollectionRequest, opts ...grpc.CallOption) (*admin_collection_v1.DeleteCollectionResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).CollectionV1().DeleteCollection(ctx, in, opts...)
}
func (c *routingImpl_CollectionV1) UpdateCollection(ctx context.Context, in *admin_collection_v1.UpdateCollectionRequest, opts ...grpc.CallOption) (*admin_collection_v1.UpdateCollectionResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).CollectionV1().UpdateCollection(ctx, in, opts...)
}
