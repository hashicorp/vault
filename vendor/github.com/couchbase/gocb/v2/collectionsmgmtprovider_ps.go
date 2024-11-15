package gocb

import (
	"time"

	"github.com/couchbase/goprotostellar/genproto/admin_collection_v1"
)

type collectionsManagementProviderPs struct {
	provider admin_collection_v1.CollectionAdminServiceClient

	managerProvider *psOpManagerProvider
	bucketName      string
}

func (cm collectionsManagementProviderPs) newOpManager(parentSpan RequestSpan, opName string, attribs map[string]interface{}) *psOpManagerDefault {
	return cm.managerProvider.NewManager(parentSpan, opName, attribs)
}

func (cm *collectionsManagementProviderPs) GetAllScopes(opts *GetAllScopesOptions) ([]ScopeSpec, error) {
	manager := cm.newOpManager(opts.ParentSpan, "manager_collections_get_all_scopes", map[string]interface{}{
		"db.name":      cm.bucketName,
		"db.operation": "ListCollections",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	req := &admin_collection_v1.ListCollectionsRequest{
		BucketName: cm.bucketName,
	}

	resp, err := wrapPSOp(manager, req, cm.provider.ListCollections)
	if err != nil {
		return nil, err
	}

	var scopes []ScopeSpec
	for _, scope := range resp.GetScopes() {
		var collections []CollectionSpec
		for _, col := range scope.Collections {

			var expiry int32
			// Since protostellar does not support negative values for MaxExpirySecs no_expiry and bucket expiry inheritance are
			// indicated with MaxExpirySecs = 0, and MaxExpirySecs = nil respectively
			if col.MaxExpirySecs == nil {
				expiry = 0
			} else if col.GetMaxExpirySecs() == 0 {
				expiry = -1
			} else {
				expiry = int32(col.GetMaxExpirySecs())
			}

			collections = append(collections, CollectionSpec{
				Name:      col.Name,
				ScopeName: scope.Name,
				MaxExpiry: time.Duration(expiry) * time.Second,
			})
		}
		scopes = append(scopes, ScopeSpec{
			Name:        scope.Name,
			Collections: collections,
		})
	}

	return scopes, nil
}

// CreateCollection creates a new collection on the bucket.
func (cm *collectionsManagementProviderPs) CreateCollection(scopeName string, collectionName string, settings *CreateCollectionSettings, opts *CreateCollectionOptions) error {
	manager := cm.newOpManager(opts.ParentSpan, "manager_collections_create_collection", map[string]interface{}{
		"db.name":                 cm.bucketName,
		"db.couchbase.scope":      scopeName,
		"db.couchbase.collection": collectionName,
		"db.operation":            "CreateCollection",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_collection_v1.CreateCollectionRequest{
		BucketName:     cm.bucketName,
		ScopeName:      scopeName,
		CollectionName: collectionName,
	}
	if settings.MaxExpiry > 0 {
		expiry := uint32(settings.MaxExpiry.Seconds())
		req.MaxExpirySecs = &expiry
	} else if settings.MaxExpiry < 0 {
		// Since protostellar does not support negative values for MaxExpirySecs no_expiry and bucket expiry inheritance are
		// indicated with MaxExpirySecs = 0, and MaxExpirySecs = nil respectively
		expiry := uint32(0)
		req.MaxExpirySecs = &expiry
	}

	_, err := wrapPSOp(manager, req, cm.provider.CreateCollection)
	if err != nil {
		return err
	}

	return nil
}

func (cm *collectionsManagementProviderPs) UpdateCollection(scopeName string, collectionName string, settings UpdateCollectionSettings, opts *UpdateCollectionOptions) error {
	manager := cm.newOpManager(opts.ParentSpan, "manager_collections_update_collection", map[string]interface{}{
		"db.name":                 cm.bucketName,
		"db.couchbase.scope":      scopeName,
		"db.couchbase.collection": collectionName,
		"db.operation":            "UpdateCollection",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_collection_v1.UpdateCollectionRequest{
		BucketName:     cm.bucketName,
		ScopeName:      scopeName,
		CollectionName: collectionName,
	}
	if settings.MaxExpiry > 0 {
		expiry := uint32(settings.MaxExpiry.Seconds())
		req.MaxExpirySecs = &expiry
	} else if settings.MaxExpiry < 0 {
		// Since protostellar does not support negative values for MaxExpirySecs no_expiry is indicated with MaxExpirySecs = 0
		expiry := uint32(0)
		req.MaxExpirySecs = &expiry
	}

	if settings.History != nil {
		req.HistoryRetentionEnabled = &settings.History.Enabled
	}

	_, err := wrapPSOp(manager, req, cm.provider.UpdateCollection)
	if err != nil {
		return err
	}

	return nil
}

// DropCollection removes a collection.
func (cm *collectionsManagementProviderPs) DropCollection(scopeName string, collectionName string, opts *DropCollectionOptions) error {
	manager := cm.newOpManager(opts.ParentSpan, "manager_collections_drop_collection", map[string]interface{}{
		"db.name":                 cm.bucketName,
		"db.couchbase.scope":      scopeName,
		"db.couchbase.collection": collectionName,
		"db.operation":            "DeleteCollection",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_collection_v1.DeleteCollectionRequest{
		BucketName:     cm.bucketName,
		ScopeName:      scopeName,
		CollectionName: collectionName,
	}

	_, err := wrapPSOp(manager, req, cm.provider.DeleteCollection)
	if err != nil {
		return err
	}

	return nil
}

// CreateScope creates a new scope on the bucket.
func (cm *collectionsManagementProviderPs) CreateScope(scopeName string, opts *CreateScopeOptions) error {
	manager := cm.newOpManager(opts.ParentSpan, "manager_collections_create_scope", map[string]interface{}{
		"db.name":            cm.bucketName,
		"db.couchbase.scope": scopeName,
		"db.operation":       "CreateScope",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_collection_v1.CreateScopeRequest{
		BucketName: cm.bucketName,
		ScopeName:  scopeName,
	}

	_, err := wrapPSOp(manager, req, cm.provider.CreateScope)
	if err != nil {
		return err
	}

	return nil
}

// DropScope removes a scope.
func (cm *collectionsManagementProviderPs) DropScope(scopeName string, opts *DropScopeOptions) error {
	manager := cm.newOpManager(opts.ParentSpan, "manager_collections_drop_scope", map[string]interface{}{
		"db.name":            cm.bucketName,
		"db.couchbase.scope": scopeName,
		"db.operation":       "DeleteScope",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_collection_v1.DeleteScopeRequest{
		BucketName: cm.bucketName,
		ScopeName:  scopeName,
	}

	_, err := wrapPSOp(manager, req, cm.provider.DeleteScope)
	if err != nil {
		return err
	}

	return nil
}
