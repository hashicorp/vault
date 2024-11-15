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

import "fmt"

// TransactionAttemptState represents the current State of a transaction
type TransactionAttemptState int

const (
	// TransactionAttemptStateNothingWritten indicates that nothing has been written yet.
	TransactionAttemptStateNothingWritten = TransactionAttemptState(1)

	// TransactionAttemptStatePending indicates that the transaction ATR has been written and
	// the transaction is currently pending.
	TransactionAttemptStatePending = TransactionAttemptState(2)

	// TransactionAttemptStateCommitting indicates that the transaction is now trying to become
	// committed, if we stay in this state, it implies ambiguity.
	TransactionAttemptStateCommitting = TransactionAttemptState(3)

	// TransactionAttemptStateCommitted indicates that the transaction is now logically committed
	// but the unstaging of documents is still underway.
	TransactionAttemptStateCommitted = TransactionAttemptState(4)

	// TransactionAttemptStateCompleted indicates that the transaction has been fully completed
	// and no longer has work to perform.
	TransactionAttemptStateCompleted = TransactionAttemptState(5)

	// TransactionAttemptStateAborted indicates that the transaction was aborted.
	TransactionAttemptStateAborted = TransactionAttemptState(6)

	// TransactionAttemptStateRolledBack indicates that the transaction was not committed and instead
	// was rolled back in its entirety.
	TransactionAttemptStateRolledBack = TransactionAttemptState(7)
)

func (state TransactionAttemptState) String() string {
	switch state {
	case TransactionAttemptStateNothingWritten:
		return "nothing_written"
	case TransactionAttemptStatePending:
		return "pending"
	case TransactionAttemptStateCommitting:
		return "committing"
	case TransactionAttemptStateCommitted:
		return "committed"
	case TransactionAttemptStateCompleted:
		return "completed"
	case TransactionAttemptStateAborted:
		return "aborted"
	case TransactionAttemptStateRolledBack:
		return "rolled_back"
	default:
		return "unknown"
	}
}

// TransactionErrorReason is the reason why a transaction should be failed.
// Internal: This should never be used and is not supported.
type TransactionErrorReason uint8

// NOTE: The errors within this section are critically ordered, as the order of
// precedence used when merging errors together is based on this.
const (
	// TransactionErrorReasonSuccess indicates the transaction succeeded and did not fail.
	TransactionErrorReasonSuccess TransactionErrorReason = iota

	// TransactionErrorReasonTransactionFailed indicates the transaction should be failed because it failed.
	TransactionErrorReasonTransactionFailed

	// TransactionErrorReasonTransactionExpired indicates the transaction should be failed because it expired.
	TransactionErrorReasonTransactionExpired

	// TransactionErrorReasonTransactionCommitAmbiguous indicates the transaction should be failed and the commit was ambiguous.
	TransactionErrorReasonTransactionCommitAmbiguous

	// TransactionErrorReasonTransactionFailedPostCommit indicates the transaction should be failed because it failed post commit.
	TransactionErrorReasonTransactionFailedPostCommit
)

func (reason TransactionErrorReason) String() string {
	switch reason {
	case TransactionErrorReasonTransactionFailed:
		return "failed"
	case TransactionErrorReasonTransactionExpired:
		return "expired"
	case TransactionErrorReasonTransactionCommitAmbiguous:
		return "commit_ambiguous"
	case TransactionErrorReasonTransactionFailedPostCommit:
		return "failed_post_commit"
	default:
		return fmt.Sprintf("unknown:%d", reason)
	}
}

// TransactionErrorClass describes the reason that a transaction error occurred.
// Internal: This should never be used and is not supported.
type TransactionErrorClass uint8

const (
	// TransactionErrorClassFailOther indicates an error occurred because it did not fit into any other reason.
	TransactionErrorClassFailOther TransactionErrorClass = iota

	// TransactionErrorClassFailTransient indicates an error occurred because of a transient reason.
	TransactionErrorClassFailTransient

	// TransactionErrorClassFailDocNotFound indicates an error occurred because of a document not found.
	TransactionErrorClassFailDocNotFound

	// TransactionErrorClassFailDocAlreadyExists indicates an error occurred because a document already exists.
	TransactionErrorClassFailDocAlreadyExists

	// TransactionErrorClassFailPathNotFound indicates an error occurred because a path was not found.
	TransactionErrorClassFailPathNotFound

	// TransactionErrorClassFailPathAlreadyExists indicates an error occurred because a path already exists.
	TransactionErrorClassFailPathAlreadyExists

	// TransactionErrorClassFailWriteWriteConflict indicates an error occurred because of a write write conflict.
	TransactionErrorClassFailWriteWriteConflict

	// TransactionErrorClassFailCasMismatch indicates an error occurred because of a cas mismatch.
	TransactionErrorClassFailCasMismatch

	// TransactionErrorClassFailHard indicates an error occurred because of a hard error.
	TransactionErrorClassFailHard

	// TransactionErrorClassFailAmbiguous indicates an error occurred leaving the transaction in an ambiguous way.
	TransactionErrorClassFailAmbiguous

	// TransactionErrorClassFailExpiry indicates an error occurred because the transaction expired.
	TransactionErrorClassFailExpiry

	// TransactionErrorClassFailOutOfSpace indicates an error occurred because the ATR is full.
	TransactionErrorClassFailOutOfSpace
)

const (
	transactionStateBitShouldNotCommit       = 1 << 0
	transactionStateBitShouldNotRollback     = 1 << 1
	transactionStateBitShouldNotRetry        = 1 << 2
	transactionStateBitHasExpired            = 1 << 3
	transactionStateBitPreExpiryAutoRollback = 1 << 4
)

const (
	transactionStateBitsMaskFinalError     = 0b1110000
	transactionStateBitsMaskBits           = 0b0001111
	transactionStateBitsPositionFinalError = 4
)
