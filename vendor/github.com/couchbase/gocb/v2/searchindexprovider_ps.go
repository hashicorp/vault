package gocb

import (
	"encoding/json"

	"github.com/couchbase/goprotostellar/genproto/admin_search_v1"
)

type searchIndexProviderPs struct {
	provider admin_search_v1.SearchAdminServiceClient

	managerProvider *psOpManagerProvider
}

var _ searchIndexProvider = (*searchIndexProviderPs)(nil)

func (sip *searchIndexProviderPs) newOpManager(parentSpan RequestSpan, opName string, attribs map[string]interface{}) *psOpManagerDefault {
	return sip.managerProvider.NewManager(parentSpan, opName, attribs)
}

func (sip *searchIndexProviderPs) GetAllIndexes(scope *Scope, opts *GetAllSearchIndexOptions) ([]SearchIndex, error) {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_get_all_indexes", map[string]interface{}{
		"db.operation": "ListIndexes",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	req := &admin_search_v1.ListIndexesRequest{}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	resp, err := wrapPSOp(manager, req, sip.provider.ListIndexes)
	if err != nil {
		return nil, err
	}

	indexes := make([]SearchIndex, len(resp.Indexes))
	for i, idx := range resp.Indexes {
		params, err := deserializeBytesMap(idx.Params)
		if err != nil {
			return nil, makeGenericError(err, nil)
		}
		sourceParams, err := deserializeBytesMap(idx.Params)
		if err != nil {
			return nil, makeGenericError(err, nil)
		}
		planParams, err := deserializeBytesMap(idx.Params)
		if err != nil {
			return nil, makeGenericError(err, nil)
		}

		index := SearchIndex{
			UUID:         idx.Uuid,
			Name:         idx.Name,
			SourceName:   idx.GetSourceName(),
			Type:         idx.Type,
			Params:       params,
			SourceUUID:   idx.GetSourceUuid(),
			SourceParams: sourceParams,
			SourceType:   idx.GetSourceType(),
			PlanParams:   planParams,
		}

		indexes[i] = index
	}

	return indexes, nil
}

func (sip *searchIndexProviderPs) GetIndex(scope *Scope, indexName string, opts *GetSearchIndexOptions) (*SearchIndex, error) {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_get_index", map[string]interface{}{
		"db.operation": "GetIndex",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	req := &admin_search_v1.GetIndexRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	resp, err := wrapPSOp(manager, req, sip.provider.GetIndex)
	if err != nil {
		return nil, err
	}

	idx := resp.Index

	params, err := deserializeBytesMap(idx.Params)
	if err != nil {
		return nil, makeGenericError(err, nil)
	}
	sourceParams, err := deserializeBytesMap(idx.Params)
	if err != nil {
		return nil, makeGenericError(err, nil)
	}
	planParams, err := deserializeBytesMap(idx.Params)
	if err != nil {
		return nil, makeGenericError(err, nil)
	}

	return &SearchIndex{
		UUID:         idx.Uuid,
		Name:         idx.Name,
		SourceName:   idx.GetSourceName(),
		Type:         idx.Type,
		Params:       params,
		SourceUUID:   idx.GetSourceUuid(),
		SourceParams: sourceParams,
		SourceType:   idx.GetSourceType(),
		PlanParams:   planParams,
	}, nil
}

func (sip *searchIndexProviderPs) UpsertIndex(scope *Scope, indexDefinition SearchIndex, opts *UpsertSearchIndexOptions) error {
	if indexDefinition.UUID == "" {
		return sip.createIndex(scope, indexDefinition, opts)
	}

	return sip.updateIndex(scope, indexDefinition, opts)
}

