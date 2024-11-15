package gocb

import (
	"context"
	"encoding/json"
	"time"
)

// SearchIndex is used to define a search index.
type SearchIndex struct {
	// UUID is required for updates. It provides a means of ensuring consistency, the UUID must match the UUID value
	// for the index on the server.
	UUID string
	// Name represents the name of this index.
	Name string
	// SourceName is the name of the source of the data for the index e.g. bucket name.
	SourceName string
	// Type is the type of index, e.g. fulltext-index or fulltext-alias.
	Type string
	// IndexParams are index properties such as store type and mappings.
	Params map[string]interface{}
	// SourceUUID is the UUID of the data source, this can be used to more tightly tie the index to a source.
	SourceUUID string
	// SourceParams are extra parameters to be defined. These are usually things like advanced connection and tuning
	// parameters.
	SourceParams map[string]interface{}
	// SourceType is the type of the data source, e.g. couchbase or nil depending on the Type field.
	SourceType string
	// PlanParams are plan properties such as number of replicas and number of partitions.
	PlanParams map[string]interface{}
}

func (si *SearchIndex) UnmarshalJSON(bytes []byte) error {
	var index jsonSearchIndex
	err := json.Unmarshal(bytes, &index)
	if err != nil {
		return err
	}

	return si.fromData(index)
}

func (si *SearchIndex) MarshalJSON() ([]byte, error) {
	index, err := si.toData()
	if err != nil {
		return nil, err
	}

	return json.Marshal(index)
}

func (si *SearchIndex) fromData(data jsonSearchIndex) error {
	si.UUID = data.UUID
	si.Name = data.Name
	si.SourceName = data.SourceName
	si.Type = data.Type
	si.Params = data.Params
	si.SourceUUID = data.SourceUUID
	si.SourceParams = data.SourceParams
	si.SourceType = data.SourceType
	si.PlanParams = data.PlanParams

	return nil
}

func (si *SearchIndex) toData() (jsonSearchIndex, error) {
	var data jsonSearchIndex

	data.UUID = si.UUID
	data.Name = si.Name
	data.SourceName = si.SourceName
	data.Type = si.Type
	data.Params = si.Params
	data.SourceUUID = si.SourceUUID
	data.SourceParams = si.SourceParams
	data.SourceType = si.SourceType
	data.PlanParams = si.PlanParams

	return data, nil
}

// SearchIndexManager provides methods for performing Couchbase search index management.
type SearchIndexManager struct {
	controller *providerController[searchIndexProvider]
}

// GetAllSearchIndexOptions is the set of options available to the search indexes GetAllIndexes operation.
type GetAllSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllIndexes retrieves all of the search indexes for the cluster.
func (sm *SearchIndexManager) GetAllIndexes(opts *GetAllSearchIndexOptions) ([]SearchIndex, error) {
	return autoOpControl(sm.controller, func(provider searchIndexProvider) ([]SearchIndex, error) {
		if opts == nil {
			opts = &GetAllSearchIndexOptions{}
		}

		return provider.GetAllIndexes(nil, opts)
	})
}

// GetSearchIndexOptions is the set of options available to the search indexes GetIndex operation.
type GetSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetIndex retrieves a specific search index by name.
func (sm *SearchIndexManager) GetIndex(indexName string, opts *GetSearchIndexOptions) (*SearchIndex, error) {
	return autoOpControl(sm.controller, func(provider searchIndexProvider) (*SearchIndex, error) {
		if opts == nil {
			opts = &GetSearchIndexOptions{}
		}

		if indexName == "" {
			return nil, invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.GetIndex(nil, indexName, opts)
	})
}

// UpsertSearchIndexOptions is the set of options available to the search index manager UpsertIndex operation.
type UpsertSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpsertIndex creates or updates a search index.
func (sm *SearchIndexManager) UpsertIndex(indexDefinition SearchIndex, opts *UpsertSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &UpsertSearchIndexOptions{}
		}

		if indexDefinition.Name == "" {
			return invalidArgumentsError{"index name cannot be empty"}
		}
		if indexDefinition.Type == "" {
			return invalidArgumentsError{"index type cannot be empty"}
		}

		return provider.UpsertIndex(nil, indexDefinition, opts)
	})
}

