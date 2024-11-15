package gocb

import "github.com/couchbase/gocbcore/v10"

// TransactionErrorReason is the reason why a transaction should be failed.
// Internal: This should never be used and is not supported.
type TransactionErrorReason uint8

const (
	// TransactionErrorReasonSuccess indicates the transaction succeeded and did not fail.
	TransactionErrorReasonSuccess TransactionErrorReason = TransactionErrorReason(gocbcore.TransactionErrorReasonSuccess)

	// TransactionErrorReasonTransactionFailed indicates the transaction should be failed because it failed.
	TransactionErrorReasonTransactionFailed = TransactionErrorReason(gocbcore.TransactionErrorReasonTransactionFailed)

	// TransactionErrorReasonTransactionExpired indicates the transaction should be failed because it expired.
	TransactionErrorReasonTransactionExpired = TransactionErrorReason(gocbcore.TransactionErrorReasonTransactionExpired)

	// TransactionErrorReasonTransactionCommitAmbiguous indicates the transaction should be failed and the commit was ambiguous.
	TransactionErrorReasonTransactionCommitAmbiguous = TransactionErrorReason(gocbcore.TransactionErrorReasonTransactionCommitAmbiguous)

	// TransactionErrorReasonTransactionFailedPostCommit indicates the transaction should be failed because it failed post commit.
	TransactionErrorReasonTransactionFailedPostCommit = TransactionErrorReason(gocbcore.TransactionErrorReasonTransactionFailedPostCommit)
)
