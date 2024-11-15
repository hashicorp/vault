package gocb

import "github.com/couchbase/gocbcore/v10"

type transactionsProvider interface {
	Run(logicFn AttemptFunc, perConfig *TransactionOptions, singleQueryMode bool) (*TransactionResult, error)

	Internal() transactionsInternal
}

type transactionsInternal interface {
	ForceCleanupQueue() []TransactionCleanupAttempt
	CleanupQueueLength() int32
	ClientCleanupEnabled() bool
	CleanupLocations() []gocbcore.TransactionLostATRLocation
}
