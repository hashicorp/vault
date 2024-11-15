package gocb

import (
	"time"
)

// Bucket represents a single bucket within a cluster.
type Bucket struct {
	bucketName string

	timeoutsConfig TimeoutsConfig

	transcoder           Transcoder
	retryStrategyWrapper *coreRetryStrategyWrapper
	compressor           *compressor

	useServerDurations bool
	useMutationTokens  bool

	bootstrapError    error
	connectionManager connectionManager
	getTransactions   func() *Transactions
}

func newBucket(c *Cluster, bucketName string) *Bucket {
	return &Bucket{
		bucketName: bucketName,

		timeoutsConfig: c.timeoutsConfig,

		transcoder: c.transcoder,

		retryStrategyWrapper: c.retryStrategyWrapper,

		compressor: c.compressor,

		useServerDurations: c.useServerDurations,
		useMutationTokens:  c.useMutationTokens,

		connectionManager: c.connectionManager,
		getTransactions:   c.Transactions,
	}
}

func (b *Bucket) setBootstrapError(err error) {
	b.bootstrapError = err
}

func (b *Bucket) getKvProvider() (kvProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	agent, err := b.connectionManager.getKvProvider(b.bucketName)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (b *Bucket) getKvCapabilitiesProvider() (kvCapabilityVerifier, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	agent, err := b.connectionManager.getKvCapabilitiesProvider(b.bucketName)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (b *Bucket) getKvBulkProvider() (kvBulkProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	agent, err := b.connectionManager.getKvBulkProvider(b.bucketName)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (b *Bucket) getQueryProvider() (queryProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	agent, err := b.connectionManager.getQueryProvider()
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (b *Bucket) getQueryIndexProvider() (queryIndexProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	provider, err := b.connectionManager.getQueryIndexProvider()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getSearchProvider() (searchProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	provider, err := b.connectionManager.getSearchProvider()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getSearchIndexProvider() (searchIndexProvider, error) {
	provider, err := b.connectionManager.getSearchIndexProvider()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getAnalyticsProvider() (analyticsProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	agent, err := b.connectionManager.getAnalyticsProvider()
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (b *Bucket) getEventingManagementProvider() (eventingManagementProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	provider, err := b.connectionManager.getEventingManagementProvider()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getViewProvider() (viewProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	provider, err := b.connectionManager.getViewProvider(b.Name())
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getViewIndexProvider() (viewIndexProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	provider, err := b.connectionManager.getViewIndexProvider(b.Name())
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getCollectionsManagementProvider() (collectionsManagementProvider, error) {
	if b.bootstrapError != nil {
		return nil, b.bootstrapError
	}

	provider, err := b.connectionManager.getCollectionsManagementProvider(b.Name())
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getDiagnosticsProvider() (diagnosticsProvider, error) {
	provider, err := b.connectionManager.getDiagnosticsProvider(b.Name())
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) getWaitUntilReadyProvider() (waitUntilReadyProvider, error) {
	provider, err := b.connectionManager.getWaitUntilReadyProvider(b.Name())
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bucket) diagnosticsController() *providerController[diagnosticsProvider] {
	return &providerController[diagnosticsProvider]{
		get:          b.getDiagnosticsProvider,
		opController: b.connectionManager,
	}
}

func (b *Bucket) viewController() *providerController[viewProvider] {
	return &providerController[viewProvider]{
		get:          b.getViewProvider,
		opController: b.connectionManager,
	}
}

func (b *Bucket) waitUntilReadyController() *providerController[waitUntilReadyProvider] {
	return &providerController[waitUntilReadyProvider]{
		get:          b.getWaitUntilReadyProvider,
		opController: b.connectionManager,
	}
}

// Name returns the name of the bucket.
func (b *Bucket) Name() string {
	return b.bucketName
}

// Scope returns an instance of a Scope.
func (b *Bucket) Scope(scopeName string) *Scope {
	return newScope(b, scopeName)
}

// DefaultScope returns an instance of the default scope.
func (b *Bucket) DefaultScope() *Scope {
	return b.Scope("_default")
}

// Collection returns an instance of a collection from within the default scope.
func (b *Bucket) Collection(collectionName string) *Collection {
	return b.DefaultScope().Collection(collectionName)
}

// DefaultCollection returns an instance of the default collection.
func (b *Bucket) DefaultCollection() *Collection {
	return b.DefaultScope().Collection("_default")
}

// ViewIndexes returns a ViewIndexManager instance for managing views.
func (b *Bucket) ViewIndexes() *ViewIndexManager {
	return &ViewIndexManager{
		controller: &providerController[viewIndexProvider]{
			get:          b.getViewIndexProvider,
			opController: b.connectionManager,
		},
	}
}

// CollectionsV2 provides functions for managing collections.
func (b *Bucket) CollectionsV2() *CollectionManagerV2 {
	return &CollectionManagerV2{
		controller: &providerController[collectionsManagementProvider]{
			get:          b.getCollectionsManagementProvider,
			opController: b.connectionManager,
		},
	}
}

// Collections provides functions for managing collections.
// Deprecated: See CollectionsV2.
func (b *Bucket) Collections() *CollectionManager {
	// TODO: return error for unsupported collections
	return &CollectionManager{
		managerV2: b.CollectionsV2(),
	}
}

// WaitUntilReady will wait for the bucket object to be ready for use.
// At present this will wait until memd connections have been established with the server and are ready
// to be used before performing a ping against the specified services (except KeyValue) which also
// exist in the cluster map.
// If no services are specified then will wait until KeyValue is ready.
// Valid service types are: ServiceTypeKeyValue, ServiceTypeManagement, ServiceTypeQuery, ServiceTypeSearch,
// ServiceTypeAnalytics, ServiceTypeViews.
func (b *Bucket) WaitUntilReady(timeout time.Duration, opts *WaitUntilReadyOptions) error {
	return autoOpControlErrorOnly(b.waitUntilReadyController(), func(provider waitUntilReadyProvider) error {
		if opts == nil {
			opts = &WaitUntilReadyOptions{}
		}

		if b.bootstrapError != nil {
			return b.bootstrapError
		}

		return provider.WaitUntilReady(
			opts.Context,
			time.Now().Add(timeout),
			opts,
		)
	})
}
