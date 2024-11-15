package gocb

import (
	"context"
	"time"
)

// AnalyticsIndexManager provides methods for performing Couchbase Analytics index management.
type AnalyticsIndexManager struct {
	controller *providerController[analyticsIndexProvider]
}

// AnalyticsDataset contains information about an analytics dataset.
type AnalyticsDataset struct {
	Name          string
	DataverseName string
	LinkName      string
	BucketName    string
}

func (ad *AnalyticsDataset) fromData(data jsonAnalyticsDataset) error {
	ad.Name = data.DatasetName
	ad.DataverseName = data.DataverseName
	ad.LinkName = data.LinkName
	ad.BucketName = data.BucketName

	return nil
}

// AnalyticsIndex contains information about an analytics index.
type AnalyticsIndex struct {
	Name          string
	DatasetName   string
	DataverseName string
	IsPrimary     bool
}

func (ai *AnalyticsIndex) fromData(data jsonAnalyticsIndex) error {
	ai.Name = data.IndexName
	ai.DatasetName = data.DatasetName
	ai.DataverseName = data.DataverseName
	ai.IsPrimary = data.IsPrimary

	return nil
}

// CreateAnalyticsDataverseOptions is the set of options available to the AnalyticsManager CreateDataverse operation.
type CreateAnalyticsDataverseOptions struct {
	IgnoreIfExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateDataverse creates a new analytics dataset.
func (am *AnalyticsIndexManager) CreateDataverse(dataverseName string, opts *CreateAnalyticsDataverseOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &CreateAnalyticsDataverseOptions{}
		}

		if dataverseName == "" {
			return invalidArgumentsError{
				message: "dataset name cannot be empty",
			}
		}

		return provider.CreateDataverse(dataverseName, opts)
	})
}

// DropAnalyticsDataverseOptions is the set of options available to the AnalyticsManager DropDataverse operation.
type DropAnalyticsDataverseOptions struct {
	IgnoreIfNotExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropDataverse drops an analytics dataset.
func (am *AnalyticsIndexManager) DropDataverse(dataverseName string, opts *DropAnalyticsDataverseOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &DropAnalyticsDataverseOptions{}
		}

		return provider.DropDataverse(dataverseName, opts)
	})
}

// CreateAnalyticsDatasetOptions is the set of options available to the AnalyticsManager CreateDataset operation.
type CreateAnalyticsDatasetOptions struct {
	IgnoreIfExists bool
	Condition      string
	DataverseName  string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateDataset creates a new analytics dataset.
func (am *AnalyticsIndexManager) CreateDataset(datasetName, bucketName string, opts *CreateAnalyticsDatasetOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &CreateAnalyticsDatasetOptions{}
		}

		if datasetName == "" {
			return invalidArgumentsError{
				message: "dataset name cannot be empty",
			}
		}

		return provider.CreateDataset(datasetName, bucketName, opts)
	})
}

// DropAnalyticsDatasetOptions is the set of options available to the AnalyticsManager DropDataset operation.
type DropAnalyticsDatasetOptions struct {
	IgnoreIfNotExists bool
	DataverseName     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropDataset drops an analytics dataset.
func (am *AnalyticsIndexManager) DropDataset(datasetName string, opts *DropAnalyticsDatasetOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &DropAnalyticsDatasetOptions{}
		}

		return provider.DropDataset(datasetName, opts)
	})
}

// GetAllAnalyticsDatasetsOptions is the set of options available to the AnalyticsManager GetAllDatasets operation.
type GetAllAnalyticsDatasetsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllDatasets gets all analytics datasets.
func (am *AnalyticsIndexManager) GetAllDatasets(opts *GetAllAnalyticsDatasetsOptions) ([]AnalyticsDataset, error) {
	return autoOpControl(am.controller, func(provider analyticsIndexProvider) ([]AnalyticsDataset, error) {
		if opts == nil {
			opts = &GetAllAnalyticsDatasetsOptions{}
		}

		return provider.GetAllDatasets(opts)
	})
}

// CreateAnalyticsIndexOptions is the set of options available to the AnalyticsManager CreateIndex operation.
type CreateAnalyticsIndexOptions struct {
	IgnoreIfExists bool
	DataverseName  string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateIndex creates a new analytics dataset.
func (am *AnalyticsIndexManager) CreateIndex(datasetName, indexName string, fields map[string]string, opts *CreateAnalyticsIndexOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &CreateAnalyticsIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{
				message: "index name cannot be empty",
			}
		}
		if len(fields) <= 0 {
			return invalidArgumentsError{
				message: "you must specify at least one field to index",
			}
		}

		return provider.CreateIndex(datasetName, indexName, fields, opts)
	})
}

