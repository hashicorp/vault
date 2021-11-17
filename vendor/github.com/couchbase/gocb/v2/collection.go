package gocb

import (
	"time"
)

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
	tracer               RequestTracer
	meter                *meterWrapper

	useMutationTokens bool

	getKvProvider func() (kvProvider, error)
}

func newCollection(scope *Scope, collectionName string) *Collection {
	return &Collection{
		collectionName: collectionName,
		scope:          scope.Name(),
		bucket:         scope.bucket,

		timeoutsConfig: kvTimeoutsConfig{
			KVTimeout:        scope.timeoutsConfig.KVTimeout,
			KVDurableTimeout: scope.timeoutsConfig.KVDurableTimeout,
		},

		transcoder:           scope.transcoder,
		retryStrategyWrapper: scope.retryStrategyWrapper,
		tracer:               scope.tracer,
		meter:                scope.meter,

		useMutationTokens: scope.useMutationTokens,

		getKvProvider: scope.getKvProvider,
	}
}

func (c *Collection) name() string {
	return c.collectionName
}

// ScopeName returns the name of the scope to which this collection belongs.
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

func (c *Collection) startKvOpTrace(operationName string, tracectx RequestSpanContext, noAttributes bool) RequestSpan {
	span := c.tracer.RequestSpan(tracectx, operationName)
	if !noAttributes {
		span.SetAttribute(spanAttribDBNameKey, c.bucket.Name())
		span.SetAttribute(spanAttribDBCollectionNameKey, c.Name())
		span.SetAttribute(spanAttribDBScopeNameKey, c.ScopeName())
		span.SetAttribute(spanAttribServiceKey, "kv")
		span.SetAttribute(spanAttribOperationKey, operationName)
	}
	span.SetAttribute(spanAttribDBSystemKey, spanAttribDBSystemValue)

	return span
}

func (c *Collection) bucketName() string {
	return c.bucket.Name()
}
