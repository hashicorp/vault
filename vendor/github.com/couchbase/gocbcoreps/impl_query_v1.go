package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/query_v1"
	"google.golang.org/grpc"
)

type routingImpl_QueryV1 struct {
	client *RoutingClient
}

// Verify that RoutingClient implements Conn
var _ query_v1.QueryServiceClient = (*routingImpl_QueryV1)(nil)

func (c *routingImpl_QueryV1) Query(ctx context.Context, in *query_v1.QueryRequest, opts ...grpc.CallOption) (query_v1.QueryService_QueryClient, error) {
	if in.BucketName != nil {
		return c.client.fetchConnForBucket(*in.BucketName).QueryV1().Query(ctx, in, opts...)
	} else {
		return c.client.fetchConn().QueryV1().Query(ctx, in, opts...)
	}
}
