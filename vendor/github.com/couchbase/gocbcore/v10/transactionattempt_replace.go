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
	"github.com/couchbase/gocbcore/v10/memd"
)

func (t *transactionAttempt) Replace(opts TransactionReplaceOptions, cb TransactionStoreCallback) error {
	return t.replace(opts, func(res *TransactionGetResult, err error) {
		if err != nil {
			t.logger.logInfof(t.id, "Replace failed")
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

func (t *transactionAttempt) replace(
	opts TransactionReplaceOptions,
	cb func(*TransactionGetResult, error),
) error {
	t.logger.logInfof(t.id, "Performing replace for %s", newLoggableDocKey(
		opts.Document.agent.BucketName(),
		opts.Document.scopeName,
		opts.Document.collectionName,
		opts.Document.key,
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

		agent := opts.Document.agent
		oboUser := opts.Document.oboUser
		scopeName := opts.Document.scopeName
		collectionName := opts.Document.collectionName
		key := opts.Document.key
		value := opts.Value
		cas := opts.Document.Cas
		meta := opts.Document.Meta

		t.checkExpiredAtomic(hookReplace, key, false, func(cerr *classifiedError) {
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
				case TransactionStagedMutationInsert:
					t.logger.logInfof(t.id, "Staged insert exists on doc, performing insert")
					t.stageInsert(
						agent, oboUser, scopeName, collectionName, key,
						value, cas,
						func(result *TransactionGetResult, err error) {
							endAndCb(result, err)
						})
					return

				case TransactionStagedMutationReplace:
					t.logger.logInfof(t.id, "Staged replace exists on doc, this is ok")
					// We can overwrite other replaces without issue, any conflicts between the mutation
					// the user passed to us and the existing mutation is caught by WriteWriteConflict.
				case TransactionStagedMutationRemove:
					endAndCb(nil, t.operationFailed(operationFailedDef{
						Cerr: classifyError(
							wrapError(ErrDocumentNotFound, "attempted to replace a document previously removed in this transaction")),
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

			t.writeWriteConflictPoll(
				forwardCompatStageWWCReplacing,
				agent, oboUser, scopeName, collectionName, key, cas,
				meta,
				existingMutation,
				func(err *TransactionOperationFailedError) {
					if err != nil {
						endAndCb(nil, err)
						return
					}

					t.confirmATRPending(agent, oboUser, scopeName, collectionName, key, func(err *TransactionOperationFailedError) {
						if err != nil {
							endAndCb(nil, err)
							return
						}

						t.stageReplace(
							agent, oboUser, scopeName, collectionName, key,
							value, cas,
							func(result *TransactionGetResult, err error) {
								endAndCb(result, err)
							})
					})
				})
		})
	})

	return nil
}

func (t *transactionAttempt) stageReplace(
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
		case TransactionErrorClassFailExpiry:
			t.setExpiryOvertimeAtomic()
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionExpired,
			}))
		case TransactionErrorClassFailDocNotFound:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrDocumentNotFound, "document not found during staging")),
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailDocAlreadyExists:
			cerr.Class = TransactionErrorClassFailCasMismatch
			fallthrough
		case TransactionErrorClassFailCasMismatch:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailTransient:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailAmbiguous:
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
		default:
			cb(nil, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	}

	t.checkExpiredAtomic(hookRemove, key, false, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(nil, cerr)
			return
		}

		t.hooks.BeforeStagedReplace(key, func(err error) {
			if err != nil {
				ecCb(nil, classifyHookError(err))
				return
			}

			stagedInfo := &transactionStagedMutation{
				OpType:         TransactionStagedMutationReplace,
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
			txnMeta.Operation.Type = jsonMutationReplace
			txnMeta.Operation.Staged = stagedInfo.Staged
			txnMeta.Restore = &jsonTxnXattrRestore{
				OriginalCAS: "",
				ExpiryTime:  0,
				RevID:       "",
			}

			txnMetaBytes, err := json.Marshal(txnMeta)
			if err != nil {
				ecCb(nil, classifyError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			_, err = stagedInfo.Agent.MutateIn(MutateInOptions{
				ScopeName:      stagedInfo.ScopeName,
				CollectionName: stagedInfo.CollectionName,
				Key:            stagedInfo.Key,
				Cas:            cas,
				Ops: []SubDocOp{
					{
						Op:    memd.SubDocOpDictSet,
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
					{
						Op:    memd.SubDocOpDictSet,
						Path:  "txn.restore.CAS",
						Flags: memd.SubdocFlagXattrPath | memd.SubdocFlagExpandMacros,
						Value: casMacro,
					},
					{
						Op:    memd.SubDocOpDictSet,
						Path:  "txn.restore.exptime",
						Flags: memd.SubdocFlagXattrPath | memd.SubdocFlagExpandMacros,
						Value: exptimeMacro,
					},
					{
						Op:    memd.SubDocOpDictSet,
						Path:  "txn.restore.revid",
						Flags: memd.SubdocFlagXattrPath | memd.SubdocFlagExpandMacros,
						Value: revidMacro,
					},
				},
				Flags:                  memd.SubdocDocFlagAccessDeleted,
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Deadline:               deadline,
				User:                   stagedInfo.OboUser,
			}, func(result *MutateInResult, err error) {
				if err != nil {
					ecCb(nil, classifyError(err))
					return
				}

				t.ReportResourceUnits(result.Internal.ResourceUnits)

				stagedInfo.Cas = result.Cas

				t.hooks.AfterStagedReplaceComplete(key, func(err error) {
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
				return
			}
		})
	})
}
