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
	"bytes"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

func transactionHasExpired(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}

func (t *transactionAttempt) beginOpAndLock(cb func(unlock func(), endOp func())) {
	t.lock.Lock(func(unlock func()) {
		t.opsWg.Add(1)

		cb(unlock, func() {
			t.opsWg.Done()
		})
	})
}

func (t *transactionAttempt) waitForOpsAndLock(cb func(unlock func())) {
	var tryWaitAndLock func()
	tryWaitAndLock = func() {
		t.opsWg.Wait(func() {
			t.lock.Lock(func(unlock func()) {
				if !t.opsWg.IsEmpty() {
					unlock()
					tryWaitAndLock()
					return
				}

				cb(unlock)
			})
		})
	}
	tryWaitAndLock()
}

func (t *transactionAttempt) checkCanPerformOpLocked() *TransactionOperationFailedError {
	switch t.state {
	case TransactionAttemptStateNothingWritten:
		fallthrough
	case TransactionAttemptStatePending:
		// Good to continue
	case TransactionAttemptStateCommitting:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, "transaction is ambiguously committed")),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	case TransactionAttemptStateCommitted:
		fallthrough
	case TransactionAttemptStateCompleted:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, "transaction already committed")),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	case TransactionAttemptStateAborted:
		fallthrough
	case TransactionAttemptStateRolledBack:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, "transaction already aborted")),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	default:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, fmt.Sprintf("invalid transaction state: %v", t.state))),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	}

	stateBits := atomic.LoadUint32(&t.stateBits)
	if (stateBits & transactionStateBitShouldNotCommit) != 0 {
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrPreviousOperationFailed, "previous operation prevents further operations")),
			ShouldNotRetry:    true,
			ShouldNotRollback: false,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	}

	return nil
}

func (t *transactionAttempt) checkCanCommitRollbackLocked() *TransactionOperationFailedError {
	switch t.state {
	case TransactionAttemptStateNothingWritten:
		fallthrough
	case TransactionAttemptStatePending:
		// Good to continue
	case TransactionAttemptStateCommitting:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, "transaction is ambiguously committed")),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	case TransactionAttemptStateCommitted:
		fallthrough
	case TransactionAttemptStateCompleted:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, "transaction already committed")),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	case TransactionAttemptStateAborted:
		fallthrough
	case TransactionAttemptStateRolledBack:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, "transaction already aborted")),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	default:
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrIllegalState, fmt.Sprintf("invalid transaction state: %v", t.state))),
			ShouldNotRetry:    true,
			ShouldNotRollback: true,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	}

	return nil
}

func (t *transactionAttempt) checkCanCommitLocked() *TransactionOperationFailedError {
	err := t.checkCanCommitRollbackLocked()
	if err != nil {
		return err
	}

	stateBits := atomic.LoadUint32(&t.stateBits)
	if (stateBits & transactionStateBitShouldNotCommit) != 0 {
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrPreviousOperationFailed, "previous operation prevents commit")),
			ShouldNotRetry:    true,
			ShouldNotRollback: false,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	}

	return nil
}

func (t *transactionAttempt) checkCanRollbackLocked() *TransactionOperationFailedError {
	err := t.checkCanCommitRollbackLocked()
	if err != nil {
		return err
	}

	stateBits := atomic.LoadUint32(&t.stateBits)
	if (stateBits & transactionStateBitShouldNotRollback) != 0 {
		return t.operationFailed(operationFailedDef{
			Cerr: classifyError(
				wrapError(ErrPreviousOperationFailed, "previous operation prevents rollback")),
			ShouldNotRetry:    true,
			ShouldNotRollback: false,
			Reason:            TransactionErrorReasonTransactionFailed,
		})
	}

	return nil
}

func (t *transactionAttempt) setExpiryOvertimeAtomic() {
	t.logger.logInfof(t.id, "Entering expiry overtime")
	t.applyStateBits(transactionStateBitHasExpired, 0)
}

func (t *transactionAttempt) isExpiryOvertimeAtomic() bool {
	stateBits := atomic.LoadUint32(&t.stateBits)
	return (stateBits & transactionStateBitHasExpired) != 0
}