// DropAnalyticsIndexOptions is the set of options available to the AnalyticsManager DropIndex operation.
type DropAnalyticsIndexOptions struct {
	IgnoreIfNotExists bool
	DataverseName     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropIndex drops an analytics index.
func (am *AnalyticsIndexManager) DropIndex(datasetName, indexName string, opts *DropAnalyticsIndexOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &DropAnalyticsIndexOptions{}
		}

		return provider.DropIndex(datasetName, indexName, opts)
	})
}

// GetAllAnalyticsIndexesOptions is the set of options available to the AnalyticsManager GetAllIndexes operation.
type GetAllAnalyticsIndexesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllIndexes gets all analytics indexes.
func (am *AnalyticsIndexManager) GetAllIndexes(opts *GetAllAnalyticsIndexesOptions) ([]AnalyticsIndex, error) {
	return autoOpControl(am.controller, func(provider analyticsIndexProvider) ([]AnalyticsIndex, error) {
		if opts == nil {
			opts = &GetAllAnalyticsIndexesOptions{}
		}

		return provider.GetAllIndexes(opts)
	})
}

// ConnectAnalyticsLinkOptions is the set of options available to the AnalyticsManager ConnectLink operation.
type ConnectAnalyticsLinkOptions struct {
	LinkName      string
	DataverseName string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ConnectLink connects an analytics link.
func (am *AnalyticsIndexManager) ConnectLink(opts *ConnectAnalyticsLinkOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &ConnectAnalyticsLinkOptions{}
		}

		return provider.ConnectLink(opts)
	})
}

// DisconnectAnalyticsLinkOptions is the set of options available to the AnalyticsManager DisconnectLink operation.
type DisconnectAnalyticsLinkOptions struct {
	LinkName      string
	DataverseName string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DisconnectLink disconnects an analytics link.
func (am *AnalyticsIndexManager) DisconnectLink(opts *DisconnectAnalyticsLinkOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &DisconnectAnalyticsLinkOptions{}
		}

		return provider.DisconnectLink(opts)
	})
}

// GetPendingMutationsAnalyticsOptions is the set of options available to the user manager GetPendingMutations operation.
type GetPendingMutationsAnalyticsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetPendingMutations returns the number of pending mutations for all indexes in the form of dataverse.dataset:mutations.
func (am *AnalyticsIndexManager) GetPendingMutations(opts *GetPendingMutationsAnalyticsOptions) (map[string]map[string]int, error) {
	return autoOpControl(am.controller, func(provider analyticsIndexProvider) (map[string]map[string]int, error) {
		if opts == nil {
			opts = &GetPendingMutationsAnalyticsOptions{}
		}

		return provider.GetPendingMutations(opts)
	})
}

// AnalyticsLink describes an external or remote analytics link, used to access data external to the cluster.
type AnalyticsLink interface {
	// Name returns the name of this link.
	Name() string
	// DataverseName returns the name of the dataverse that this link belongs to.
	DataverseName() string
	// FormEncode encodes the link into a form data representation, to be sent as the body of a CreateLink or ReplaceLink
	// request.
	FormEncode() ([]byte, error)
	// Validate is used by CreateLink and ReplaceLink to ensure that the link is valid.
	Validate() error
	// LinkType returns the type of analytics type this link is.
	LinkType() AnalyticsLinkType
}

// NewCouchbaseRemoteAnalyticsLinkOptions are the options available when creating a new CouchbaseRemoteAnalyticsLink.
type NewCouchbaseRemoteAnalyticsLinkOptions struct {
	Encryption CouchbaseRemoteAnalyticsEncryptionSettings
	Username   string
	Password   string
}

// NewCouchbaseRemoteAnalyticsLink creates a new CouchbaseRemoteAnalyticsLink.
// Scope is the analytics scope in the form of "bucket/scope".
func NewCouchbaseRemoteAnalyticsLink(linkName, hostname, dataverseName string,
	opts *NewCouchbaseRemoteAnalyticsLinkOptions) *CouchbaseRemoteAnalyticsLink {
	if opts == nil {
		opts = &NewCouchbaseRemoteAnalyticsLinkOptions{}
	}
	return &CouchbaseRemoteAnalyticsLink{
		Dataverse:  dataverseName,
		LinkName:   linkName,
		Hostname:   hostname,
		Encryption: opts.Encryption,
		Username:   opts.Username,
		Password:   opts.Password,
	}
}

// NewS3ExternalAnalyticsLinkOptions are the options available when creating a new S3ExternalAnalyticsLink.
type NewS3ExternalAnalyticsLinkOptions struct {
	SessionToken    string
	ServiceEndpoint string
}

