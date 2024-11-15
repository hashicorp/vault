package gocb

import (
	"context"
	"time"
)

// QueryIndexManager provides methods for performing Couchbase query index management.
type QueryIndexManager struct {
	controller *providerController[queryIndexProvider]
}

func (qm *QueryIndexManager) validateScopeCollection(scope, collection string) error {
	if scope == "" && collection != "" {
		return makeInvalidArgumentsError("if collection is set then scope must be set")
	} else if scope != "" && collection == "" {
		return makeInvalidArgumentsError("if scope is set then collection must be set")
	}

	return nil
}

// CreateQueryIndexOptions is the set of options available to the query indexes CreateIndex operation.
type CreateQueryIndexOptions struct {
	IgnoreIfExists bool
	Deferred       bool
	NumReplicas    int

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateIndex creates an index over the specified fields.
// The SDK will automatically escape the provided index keys. For more advanced use cases like index keys using keywords
// cluster.Query or scope.Query should be used with the query directly.
func (qm *QueryIndexManager) CreateIndex(bucketName, indexName string, keys []string, opts *CreateQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &CreateQueryIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{
				message: "an invalid index name was specified",
			}
		}
		if len(keys) <= 0 {
			return invalidArgumentsError{
				message: "you must specify at least one index-key to index",
			}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.CreateIndex(nil, bucketName, indexName, keys, opts)
	})
}

// CreatePrimaryQueryIndexOptions is the set of options available to the query indexes CreatePrimaryIndex operation.
type CreatePrimaryQueryIndexOptions struct {
	IgnoreIfExists bool
	Deferred       bool
	CustomName     string
	NumReplicas    int

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreatePrimaryIndex creates a primary index.  An empty customName uses the default naming.
func (qm *QueryIndexManager) CreatePrimaryIndex(bucketName string, opts *CreatePrimaryQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &CreatePrimaryQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.CreatePrimaryIndex(nil, bucketName, opts)
	})
}

// DropQueryIndexOptions is the set of options available to the query indexes DropIndex operation.
type DropQueryIndexOptions struct {
	IgnoreIfNotExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropIndex drops a specific index by name.
func (qm *QueryIndexManager) DropIndex(bucketName, indexName string, opts *DropQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &DropQueryIndexOptions{}
		}

		if indexName == "" {
			return invalidArgumentsError{
				message: "an invalid index name was specified",
			}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.DropIndex(nil, bucketName, indexName, opts)
	})
}

// DropPrimaryQueryIndexOptions is the set of options available to the query indexes DropPrimaryIndex operation.
type DropPrimaryQueryIndexOptions struct {
	IgnoreIfNotExists bool
	CustomName        string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropPrimaryIndex drops the primary index.  Pass an empty customName for unnamed primary indexes.
func (qm *QueryIndexManager) DropPrimaryIndex(bucketName string, opts *DropPrimaryQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &DropPrimaryQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.DropPrimaryIndex(nil, bucketName, opts)
	})
}

// GetAllQueryIndexesOptions is the set of options available to the query indexes GetAllIndexes operation.
type GetAllQueryIndexesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllIndexes returns a list of all currently registered indexes.
func (qm *QueryIndexManager) GetAllIndexes(bucketName string, opts *GetAllQueryIndexesOptions) ([]QueryIndex, error) {
	return autoOpControl(qm.controller, func(provider queryIndexProvider) ([]QueryIndex, error) {
		if opts == nil {
			opts = &GetAllQueryIndexesOptions{}
		}

		return provider.GetAllIndexes(nil, bucketName, opts)
	})
}

// BuildDeferredQueryIndexOptions is the set of options available to the query indexes BuildDeferredIndexes operation.
type BuildDeferredQueryIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// BuildDeferredIndexes builds all indexes which are currently in deferred state.
// If no collection and scope names are specified in the options then *only* indexes created on the bucket directly
// will be built.
func (qm *QueryIndexManager) BuildDeferredIndexes(bucketName string, opts *BuildDeferredQueryIndexOptions) ([]string, error) {
	return autoOpControl(qm.controller, func(provider queryIndexProvider) ([]string, error) {
		if opts == nil {
			opts = &BuildDeferredQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return nil, err
		}

		return provider.BuildDeferredIndexes(nil, bucketName, opts)
	})
}

// WatchQueryIndexOptions is the set of options available to the query indexes Watch operation.
type WatchQueryIndexOptions struct {
	WatchPrimary bool

	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Deprecated: See CollectionQueryIndexManager.
	ScopeName string
	// Deprecated: See CollectionQueryIndexManager.
	CollectionName string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// WatchIndexes waits for a set of indexes to come online.
func (qm *QueryIndexManager) WatchIndexes(bucketName string, watchList []string, timeout time.Duration, opts *WatchQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &WatchQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}
		return provider.WatchIndexes(nil, bucketName, watchList, timeout, opts)
	})
}
