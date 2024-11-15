package gocb

type analyticsIndexProvider interface {
	CreateDataverse(dataverseName string, opts *CreateAnalyticsDataverseOptions) error
	DropDataverse(dataverseName string, opts *DropAnalyticsDataverseOptions) error
	CreateDataset(datasetName, bucketName string, opts *CreateAnalyticsDatasetOptions) error
	DropDataset(datasetName string, opts *DropAnalyticsDatasetOptions) error
	GetAllDatasets(opts *GetAllAnalyticsDatasetsOptions) ([]AnalyticsDataset, error)
	CreateIndex(datasetName, indexName string, fields map[string]string, opts *CreateAnalyticsIndexOptions) error
	DropIndex(datasetName, indexName string, opts *DropAnalyticsIndexOptions) error
	GetAllIndexes(opts *GetAllAnalyticsIndexesOptions) ([]AnalyticsIndex, error)
	ConnectLink(opts *ConnectAnalyticsLinkOptions) error
	DisconnectLink(opts *DisconnectAnalyticsLinkOptions) error
	GetPendingMutations(opts *GetPendingMutationsAnalyticsOptions) (map[string]map[string]int, error)
	CreateLink(link AnalyticsLink, opts *CreateAnalyticsLinkOptions) error
	ReplaceLink(link AnalyticsLink, opts *ReplaceAnalyticsLinkOptions) error
	DropLink(linkName, dataverseName string, opts *DropAnalyticsLinkOptions) error
	GetLinks(opts *GetAnalyticsLinksOptions) ([]AnalyticsLink, error)
}
