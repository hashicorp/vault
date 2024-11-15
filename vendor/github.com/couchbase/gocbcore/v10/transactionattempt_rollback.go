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
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

func (t *transactionAttempt) Rollback(cb TransactionRollbackCallback) error {
	return t.rollback(func(err *TransactionOperationFailedError) {
		if err != nil {
			t.logger.logInfof(t.id, "Rollback failed")
			t.ensureCleanUpRequest()
			cb(err)
			return
		}

		t.ensureCleanUpRequest()
		cb(nil)
	})
}

func (t *transactionAttempt) rollback(
	cb func(*TransactionOperationFailedError),
) error {
	t.logger.logInfof(t.id, "Rolling back")
	t.waitForOpsAndLock(func(unlock func()) {
		unlockAndCb := func(err *TransactionOperationFailedError) {
			unlock()
			cb(err)
		}

		err := t.checkCanRollbackLocked()
		if err != nil {
			unlockAndCb(err)
			return
		}

		t.applyStateBits(transactionStateBitShouldNotCommit|transactionStateBitShouldNotRollback, 0)

		if t.state == TransactionAttemptStateNothingWritten {
			unlockAndCb(nil)
			return
		}

		t.checkExpiredAtomic(hookRollback, []byte{}, true, func(cerr *classifiedError) {
			if cerr != nil {
				t.setExpiryOvertimeAtomic()
			}

			t.setATRAbortedLocked(func(err *TransactionOperationFailedError) {
				if err != nil {
					unlockAndCb(err)
					return
				}

				t.state = TransactionAttemptStateAborted

				go func() {
					removeStagedMutation := func(
						mutation *transactionStagedMutation,
						unstageCb func(*TransactionOperationFailedError),
					) {
						switch mutation.OpType {
						case TransactionStagedMutationInsert:
							t.removeStagedInsert(*mutation, unstageCb)
						case TransactionStagedMutationReplace:
							fallthrough
						case TransactionStagedMutationRemove:
							t.removeStagedRemoveReplace(*mutation, unstageCb)
						default:
							unstageCb(t.operationFailed(operationFailedDef{
								Cerr: classifyError(
									wrapError(ErrIllegalState, "unexpected staged mutation type")),
								ShouldNotRetry:    true,
								ShouldNotRollback: true,
							}))
						}
					}

					var mutErrs []*TransactionOperationFailedError
					if !t.enableParallelUnstaging {
						for _, mutation := range t.stagedMutations {
							waitCh := make(chan struct{}, 1)

							removeStagedMutation(mutation, func(err *TransactionOperationFailedError) {
								if err != nil {
									mutErrs = append(mutErrs, err)
									waitCh <- struct{}{}
									return
								}

								waitCh <- struct{}{}
							})

							<-waitCh
							if len(mutErrs) > 0 {
								break
							}
						}
					} else {
						type mutResult struct {
							Err *TransactionOperationFailedError
						}

						numMutations := len(t.stagedMutations)
						waitCh := make(chan mutResult, numMutations)

						// Unlike the RFC we do insert and replace separately. We have a bug in gocbcore where subdocs
						// will raise doc exists rather than a cas mismatch so we need to do these ops separately to tell
						// how to handle that error.
						for _, mutation := range t.stagedMutations {
							removeStagedMutation(mutation, func(err *TransactionOperationFailedError) {
								waitCh <- mutResult{
									Err: err,
								}
							})
						}

						for i := 0; i < numMutations; i++ {
							res := <-waitCh

							if res.Err != nil {
								mutErrs = append(mutErrs, res.Err)
								continue
							}
						}
					}
					err = mergeOperationFailedErrors(mutErrs)
					if err != nil {
						unlockAndCb(err)
						return
					}

					t.setATRRolledBackLocked(func(err *TransactionOperationFailedError) {
						if err != nil {
							unlockAndCb(err)
							return
						}

						t.state = TransactionAttemptStateRolledBack

						unlockAndCb(nil)
					})
				}()
			})
		})
	})

	return nil
}

