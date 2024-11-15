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
	"encoding/json"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

func (t *transactionAttempt) Insert(opts TransactionInsertOptions, cb TransactionStoreCallback) error {
	return t.insert(opts, func(res *TransactionGetResult, err error) {
		if err != nil {
			var e *TransactionOperationFailedError
			if errors.As(err, &e) {
				if e.shouldNotRollback {
					t.ensureCleanUpRequest()
				}
			}

			cb(nil, err)
			return
		}

		cb(res, nil)
	})
}

func (t *transactionAttempt) insert(
	opts TransactionInsertOptions,
	cb func(*TransactionGetResult, error),
) error {
	t.logger.logInfof(t.id, "Performing insert for %s", newLoggableDocKey(
		opts.Agent.BucketName(),
		opts.ScopeName,
		opts.CollectionName,
		opts.Key,
	))

	t.beginOpAndLock(func(unlock func(), endOp func()) {
		endAndCb := func(result *TransactionGetResult, err error) {
			endOp()
			cb(result, err)
		}

		err := t.checkCanPerformOpLocked()
		if err != nil {
			unlock()
			endAndCb(nil, err)
			return
		}

		agent := opts.Agent
		oboUser := opts.OboUser
		scopeName := opts.ScopeName
		collectionName := opts.CollectionName
		key := opts.Key
		value := opts.Value

		t.checkExpiredAtomic(hookInsert, key, false, func(cerr *classifiedError) {
			if cerr != nil {
				unlock()
				endAndCb(nil, t.operationFailed(operationFailedDef{
					Cerr:              cerr,
					ShouldNotRetry:    true,
					ShouldNotRollback: false,
					Reason:            TransactionErrorReasonTransactionExpired,
				}))
				return
			}

			_, existingMutation := t.getStagedMutationLocked(agent.BucketName(), scopeName, collectionName, key)
			unlock()

			if existingMutation != nil {
				switch existingMutation.OpType {
				case TransactionStagedMutationRemove:
					t.logger.logInfof(t.id, "Staged remove exists on doc, performing replace")
					t.stageReplace(
						agent, oboUser, scopeName, collectionName, key,
						value, existingMutation.Cas,
						func(result *TransactionGetResult, err error) {
							endAndCb(result, err)
						})
					return
				case TransactionStagedMutationInsert:
					endAndCb(nil, t.operationFailed(operationFailedDef{
						Cerr: classifyError(
							wrapError(ErrDocumentExists, "attempted to insert a document previously inserted in this transaction")),
						ShouldNotRetry:    true,
						ShouldNotRollback: false,
						Reason:            TransactionErrorReasonTransactionFailed,
					}))
					return
				case TransactionStagedMutationReplace:
					endAndCb(nil, t.operationFailed(operationFailedDef{
						Cerr: classifyError(
							wrapError(ErrDocumentExists, "attempted to insert a document previously replaced in this transaction")),
						ShouldNotRetry:    true,
						ShouldNotRollback: false,
						Reason:            TransactionErrorReasonTransactionFailed,
					}))
					return
				default:
					endAndCb(nil, t.operationFailed(operationFailedDef{
						Cerr: classifyError(
							wrapError(ErrIllegalState, "unexpected staged mutation type")),
						ShouldNotRetry:    true,
						ShouldNotRollback: false,
						Reason:            TransactionErrorReasonTransactionFailed,
					}))
					return
				}
			}

			t.confirmATRPending(agent, oboUser, scopeName, collectionName, key, func(err *TransactionOperationFailedError) {
				if err != nil {
					endAndCb(nil, err)
					return
				}

				t.stageInsert(
					agent, oboUser, scopeName, collectionName, key,
					value, 0,
					func(result *TransactionGetResult, err error) {
						endAndCb(result, err)
					})
			})
		})
	})

	return nil
}

