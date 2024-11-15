package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/routing_v1"
	"google.golang.org/grpc"
)

type routingImpl_RoutingV1 struct {
	client *RoutingClient
}

// Verify that RoutingClient implements Conn
var _ routing_v1.RoutingServiceClient = (*routingImpl_RoutingV1)(nil)

func (c *routingImpl_RoutingV1) WatchRouting(
	ctx context.Context,
	in *routing_v1.WatchRoutingRequest,
	opts ...grpc.CallOption,
) (routing_v1.RoutingService_WatchRoutingClient, error) {
	// We intentionally ignore the bucket name in this request due to the fact
	// that technically routing of a bucket isn't part of the bucket itself.  If
	// we used routing for the bucket routing, it's a circular dependancy.
	return c.client.fetchConn().RoutingV1().WatchRouting(ctx, in, opts...)
}
