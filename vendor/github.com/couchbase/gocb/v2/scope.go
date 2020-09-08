package gocb

// Scope represents a single scope within a bucket.
// VOLATILE: This API is subject to change at any time.
type Scope struct {
	scopeName string
	bucket    *Bucket

	timeoutsConfig kvTimeoutsConfig

	transcoder           Transcoder
	retryStrategyWrapper *retryStrategyWrapper
	tracer               requestTracer

	useMutationTokens bool

	getKvProvider func() (kvProvider, error)
}

func newScope(bucket *Bucket, scopeName string) *Scope {
	return &Scope{
		scopeName: scopeName,
		bucket:    bucket,

		timeoutsConfig: kvTimeoutsConfig{
			KVTimeout:        bucket.timeoutsConfig.KVTimeout,
			KVDurableTimeout: bucket.timeoutsConfig.KVDurableTimeout,
		},

		transcoder:           bucket.transcoder,
		retryStrategyWrapper: bucket.retryStrategyWrapper,
		tracer:               bucket.tracer,

		useMutationTokens: bucket.useMutationTokens,

		getKvProvider: bucket.getKvProvider,
	}
}

// Name returns the name of the scope.
func (s *Scope) Name() string {
	return s.scopeName
}

// BucketName returns the name of the bucket to which this collection belongs.
// UNCOMMITTED: This API may change in the future.
func (s *Scope) BucketName() string {
	return s.bucket.Name()
}

// Collection returns an instance of a collection.
// VOLATILE: This API is subject to change at any time.
func (s *Scope) Collection(collectionName string) *Collection {
	return newCollection(s, collectionName)
}
