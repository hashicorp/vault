package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/admin_bucket_v1"
	"google.golang.org/grpc"
)

type routingImpl_BucketV1 struct {
	client *RoutingClient
}

var _ admin_bucket_v1.BucketAdminServiceClient = (*routingImpl_BucketV1)(nil)

func (c *routingImpl_BucketV1) ListBuckets(ctx context.Context, in *admin_bucket_v1.ListBucketsRequest, opts ...grpc.CallOption) (*admin_bucket_v1.ListBucketsResponse, error) {
	return c.client.fetchConn().BucketV1().ListBuckets(ctx, in, opts...)
}
func (c *routingImpl_BucketV1) CreateBucket(ctx context.Context, in *admin_bucket_v1.CreateBucketRequest, opts ...grpc.CallOption) (*admin_bucket_v1.CreateBucketResponse, error) {
	return c.client.fetchConn().BucketV1().CreateBucket(ctx, in, opts...)
}
func (c *routingImpl_BucketV1) UpdateBucket(ctx context.Context, in *admin_bucket_v1.UpdateBucketRequest, opts ...grpc.CallOption) (*admin_bucket_v1.UpdateBucketResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).bucketV1.UpdateBucket(ctx, in, opts...)
}
func (c *routingImpl_BucketV1) DeleteBucket(ctx context.Context, in *admin_bucket_v1.DeleteBucketRequest, opts ...grpc.CallOption) (*admin_bucket_v1.DeleteBucketResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).bucketV1.DeleteBucket(ctx, in, opts...)
}
func (c *routingImpl_BucketV1) FlushBucket(ctx context.Context, in *admin_bucket_v1.FlushBucketRequest, opts ...grpc.CallOption) (*admin_bucket_v1.FlushBucketResponse, error) {
	return c.client.fetchConnForBucket(in.BucketName).bucketV1.FlushBucket(ctx, in, opts...)
}
