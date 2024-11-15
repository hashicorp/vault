package gocb

import (
	"encoding/json"
	"errors"
	"github.com/couchbase/gocbcore/v10"
)

func queryErrorCodeToError(code uint32, c *TransactionAttemptContext) error {
	switch code {
	case 1065:
		return operationFailed(transactionQueryOperationFailedDef{
			ShouldNotRetry: true,
			Reason:         gocbcore.TransactionErrorReasonTransactionFailed,
			ErrorCause:     ErrFeatureNotAvailable,
		}, c)
	case 1080:
		return operationFailed(transactionQueryOperationFailedDef{
			ShouldNotRetry:    true,
			Reason:            gocbcore.TransactionErrorReasonTransactionExpired,
			ErrorCause:        gocbcore.ErrAttemptExpired,
			ErrorClass:        gocbcore.TransactionErrorClassFailExpiry,
			ShouldNotRollback: true,
		}, c)
	case 17004:
		return ErrAttemptNotFoundOnQuery
	case 17010:
		return operationFailed(transactionQueryOperationFailedDef{
			ShouldNotRetry:    true,
			Reason:            gocbcore.TransactionErrorReasonTransactionExpired,
			ErrorCause:        gocbcore.ErrAttemptExpired,
			ErrorClass:        gocbcore.TransactionErrorClassFailExpiry,
			ShouldNotRollback: true,
		}, c)
	case 17012:
		return ErrDocumentExists
	case 17014:
		return ErrDocumentNotFound
	case 17015:
		return ErrCasMismatch
	default:
		return nil
	}
}

func queryCauseToOperationFailedError(queryErr *QueryError, c *TransactionAttemptContext) error {

	var operationFailedErrs []jsonQueryTransactionOperationFailedCause
	if err := json.Unmarshal([]byte(queryErr.ErrorText), &operationFailedErrs); err == nil {
		for _, operationFailedErr := range operationFailedErrs {
			if operationFailedErr.Cause != nil {
				if operationFailedErr.Code >= 17000 && operationFailedErr.Code <= 18000 {
					if err := queryErrorCodeToError(operationFailedErr.Code, c); err != nil {
						return err
					}
				}

				return operationFailed(transactionQueryOperationFailedDef{
					ShouldNotRetry:    !operationFailedErr.Cause.Retry,
					ShouldNotRollback: !operationFailedErr.Cause.Rollback,
					Reason:            errorReasonFromString(operationFailedErr.Cause.Raise),
					ErrorCause:        queryErr,
					ShouldNotCommit:   true,
				}, c)
			}
		}
	}
	return nil
}

func queryMaybeTranslateToTransactionsError(err error, c *TransactionAttemptContext) error {
	if errors.Is(err, ErrTimeout) {
		return operationFailed(transactionQueryOperationFailedDef{
			ShouldNotRetry: true,
			Reason:         gocbcore.TransactionErrorReasonTransactionExpired,
			ErrorCause:     err,
		}, c)
	}

	var queryErr *QueryError
	if !errors.As(err, &queryErr) {
		return err
	}

	if len(queryErr.Errors) == 0 {
		return queryErr
	}

	// If an error contains a cause field, use that error.
	// Otherwise, if an error has code between 17000 and 18000 inclusive, it is a transactions-related error. Use that.
	// Otherwise, fallback to using the first error.
	if err := queryCauseToOperationFailedError(queryErr, c); err != nil {
		return err
	}

	for _, e := range queryErr.Errors {
		if e.Code >= 17000 && e.Code <= 18000 {
			if err := queryErrorCodeToError(e.Code, c); err != nil {
				return err
			}
		}
	}

	if err := queryErrorCodeToError(queryErr.Errors[0].Code, c); err != nil {
		return err
	}

	return queryErr
}

type transactionQueryOperationFailedDef struct {
	ShouldNotRetry    bool
	ShouldNotRollback bool
	Reason            gocbcore.TransactionErrorReason
	ErrorCause        error
	ErrorClass        gocbcore.TransactionErrorClass
	ShouldNotCommit   bool
}

func operationFailed(def transactionQueryOperationFailedDef, c *TransactionAttemptContext) *TransactionOperationFailedError {
	err := &TransactionOperationFailedError{
		shouldRetry:       !def.ShouldNotRetry,
		shouldNotRollback: def.ShouldNotRollback,
		errorCause:        def.ErrorCause,
		shouldRaise:       def.Reason,
		errorClass:        def.ErrorClass,
	}

	if c != nil {
		c.logger.logInfof(c.attemptID, "Operation failed: can still commit: %t, should not rollback: %t, should not retry: %t, "+
			"reason: %s", !def.ShouldNotCommit, def.ShouldNotRollback, def.ShouldNotRetry, def.Reason)
		c.updateState(def)
	}
	return err
}

func (c *TransactionAttemptContext) updateState(def transactionQueryOperationFailedDef) {
	opts := gocbcore.TransactionUpdateStateOptions{}
	if def.ShouldNotRollback {
		opts.ShouldNotRollback = true
	}
	if def.ShouldNotRetry {
		opts.ShouldNotRetry = true
	}
	if def.ShouldNotCommit {
		opts.ShouldNotCommit = true
	}
	opts.Reason = def.Reason
	c.txn.UpdateState(opts)
}

func singleQueryErrToTransactionError(err error, txnID string) error {
	err = queryMaybeTranslateToTransactionsError(err, nil)

	var tErr *TransactionOperationFailedError
	if errors.As(err, &tErr) {
		switch tErr.shouldRaise {
		case gocbcore.TransactionErrorReasonTransactionFailed:
			return &TransactionFailedError{
				cause: tErr.errorCause,
				result: &TransactionResult{
					TransactionID:     txnID,
					UnstagingComplete: false,
				},
			}
		case gocbcore.TransactionErrorReasonTransactionCommitAmbiguous:
			return &TransactionCommitAmbiguousError{
				cause: tErr.errorCause,
				result: &TransactionResult{
					TransactionID:     txnID,
					UnstagingComplete: false,
				},
			}
		case gocbcore.TransactionErrorReasonTransactionExpired:
			return &TransactionExpiredError{
				result: &TransactionResult{
					TransactionID:     txnID,
					UnstagingComplete: false,
				},
			}
		}
	}

	return &TransactionFailedError{
		cause: err,
		result: &TransactionResult{
			TransactionID:     txnID,
			UnstagingComplete: false,
		},
	}
}
