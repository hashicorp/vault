package gocbcoreps

import (
	"github.com/couchbase/goprotostellar/genproto/admin_bucket_v1"
	"github.com/couchbase/goprotostellar/genproto/admin_collection_v1"
	"github.com/couchbase/goprotostellar/genproto/analytics_v1"
	"github.com/couchbase/goprotostellar/genproto/kv_v1"
	"github.com/couchbase/goprotostellar/genproto/query_v1"
	"github.com/couchbase/goprotostellar/genproto/routing_v1"
	"github.com/couchbase/goprotostellar/genproto/search_v1"
)

type Conn interface {
	RoutingV1() routing_v1.RoutingServiceClient
	KvV1() kv_v1.KvServiceClient
	QueryV1() query_v1.QueryServiceClient
	CollectionV1() admin_collection_v1.CollectionAdminServiceClient
	BucketV1() admin_bucket_v1.BucketAdminServiceClient
	AnalyticsV1() analytics_v1.AnalyticsServiceClient
	SearchV1() search_v1.SearchServiceClient
}