func (sip *searchIndexProviderPs) updateIndex(scope *Scope, index SearchIndex, opts *UpsertSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_upsert_index", map[string]interface{}{
		"db.operation": "UpdateIndex",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_search_v1.UpdateIndexRequest{}
	var err error
	req.Index, err = sip.makeIndex(index)
	if err != nil {
		return err
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err = wrapPSOp(manager, req, sip.provider.UpdateIndex)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) createIndex(scope *Scope, index SearchIndex, opts *UpsertSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_upsert_index", map[string]interface{}{
		"db.operation": "CreateIndex",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_search_v1.CreateIndexRequest{
		Name: index.Name,
		Type: index.Type,
	}

	if index.SourceName != "" {
		req.SourceName = &index.SourceName
	}

	if index.SourceType != "" {
		req.SourceType = &index.SourceType
	}

	if index.SourceUUID != "" {
		req.SourceName = &index.SourceUUID
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	var err error
	req.Params, err = serializeBytesMap(index.Params)
	if err != nil {
		return err
	}

	req.PlanParams, err = serializeBytesMap(index.PlanParams)
	if err != nil {
		return err
	}

	req.SourceParams, err = serializeBytesMap(index.SourceParams)
	if err != nil {
		return err
	}

	_, err = wrapPSOp(manager, req, sip.provider.CreateIndex)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) DropIndex(scope *Scope, indexName string, opts *DropSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_drop_index", map[string]interface{}{
		"db.operation": "DeleteIndex",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_search_v1.DeleteIndexRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.DeleteIndex)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) AnalyzeDocument(scope *Scope, indexName string, doc interface{}, opts *AnalyzeDocumentOptions) ([]interface{}, error) {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_analyze_document", map[string]interface{}{
		"db.operation": "AnalyzeDocument",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	b, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	req := &admin_search_v1.AnalyzeDocumentRequest{
		Name: indexName,
		Doc:  b,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	resp, err := wrapPSOp(manager, req, sip.provider.AnalyzeDocument)
	if err != nil {
		return nil, err
	}

	var analyzed []interface{}
	err = json.Unmarshal(resp.Analyzed, &analyzed)
	if err != nil {
		return nil, err
	}

	return analyzed, nil
}

func (sip *searchIndexProviderPs) GetIndexedDocumentsCount(scope *Scope, indexName string, opts *GetIndexedDocumentsCountOptions) (uint64, error) {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_get_indexed_documents_count", map[string]interface{}{
		"db.operation": "GetIndexedDocumentsCount",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.GetIndexedDocumentsCountRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	resp, err := wrapPSOp(manager, req, sip.provider.GetIndexedDocumentsCount)
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (sip *searchIndexProviderPs) PauseIngest(scope *Scope, indexName string, opts *PauseIngestSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_pause_ingest", map[string]interface{}{
		"db.operation": "PauseIndexIngest",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.PauseIndexIngestRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.PauseIndexIngest)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) ResumeIngest(scope *Scope, indexName string, opts *ResumeIngestSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_resume_ingest", map[string]interface{}{
		"db.operation": "ResumeIndexIngest",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.ResumeIndexIngestRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.ResumeIndexIngest)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) AllowQuerying(scope *Scope, indexName string, opts *AllowQueryingSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_allow_querying", map[string]interface{}{
		"db.operation": "AllowIndexQuerying",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.AllowIndexQueryingRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.AllowIndexQuerying)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) DisallowQuerying(scope *Scope, indexName string, opts *DisallowQueryingSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_disallow_querying", map[string]interface{}{
		"db.operation": "DisallowIndexQuerying",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.DisallowIndexQueryingRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.DisallowIndexQuerying)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) FreezePlan(scope *Scope, indexName string, opts *FreezePlanSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_freeze_plan", map[string]interface{}{
		"db.operation": "FreezeIndexPlan",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.FreezeIndexPlanRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.FreezeIndexPlan)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) UnfreezePlan(scope *Scope, indexName string, opts *UnfreezePlanSearchIndexOptions) error {
	manager := sip.newOpManager(opts.ParentSpan, "manager_search_unfreeze_plan", map[string]interface{}{
		"db.operation": "UnfreezeIndexPlan",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	req := &admin_search_v1.UnfreezeIndexPlanRequest{
		Name: indexName,
	}

	if scope != nil {
		req.ScopeName = &scope.scopeName
		req.BucketName = &scope.bucket.bucketName
	}

	_, err := wrapPSOp(manager, req, sip.provider.UnfreezeIndexPlan)
	if err != nil {
		return err
	}

	return nil
}

func (sip *searchIndexProviderPs) makeIndex(idx SearchIndex) (*admin_search_v1.Index, error) {
	newIdx := &admin_search_v1.Index{
		Name: idx.Name,
		Type: idx.Type,
		Uuid: idx.UUID,
	}

	if idx.SourceName != "" {
		newIdx.SourceName = &idx.SourceName
	}

	if idx.SourceType != "" {
		newIdx.SourceType = &idx.SourceType
	}

	if idx.SourceUUID != "" {
		newIdx.SourceUuid = &idx.SourceUUID
	}

	var err error
	newIdx.Params, err = serializeBytesMap(idx.Params)
	if err != nil {
		return nil, err
	}

	newIdx.PlanParams, err = serializeBytesMap(idx.PlanParams)
	if err != nil {
		return nil, err
	}

	newIdx.SourceParams, err = serializeBytesMap(idx.SourceParams)
	if err != nil {
		return nil, err
	}

	return newIdx, nil
}

func serializeBytesMap(m map[string]interface{}) (map[string][]byte, error) {
	deserialized := make(map[string][]byte, len(m))
	for k, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		deserialized[k] = b
	}

	return deserialized, nil
}

func deserializeBytesMap(m map[string][]byte) (map[string]interface{}, error) {
	deserialized := make(map[string]interface{}, len(m))
	for k, v := range m {
		var d interface{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			return nil, err
		}

		deserialized[k] = d
	}

	return deserialized, nil
}
