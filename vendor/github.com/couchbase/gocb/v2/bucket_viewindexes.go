package gocb

import (
	"context"
	"time"
)

// DesignDocumentNamespace represents which namespace a design document resides in.
type DesignDocumentNamespace uint

const (
	// DesignDocumentNamespaceProduction means that a design document resides in the production namespace.
	DesignDocumentNamespaceProduction DesignDocumentNamespace = iota

	// DesignDocumentNamespaceDevelopment means that a design document resides in the development namespace.
	DesignDocumentNamespaceDevelopment
)

// View represents a Couchbase view within a design document.
type View struct {
	Map    string
	Reduce string
}

func (v *View) fromData(data jsonView) error {
	v.Map = data.Map
	v.Reduce = data.Reduce

	return nil
}

func (v *View) toData() (jsonView, error) {
	var data jsonView

	data.Map = v.Map
	data.Reduce = v.Reduce

	return data, nil
}

// DesignDocument represents a Couchbase design document containing multiple views.
type DesignDocument struct {
	Name  string
	Views map[string]View
}

func (dd *DesignDocument) fromData(data jsonDesignDocument, name string) error {
	dd.Name = name

	views := make(map[string]View)
	for viewName, viewData := range data.Views {
		var view View
		err := view.fromData(viewData)
		if err != nil {
			return err
		}

		views[viewName] = view
	}
	dd.Views = views

	return nil
}

func (dd *DesignDocument) toData() (jsonDesignDocument, string, error) {
	var data jsonDesignDocument

	views := make(map[string]jsonView)
	for viewName, view := range dd.Views {
		viewData, err := view.toData()
		if err != nil {
			return jsonDesignDocument{}, "", err
		}

		views[viewName] = viewData
	}
	data.Views = views

	return data, dd.Name, nil
}

// ViewIndexManager provides methods for performing View management.
type ViewIndexManager struct {
	controller *providerController[viewIndexProvider]
}

// GetDesignDocumentOptions is the set of options available to the ViewIndexManager GetDesignDocument operation.
type GetDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetDesignDocument retrieves a single design document for the given bucket.
func (vm *ViewIndexManager) GetDesignDocument(name string, namespace DesignDocumentNamespace, opts *GetDesignDocumentOptions) (*DesignDocument, error) {
	return autoOpControl(vm.controller, func(provider viewIndexProvider) (*DesignDocument, error) {
		if opts == nil {
			opts = &GetDesignDocumentOptions{}
		}

		return provider.GetDesignDocument(name, namespace, opts)
	})
}

// GetAllDesignDocumentsOptions is the set of options available to the ViewIndexManager GetAllDesignDocuments operation.
type GetAllDesignDocumentsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllDesignDocuments will retrieve all design documents for the given bucket.
func (vm *ViewIndexManager) GetAllDesignDocuments(namespace DesignDocumentNamespace, opts *GetAllDesignDocumentsOptions) ([]DesignDocument, error) {
	return autoOpControl(vm.controller, func(provider viewIndexProvider) ([]DesignDocument, error) {
		if opts == nil {
			opts = &GetAllDesignDocumentsOptions{}
		}

		return provider.GetAllDesignDocuments(namespace, opts)
	})
}

// UpsertDesignDocumentOptions is the set of options available to the ViewIndexManager UpsertDesignDocument operation.
type UpsertDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpsertDesignDocument will insert a design document to the given bucket, or update
// an existing design document with the same name.
func (vm *ViewIndexManager) UpsertDesignDocument(ddoc DesignDocument, namespace DesignDocumentNamespace, opts *UpsertDesignDocumentOptions) error {
	return autoOpControlErrorOnly(vm.controller, func(provider viewIndexProvider) error {
		if opts == nil {
			opts = &UpsertDesignDocumentOptions{}
		}

		return provider.UpsertDesignDocument(ddoc, namespace, opts)
	})
}

// DropDesignDocumentOptions is the set of options available to the ViewIndexManager Upsert operation.
type DropDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropDesignDocument will remove a design document from the given bucket.
func (vm *ViewIndexManager) DropDesignDocument(name string, namespace DesignDocumentNamespace, opts *DropDesignDocumentOptions) error {
	return autoOpControlErrorOnly(vm.controller, func(provider viewIndexProvider) error {
		if opts == nil {
			opts = &DropDesignDocumentOptions{}
		}

		return provider.DropDesignDocument(name, namespace, opts)
	})
}

// PublishDesignDocumentOptions is the set of options available to the ViewIndexManager PublishDesignDocument operation.
type PublishDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// PublishDesignDocument publishes a design document to the given bucket.
func (vm *ViewIndexManager) PublishDesignDocument(name string, opts *PublishDesignDocumentOptions) error {
	return autoOpControlErrorOnly(vm.controller, func(provider viewIndexProvider) error {
		if opts == nil {
			opts = &PublishDesignDocumentOptions{}
		}

		return provider.PublishDesignDocument(name, opts)
	})
}