func (t *transactionAttempt) resolveConflictedInsert(
	agent *Agent,
	oboUser string,
	scopeName string,
	collectionName string,
	key []byte,
	value json.RawMessage,
	cb func(*TransactionGetResult, error),
) {
	t.getMetaForConflictedInsert(agent, oboUser, scopeName, collectionName, key,
		func(isTombstone bool, txnMeta *jsonTxnXattr, cas Cas, err error) {
			if err != nil {
				cb(nil, err)
				return
			}

			if txnMeta == nil {
				// This doc isn't in a transaction
				if !isTombstone {
					cb(nil, ErrDocumentExists)
					return
				}

				// There wasn't actually a staged mutation there.
				t.stageInsert(agent, oboUser, scopeName, collectionName, key, value, cas, cb)
				return
			}

			meta := &TransactionMutableItemMeta{
				TransactionID: txnMeta.ID.Transaction,
				AttemptID:     txnMeta.ID.Attempt,
				ATR: TransactionMutableItemMetaATR{
					BucketName:     txnMeta.ATR.BucketName,
					ScopeName:      txnMeta.ATR.ScopeName,
					CollectionName: txnMeta.ATR.CollectionName,
					DocID:          txnMeta.ATR.DocID,
				},
				ForwardCompat: jsonForwardCompatToForwardCompat(txnMeta.ForwardCompat),
			}

			t.checkForwardCompatability(
				key,
				agent.BucketName(),
				scopeName,
				collectionName,
				forwardCompatStageWWCInsertingGet, meta.ForwardCompat, false, func(err *TransactionOperationFailedError) {
					if err != nil {
						cb(nil, err)
						return
					}

					if txnMeta.Operation.Type != jsonMutationInsert {
						cb(nil, t.operationFailed(operationFailedDef{
							Cerr: classifyError(
								wrapError(ErrDocumentExists, "found staged non-insert mutation")),
							ShouldNotRetry:    true,
							ShouldNotRollback: false,
							Reason:            TransactionErrorReasonTransactionFailed,
						}))
						return
					}

					// We have guards in place within the write write conflict polling to prevent miss-use when
					// an existing mutation must have been discovered before it's safe to overwrite.  This logic
					// is unneccessary, as is the forwards compatibility check when resolving conflicted inserts
					// so we can safely just ignore it.
					if meta.TransactionID == t.transactionID && meta.AttemptID == t.id {
						t.stageInsert(agent, oboUser, scopeName, collectionName, key, value, cas, cb)
						return
					}

					t.writeWriteConflictPoll(forwardCompatStageWWCInserting, agent, oboUser, scopeName, collectionName, key, cas, meta, nil, func(err *TransactionOperationFailedError) {
						if err != nil {
							cb(nil, err)
							return
						}

						t.cleanupStagedInsert(agent, oboUser, scopeName, collectionName, key, cas, isTombstone, func(cas Cas, err *TransactionOperationFailedError) {
							if err != nil {
								cb(nil, err)
								return
							}

							t.stageInsert(agent, oboUser, scopeName, collectionName, key, value, cas, cb)
						})
					})
				})
		})
}

func (t *transactionAttempt) stageInsert(
	agent *Agent,
	oboUser string,
	scopeName string,
	collectionName string,
	key []byte,
	value json.RawMessage,
	cas Cas,
	cb func(*TransactionGetResult, error),
) {
	ecCb := func(result *TransactionGetResult, cerr *classifiedError) {
		if cerr == nil {
			cb(result, nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		switch cerr.Class {
		case TransactionErrorClassFailAmbiguous:
			time.AfterFunc(3*time.Millisecond, func() {
				t.stageInsert(agent, oboUser, scopeName, collectionName, key, value, cas, cb)
			})
		case TransactionErrorClassFailExpiry:
			t.setExpiryOvertimeAtomic()
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionExpired,
			}))
		case TransactionErrorClassFailTransient:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailHard:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailDocAlreadyExists:
			fallthrough
		case TransactionErrorClassFailCasMismatch:
			t.resolveConflictedInsert(agent, oboUser, scopeName, collectionName, key, value, cb)
		default:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	}

	t.checkExpiredAtomic(hookInsert, key, false, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(nil, cerr)
			return
		}

		t.hooks.BeforeStagedInsert(key, func(err error) {
			if err != nil {
				ecCb(nil, classifyHookError(err))
				return
			}

			stagedInfo := &transactionStagedMutation{
				OpType:         TransactionStagedMutationInsert,
				Agent:          agent,
				OboUser:        oboUser,
				ScopeName:      scopeName,
				CollectionName: collectionName,
				Key:            key,
				Staged:         value,
			}

			var txnMeta jsonTxnXattr
			txnMeta.ID.Transaction = t.transactionID
			txnMeta.ID.Attempt = t.id
			txnMeta.ATR.CollectionName = t.atrCollectionName
			txnMeta.ATR.ScopeName = t.atrScopeName
			txnMeta.ATR.BucketName = t.atrAgent.BucketName()
			txnMeta.ATR.DocID = string(t.atrKey)
			txnMeta.Operation.Type = jsonMutationInsert
			txnMeta.Operation.Staged = stagedInfo.Staged

			txnMetaBytes, err := json.Marshal(txnMeta)
			if err != nil {
				ecCb(nil, classifyError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			flags := memd.SubdocDocFlagCreateAsDeleted | memd.SubdocDocFlagAccessDeleted
			var txnOp memd.SubDocOpType
			if cas == 0 {
				flags |= memd.SubdocDocFlagAddDoc
				txnOp = memd.SubDocOpDictAdd
			} else {
				txnOp = memd.SubDocOpDictSet
			}

			_, err = stagedInfo.Agent.MutateIn(MutateInOptions{
				ScopeName:      stagedInfo.ScopeName,
				CollectionName: stagedInfo.CollectionName,
				Key:            stagedInfo.Key,
				Cas:            cas,
				Ops: []SubDocOp{
					{
						Op:    txnOp,
						Path:  "txn",
						Flags: memd.SubdocFlagMkDirP | memd.SubdocFlagXattrPath,
						Value: txnMetaBytes,
					},
					{
						Op:    memd.SubDocOpDictSet,
						Path:  "txn.op.crc32",
						Flags: memd.SubdocFlagXattrPath | memd.SubdocFlagExpandMacros,
						Value: crc32cMacro,
					},
				},
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Deadline:               deadline,
				Flags:                  flags,
				User:                   stagedInfo.OboUser,
			}, func(result *MutateInResult, err error) {
				if err != nil {
					ecCb(nil, classifyError(err))
					return
				}

				t.ReportResourceUnits(result.Internal.ResourceUnits)

				stagedInfo.Cas = result.Cas

				t.hooks.AfterStagedInsertComplete(key, func(err error) {
					if err != nil {
						ecCb(nil, classifyHookError(err))
						return
					}

					t.recordStagedMutation(stagedInfo, func() {

						ecCb(&TransactionGetResult{
							agent:          stagedInfo.Agent,
							oboUser:        stagedInfo.OboUser,
							scopeName:      stagedInfo.ScopeName,
							collectionName: stagedInfo.CollectionName,
							key:            stagedInfo.Key,
							Value:          stagedInfo.Staged,
							Cas:            stagedInfo.Cas,
							Meta:           nil,
						}, nil)
					})
				})
			})
			if err != nil {
				ecCb(nil, classifyError(err))
			}
		})
	})
}

