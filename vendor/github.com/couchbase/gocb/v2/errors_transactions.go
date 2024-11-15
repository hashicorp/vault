package gocb

import (
	"errors"
	"github.com/couchbase/gocbcore/v10"
)

var (
	// ErrOther indicates an non-specific error has occured.
	ErrOther = gocbcore.ErrOther

	// ErrTransient indicates a transient error occured which may succeed at a later point in time.
	ErrTransient = gocbcore.ErrTransient

	// ErrWriteWriteConflict indicates that another transaction conflicted with this one.
	ErrWriteWriteConflict = gocbcore.ErrWriteWriteConflict

	// ErrHard indicates that an unrecoverable error occured.
	ErrHard = gocbcore.ErrHard

	// ErrAmbiguous indicates that a failure occured but the outcome was not known.
	ErrAmbiguous = gocbcore.ErrAmbiguous

	// ErrAtrFull indicates that the ATR record was too full to accept a new mutation.
	ErrAtrFull = gocbcore.ErrAtrFull

	// ErrAttemptExpired indicates an transactionAttempt expired
	ErrAttemptExpired = gocbcore.ErrAttemptExpired

	// ErrAtrNotFound indicates that an expected ATR document was missing
	ErrAtrNotFound = gocbcore.ErrAtrNotFound

	// ErrAtrEntryNotFound indicates that an expected ATR entry was missing
	ErrAtrEntryNotFound = gocbcore.ErrAtrEntryNotFound

	// ErrDocAlreadyInTransaction indicates that a document is already in a transaction.
	ErrDocAlreadyInTransaction = gocbcore.ErrDocAlreadyInTransaction

	// ErrTransactionAbortedExternally indicates the transaction was aborted externally.
	ErrTransactionAbortedExternally = gocbcore.ErrTransactionAbortedExternally

	// ErrPreviousOperationFailed indicates a previous operation already failed.
	ErrPreviousOperationFailed = gocbcore.ErrPreviousOperationFailed

	// ErrForwardCompatibilityFailure indicates an operation failed due to involving a document in another transaction
	// which contains features this transaction does not support.
	ErrForwardCompatibilityFailure = gocbcore.ErrForwardCompatibilityFailure

	// ErrIllegalState is used for when a transaction enters an illegal State.
	ErrIllegalState = gocbcore.ErrIllegalState

	ErrAttemptNotFoundOnQuery = errors.New("transactionAttempt not found on query")
)

type TransactionFailedError struct {
	cause  error
	result *TransactionResult
}

func (tfe TransactionFailedError) Error() string {
	if tfe.cause == nil {
		return "transaction failed"
	}
	return "transaction failed | " + tfe.cause.Error()
}

func (tfe TransactionFailedError) Unwrap() error {
	return tfe.cause
}

// Internal: This should never be used and is not supported.
func (tfe TransactionFailedError) Result() *TransactionResult {
	return tfe.result
}

type TransactionExpiredError struct {
	result *TransactionResult
}

func (tfe TransactionExpiredError) Error() string {
	return ErrAttemptExpired.Error()
}

func (tfe TransactionExpiredError) Unwrap() error {
	return ErrAttemptExpired
}

// Internal: This should never be used and is not supported.
func (tfe TransactionExpiredError) Result() *TransactionResult {
	return tfe.result
}

type TransactionCommitAmbiguousError struct {
	cause  error
	result *TransactionResult
}

func (tfe TransactionCommitAmbiguousError) Error() string {
	if tfe.cause == nil {
		return "transaction commit ambiguous"
	}
	return "transaction failed | " + tfe.cause.Error()
}

func (tfe TransactionCommitAmbiguousError) Unwrap() error {
	return tfe.cause
}

// Internal: This should never be used and is not supported.
func (tfe TransactionCommitAmbiguousError) Result() *TransactionResult {
	return tfe.result
}

type TransactionFailedPostCommit struct {
	cause  error
	result *TransactionResult
}

func (tfe TransactionFailedPostCommit) Error() string {
	if tfe.cause == nil {
		return "transaction failed post commit"
	}
	return "transaction failed | " + tfe.cause.Error()
}

func (tfe TransactionFailedPostCommit) Unwrap() error {
	return tfe.cause
}

// Internal: This should never be used and is not supported.
func (tfe TransactionFailedPostCommit) Result() *TransactionResult {
	return tfe.result
}

// TransactionOperationFailedError is used when a transaction operation fails.
// Internal: This should never be used and is not supported.
type TransactionOperationFailedError struct {
	shouldRetry       bool
	shouldNotRollback bool
	errorCause        error
	shouldRaise       gocbcore.TransactionErrorReason
	errorClass        gocbcore.TransactionErrorClass
}

func (tfe TransactionOperationFailedError) Error() string {
	if tfe.errorCause == nil {
		return "transaction operation failed"
	}
	return "transaction operation failed | " + tfe.errorCause.Error()
}

// InternalUnwrap returns the underlying error for this error.
func (tfe TransactionOperationFailedError) InternalUnwrap() error {
	return tfe.errorCause
}

// Retry signals whether a new transactionAttempt should be made at rollback.
func (tfe TransactionOperationFailedError) Retry() bool {
	return tfe.shouldRetry
}

// Rollback signals whether the transactionAttempt should be auto-rolled back.
func (tfe TransactionOperationFailedError) Rollback() bool {
	return !tfe.shouldNotRollback
}

// ToRaise signals which error type should be raised to the application.
func (tfe TransactionOperationFailedError) ToRaise() TransactionErrorReason {
	return TransactionErrorReason(tfe.shouldRaise)
}

func createTransactionOperationFailedError(err error) error {
	if err == nil {
		return nil
	}

	var txnErr *gocbcore.TransactionOperationFailedError
	if errors.As(err, &txnErr) {
		return &TransactionOperationFailedError{
			shouldRetry:       txnErr.Retry(),
			shouldNotRollback: !txnErr.Rollback(),
			errorCause:        txnErr.InternalUnwrap(),
			shouldRaise:       txnErr.ToRaise(),
			errorClass:        txnErr.ErrorClass(),
		}
	}

	return &TransactionOperationFailedError{
		errorCause: err,
		errorClass: gocbcore.TransactionErrorClassFailOther,
	}
}

func errorReasonFromString(reason string) gocbcore.TransactionErrorReason {
	switch reason {
	case "failed":
		return gocbcore.TransactionErrorReasonTransactionFailed
	case "expired":
		return gocbcore.TransactionErrorReasonTransactionExpired
	case "commit_ambiguous":
		return gocbcore.TransactionErrorReasonTransactionCommitAmbiguous
	case "failed_post_commit":
		return gocbcore.TransactionErrorReasonTransactionFailedPostCommit
	default:
		return gocbcore.TransactionErrorReasonTransactionFailed
	}
}