func (t *transactionAttempt) checkExpiredAtomic(stage string, id []byte, proceedInOvertime bool, cb func(*classifiedError)) {
	if proceedInOvertime && t.isExpiryOvertimeAtomic() {
		cb(nil)
		return
	}

	t.hooks.HasExpiredClientSideHook(stage, id, func(expired bool, err error) {
		if err != nil {
			cb(classifyError(wrapError(err, "HasExpired hook returned an unexpected error")))
			return
		}

		if expired {
			cb(classifyError(wrapError(ErrAttemptExpired, "a hook has marked this attempt expired")))
			return
		} else if transactionHasExpired(t.expiryTime) {
			cb(classifyError(wrapError(ErrAttemptExpired, "the expiry for the attempt was reached")))
			return
		}

		cb(nil)
	})
}

func (t *transactionAttempt) confirmATRPending(
	firstAgent *Agent,
	firstOboUser string,
	firstScopeName string,
	firstCollectionName string,
	firstKey []byte,
	cb func(*TransactionOperationFailedError),
) {
	t.lock.Lock(func(unlock func()) {
		unlockAndCb := func(err *TransactionOperationFailedError) {
			unlock()
			cb(err)
		}

		if t.state != TransactionAttemptStateNothingWritten {
			unlockAndCb(nil)
			return
		}

		t.selectAtrLocked(
			firstAgent,
			firstOboUser,
			firstScopeName,
			firstCollectionName,
			firstKey,
			func(err *TransactionOperationFailedError) {
				if err != nil {
					unlockAndCb(err)
					return
				}

				t.setATRPendingLocked(func(err *TransactionOperationFailedError) {
					if err != nil {
						unlockAndCb(err)
						return
					}

					t.state = TransactionAttemptStatePending

					unlockAndCb(nil)
				})
			})
	})
}

func (t *transactionAttempt) getStagedMutationLocked(
	bucketName, scopeName, collectionName string, key []byte,
) (int, *transactionStagedMutation) {
	for i, mutation := range t.stagedMutations {
		if mutation.Agent.BucketName() == bucketName &&
			mutation.ScopeName == scopeName &&
			mutation.CollectionName == collectionName &&
			bytes.Equal(mutation.Key, key) {
			return i, mutation
		}
	}

	return -1, nil
}

func (t *transactionAttempt) removeStagedMutation(
	bucketName, scopeName, collectionName string, key []byte,
	cb func(),
) {
	t.lock.Lock(func(unlock func()) {
		mutIdx, _ := t.getStagedMutationLocked(bucketName, scopeName, collectionName, key)
		if mutIdx >= 0 {
			// Not finding the item should be basically impossible, but we wrap it just in case...
			t.stagedMutations = append(t.stagedMutations[:mutIdx], t.stagedMutations[mutIdx+1:]...)
		}

		unlock()
		cb()
	})
}

func (t *transactionAttempt) recordStagedMutation(
	stagedInfo *transactionStagedMutation,
	cb func(),
) {
	if !t.enableMutationCaching {
		stagedInfo.Staged = nil
	}

	t.lock.Lock(func(unlock func()) {
		mutIdx, _ := t.getStagedMutationLocked(
			stagedInfo.Agent.BucketName(),
			stagedInfo.ScopeName,
			stagedInfo.CollectionName,
			stagedInfo.Key)
		if mutIdx >= 0 {
			t.stagedMutations[mutIdx] = stagedInfo
		} else {
			t.stagedMutations = append(t.stagedMutations, stagedInfo)
		}

		unlock()
		cb()
	})
}