func (t *transactionAttempt) getMetaForConflictedInsert(
	agent *Agent,
	oboUser string,
	scopeName string,
	collectionName string,
	key []byte,
	cb func(bool, *jsonTxnXattr, Cas, error),
) {
	ecCb := func(isTombstone bool, meta *jsonTxnXattr, cas Cas, cerr *classifiedError) {
		if cerr == nil {
			cb(isTombstone, meta, cas, nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		switch cerr.Class {
		case TransactionErrorClassFailDocNotFound:
			fallthrough
		case TransactionErrorClassFailPathNotFound:
			fallthrough
		case TransactionErrorClassFailTransient:
			cb(isTombstone, nil, 0, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		default:
			cb(isTombstone, nil, 0, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	}

	t.hooks.BeforeGetDocInExistsDuringStagedInsert(key, func(err error) {
		if err != nil {
			ecCb(false, nil, 0, classifyHookError(err))
			return
		}

		var deadline time.Time
		if t.keyValueTimeout > 0 {
			deadline = time.Now().Add(t.keyValueTimeout)
		}

		_, err = agent.LookupIn(LookupInOptions{
			ScopeName:      scopeName,
			CollectionName: collectionName,
			Key:            key,
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpGet,
					Path:  "txn",
					Flags: memd.SubdocFlagXattrPath,
				},
			},
			Deadline: deadline,
			Flags:    memd.SubdocDocFlagAccessDeleted,
			User:     oboUser,
		}, func(result *LookupInResult, err error) {
			if err != nil {
				ecCb(false, nil, 0, classifyError(err))
				return
			}

			t.ReportResourceUnits(result.Internal.ResourceUnits)

			var txnMeta *jsonTxnXattr
			if result.Ops[0].Err == nil {
				var txnMetaVal jsonTxnXattr
				if err := json.Unmarshal(result.Ops[0].Value, &txnMetaVal); err != nil {
					ecCb(false, nil, 0, classifyError(err))
					return
				}
				txnMeta = &txnMetaVal
			}

			isTombstone := result.Internal.IsDeleted
			ecCb(isTombstone, txnMeta, result.Cas, nil)
		})
		if err != nil {
			ecCb(false, nil, 0, classifyError(err))
			return
		}
	})
}

func (t *transactionAttempt) cleanupStagedInsert(
	agent *Agent,
	oboUser string,
	scopeName string,
	collectionName string,
	key []byte,
	cas Cas,
	isTombstone bool,
	cb func(Cas, *TransactionOperationFailedError),
) {
	ecCb := func(cas Cas, cerr *classifiedError) {
		if cerr == nil {
			cb(cas, nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		switch cerr.Class {
		case TransactionErrorClassFailDocNotFound:
			fallthrough
		case TransactionErrorClassFailCasMismatch:
			fallthrough
		case TransactionErrorClassFailTransient:
			cb(0, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		default:
			cb(0, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	}

	if isTombstone {
		// This is already a tombstone, so we can just proceed.
		ecCb(cas, nil)
		return
	}

	t.hooks.BeforeRemovingDocDuringStagedInsert(key, func(err error) {
		if err != nil {
			ecCb(0, classifyHookError(err))
			return
		}

		var deadline time.Time
		if t.keyValueTimeout > 0 {
			deadline = time.Now().Add(t.keyValueTimeout)
		}

		_, err = agent.Delete(DeleteOptions{
			ScopeName:      scopeName,
			CollectionName: collectionName,
			Key:            key,
			Deadline:       deadline,
			User:           oboUser,
		}, func(result *DeleteResult, err error) {
			if err != nil {
				ecCb(0, classifyError(err))
				return
			}

			t.ReportResourceUnits(result.Internal.ResourceUnits)

			ecCb(result.Cas, nil)
		})
		if err != nil {
			ecCb(0, classifyError(err))
			return
		}
	})
}
