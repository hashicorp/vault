package gocb

import (
	"time"

	"github.com/couchbase/gocbcore/v9"
)

// Bucket represents a single bucket within a cluster.
type Bucket struct {
	bucketName string

	timeoutsConfig TimeoutsConfig

	transcoder           Transcoder
	retryStrategyWrapper *retryStrategyWrapper
	tracer               requestTracer

	useServerDurations bool
	useMutationTokens  bool

	bootstrapError    error
	connectionManager connectionManager
}

func newBucket(c *Cluster, bucketName string) *Bucket {
	return &Bucket{
		bucketName: bucketName,

		timeoutsConfig: c.timeoutsConfig,

		transcoder: c.transcoder,

		retryStrategyWrapper: c.retryStrategyWrapper,

		tracer: c.tracer,

		useServerDurations: c.useServerDurations,
		useMutationTokens:  c.useMutationTokens,

		connectionManager: c.connectionManager,
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

// Name returns the name of the bucket.
func (b *Bucket) Name() string {
	return b.bucketName
}

// Scope returns an instance of a Scope.
// VOLATILE: This API is subject to change at any time.
func (b *Bucket) Scope(scopeName string) *Scope {
	return newScope(b, scopeName)
}

// DefaultScope returns an instance of the default scope.
// VOLATILE: This API is subject to change at any time.
func (b *Bucket) DefaultScope() *Scope {
	return b.Scope("_default")
}

// Collection returns an instance of a collection from within the default scope.
// VOLATILE: This API is subject to change at any time.
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
		mgmtProvider: b,
		bucketName:   b.Name(),
		tracer:       b.tracer,
	}
}

// Collections provides functions for managing collections.
func (b *Bucket) Collections() *CollectionManager {
	// TODO: return error for unsupported collections
	return &CollectionManager{
		mgmtProvider: b,
		bucketName:   b.Name(),
		tracer:       b.tracer,
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
	if opts == nil {
		opts = &WaitUntilReadyOptions{}
	}

	if b.bootstrapError != nil {
		return b.bootstrapError
	}

	provider, err := b.connectionManager.getWaitUntilReadyProvider(b.bucketName)
	if err != nil {
		return err
	}

	desiredState := opts.DesiredState
	if desiredState == 0 {
		desiredState = ClusterStateOnline
	}

	services := opts.ServiceTypes
	gocbcoreServices := make([]gocbcore.ServiceType, len(services))
	for i, svc := range services {
		gocbcoreServices[i] = gocbcore.ServiceType(svc)
	}

	err = provider.WaitUntilReady(
		time.Now().Add(timeout),
		gocbcore.WaitUntilReadyOptions{
			DesiredState: gocbcore.ClusterState(desiredState),
			ServiceTypes: gocbcoreServices,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