// NewS3ExternalAnalyticsLink creates a new S3ExternalAnalyticsLink with the scope field populated.
// Scope is the analytics scope in the form of "bucket/scope".
func NewS3ExternalAnalyticsLink(linkName, dataverseName, accessKeyID, secretAccessKey, region string,
	opts *NewS3ExternalAnalyticsLinkOptions) *S3ExternalAnalyticsLink {
	if opts == nil {
		opts = &NewS3ExternalAnalyticsLinkOptions{}
	}
	return &S3ExternalAnalyticsLink{
		Dataverse:       dataverseName,
		LinkName:        linkName,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    opts.SessionToken,
		Region:          region,
		ServiceEndpoint: opts.ServiceEndpoint,
	}
}

// NewAzureBlobExternalAnalyticsLinkOptions are the options available when creating a new AzureBlobExternalAnalyticsLink.
// VOLATILE: This API is subject to change at any time.
type NewAzureBlobExternalAnalyticsLinkOptions struct {
	ConnectionString      string
	AccountName           string
	AccountKey            string
	SharedAccessSignature string
	BlobEndpoint          string
	EndpointSuffix        string
}

// NewAzureBlobExternalAnalyticsLink creates a new AzureBlobExternalAnalyticsLink.
// VOLATILE: This API is subject to change at any time.
func NewAzureBlobExternalAnalyticsLink(linkName, dataverseName string,
	opts *NewAzureBlobExternalAnalyticsLinkOptions) *AzureBlobExternalAnalyticsLink {
	if opts == nil {
		opts = &NewAzureBlobExternalAnalyticsLinkOptions{}
	}
	return &AzureBlobExternalAnalyticsLink{
		Dataverse:             dataverseName,
		LinkName:              linkName,
		ConnectionString:      opts.ConnectionString,
		AccountName:           opts.AccountName,
		AccountKey:            opts.AccountKey,
		SharedAccessSignature: opts.SharedAccessSignature,
		BlobEndpoint:          opts.BlobEndpoint,
		EndpointSuffix:        opts.EndpointSuffix,
	}
}

// CreateAnalyticsLinkOptions is the set of options available to the analytics manager CreateLink
// function.
type CreateAnalyticsLinkOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateLink creates an analytics link.
func (am *AnalyticsIndexManager) CreateLink(link AnalyticsLink, opts *CreateAnalyticsLinkOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &CreateAnalyticsLinkOptions{}
		}

		return provider.CreateLink(link, opts)
	})
}

// ReplaceAnalyticsLinkOptions is the set of options available to the analytics manager ReplaceLink
// function.
type ReplaceAnalyticsLinkOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ReplaceLink modifies an existing analytics link.
func (am *AnalyticsIndexManager) ReplaceLink(link AnalyticsLink, opts *ReplaceAnalyticsLinkOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &ReplaceAnalyticsLinkOptions{}
		}

		return provider.ReplaceLink(link, opts)
	})
}

// DropAnalyticsLinkOptions is the set of options available to the analytics manager DropLink
// function.
type DropAnalyticsLinkOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropLink removes an existing external analytics link from specified scope.
// dataverseName can be given in the form of "namepart" or "namepart1/namepart2".
// Only available against Couchbase Server 7.0+.
func (am *AnalyticsIndexManager) DropLink(linkName, dataverseName string, opts *DropAnalyticsLinkOptions) error {
	return autoOpControlErrorOnly(am.controller, func(provider analyticsIndexProvider) error {
		if opts == nil {
			opts = &DropAnalyticsLinkOptions{}
		}

		return provider.DropLink(linkName, dataverseName, opts)
	})
}

// GetAnalyticsLinksOptions are the options available to the AnalyticsManager GetLinks function.
type GetAnalyticsLinksOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Dataverse restricts the results to a given dataverse, can be given in the form of "namepart" or "namepart1/namepart2".
	Dataverse string
	// LinkType restricts the results to the given link type.
	LinkType AnalyticsLinkType
	// Name restricts the results to the link with the specified name.
	// If set then `Scope` must also be set.
	Name string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetLinks retrieves all external or remote analytics links.
func (am *AnalyticsIndexManager) GetLinks(opts *GetAnalyticsLinksOptions) ([]AnalyticsLink, error) {
	return autoOpControl(am.controller, func(provider analyticsIndexProvider) ([]AnalyticsLink, error) {
		if opts == nil {
			opts = &GetAnalyticsLinksOptions{}
		}

		if opts.Name != "" && opts.Dataverse == "" {
			return nil, makeInvalidArgumentsError("when name is set then dataverse must also be set")
		}

		return provider.GetLinks(opts)
	})
}
