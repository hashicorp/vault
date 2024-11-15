package gocbcoreps

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/couchbase/goprotostellar/genproto/view_v1"

	"github.com/couchbase/goprotostellar/genproto/admin_search_v1"

	"google.golang.org/grpc/connectivity"

	"github.com/couchbase/goprotostellar/genproto/admin_query_v1"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/couchbase/goprotostellar/genproto/admin_bucket_v1"
	"github.com/couchbase/goprotostellar/genproto/admin_collection_v1"
	"github.com/couchbase/goprotostellar/genproto/analytics_v1"
	"github.com/couchbase/goprotostellar/genproto/kv_v1"
	"github.com/couchbase/goprotostellar/genproto/query_v1"
	"github.com/couchbase/goprotostellar/genproto/routing_v1"
	"github.com/couchbase/goprotostellar/genproto/search_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type routingConnOptions struct {
	InsecureSkipVerify bool // used for enabling TLS, but skipping verification
	ClientCertificate  *x509.CertPool
	Username           string
	Password           string
	TracerProvider     trace.TracerProvider
	MeterProvider      metric.MeterProvider
}

type routingConn struct {
	conn          *grpc.ClientConn
	routingV1     routing_v1.RoutingServiceClient
	kvV1          kv_v1.KvServiceClient
	queryV1       query_v1.QueryServiceClient
	collectionV1  admin_collection_v1.CollectionAdminServiceClient
	bucketV1      admin_bucket_v1.BucketAdminServiceClient
	analyticsV1   analytics_v1.AnalyticsServiceClient
	searchV1      search_v1.SearchServiceClient
	viewV1        view_v1.ViewServiceClient
	queryAdminV1  admin_query_v1.QueryAdminServiceClient
	searchAdminV1 admin_search_v1.SearchAdminServiceClient
}

// Verify that routingConn implements Conn
var _ Conn = (*routingConn)(nil)

const maxMsgSize = 26214400 // 25MiB

func dialRoutingConn(ctx context.Context, address string, opts *routingConnOptions) (*routingConn, error) {
	var transportDialOpt grpc.DialOption
	var perRpcDialOpt grpc.DialOption

	if opts.ClientCertificate != nil { // use tls
		transportDialOpt = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(opts.ClientCertificate, ""))
	} else if opts.InsecureSkipVerify { // use tls, but skip verification
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		transportDialOpt = grpc.WithTransportCredentials(creds)
	} else { // plain text
		transportDialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	// setup basic auth.
	if opts.Username != "" && opts.Password != "" {
		basicAuthCreds, err := NewGrpcBasicAuth(opts.Username, opts.Password)
		if err != nil {
			return nil, err
		}
		perRpcDialOpt = grpc.WithPerRPCCredentials(basicAuthCreds)
	} else {
		perRpcDialOpt = nil
	}

	dialOpts := []grpc.DialOption{transportDialOpt}
	if perRpcDialOpt != nil {
		dialOpts = append(dialOpts, perRpcDialOpt)
	}

	clientOpts := []otelgrpc.Option{
		otelgrpc.WithPropagators(propagation.TraceContext{}),
	}
	if opts.TracerProvider != nil {
		clientOpts = append(clientOpts, otelgrpc.WithTracerProvider(opts.TracerProvider))
	}
	if opts.MeterProvider != nil {
		clientOpts = append(clientOpts, otelgrpc.WithMeterProvider(opts.MeterProvider))
	}
	dialOpts = append(dialOpts, grpc.WithStatsHandler(otelgrpc.NewClientHandler(clientOpts...)))
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxRecvMsgSizeCallOption{MaxRecvMsgSize: maxMsgSize}))

	conn, err := grpc.DialContext(ctx, address, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &routingConn{
		conn:          conn,
		routingV1:     routing_v1.NewRoutingServiceClient(conn),
		kvV1:          kv_v1.NewKvServiceClient(conn),
		queryV1:       query_v1.NewQueryServiceClient(conn),
		collectionV1:  admin_collection_v1.NewCollectionAdminServiceClient(conn),
		bucketV1:      admin_bucket_v1.NewBucketAdminServiceClient(conn),
		analyticsV1:   analytics_v1.NewAnalyticsServiceClient(conn),
		queryAdminV1:  admin_query_v1.NewQueryAdminServiceClient(conn),
		searchV1:      search_v1.NewSearchServiceClient(conn),
		viewV1:        view_v1.NewViewServiceClient(conn),
		searchAdminV1: admin_search_v1.NewSearchAdminServiceClient(conn),
	}, nil
}

func (c *routingConn) RoutingV1() routing_v1.RoutingServiceClient {
	return c.routingV1
}

func (c *routingConn) KvV1() kv_v1.KvServiceClient {
	return c.kvV1
}

func (c *routingConn) QueryV1() query_v1.QueryServiceClient {
	return c.queryV1
}

func (c *routingConn) CollectionV1() admin_collection_v1.CollectionAdminServiceClient {
	return c.collectionV1
}

func (c *routingConn) BucketV1() admin_bucket_v1.BucketAdminServiceClient {
	return c.bucketV1
}

func (c *routingConn) AnalyticsV1() analytics_v1.AnalyticsServiceClient {
	return c.analyticsV1
}

func (c *routingConn) SearchV1() search_v1.SearchServiceClient {
	return c.searchV1
}

func (c *routingConn) ViewV1() view_v1.ViewServiceClient {
	return c.viewV1
}

func (c *routingConn) QueryAdminV1() admin_query_v1.QueryAdminServiceClient {
	return c.queryAdminV1
}

func (c *routingConn) SearchAdminV1() admin_search_v1.SearchAdminServiceClient {
	return c.searchAdminV1
}

func (c *routingConn) Close() error {
	return c.conn.Close()
}

func (c *routingConn) State() ConnState {
	switch c.conn.GetState() {
	case connectivity.Connecting:
		return ConnStateOffline
	case connectivity.Shutdown:
		return ConnStateOffline
	case connectivity.TransientFailure:
		return ConnStateOffline
	case connectivity.Idle:
		return ConnStateOffline
	case connectivity.Ready:
		return ConnStateOnline
	}

	// This connection is in an unknown state so let's assume offline.
	return ConnStateOffline
}
