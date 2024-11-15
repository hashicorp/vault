package gocbcoreps

import (
	"context"

	"github.com/couchbase/goprotostellar/genproto/admin_search_v1"

	"google.golang.org/grpc"
)

type routingImpl_SearchAdminV1 struct {
	client *RoutingClient
}

var _ admin_search_v1.SearchAdminServiceClient = (*routingImpl_SearchAdminV1)(nil)

func (r routingImpl_SearchAdminV1) GetIndex(ctx context.Context, in *admin_search_v1.GetIndexRequest, opts ...grpc.CallOption) (*admin_search_v1.GetIndexResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().GetIndex(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().GetIndex(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) ListIndexes(ctx context.Context, in *admin_search_v1.ListIndexesRequest, opts ...grpc.CallOption) (*admin_search_v1.ListIndexesResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().ListIndexes(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().ListIndexes(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) CreateIndex(ctx context.Context, in *admin_search_v1.CreateIndexRequest, opts ...grpc.CallOption) (*admin_search_v1.CreateIndexResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().CreateIndex(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().CreateIndex(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) UpdateIndex(ctx context.Context, in *admin_search_v1.UpdateIndexRequest, opts ...grpc.CallOption) (*admin_search_v1.UpdateIndexResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().UpdateIndex(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().UpdateIndex(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) DeleteIndex(ctx context.Context, in *admin_search_v1.DeleteIndexRequest, opts ...grpc.CallOption) (*admin_search_v1.DeleteIndexResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().DeleteIndex(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().DeleteIndex(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) AnalyzeDocument(ctx context.Context, in *admin_search_v1.AnalyzeDocumentRequest, opts ...grpc.CallOption) (*admin_search_v1.AnalyzeDocumentResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().AnalyzeDocument(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().AnalyzeDocument(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) GetIndexedDocumentsCount(ctx context.Context, in *admin_search_v1.GetIndexedDocumentsCountRequest, opts ...grpc.CallOption) (*admin_search_v1.GetIndexedDocumentsCountResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().GetIndexedDocumentsCount(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().GetIndexedDocumentsCount(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) PauseIndexIngest(ctx context.Context, in *admin_search_v1.PauseIndexIngestRequest, opts ...grpc.CallOption) (*admin_search_v1.PauseIndexIngestResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().PauseIndexIngest(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().PauseIndexIngest(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) ResumeIndexIngest(ctx context.Context, in *admin_search_v1.ResumeIndexIngestRequest, opts ...grpc.CallOption) (*admin_search_v1.ResumeIndexIngestResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().ResumeIndexIngest(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().ResumeIndexIngest(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) AllowIndexQuerying(ctx context.Context, in *admin_search_v1.AllowIndexQueryingRequest, opts ...grpc.CallOption) (*admin_search_v1.AllowIndexQueryingResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().AllowIndexQuerying(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().AllowIndexQuerying(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) DisallowIndexQuerying(ctx context.Context, in *admin_search_v1.DisallowIndexQueryingRequest, opts ...grpc.CallOption) (*admin_search_v1.DisallowIndexQueryingResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().DisallowIndexQuerying(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().DisallowIndexQuerying(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) FreezeIndexPlan(ctx context.Context, in *admin_search_v1.FreezeIndexPlanRequest, opts ...grpc.CallOption) (*admin_search_v1.FreezeIndexPlanResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().FreezeIndexPlan(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().FreezeIndexPlan(ctx, in, opts...)
	}
}

func (r routingImpl_SearchAdminV1) UnfreezeIndexPlan(ctx context.Context, in *admin_search_v1.UnfreezeIndexPlanRequest, opts ...grpc.CallOption) (*admin_search_v1.UnfreezeIndexPlanResponse, error) {
	if in.BucketName != nil {
		return r.client.fetchConnForBucket(*in.BucketName).SearchAdminV1().UnfreezeIndexPlan(ctx, in, opts...)
	} else {
		return r.client.fetchConn().SearchAdminV1().UnfreezeIndexPlan(ctx, in, opts...)
	}
}