// DropSearchIndexOptions is the set of options available to the search index DropIndex operation.
type DropSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropIndex removes the search index with the specific name.
func (sm *SearchIndexManager) DropIndex(indexName string, opts *DropSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &DropSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.DropIndex(nil, indexName, opts)
	})
}

// AnalyzeDocumentOptions is the set of options available to the search index AnalyzeDocument operation.
type AnalyzeDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// AnalyzeDocument returns how a doc is analyzed against a specific index.
func (sm *SearchIndexManager) AnalyzeDocument(indexName string, doc interface{}, opts *AnalyzeDocumentOptions) ([]interface{}, error) {
	return autoOpControl(sm.controller, func(provider searchIndexProvider) ([]interface{}, error) {
		if opts == nil {
			opts = &AnalyzeDocumentOptions{}
		}

		if indexName == "" {
			return nil, invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.AnalyzeDocument(nil, indexName, doc, opts)
	})
}

// GetIndexedDocumentsCountOptions is the set of options available to the search index GetIndexedDocumentsCount operation.
type GetIndexedDocumentsCountOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetIndexedDocumentsCount retrieves the document count for a search index.
func (sm *SearchIndexManager) GetIndexedDocumentsCount(indexName string, opts *GetIndexedDocumentsCountOptions) (uint64, error) {
	return autoOpControl(sm.controller, func(provider searchIndexProvider) (uint64, error) {
		if opts == nil {
			opts = &GetIndexedDocumentsCountOptions{}
		}

		if indexName == "" {
			return 0, invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.GetIndexedDocumentsCount(nil, indexName, opts)
	})
}

// PauseIngestSearchIndexOptions is the set of options available to the search index PauseIngest operation.
type PauseIngestSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// PauseIngest pauses updates and maintenance for an index.
func (sm *SearchIndexManager) PauseIngest(indexName string, opts *PauseIngestSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &PauseIngestSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.PauseIngest(nil, indexName, opts)
	})
}

// ResumeIngestSearchIndexOptions is the set of options available to the search index ResumeIngest operation.
type ResumeIngestSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ResumeIngest resumes updates and maintenance for an index.
func (sm *SearchIndexManager) ResumeIngest(indexName string, opts *ResumeIngestSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &ResumeIngestSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.ResumeIngest(nil, indexName, opts)
	})
}

// AllowQueryingSearchIndexOptions is the set of options available to the search index AllowQuerying operation.
type AllowQueryingSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// AllowQuerying allows querying against an index.
func (sm *SearchIndexManager) AllowQuerying(indexName string, opts *AllowQueryingSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &AllowQueryingSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.AllowQuerying(nil, indexName, opts)
	})
}

// DisallowQueryingSearchIndexOptions is the set of options available to the search index DisallowQuerying operation.
type DisallowQueryingSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DisallowQuerying disallows querying against an index.
func (sm *SearchIndexManager) DisallowQuerying(indexName string, opts *AllowQueryingSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &AllowQueryingSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.DisallowQuerying(nil, indexName, &DisallowQueryingSearchIndexOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	})
}

// FreezePlanSearchIndexOptions is the set of options available to the search index FreezePlan operation.
type FreezePlanSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// FreezePlan freezes the assignment of index partitions to nodes.
func (sm *SearchIndexManager) FreezePlan(indexName string, opts *AllowQueryingSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &AllowQueryingSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.FreezePlan(nil, indexName, &FreezePlanSearchIndexOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	})
}

// UnfreezePlanSearchIndexOptions is the set of options available to the search index UnfreezePlan operation.
type UnfreezePlanSearchIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UnfreezePlan unfreezes the assignment of index partitions to nodes.
func (sm *SearchIndexManager) UnfreezePlan(indexName string, opts *AllowQueryingSearchIndexOptions) error {
	return autoOpControlErrorOnly(sm.controller, func(provider searchIndexProvider) error {
		if opts == nil {
			opts = &AllowQueryingSearchIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{"indexName cannot be empty"}
		}

		return provider.UnfreezePlan(nil, indexName, &UnfreezePlanSearchIndexOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	})
}
