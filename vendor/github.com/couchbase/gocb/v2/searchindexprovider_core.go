package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type searchIndexProviderCore struct {
	mgmtProvider      mgmtProvider
	searchCapVerifier searchCapabilityVerifier

	tracer RequestTracer
	meter  *meterWrapper
}

var _ searchIndexProvider = (*searchIndexProviderCore)(nil)

func (sm *searchIndexProviderCore) GetAllIndexes(scope *Scope, opts *GetAllSearchIndexOptions) ([]SearchIndex, error) {
	if opts == nil {
		opts = &GetAllSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return nil, wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	start := time.Now()
	defer sm.meter.ValueRecord(meterValueServiceManagement, "manager_search_get_all_indexes", start)

	path := sm.pathPrefix(scope)

	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_get_all_indexes", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeSearch,
		Method:        "GET",
		Path:          path,
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := sm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return nil, idxErr
		}

		return nil, makeMgmtBadStatusError("failed to get index", &req, resp)
	}

	var indexesResp jsonSearchIndexesResp
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&indexesResp)
	if err != nil {
		return nil, err
	}

	indexDefs := indexesResp.IndexDefs.IndexDefs
	var indexes []SearchIndex
	for _, indexData := range indexDefs {
		var index SearchIndex
		err := index.fromData(indexData)
		if err != nil {
			return nil, err
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (sm *searchIndexProviderCore) GetIndex(scope *Scope, indexName string, opts *GetSearchIndexOptions) (*SearchIndex, error) {
	if opts == nil {
		opts = &GetSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return nil, wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	start := time.Now()
	defer sm.meter.ValueRecord(meterValueServiceManagement, "manager_search_get_index", start)

	path := fmt.Sprintf("%s/%s", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_get_index", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeSearch,
		Method:        "GET",
		Path:          path,
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := sm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return nil, idxErr
		}

		return nil, makeMgmtBadStatusError("failed to get index", &req, resp)
	}

	var indexResp jsonSearchIndexResp
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&indexResp)
	if err != nil {
		return nil, err
	}

	var indexDef SearchIndex
	err = indexDef.fromData(*indexResp.IndexDef)
	if err != nil {
		return nil, err
	}

	return &indexDef, nil
}

func (sm *searchIndexProviderCore) validateIndexVectorMappingsProperties(vectorMappingProperties map[string]interface{}) error {
	for _, val := range vectorMappingProperties {
		if v, ok := val.(map[string]interface{}); ok {
			if fields, ok := v["fields"]; ok {
				if f, ok := fields.([]interface{}); ok {
					for _, field := range f {
						if item, ok := field.(map[string]interface{}); ok {
							if typ, ok := item["type"]; ok {
								if t, ok := typ.(string); ok {
									if (t == "vector" || t == "vector_base64") && sm.vectorSearchUnsupported() {
										return wrapError(ErrFeatureNotAvailable, "indexes with vectors cannot be used with this server version")
									}
								}
							}
						}
					}
				}
			} else if properties, ok := v["properties"]; ok {
				if p, ok := properties.(map[string]interface{}); ok {
					if err := sm.validateIndexVectorMappingsProperties(p); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (sm *searchIndexProviderCore) validateIndexVectorMappings(indexDefinition SearchIndex) error {
	if mapping, ok := indexDefinition.Params["mapping"]; ok {
		if m, ok := mapping.(map[string]interface{}); ok {
			if types, ok := m["types"]; ok {
				if t, ok := types.(map[string]interface{}); ok {
					for _, val := range t {
						if v, ok := val.(map[string]interface{}); ok {
							if properties, ok := v["properties"]; ok {
								if p, ok := properties.(map[string]interface{}); ok {
									if err := sm.validateIndexVectorMappingsProperties(p); err != nil {
										return err
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (sm *searchIndexProviderCore) UpsertIndex(scope *Scope, indexDefinition SearchIndex, opts *UpsertSearchIndexOptions) error {
	if opts == nil {
		opts = &UpsertSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if err := sm.validateIndexVectorMappings(indexDefinition); err != nil {
		return err
	}

	if indexDefinition.Name == "" {
		return invalidArgumentsError{"index name cannot be empty"}
	}
	if indexDefinition.Type == "" {
		return invalidArgumentsError{"index type cannot be empty"}
	}

	start := time.Now()
	defer sm.meter.ValueRecord(meterValueServiceManagement, "manager_search_upsert_index", start)

	path := fmt.Sprintf("%s/%s", sm.pathPrefix(scope), url.PathEscape(indexDefinition.Name))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_upsert_index", "management")
	span.SetAttribute("db.operation", "PUT "+path)
	defer span.End()

	indexData, err := indexDefinition.toData()
	if err != nil {
		return err
	}

	b, err := json.Marshal(indexData)
	if err != nil {
		return err
	}

	req := mgmtRequest{
		Service: ServiceTypeSearch,
		Method:  "PUT",
		Path:    path,
		Headers: map[string]string{
			"cache-control": "no-cache",
		},
		Body:          b,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := sm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return idxErr
		}

		return makeMgmtBadStatusError("failed to create index", &req, resp)
	}

	return nil
}

func (sm *searchIndexProviderCore) DropIndex(scope *Scope, indexName string, opts *DropSearchIndexOptions) error {
	if opts == nil {
		opts = &DropSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	start := time.Now()
	defer sm.meter.ValueRecord(meterValueServiceManagement, "manager_search_drop_index", start)

	path := fmt.Sprintf("%s/%s", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_drop_index", "management")
	span.SetAttribute("db.operation", "DELETE "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeSearch,
		Method:        "DELETE",
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := sm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return idxErr
		}

		return makeMgmtBadStatusError("failed to drop the index", &req, resp)
	}

	return nil
}

func (sm *searchIndexProviderCore) AnalyzeDocument(scope *Scope, indexName string, doc interface{}, opts *AnalyzeDocumentOptions) ([]interface{}, error) {
	if opts == nil {
		opts = &AnalyzeDocumentOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return nil, wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return nil, invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/analyzeDoc", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_analyze_document", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	b, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	req := mgmtRequest{
		Service:       ServiceTypeSearch,
		Method:        "POST",
		Path:          path,
		Body:          b,
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := sm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return nil, idxErr
		}

		return nil, makeMgmtBadStatusError("failed to analyze document", &req, resp)
	}

	var analysis struct {
		Status   string        `json:"status"`
		Analyzed []interface{} `json:"analyzed"`
	}
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&analysis)
	if err != nil {
		return nil, err
	}

	return analysis.Analyzed, nil
}

func (sm *searchIndexProviderCore) GetIndexedDocumentsCount(scope *Scope, indexName string, opts *GetIndexedDocumentsCountOptions) (uint64, error) {
	if opts == nil {
		opts = &GetIndexedDocumentsCountOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return 0, wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return 0, invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/count", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_get_indexed_documents_count", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeSearch,
		Method:        "GET",
		Path:          path,
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := sm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return 0, err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return 0, idxErr
		}

		return 0, makeMgmtBadStatusError("failed to get the indexed documents count", &req, resp)
	}

	var count struct {
		Count uint64 `json:"count"`
	}
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&count)
	if err != nil {
		return 0, err
	}

	return count.Count, nil
}

func (sm *searchIndexProviderCore) PauseIngest(scope *Scope, indexName string, opts *PauseIngestSearchIndexOptions) error {
	if opts == nil {
		opts = &PauseIngestSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/ingestControl/pause", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_pause_ingest", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	return sm.performControlRequest(
		opts.Context,
		span.Context(),
		"POST",
		path,
		opts.Timeout,
		opts.RetryStrategy)
}

func (sm *searchIndexProviderCore) ResumeIngest(scope *Scope, indexName string, opts *ResumeIngestSearchIndexOptions) error {
	if opts == nil {
		opts = &ResumeIngestSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/ingestControl/resume", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_resume_ingest", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	return sm.performControlRequest(
		opts.Context,
		span.Context(),
		"POST",
		path,
		opts.Timeout,
		opts.RetryStrategy)
}

func (sm *searchIndexProviderCore) AllowQuerying(scope *Scope, indexName string, opts *AllowQueryingSearchIndexOptions) error {
	if opts == nil {
		opts = &AllowQueryingSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/queryControl/allow", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_allow_querying", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	return sm.performControlRequest(
		opts.Context,
		span.Context(),
		"POST",
		path,
		opts.Timeout,
		opts.RetryStrategy)
}

func (sm *searchIndexProviderCore) DisallowQuerying(scope *Scope, indexName string, opts *DisallowQueryingSearchIndexOptions) error {
	if opts == nil {
		opts = &DisallowQueryingSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/queryControl/disallow", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_disallow_querying", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	return sm.performControlRequest(
		opts.Context,
		span.Context(),
		"POST",
		path,
		opts.Timeout,
		opts.RetryStrategy)
}

func (sm *searchIndexProviderCore) FreezePlan(scope *Scope, indexName string, opts *FreezePlanSearchIndexOptions) error {
	if opts == nil {
		opts = &FreezePlanSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/planFreezeControl/freeze", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_freeze_plan", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	return sm.performControlRequest(
		opts.Context,
		span.Context(),
		"POST",
		path,
		opts.Timeout,
		opts.RetryStrategy)
}

func (sm *searchIndexProviderCore) UnfreezePlan(scope *Scope, indexName string, opts *UnfreezePlanSearchIndexOptions) error {
	if opts == nil {
		opts = &UnfreezePlanSearchIndexOptions{}
	}

	if scope != nil && sm.scopedIndexesUnsupported() {
		return wrapError(ErrFeatureNotAvailable, "scoped indexes cannot be used with this server version")
	}

	if indexName == "" {
		return invalidArgumentsError{"indexName cannot be empty"}
	}

	path := fmt.Sprintf("%s/%s/planFreezeControl/unfreeze", sm.pathPrefix(scope), url.PathEscape(indexName))
	span := createSpan(sm.tracer, opts.ParentSpan, "manager_search_unfreeze_plan", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	return sm.performControlRequest(
		opts.Context,
		span.Context(),
		"POST",
		path,
		opts.Timeout,
		opts.RetryStrategy)
}

func (sm *searchIndexProviderCore) pathPrefix(scope *Scope) string {
	if scope == nil {
		return "/api/index"
	} else {
		return fmt.Sprintf("/api/bucket/%s/scope/%s/index", url.PathEscape(scope.bucket.bucketName), url.PathEscape(scope.scopeName))
	}
}

func (sm *searchIndexProviderCore) scopedIndexesUnsupported() bool {
	return sm.searchCapVerifier.SearchCapabilityStatus(gocbcore.SearchCapabilityScopedIndexes) == gocbcore.CapabilityStatusUnsupported
}

func (sm *searchIndexProviderCore) vectorSearchUnsupported() bool {
	return sm.searchCapVerifier.SearchCapabilityStatus(gocbcore.SearchCapabilityVectorSearch) == gocbcore.CapabilityStatusUnsupported
}

func (sm *searchIndexProviderCore) performControlRequest(
	ctx context.Context,
	tracectx RequestSpanContext,
	method, uri string,
	timeout time.Duration,
	retryStrategy RetryStrategy,
) error {
	req := mgmtRequest{
		Service:       ServiceTypeSearch,
		Method:        method,
		Path:          uri,
		IsIdempotent:  true,
		Timeout:       timeout,
		RetryStrategy: retryStrategy,
		parentSpanCtx: tracectx,
	}

	resp, err := sm.doMgmtRequest(ctx, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := sm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return idxErr
		}

		return makeMgmtBadStatusError("failed to perform the control request", &req, resp)
	}

	return nil
}

func (sm *searchIndexProviderCore) checkForRateLimitError(statusCode uint32, errMsg string) error {
	errMsg = strings.ToLower(errMsg)

	var err error
	if statusCode == 400 && strings.Contains(errMsg, "num_fts_indexes") {
		err = ErrQuotaLimitedFailure
	} else if statusCode == 429 {
		if strings.Contains(errMsg, "num_concurrent_requests") {
			err = ErrRateLimitedFailure
		} else if strings.Contains(errMsg, "num_queries_per_min") {
			err = ErrRateLimitedFailure
		} else if strings.Contains(errMsg, "ingress_mib_per_min") {
			err = ErrRateLimitedFailure
		} else if strings.Contains(errMsg, "egress_mib_per_min") {
			err = ErrRateLimitedFailure
		}
	}

	return err
}

func (sm *searchIndexProviderCore) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read search index response body: %s", err)
		return nil
	}

	if err := sm.checkForRateLimitError(resp.StatusCode, string(b)); err != nil {
		return makeGenericMgmtError(err, req, resp, string(b))
	}

	if resp.StatusCode == 404 {
		return makeGenericMgmtError(ErrFeatureNotAvailable, req, resp, "scoped indexes cannot be used with this server version")
	}

	var bodyErr error
	if strings.Contains(strings.ToLower(string(b)), "index not found") {
		bodyErr = ErrIndexNotFound
	} else if strings.Contains(strings.ToLower(string(b)), "index with the same name already exists") {
		bodyErr = ErrIndexExists
	} else {
		bodyErr = errors.New(string(b))
	}

	return makeGenericMgmtError(bodyErr, req, resp, string(b))
}

func (sm *searchIndexProviderCore) doMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error) {
	resp, err := sm.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type jsonSearchIndexResp struct {
	Status   string           `json:"status"`
	IndexDef *jsonSearchIndex `json:"indexDef"`
}

type jsonSearchIndexDefs struct {
	IndexDefs   map[string]jsonSearchIndex `json:"indexDefs"`
	ImplVersion string                     `json:"implVersion"`
}

type jsonSearchIndexesResp struct {
	Status    string              `json:"status"`
	IndexDefs jsonSearchIndexDefs `json:"indexDefs"`
}

type jsonSearchIndex struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	SourceName   string                 `json:"sourceName"`
	Type         string                 `json:"type"`
	Params       map[string]interface{} `json:"params"`
	SourceUUID   string                 `json:"sourceUUID"`
	SourceParams map[string]interface{} `json:"sourceParams"`
	SourceType   string                 `json:"sourceType"`
	PlanParams   map[string]interface{} `json:"planParams"`
}
