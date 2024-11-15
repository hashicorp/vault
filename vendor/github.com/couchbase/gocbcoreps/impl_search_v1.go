package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/search_v1"
	"google.golang.org/grpc"
)

type routingImpl_SearchV1 struct {
	client *RoutingClient
}

// Verify that RoutingClient implements Conn
var _ search_v1.SearchServiceClient = (*routingImpl_SearchV1)(nil)

func (c *routingImpl_SearchV1) SearchQuery(ctx context.Context, in *search_v1.SearchQueryRequest, opts ...grpc.CallOption) (search_v1.SearchService_SearchQueryClient, error) {
	if in.BucketName != nil {
		return c.client.fetchConnForBucket(*in.BucketName).SearchV1().SearchQuery(ctx, in, opts...)
	}

	return c.client.fetchConn().SearchV1().SearchQuery(ctx, in, opts...)
}
