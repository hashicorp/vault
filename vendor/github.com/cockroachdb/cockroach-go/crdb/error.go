package crdb

import "fmt"

type txError struct {
	cause error
}

// Error implements the error interface
func (e *txError) Error() string { return e.cause.Error() }

// Cause returns the error encountered by the "ROLLBACK TO SAVEPOINT
// cockroach_restart" statement. This method also implements the internal
// pkg/errors.causer interface, so TxnRestartError works with pkg/errors.Cause().
func (e *txError) Cause() error { return e.cause }

// AmbiguousCommitError represents an error that left a transaction in an
// ambiguous state: unclear if it committed or not.
type AmbiguousCommitError struct {
	txError
}

func newAmbiguousCommitError(err error) *AmbiguousCommitError {
	return &AmbiguousCommitError{txError{cause: err}}
}

// TxnRestartError represents an error when restarting a transaction. `cause` is
// the error from restarting the txn and `retryCause` is the original error which
// triggered the restart.
type TxnRestartError struct {
	txError
	retryCause error
	msg        string
}

func newTxnRestartError(err error, retryErr error) *TxnRestartError {
	const msgPattern = "restarting txn failed. ROLLBACK TO SAVEPOINT " +
		"encountered error: %s. Original error: %s."
	return &TxnRestartError{
		txError:    txError{cause: err},
		retryCause: retryErr,
		msg:        fmt.Sprintf(msgPattern, err, retryErr),
	}
}

// Error implements the error interface
func (e *TxnRestartError) Error() string { return e.msg }

// RetryCause returns the error encountered by the transaction, which caused the
// transaction to be restarted.
func (e *TxnRestartError) RetryCause() error { return e.retryCause }
