package gocb

import "time"

type kvTimeoutsConfig struct {
	KVTimeout        time.Duration
	KVDurableTimeout time.Duration
}

// Collection represents a single collection.
type Collection struct {
	collectionName string
	scope          string
	bucket         *Bucket

	timeoutsConfig kvTimeoutsConfig

	transcoder           Transcoder
	retryStrategyWrapper *retryStrategyWrapper
	tracer               requestTracer

	useMutationTokens bool

	getKvProvider func() (kvProvider, error)
}

func newCollection(scope *Scope, collectionName string) *Collection {
	return &Collection{
		collectionName: collectionName,
		scope:          scope.Name(),
		bucket:         scope.bucket,

		timeoutsConfig: scope.timeoutsConfig,

		transcoder:           scope.transcoder,
		retryStrategyWrapper: scope.retryStrategyWrapper,
		tracer:               scope.tracer,

		useMutationTokens: scope.useMutationTokens,

		getKvProvider: scope.getKvProvider,
	}
}

func (c *Collection) name() string {
	return c.collectionName
}

// ScopeName returns the name of the scope to which this collection belongs.
// UNCOMMITTED: This API may change in the future.
func (c *Collection) ScopeName() string {
	return c.scope
}

// Bucket returns the name of the bucket to which this collection belongs.
// UNCOMMITTED: This API may change in the future.
func (c *Collection) Bucket() *Bucket {
	return c.bucket
}

// Name returns the name of the collection.
func (c *Collection) Name() string {
	return c.collectionName
}

func (c *Collection) startKvOpTrace(operationName string, tracectx requestSpanContext) requestSpan {
	return c.tracer.StartSpan(operationName, tracectx).
		SetTag("couchbase.bucket", c.bucket).
		SetTag("couchbase.collection", c.collectionName).
		SetTag("couchbase.service", "kv")
}

func (c *Collection) bucketName() string {
	return c.bucket.Name()
}
