package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/kv_v1"
	"google.golang.org/grpc"
)

type routingImpl_KvV1 struct {
	client *RoutingClient
}

// Verify that RoutingClient implements Conn
var _ kv_v1.KvServiceClient = (*routingImpl_KvV1)(nil)

func (c *routingImpl_KvV1) Get(ctx context.Context, in *kv_v1.GetRequest, opts ...grpc.CallOption) (*kv_v1.GetResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Get(ctx, in, opts...)
}

func (c *routingImpl_KvV1) GetAndTouch(ctx context.Context, in *kv_v1.GetAndTouchRequest, opts ...grpc.CallOption) (*kv_v1.GetAndTouchResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().GetAndTouch(ctx, in, opts...)
}

func (c *routingImpl_KvV1) GetAndLock(ctx context.Context, in *kv_v1.GetAndLockRequest, opts ...grpc.CallOption) (*kv_v1.GetAndLockResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().GetAndLock(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Unlock(ctx context.Context, in *kv_v1.UnlockRequest, opts ...grpc.CallOption) (*kv_v1.UnlockResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Unlock(ctx, in, opts...)
}

func (c *routingImpl_KvV1) GetAllReplicas(ctx context.Context, in *kv_v1.GetAllReplicasRequest, opts ...grpc.CallOption) (kv_v1.KvService_GetAllReplicasClient, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().GetAllReplicas(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Touch(ctx context.Context, in *kv_v1.TouchRequest, opts ...grpc.CallOption) (*kv_v1.TouchResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Touch(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Exists(ctx context.Context, in *kv_v1.ExistsRequest, opts ...grpc.CallOption) (*kv_v1.ExistsResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Exists(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Insert(ctx context.Context, in *kv_v1.InsertRequest, opts ...grpc.CallOption) (*kv_v1.InsertResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Insert(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Upsert(ctx context.Context, in *kv_v1.UpsertRequest, opts ...grpc.CallOption) (*kv_v1.UpsertResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Upsert(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Replace(ctx context.Context, in *kv_v1.ReplaceRequest, opts ...grpc.CallOption) (*kv_v1.ReplaceResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Replace(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Remove(ctx context.Context, in *kv_v1.RemoveRequest, opts ...grpc.CallOption) (*kv_v1.RemoveResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Remove(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Increment(ctx context.Context, in *kv_v1.IncrementRequest, opts ...grpc.CallOption) (*kv_v1.IncrementResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Increment(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Decrement(ctx context.Context, in *kv_v1.DecrementRequest, opts ...grpc.CallOption) (*kv_v1.DecrementResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Decrement(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Append(ctx context.Context, in *kv_v1.AppendRequest, opts ...grpc.CallOption) (*kv_v1.AppendResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Append(ctx, in, opts...)
}

func (c *routingImpl_KvV1) Prepend(ctx context.Context, in *kv_v1.PrependRequest, opts ...grpc.CallOption) (*kv_v1.PrependResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().Prepend(ctx, in, opts...)
}

func (c *routingImpl_KvV1) LookupIn(ctx context.Context, in *kv_v1.LookupInRequest, opts ...grpc.CallOption) (*kv_v1.LookupInResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().LookupIn(ctx, in, opts...)
}

func (c *routingImpl_KvV1) MutateIn(ctx context.Context, in *kv_v1.MutateInRequest, opts ...grpc.CallOption) (*kv_v1.MutateInResponse, error) {
	return c.client.fetchConnForKey(in.BucketName, in.Key).KvV1().MutateIn(ctx, in, opts...)
}
