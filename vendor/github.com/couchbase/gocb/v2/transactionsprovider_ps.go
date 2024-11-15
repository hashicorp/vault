package gocb

import "github.com/couchbase/gocbcore/v10"

type transactionsProviderPs struct{}

type transactionsInternalPs struct{}

func (t *transactionsProviderPs) Run(logicFn AttemptFunc, perConfig *TransactionOptions, singleQueryMode bool) (*TransactionResult, error) {
	return nil, wrapError(ErrFeatureNotAvailable, "transactions are not currently supported against the couchbase2 protocol")
}

func (t *transactionsProviderPs) Internal() transactionsInternal {
	return &transactionsInternalPs{}
}

func (t *transactionsInternalPs) ForceCleanupQueue() []TransactionCleanupAttempt {
	return nil
}

func (t *transactionsInternalPs) CleanupQueueLength() int32 {
	return 0
}

func (t *transactionsInternalPs) ClientCleanupEnabled() bool {
	return false
}

func (t *transactionsInternalPs) CleanupLocations() []gocbcore.TransactionLostATRLocation {
	return nil
}