func (t *transactionAttempt) checkForwardCompatability(
	key []byte,
	bucket, scope, collection string,
	stage forwardCompatStage,
	fc map[string][]TransactionForwardCompatibilityEntry,
	forceNonFatal bool,
	cb func(*TransactionOperationFailedError),
) {
	t.logger.logInfof(t.id, "Checking forward compatibility")
	isCompat, shouldRetry, retryWait, err := checkForwardCompatability(stage, fc)
	if err != nil {
		t.logger.logInfof(t.id, "Forward compatability error")
		cb(t.operationFailed(operationFailedDef{
			Cerr:              classifyError(err),
			CanStillCommit:    forceNonFatal,
			ShouldNotRetry:    false,
			ShouldNotRollback: false,
			Reason:            TransactionErrorReasonTransactionFailed,
		}))
		return
	}

	if !isCompat {
		if shouldRetry {
			cbRetryError := func() {
				t.logger.logInfof(t.id, "Forward compatability failed - incompatible, should retry")
				cb(t.operationFailed(operationFailedDef{
					Cerr: classifyError(forwardCompatError{
						BucketName:     bucket,
						ScopeName:      scope,
						CollectionName: collection,
						DocumentKey:    key,
					}),
					CanStillCommit:    forceNonFatal,
					ShouldNotRetry:    false,
					ShouldNotRollback: false,
					Reason:            TransactionErrorReasonTransactionFailed,
				}))
			}

			if retryWait > 0 {
				time.AfterFunc(retryWait, cbRetryError)
			} else {
				cbRetryError()
			}

			return
		}

		t.logger.logInfof(t.id, "Forward compatability failed - incompatible")
		cb(t.operationFailed(operationFailedDef{
			Cerr: classifyError(forwardCompatError{
				BucketName:     bucket,
				ScopeName:      scope,
				CollectionName: collection,
				DocumentKey:    key,
			}),
			CanStillCommit:    forceNonFatal,
			ShouldNotRetry:    true,
			ShouldNotRollback: false,
			Reason:            TransactionErrorReasonTransactionFailed,
		}))
		return
	}

	cb(nil)
}

