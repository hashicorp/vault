package gocb

import (
	"github.com/couchbase/gocbcore/v10"
)

// TransactionAttemptState represents the current state of a transaction attempt.
// Internal: This should never be used and is not supported.
type TransactionAttemptState int

const (
	// TransactionAttemptStateNothingWritten indicates that nothing has been written in this attempt.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStateNothingWritten = TransactionAttemptState(gocbcore.TransactionAttemptStateNothingWritten)

	// TransactionAttemptStatePending indicates that this attempt is in pending state.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStatePending = TransactionAttemptState(gocbcore.TransactionAttemptStatePending)

	// TransactionAttemptStateCommitting indicates that this attempt is in committing state.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStateCommitting = TransactionAttemptState(gocbcore.TransactionAttemptStateCommitting)

	// TransactionAttemptStateCommitted indicates that this attempt is in committed state.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStateCommitted = TransactionAttemptState(gocbcore.TransactionAttemptStateCommitted)

	// TransactionAttemptStateCompleted indicates that this attempt is in completed state.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStateCompleted = TransactionAttemptState(gocbcore.TransactionAttemptStateCompleted)

	// TransactionAttemptStateAborted indicates that this attempt is in aborted state.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStateAborted = TransactionAttemptState(gocbcore.TransactionAttemptStateAborted)

	// TransactionAttemptStateRolledBack indicates that this attempt is in rolled back state.
	// Internal: This should never be used and is not supported.
	TransactionAttemptStateRolledBack = TransactionAttemptState(gocbcore.TransactionAttemptStateRolledBack)
)

// TransactionResult represents the result of a transaction which was executed.
type TransactionResult struct {
	// TransactionID represents the UUID assigned to this transaction
	TransactionID string

	// UnstagingComplete indicates whether the transaction was succesfully
	// unstaged, or if a later cleanup job will be responsible.
	UnstagingComplete bool

	// Logs returns the set of logs that were created during this transaction.
	// UNCOMMITTED: This API may change in the future.
	Logs []TransactionLogItem
}
