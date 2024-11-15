package gocb

import (
	"time"
)

// CollectionQueryIndexManager provides methods for performing Couchbase query index management against collections.
// UNCOMMITTED: This API may change in the future.
type CollectionQueryIndexManager struct {
	controller *providerController[queryIndexProvider]

	c *Collection
}

func (qm *CollectionQueryIndexManager) validateScopeCollection(scope, collection string) error {
	if scope != "" || collection != "" {
		return makeInvalidArgumentsError("cannot use scope or collection with collection query index manager")
	}
	return nil
}

// CreateIndex creates an index over the specified fields.
// The SDK will automatically escape the provided index keys. For more advanced use cases like index keys using keywords
// scope.Query should be used with the query directly.
func (qm *CollectionQueryIndexManager) CreateIndex(indexName string, keys []string, opts *CreateQueryIndexOptions) error {
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

		return provider.CreateIndex(qm.c, "", indexName, keys, opts)
	})
}

// CreatePrimaryIndex creates a primary index.  An empty customName uses the default naming.
func (qm *CollectionQueryIndexManager) CreatePrimaryIndex(opts *CreatePrimaryQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &CreatePrimaryQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.CreatePrimaryIndex(qm.c, "", opts)
	})
}

// DropIndex drops a specific index by name.
func (qm *CollectionQueryIndexManager) DropIndex(indexName string, opts *DropQueryIndexOptions) error {
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

		return provider.DropIndex(qm.c, "", indexName, opts)
	})
}

// DropPrimaryIndex drops the primary index.  Pass an empty customName for unnamed primary indexes.
func (qm *CollectionQueryIndexManager) DropPrimaryIndex(opts *DropPrimaryQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &DropPrimaryQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.DropPrimaryIndex(qm.c, "", opts)
	})
}

// GetAllIndexes returns a list of all currently registered indexes.
func (qm *CollectionQueryIndexManager) GetAllIndexes(opts *GetAllQueryIndexesOptions) ([]QueryIndex, error) {
	return autoOpControl(qm.controller, func(provider queryIndexProvider) ([]QueryIndex, error) {
		if opts == nil {
			opts = &GetAllQueryIndexesOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return nil, err
		}

		return provider.GetAllIndexes(qm.c, "", opts)
	})
}

// BuildDeferredIndexes builds all indexes which are currently in deferred state.
// If no collection and scope names are specified in the options then *only* indexes created on the bucket directly
// will be built.
func (qm *CollectionQueryIndexManager) BuildDeferredIndexes(opts *BuildDeferredQueryIndexOptions) ([]string, error) {
	return autoOpControl(qm.controller, func(provider queryIndexProvider) ([]string, error) {
		if opts == nil {
			opts = &BuildDeferredQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return nil, err
		}

		return provider.BuildDeferredIndexes(qm.c, "", opts)
	})
}

// WatchIndexes waits for a set of indexes to come online.
func (qm *CollectionQueryIndexManager) WatchIndexes(watchList []string, timeout time.Duration, opts *WatchQueryIndexOptions) error {
	return autoOpControlErrorOnly(qm.controller, func(provider queryIndexProvider) error {
		if opts == nil {
			opts = &WatchQueryIndexOptions{}
		}
		if err := qm.validateScopeCollection(opts.ScopeName, opts.CollectionName); err != nil {
			return err
		}

		return provider.WatchIndexes(qm.c, "", watchList, timeout, opts)
	})
}
