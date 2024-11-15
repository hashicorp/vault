package gocb

import (
	"context"
	"time"
)

// CollectionHistorySettings specifies settings for whether history retention should be enabled or disabled
// for this collection.
type CollectionHistorySettings struct {
	Enabled bool
}

// CollectionSpec describes the specification of a collection.
type CollectionSpec struct {
	Name      string
	ScopeName string

	// The maximum expiry all documents in the collection can have. Defaults to the bucket-level setting.
	// Value of -1 seconds (time.Duration(-1) * time.Second)  denotes 'no expiry'.
	MaxExpiry time.Duration

	History *CollectionHistorySettings
}

// ScopeSpec describes the specification of a scope.
type ScopeSpec struct {
	Name        string
	Collections []CollectionSpec
}

// CollectionManager provides methods for performing collections management.
// Deprecated: See CollectionsV2 and CollectionManagerV2.
type CollectionManager struct {
	managerV2 *CollectionManagerV2
}

// GetAllScopesOptions is the set of options available to the GetAllScopes operation.
type GetAllScopesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllScopes gets all scopes from the bucket.
// Will be deprecated in favor of CollectionManagerV2.GetAllScopes in the next minor release.
func (cm *CollectionManager) GetAllScopes(opts *GetAllScopesOptions) ([]ScopeSpec, error) {
	return cm.managerV2.GetAllScopes(opts)
}

// CreateCollectionOptions is the set of options available to the CreateCollection operation.
type CreateCollectionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateCollection creates a new collection on the bucket.
// Will be deprecated in favor of CollectionManagerV2.CreateCollection in the next minor release.
func (cm *CollectionManager) CreateCollection(spec CollectionSpec, opts *CreateCollectionOptions) error {
	settings := &CreateCollectionSettings{
		MaxExpiry: spec.MaxExpiry,
		History:   spec.History,
	}

	return cm.managerV2.CreateCollection(spec.ScopeName, spec.Name, settings, opts)
}

// UpdateCollectionOptions is the set of options available to the UpdateCollection operation.
type UpdateCollectionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpdateCollection updates the settings of an existing collection.
// Will be deprecated in favor of CollectionManagerV2.UpdateCollection in the next minor release.
func (cm *CollectionManager) UpdateCollection(spec CollectionSpec, opts *UpdateCollectionOptions) error {
	settings := UpdateCollectionSettings{
		MaxExpiry: spec.MaxExpiry,
		History:   spec.History,
	}

	return cm.managerV2.UpdateCollection(spec.ScopeName, spec.Name, settings, opts)
}

// DropCollectionOptions is the set of options available to the DropCollection operation.
type DropCollectionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropCollection removes a collection.
// Will be deprecated in favor of CollectionManagerV2.DropCollection in the next minor release.
func (cm *CollectionManager) DropCollection(spec CollectionSpec, opts *DropCollectionOptions) error {
	return cm.managerV2.DropCollection(spec.ScopeName, spec.Name, opts)
}

// CreateScopeOptions is the set of options available to the CreateScope operation.
type CreateScopeOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateScope creates a new scope on the bucket.
// Will be deprecated in favor of CollectionManagerV2.CreateScope in the next minor release.
func (cm *CollectionManager) CreateScope(scopeName string, opts *CreateScopeOptions) error {
	return cm.managerV2.CreateScope(scopeName, opts)
}

// DropScopeOptions is the set of options available to the DropScope operation.
type DropScopeOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropScope removes a scope.
// Will be deprecated in favor of CollectionManagerV2.DropScope in the next minor release.
func (cm *CollectionManager) DropScope(scopeName string, opts *DropScopeOptions) error {
	return cm.managerV2.DropScope(scopeName, opts)
}