func (t *transactionAttempt) removeStagedInsert(
	mutation transactionStagedMutation,
	cb func(*TransactionOperationFailedError),
) {
	ecCb := func(cerr *classifiedError) {
		if cerr == nil {
			cb(nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		if t.isExpiryOvertimeAtomic() {
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrAttemptExpired, "removing a staged insert failed during overtime")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
			return
		}

		switch cerr.Class {
		case TransactionErrorClassFailAmbiguous:
			time.AfterFunc(3*time.Millisecond, func() {
				t.removeStagedInsert(mutation, cb)
			})
		case TransactionErrorClassFailExpiry:
			t.setExpiryOvertimeAtomic()
			time.AfterFunc(3*time.Millisecond, func() {
				t.removeStagedInsert(mutation, cb)
			})
		case TransactionErrorClassFailDocNotFound:
			cb(nil)
			return
		case TransactionErrorClassFailPathNotFound:
			cb(nil)
			return
		case TransactionErrorClassFailDocAlreadyExists:
			cerr.Class = TransactionErrorClassFailCasMismatch
			fallthrough
		case TransactionErrorClassFailCasMismatch:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		case TransactionErrorClassFailHard:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		default:
			time.AfterFunc(3*time.Millisecond, func() {
				t.removeStagedInsert(mutation, cb)
			})
		}
	}

	t.checkExpiredAtomic(hookDeleteInserted, mutation.Key, true, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeRollbackDeleteInserted(mutation.Key, func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			_, err = mutation.Agent.MutateIn(MutateInOptions{
				ScopeName:      mutation.ScopeName,
				CollectionName: mutation.CollectionName,
				Key:            mutation.Key,
				Cas:            mutation.Cas,
				Flags:          memd.SubdocDocFlagAccessDeleted,
				Ops: []SubDocOp{
					{
						Op:    memd.SubDocOpDictSet,
						Path:  "txn",
						Flags: memd.SubdocFlagXattrPath,
						Value: []byte{110, 117, 108, 108}, // null
					},
					{
						Op:    memd.SubDocOpDelete,
						Path:  "txn",
						Flags: memd.SubdocFlagXattrPath,
					},
				},
				User: mutation.OboUser,
			}, func(result *MutateInResult, err error) {
				if err != nil {
					ecCb(classifyError(err))
					return
				}

				t.ReportResourceUnits(result.Internal.ResourceUnits)

				for _, op := range result.Ops {
					if op.Err != nil {
						ecCb(classifyError(op.Err))
						return
					}
				}

				t.hooks.AfterRollbackDeleteInserted(mutation.Key, func(err error) {
					if err != nil {
						ecCb(classifyHookError(err))
						return
					}

					ecCb(nil)
				})
			})
			if err != nil {
				ecCb(classifyError(err))
				return
			}
		})
	})
}

func (t *transactionAttempt) removeStagedRemoveReplace(
	mutation transactionStagedMutation,
	cb func(*TransactionOperationFailedError),
) {
	ecCb := func(cerr *classifiedError) {
		if cerr == nil {
			cb(nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		if t.isExpiryOvertimeAtomic() {
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrAttemptExpired, "removing a staged remove or replace failed during overtime")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
			return
		}

		switch cerr.Class {
		case TransactionErrorClassFailAmbiguous:
			time.AfterFunc(3*time.Millisecond, func() {
				t.removeStagedRemoveReplace(mutation, cb)
			})
		case TransactionErrorClassFailExpiry:
			t.setExpiryOvertimeAtomic()
			time.AfterFunc(3*time.Millisecond, func() {
				t.removeStagedRemoveReplace(mutation, cb)
			})
		case TransactionErrorClassFailPathNotFound:
			cb(nil)
			return
		case TransactionErrorClassFailDocNotFound:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		case TransactionErrorClassFailDocAlreadyExists:
			cerr.Class = TransactionErrorClassFailCasMismatch
			fallthrough
		case TransactionErrorClassFailCasMismatch:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		case TransactionErrorClassFailHard:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		default:
			time.AfterFunc(3*time.Millisecond, func() {
				t.removeStagedRemoveReplace(mutation, cb)
			})
		}
	}

	t.checkExpiredAtomic(hookRollbackDoc, mutation.Key, true, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeDocRolledBack(mutation.Key, func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			_, err = mutation.Agent.MutateIn(MutateInOptions{
				ScopeName:      mutation.ScopeName,
				CollectionName: mutation.CollectionName,
				Key:            mutation.Key,
				Cas:            mutation.Cas,
				Ops: []SubDocOp{
					{
						Op:    memd.SubDocOpDictSet,
						Path:  "txn",
						Flags: memd.SubdocFlagXattrPath,
						Value: []byte{110, 117, 108, 108}, // null
					},
					{
						Op:    memd.SubDocOpDelete,
						Path:  "txn",
						Flags: memd.SubdocFlagXattrPath,
					},
				},
				User: mutation.OboUser,
			}, func(result *MutateInResult, err error) {
				if err != nil {
					ecCb(classifyError(err))
					return
				}

				t.ReportResourceUnits(result.Internal.ResourceUnits)

				for _, op := range result.Ops {
					if op.Err != nil {
						ecCb(classifyError(op.Err))
						return
					}
				}

				t.hooks.AfterRollbackReplaceOrRemove(mutation.Key, func(err error) {
					if err != nil {
						ecCb(classifyHookError(err))
						return
					}

					ecCb(nil)
				})
			})
			if err != nil {
				ecCb(classifyError(err))
				return
			}
		})
	})
}
