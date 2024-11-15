package gocb

import (
	"errors"
	"github.com/couchbase/gocbcore/v10"
)

// AttemptFunc represents the lambda used by the Transactions Run function.
type AttemptFunc func(*TransactionAttemptContext) error

// Transactions can be used to perform transactions.
type Transactions struct {
	controller *providerController[transactionsProvider]
}

// initTransactions will initialize the transactions library and return a Transactions
// object which can be used to perform transactions.
func (c *Cluster) initTransactions(config TransactionsConfig) (*Transactions, error) {
	err := c.connectionManager.initTransactions(config, c)
	if err != nil {
		return nil, err
	}

	return &Transactions{
		controller: c.transactionsController(),
	}, nil
}

// Run runs a lambda to perform a number of operations as part of a
// singular transaction.
func (t *Transactions) Run(logicFn AttemptFunc, perConfig *TransactionOptions) (*TransactionResult, error) {
	return autoOpControl(t.controller, func(provider transactionsProvider) (*TransactionResult, error) {
		return provider.Run(logicFn, perConfig, false)
	})
}

func (t *Transactions) singleQuery(statement string, scope *Scope, opts QueryOptions) (*QueryResult, error) {
	return autoOpControl(t.controller, func(provider transactionsProvider) (*QueryResult, error) {
		if opts.Context != nil {
			return nil, makeInvalidArgumentsError("cannot use context and transactions together")
		}

		config := &TransactionOptions{
			DurabilityLevel: opts.AsTransaction.DurabilityLevel,
			Timeout:         opts.Timeout,
		}
		config.Internal.Hooks = opts.AsTransaction.Internal.Hooks

		var queryRes *QueryResult
		res, err := provider.Run(func(context *TransactionAttemptContext) error {
			// We need to tell the core loop that autocommit and autorollback are disabled.
			// context.txn.UpdateState(gocbcore.TransactionUpdateStateOptions{
			// 	ShouldNotCommit:   true,
			// 	ShouldNotRollback: true,
			// })
			qRes, err := context.queryWrapper(scope, statement, opts, "query", false, false,
				nil, true)
			if err != nil {
				return err
			}

			queryRes = qRes
			// If the result contains rows then we can't immediately check for errors, so we need to return here.
			if len(queryRes.peekNext()) > 0 {
				// We consider this success so tell the core to not retry - any errors on stream will happen outside the
				// context of the core loop.
				// context.txn.UpdateState(gocbcore.TransactionUpdateStateOptions{
				// 	ShouldNotRetry: true,
				// })

				return nil
			}

			if err := qRes.Err(); err != nil {
				return queryMaybeTranslateToTransactionsError(err, context)
			}

			meta, err := qRes.MetaData()
			if err != nil {
				return queryMaybeTranslateToTransactionsError(err, context)
			}

			if meta.Status == QueryStatusFatal {
				return operationFailed(transactionQueryOperationFailedDef{
					ShouldNotRetry:  true,
					Reason:          gocbcore.TransactionErrorReasonTransactionFailed,
					ShouldNotCommit: true,
				}, context)
			}

			// We won't do autocommit or autorollback so tell the core loop to not retry.
			// context.txn.UpdateState(gocbcore.TransactionUpdateStateOptions{
			// 	ShouldNotRetry: true,
			// })

			return nil
		}, config, true)
		if err != nil {
			var expiredErr *TransactionExpiredError
			if errors.As(err, &expiredErr) {
				return nil, ErrUnambiguousTimeout
			}
			return nil, err
		}

		queryRes.transactionID = res.TransactionID

		return queryRes, nil
	})
}

// TransactionsInternal exposes internal methods that are useful for testing and/or
// other forms of internal use.
type TransactionsInternal struct {
	parent *Transactions
}

// Internal returns an TransactionsInternal object which can be used for specialized
// internal use cases.
func (t *Transactions) Internal() *TransactionsInternal {
	return &TransactionsInternal{
		parent: t,
	}
}

// ForceCleanupQueue forces the transactions client cleanup queue to drain without waiting for expirations.
func (t *TransactionsInternal) ForceCleanupQueue() []TransactionCleanupAttempt {
	attempts, err := autoOpControl(t.parent.controller, func(provider transactionsProvider) ([]TransactionCleanupAttempt, error) {
		return provider.Internal().ForceCleanupQueue(), nil
	})
	if err != nil {
		return nil
	}

	return attempts
}

// CleanupQueueLength returns the current length of the client cleanup queue.
func (t *TransactionsInternal) CleanupQueueLength() int32 {
	length, err := autoOpControl(t.parent.controller, func(provider transactionsProvider) (int32, error) {
		return provider.Internal().CleanupQueueLength(), nil
	})
	if err != nil {
		return 0
	}

	return length
}

// ClientCleanupEnabled returns whether the client cleanup process is enabled.
func (t *TransactionsInternal) ClientCleanupEnabled() bool {
	enabled, err := autoOpControl(t.parent.controller, func(provider transactionsProvider) (bool, error) {
		return provider.Internal().ClientCleanupEnabled(), nil
	})
	if err != nil {
		return false
	}

	return enabled
}

// CleanupLocations returns the set of locations currently being watched by the lost transactions process.
func (t *TransactionsInternal) CleanupLocations() []gocbcore.TransactionLostATRLocation {
	locs, err := autoOpControl(t.parent.controller, func(provider transactionsProvider) ([]gocbcore.TransactionLostATRLocation, error) {
		return provider.Internal().CleanupLocations(), nil
	})
	if err != nil {
		return nil
	}

	return locs
}
