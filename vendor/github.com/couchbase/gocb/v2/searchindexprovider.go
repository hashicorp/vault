package gocb

type searchIndexProvider interface {
	GetAllIndexes(scope *Scope, opts *GetAllSearchIndexOptions) ([]SearchIndex, error)
	GetIndex(scope *Scope, indexName string, opts *GetSearchIndexOptions) (*SearchIndex, error)
	UpsertIndex(scope *Scope, indexDefinition SearchIndex, opts *UpsertSearchIndexOptions) error
	DropIndex(scope *Scope, indexName string, opts *DropSearchIndexOptions) error
	AnalyzeDocument(scope *Scope, indexName string, doc interface{}, opts *AnalyzeDocumentOptions) ([]interface{}, error)
	GetIndexedDocumentsCount(scope *Scope, indexName string, opts *GetIndexedDocumentsCountOptions) (uint64, error)
	PauseIngest(scope *Scope, indexName string, opts *PauseIngestSearchIndexOptions) error
	ResumeIngest(scope *Scope, indexName string, opts *ResumeIngestSearchIndexOptions) error
	AllowQuerying(scope *Scope, indexName string, opts *AllowQueryingSearchIndexOptions) error
	DisallowQuerying(scope *Scope, indexName string, opts *DisallowQueryingSearchIndexOptions) error
	FreezePlan(scope *Scope, indexName string, opts *FreezePlanSearchIndexOptions) error
	UnfreezePlan(scope *Scope, indexName string, opts *UnfreezePlanSearchIndexOptions) error
}
