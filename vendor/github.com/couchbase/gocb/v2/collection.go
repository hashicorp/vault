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
	bucket         string

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
		bucket:         scope.bucketName,

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

func (c *Collection) scopeName() string {
	return c.scope
}

func (c *Collection) clone() *Collection {
	newC := *c
	return &newC
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
