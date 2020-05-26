package gocb

// Scope represents a single scope within a bucket.
// VOLATILE: This API is subject to change at any time.
type Scope struct {
	scopeName  string
	bucketName string

	timeoutsConfig kvTimeoutsConfig

	transcoder           Transcoder
	retryStrategyWrapper *retryStrategyWrapper
	tracer               requestTracer

	useMutationTokens bool

	getKvProvider func() (kvProvider, error)
}

func newScope(bucket *Bucket, scopeName string) *Scope {
	return &Scope{
		scopeName:  scopeName,
		bucketName: bucket.Name(),

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

func (s *Scope) clone() *Scope {
	newS := *s
	return &newS
}

// Name returns the name of the scope.
func (s *Scope) Name() string {
	return s.scopeName
}

// Collection returns an instance of a collection.
// VOLATILE: This API is subject to change at any time.
func (s *Scope) Collection(collectionName string) *Collection {
	return newCollection(s, collectionName)
}
