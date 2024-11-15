package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/view_v1"

	"google.golang.org/grpc"
)

type routingImpl_ViewV1 struct {
	client *RoutingClient
}

// Verify that RoutingClient implements Conn
var _ view_v1.ViewServiceClient = (*routingImpl_ViewV1)(nil)

func (c *routingImpl_ViewV1) ViewQuery(ctx context.Context, in *view_v1.ViewQueryRequest, opts ...grpc.CallOption) (view_v1.ViewService_ViewQueryClient, error) {
	return c.client.fetchConnForBucket(in.BucketName).ViewV1().ViewQuery(ctx, in, opts...)
}
