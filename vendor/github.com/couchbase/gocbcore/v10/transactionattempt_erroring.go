// Copyright 2021 Couchbase
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocbcore

import (
	"errors"
	"sync/atomic"
)

func mergeOperationFailedErrors(errs []*TransactionOperationFailedError) *TransactionOperationFailedError {
	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	shouldNotRetry := false
	shouldNotRollback := false
	aggCauses := aggregateError{}
	shouldRaise := TransactionErrorReasonTransactionFailed

	for errIdx := 0; errIdx < len(errs); errIdx++ {
		tErr := errs[errIdx]

		aggCauses = append(aggCauses, tErr)

		if tErr.shouldNotRetry {
			shouldNotRetry = true
		}
		if tErr.shouldNotRollback {
			shouldNotRollback = true
		}
		if tErr.shouldRaise > shouldRaise {
			shouldRaise = tErr.shouldRaise
		}
	}

	return &TransactionOperationFailedError{
		shouldNotRetry:    shouldNotRetry,
		shouldNotRollback: shouldNotRollback,
		errorCause:        aggCauses,
		shouldRaise:       shouldRaise,
		errorClass:        TransactionErrorClassFailOther,
	}
}

type operationFailedDef struct {
	Cerr              *classifiedError
	ShouldNotRetry    bool
	ShouldNotRollback bool
	CanStillCommit    bool
	Reason            TransactionErrorReason
}

func (t *transactionAttempt) applyStateBits(stateBits uint32, errorBits uint32) {
	// This is a bit dirty, but its maximum going to do one retry per bit.
	for {
		oldStateBits := atomic.LoadUint32(&t.stateBits)
		newStateBits := oldStateBits | stateBits
		if errorBits > ((oldStateBits & transactionStateBitsMaskFinalError) >> transactionStateBitsPositionFinalError) {
			newStateBits = (newStateBits & transactionStateBitsMaskBits) | (errorBits << transactionStateBitsPositionFinalError)
		}
		t.logger.logInfof(t.id, "Applying state bits: %08b, error bits: %08b, old: %08b, new: %08b",
			stateBits, errorBits, oldStateBits, newStateBits)
		if atomic.CompareAndSwapUint32(&t.stateBits, oldStateBits, newStateBits) {
			break
		}
	}
}

func (t *transactionAttempt) operationFailed(def operationFailedDef) *TransactionOperationFailedError {
	t.logger.logInfof(t.id, "Operation failed: can still commit: %t, should not rollback: %t, should not retry: %t, "+
		"reason: %s", def.CanStillCommit, def.ShouldNotRollback, def.ShouldNotRetry, def.Reason)
	err := &TransactionOperationFailedError{
		shouldNotRetry:    def.ShouldNotRetry,
		shouldNotRollback: def.ShouldNotRollback,
		errorCause:        def.Cerr.Source,
		errorClass:        def.Cerr.Class,
		shouldRaise:       def.Reason,
	}

	stateBits := uint32(0)
	if !def.CanStillCommit {
		stateBits |= transactionStateBitShouldNotCommit
	}
	if def.ShouldNotRollback {
		stateBits |= transactionStateBitShouldNotRollback
	}
	if def.ShouldNotRetry {
		stateBits |= transactionStateBitShouldNotRetry
	}
	if def.Reason == TransactionErrorReasonTransactionExpired {
		stateBits |= transactionStateBitHasExpired
	}
	t.applyStateBits(stateBits, uint32(def.Reason))

	return err
}

func classifyHookError(err error) *classifiedError {
	// We currently have to classify the errors that are returned from the hooks, but
	// we should really just directly return the classifications and make the source
	// some special internal source showing it came from a hook...
	return classifyError(err)
}

func classifyError(err error) *classifiedError {
	ec := TransactionErrorClassFailOther
	if errors.Is(err, ErrDocAlreadyInTransaction) || errors.Is(err, ErrWriteWriteConflict) {
		ec = TransactionErrorClassFailWriteWriteConflict
	} else if errors.Is(err, ErrHard) {
		ec = TransactionErrorClassFailHard
	} else if errors.Is(err, ErrAttemptExpired) {
		ec = TransactionErrorClassFailExpiry
	} else if errors.Is(err, ErrTransient) {
		ec = TransactionErrorClassFailTransient
	} else if errors.Is(err, ErrDocumentNotFound) {
		ec = TransactionErrorClassFailDocNotFound
	} else if errors.Is(err, ErrAmbiguous) {
		ec = TransactionErrorClassFailAmbiguous
	} else if errors.Is(err, ErrCasMismatch) {
		ec = TransactionErrorClassFailCasMismatch
	} else if errors.Is(err, ErrDocumentNotFound) {
		ec = TransactionErrorClassFailDocNotFound
	} else if errors.Is(err, ErrDocumentExists) {
		ec = TransactionErrorClassFailDocAlreadyExists
	} else if errors.Is(err, ErrPathExists) {
		ec = TransactionErrorClassFailPathAlreadyExists
	} else if errors.Is(err, ErrPathNotFound) {
		ec = TransactionErrorClassFailPathNotFound
	} else if errors.Is(err, ErrCasMismatch) {
		ec = TransactionErrorClassFailCasMismatch
	} else if errors.Is(err, ErrUnambiguousTimeout) {
		ec = TransactionErrorClassFailTransient
	} else if errors.Is(err, ErrDurabilityAmbiguous) ||
		errors.Is(err, ErrAmbiguousTimeout) ||
		errors.Is(err, ErrRequestCanceled) {
		ec = TransactionErrorClassFailAmbiguous
	} else if errors.Is(err, ErrMemdTooBig) || errors.Is(err, ErrValueTooLarge) {
		ec = TransactionErrorClassFailOutOfSpace
	}

	return &classifiedError{
		Source: err,
		Class:  ec,
	}
}