func (t *transactionAttempt) getTxnState(
	srcBucketName string,
	srcScopeName string,
	srcCollectionName string,
	srcDocID []byte,
	atrBucketName string,
	atrScopeName string,
	atrCollectionName string,
	atrDocID string,
	attemptID string,
	forceNonFatal bool,
	cb func(*jsonAtrAttempt, time.Time, *TransactionOperationFailedError),
) {
	ecCb := func(res *jsonAtrAttempt, txnExp time.Time, cerr *classifiedError) {
		if cerr == nil {
			cb(res, txnExp, nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		switch cerr.Class {
		case TransactionErrorClassFailPathNotFound:
			t.logger.logInfof(t.id, "Attempt entry not found")
			// If the path is not found, we just return as if there was no
			// entry data available for that atr entry.
			cb(nil, time.Time{}, nil)
		case TransactionErrorClassFailDocNotFound:
			t.logger.logInfof(t.id, "ATR doc not found")
			// If the ATR is not found, we just return as if there was no
			// entry data available for that atr entry.
			cb(nil, time.Time{}, nil)
		default:
			cb(nil, time.Time{}, t.operationFailed(operationFailedDef{
				Cerr: classifyError(&writeWriteConflictError{
					Source:         cerr.Source,
					BucketName:     srcBucketName,
					ScopeName:      srcScopeName,
					CollectionName: srcCollectionName,
					DocumentKey:    srcDocID,
				}),
				CanStillCommit:    forceNonFatal,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	}

	t.logger.logInfof(t.id, "Getting txn state")

	atrAgent, atrOboUser, err := t.bucketAgentProvider(atrBucketName)
	if err != nil {
		t.logger.logInfof(t.id, "Failed to get atr agent")
		ecCb(nil, time.Time{}, classifyError(err))
		return
	}

	t.hooks.BeforeCheckATREntryForBlockingDoc([]byte(atrDocID), func(err error) {
		if err != nil {
			ecCb(nil, time.Time{}, classifyHookError(err))
			return
		}

		var deadline time.Time
		if t.keyValueTimeout > 0 {
			deadline = time.Now().Add(t.keyValueTimeout)
		}

		_, err = atrAgent.LookupIn(LookupInOptions{
			ScopeName:      atrScopeName,
			CollectionName: atrCollectionName,
			Key:            []byte(atrDocID),
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpGet,
					Path:  "attempts." + attemptID,
					Flags: memd.SubdocFlagXattrPath,
				},
				{
					Op:    memd.SubDocOpGet,
					Path:  hlcMacro,
					Flags: memd.SubdocFlagXattrPath,
				},
			},
			Deadline: deadline,
			User:     atrOboUser,
		}, func(result *LookupInResult, err error) {
			if err != nil {
				ecCb(nil, time.Time{}, classifyError(err))
				return
			}

			t.ReportResourceUnits(result.Internal.ResourceUnits)

			for _, op := range result.Ops {
				if op.Err != nil {
					ecCb(nil, time.Time{}, classifyError(op.Err))
					return
				}
			}

			var txnAttempt *jsonAtrAttempt
			if err := json.Unmarshal(result.Ops[0].Value, &txnAttempt); err != nil {
				ecCb(nil, time.Time{}, classifyError(err))
				return
			}

			var hlc *jsonHLC
			if err := json.Unmarshal(result.Ops[1].Value, &hlc); err != nil {
				ecCb(nil, time.Time{}, classifyError(err))
				return
			}

			nowSecs, err := parseHLCToSeconds(*hlc)
			if err != nil {
				ecCb(nil, time.Time{}, classifyError(err))
				return
			}

			txnStartMs, err := parseCASToMilliseconds(txnAttempt.PendingCAS)
			if err != nil {
				ecCb(nil, time.Time{}, classifyError(err))
				return
			}

			nowTime := time.Duration(nowSecs) * time.Second
			txnStartTime := time.Duration(txnStartMs) * time.Millisecond
			txnExpiryTime := time.Duration(txnAttempt.ExpiryTime) * time.Millisecond

			txnElapsedTime := nowTime - txnStartTime
			txnExpiry := time.Now().Add(txnExpiryTime - txnElapsedTime)

			ecCb(txnAttempt, txnExpiry, nil)
		})
		if err != nil {
			ecCb(nil, time.Time{}, classifyError(err))
			return
		}
	})
}

func (t *transactionAttempt) writeWriteConflictPoll(
	stage forwardCompatStage,
	agent *Agent,
	oboUser string,
	scopeName string,
	collectionName string,
	key []byte,
	cas Cas,
	meta *TransactionMutableItemMeta,
	existingMutation *transactionStagedMutation,
	cb func(*TransactionOperationFailedError),
) {
	if meta == nil {
		t.logger.logInfof(t.id, "Meta is nil, no write-write conflict")
		// There is no write-write conflict.
		cb(nil)
		return
	}

	if meta.TransactionID == t.transactionID {
		if meta.AttemptID == t.id {
			if existingMutation != nil {
				if cas != existingMutation.Cas {
					// There was an existing mutation but it doesn't match the expected
					// CAS.  We throw a CAS mismatch to early detect this.
					cb(t.operationFailed(operationFailedDef{
						Cerr: classifyError(
							wrapError(ErrCasMismatch, "cas mismatch occured against local staged mutation")),
						ShouldNotRetry:    false,
						ShouldNotRollback: false,
						Reason:            TransactionErrorReasonTransactionFailed,
					}))
					return
				}

				cb(nil)
				return
			}

			// This means that we are trying to overwrite a previous write this specific
			// attempt has performed without actually having found the existing mutation,
			// this is never going to work correctly.
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrIllegalState, "attempted to overwrite local staged mutation but couldn't find it")),
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
			return
		}

		t.logger.logInfof(t.id, "Transaction meta matches ours, no write-write conflict")
		// The transaction matches our transaction.  We can safely overwrite the existing
		// data in the txn meta and continue.
		cb(nil)
		return
	}

	deadline := time.Now().Add(1 * time.Second)

	var onePoll func()
	onePoll = func() {
		t.logger.logInfof(t.id, "Performing write-write conflict poll")
		if !time.Now().Before(deadline) {
			t.logger.logInfof(t.id, "Deadline expired during write-write poll")
			// If the deadline expired, lets just immediately return.
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(&writeWriteConflictError{
					Source: fmt.Errorf(
						"deadline expired before WWC was resolved on %s.%s.%s.%s",
						meta.ATR.BucketName,
						meta.ATR.ScopeName,
						meta.ATR.CollectionName,
						meta.ATR.DocID),
					BucketName:     agent.BucketName(),
					ScopeName:      scopeName,
					CollectionName: collectionName,
					DocumentKey:    key,
				}),
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
			return
		}

		t.checkForwardCompatability(key, agent.BucketName(), scopeName, collectionName, stage, meta.ForwardCompat, false,
			func(err *TransactionOperationFailedError) {
				if err != nil {
					cb(err)
					return
				}

				t.checkExpiredAtomic(hookWWC, key, false, func(cerr *classifiedError) {
					if cerr != nil {
						cb(t.operationFailed(operationFailedDef{
							Cerr:              cerr,
							ShouldNotRetry:    true,
							ShouldNotRollback: false,
							Reason:            TransactionErrorReasonTransactionExpired,
						}))
						return
					}

					t.getTxnState(
						agent.BucketName(),
						scopeName,
						collectionName,
						key,
						meta.ATR.BucketName,
						meta.ATR.ScopeName,
						meta.ATR.CollectionName,
						meta.ATR.DocID,
						meta.AttemptID,
						false,
						func(attempt *jsonAtrAttempt, expiry time.Time, err *TransactionOperationFailedError) {
							if err != nil {
								cb(err)
								return
							}

							if attempt == nil {
								t.logger.logInfof(t.id, "ATR entry missing, completing write-write conflict poll")
								// The ATR entry is missing, which counts as it being completed.
								cb(nil)
								return
							}

							state := jsonAtrState(attempt.State)
							if state == jsonAtrStateCompleted || state == jsonAtrStateRolledBack {
								t.logger.logInfof(t.id, "Attempt state %s, completing write-write conflict poll", state)
								// If we have progressed enough to continue, let's do that.
								cb(nil)
								return
							}

							time.AfterFunc(200*time.Millisecond, onePoll)
						})
				})
			})
	}
	onePoll()
}

func (t *transactionAttempt) ensureCleanUpRequest() {
	// BUG(TXNG-59): Do not use a synchronous lock for cleanup requests.
	// Because of the need to include the state of the transaction within the cleanup
	// request, we are not able to do registration until the end of commit/rollback,
	// which means that we no longer have the lock on the transaction, and need to
	// relock it.
	t.lock.LockSync()

	if t.state == TransactionAttemptStateCompleted || t.state == TransactionAttemptStateRolledBack {
		t.lock.UnlockSync()
		t.logger.logInfof(t.id, "Attempt state completed or rolled back, will not add cleanup request")
		return
	}

	if t.hasCleanupRequest {
		t.lock.UnlockSync()
		t.logger.logInfof(t.id, "Attempt already created cleanup request, will not add cleanup request")
		return
	}

	t.hasCleanupRequest = true

	var inserts []TransactionsDocRecord
	var replaces []TransactionsDocRecord
	var removes []TransactionsDocRecord
	for _, staged := range t.stagedMutations {
		dr := TransactionsDocRecord{
			CollectionName: staged.CollectionName,
			ScopeName:      staged.ScopeName,
			BucketName:     staged.Agent.BucketName(),
			ID:             staged.Key,
		}

		switch staged.OpType {
		case TransactionStagedMutationInsert:
			inserts = append(inserts, dr)
		case TransactionStagedMutationReplace:
			replaces = append(replaces, dr)
		case TransactionStagedMutationRemove:
			removes = append(removes, dr)
		}
	}

	var bucketName string
	if t.atrAgent != nil {
		bucketName = t.atrAgent.BucketName()
	}

	cleanupState := t.state
	if cleanupState == TransactionAttemptStateCommitting {
		cleanupState = TransactionAttemptStatePending
	}

	req := &TransactionsCleanupRequest{
		AttemptID:         t.id,
		AtrID:             t.atrKey,
		AtrCollectionName: t.atrCollectionName,
		AtrScopeName:      t.atrScopeName,
		AtrBucketName:     bucketName,
		Inserts:           inserts,
		Replaces:          replaces,
		Removes:           removes,
		State:             cleanupState,
		ForwardCompat:     nil, // Let's just be explicit about this, it'll change in the future anyway.
		DurabilityLevel:   t.durabilityLevel,
		Age:               time.Since(t.txnStartTime),
	}

	t.lock.UnlockSync()

	t.logger.logInfof(t.id, "Adding cleanup request for atr %s, cleanup state: %s", newLoggableATRKey(
		bucketName,
		t.atrScopeName,
		t.atrCollectionName,
		t.atrKey,
	), cleanupState)

	t.addCleanupRequest(req)
}
