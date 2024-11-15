package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/analytics_v1"
	"google.golang.org/grpc"
)

type routingImpl_AnalyticsV1 struct {
	client *RoutingClient
}

var _ analytics_v1.AnalyticsServiceClient = (*routingImpl_AnalyticsV1)(nil)

func (c *routingImpl_AnalyticsV1) AnalyticsQuery(ctx context.Context, in *analytics_v1.AnalyticsQueryRequest, opts ...grpc.CallOption) (analytics_v1.AnalyticsService_AnalyticsQueryClient, error) {
	return c.client.fetchConn().AnalyticsV1().AnalyticsQuery(ctx, in, opts...)
}
