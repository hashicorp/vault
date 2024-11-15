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
	"fmt"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

func (t *transactionAttempt) selectAtrLocked(
	firstAgent *Agent,
	firstOboUser string,
	firstScopeName string,
	firstCollectionName string,
	firstKey []byte,
	cb func(*TransactionOperationFailedError),
) {
	atrID := int(cbcVbMap(firstKey, 1024))
	atrKey := []byte(transactionAtrIDList[atrID])

	t.hooks.RandomATRIDForVbucket(func(s string, err error) {
		if err != nil {
			cb(t.operationFailed(operationFailedDef{
				Cerr:              classifyHookError(err),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
			return
		}

		if s != "" {
			atrKey = []byte(s)
		}

		atrAgent := firstAgent
		atrOboUser := firstOboUser
		atrScopeName := "_default"
		atrCollectionName := "_default"
		if t.atrLocation.Agent != nil {
			atrAgent = t.atrLocation.Agent
			atrOboUser = t.atrLocation.OboUser
			atrScopeName = t.atrLocation.ScopeName
			atrCollectionName = t.atrLocation.CollectionName
		} else {
			if t.enableExplicitATRs {
				cb(t.operationFailed(operationFailedDef{
					Cerr:              classifyError(errors.New("atrs must be explicitly defined")),
					ShouldNotRetry:    true,
					ShouldNotRollback: true,
					Reason:            TransactionErrorReasonTransactionFailed,
				}))
				return
			}
		}

		t.atrAgent = atrAgent
		t.atrOboUser = atrOboUser
		t.atrScopeName = atrScopeName
		t.atrCollectionName = atrCollectionName
		t.atrKey = atrKey

		cb(nil)
	})
}

func (t *transactionAttempt) setATRPendingLocked(
	cb func(*TransactionOperationFailedError),
) {
	ecCb := func(cerr *classifiedError) {
		if cerr == nil {
			cb(nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		switch cerr.Class {
		case TransactionErrorClassFailAmbiguous:
			time.AfterFunc(3*time.Millisecond, func() {
				t.setATRPendingLocked(cb)
			})
			return
		case TransactionErrorClassFailPathAlreadyExists:
			cb(nil)
			return
		case TransactionErrorClassFailExpiry:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionExpired,
			}))
		case TransactionErrorClassFailOutOfSpace:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrFull),
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailTransient:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case TransactionErrorClassFailHard:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		default:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	}

	t.checkExpiredAtomic(hookATRPending, []byte{}, false, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeATRPending(func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			var marshalErr error
			atrFieldOp := func(fieldName string, data interface{}, flags memd.SubdocFlag) SubDocOp {
				b, err := json.Marshal(data)
				if err != nil {
					marshalErr = err
					return SubDocOp{}
				}

				return SubDocOp{
					Op:    memd.SubDocOpDictAdd,
					Flags: memd.SubdocFlagMkDirP | flags,
					Path:  "attempts." + t.id + "." + fieldName,
					Value: b,
				}
			}

			atrOps := []SubDocOp{
				atrFieldOp("tst", "${Mutation.CAS}", memd.SubdocFlagXattrPath|memd.SubdocFlagExpandMacros),
				atrFieldOp("tid", t.transactionID, memd.SubdocFlagXattrPath),
				atrFieldOp("st", jsonAtrStatePending, memd.SubdocFlagXattrPath),
				atrFieldOp("exp", time.Until(t.expiryTime)/time.Millisecond, memd.SubdocFlagXattrPath),
				{
					Op:    memd.SubDocOpSetDoc,
					Flags: memd.SubdocFlagNone,
					Path:  "",
					Value: []byte{0},
				},
				atrFieldOp("d", transactionsDurabilityLevelToShorthand(t.durabilityLevel), memd.SubdocFlagXattrPath),
			}
			if marshalErr != nil {
				ecCb(classifyError(marshalErr))
				return
			}
			t.logger.logInfof(t.id, "Setting ATR %s pending", newLoggableATRKey(
				t.atrAgent.BucketName(),
				t.atrScopeName,
				t.atrCollectionName,
				t.atrKey,
			))

			_, err = t.atrAgent.MutateIn(MutateInOptions{
				ScopeName:              t.atrScopeName,
				CollectionName:         t.atrCollectionName,
				Key:                    t.atrKey,
				Ops:                    atrOps,
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Deadline:               deadline,
				Flags:                  memd.SubdocDocFlagMkDoc,
				User:                   t.atrOboUser,
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

				t.hooks.AfterATRPending(func(err error) {
					if err != nil {
						ecCb(classifyHookError(err))
						return
					}

					t.addLostCleanupLocation(t.atrAgent.BucketName(), t.atrScopeName, t.atrCollectionName)

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

func (t *transactionAttempt) fetchATRCommitConflictLocked(
	cb func(jsonAtrState, *TransactionOperationFailedError),
) {
	ecCb := func(st jsonAtrState, cerr *classifiedError) {
		if cerr == nil {
			cb(st, nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		switch cerr.Class {
		case TransactionErrorClassFailTransient:
			fallthrough
		case TransactionErrorClassFailOther:
			time.AfterFunc(3*time.Millisecond, func() {
				t.fetchATRCommitConflictLocked(cb)
			})
			return
		case TransactionErrorClassFailDocNotFound:
			cb(jsonAtrStateUnknown, t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrNotFound),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionCommitAmbiguous,
			}))
		case TransactionErrorClassFailPathNotFound:
			cb(jsonAtrStateUnknown, t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrEntryNotFound),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionCommitAmbiguous,
			}))
		case TransactionErrorClassFailExpiry:
			cb(jsonAtrStateUnknown, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionCommitAmbiguous,
			}))
		case TransactionErrorClassFailHard:
			cb(jsonAtrStateUnknown, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionCommitAmbiguous,
			}))
		default:
			cb(jsonAtrStateUnknown, t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionCommitAmbiguous,
			}))
			return
		}
	}

	t.checkExpiredAtomic(hookATRCommitAmbiguityResolution, []byte{}, false, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(jsonAtrStateUnknown, cerr)
			return
		}

		t.hooks.BeforeATRCommitAmbiguityResolution(func(err error) {
			if err != nil {
				ecCb(jsonAtrStateUnknown, classifyHookError(err))
				return
			}

			var deadline time.Time
			if t.keyValueTimeout > 0 {
				deadline = time.Now().Add(t.keyValueTimeout)
			}

			_, err = t.atrAgent.LookupIn(LookupInOptions{
				ScopeName:      t.atrScopeName,
				CollectionName: t.atrCollectionName,
				Key:            t.atrKey,
				Ops: []SubDocOp{
					{
						Op:    memd.SubDocOpGet,
						Path:  "attempts." + t.id + ".st",
						Flags: memd.SubdocFlagXattrPath,
					},
				},
				Deadline: deadline,
				Flags:    memd.SubdocDocFlagNone,
				User:     t.atrOboUser,
			}, func(result *LookupInResult, err error) {
				if err != nil {
					ecCb(jsonAtrStateUnknown, classifyError(err))
					return
				}

				t.ReportResourceUnits(result.Internal.ResourceUnits)

				if result.Ops[0].Err != nil {
					ecCb(jsonAtrStateUnknown, classifyError(err))
					return
				}

				var st jsonAtrState
				if err := json.Unmarshal(result.Ops[0].Value, &st); err != nil {
					ecCb(jsonAtrStateUnknown, classifyError(err))
					return
				}

				ecCb(st, nil)
			})
			if err != nil {
				ecCb(jsonAtrStateUnknown, classifyError(err))
				return
			}
		})
	})
}

func (t *transactionAttempt) resolveATRCommitConflictLocked(
	cb func(*TransactionOperationFailedError),
) {
	t.fetchATRCommitConflictLocked(func(st jsonAtrState, err *TransactionOperationFailedError) {
		if err != nil {
			cb(err)
			return
		}

		switch st {
		case jsonAtrStatePending:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrIllegalState, "transaction still pending even with p set during commit")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case jsonAtrStateCommitted:
			cb(nil)
		case jsonAtrStateCompleted:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrIllegalState, "transaction already completed during commit")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case jsonAtrStateAborted:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrIllegalState, "transaction already aborted during commit")),
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		case jsonAtrStateRolledBack:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrIllegalState, "transaction already rolled back during commit")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		default:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrIllegalState, fmt.Sprintf("illegal transaction state during commit: %s", st))),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailed,
			}))
		}
	})
}

func (t *transactionAttempt) setATRCommittedLocked(
	ambiguityResolution bool,
	cb func(*TransactionOperationFailedError),
) {
	ecCb := func(cerr *classifiedError) {
		if cerr == nil {
			cb(nil)
			return
		}

		t.ReportResourceUnitsError(cerr.Source)

		errorReason := TransactionErrorReasonTransactionFailed
		if ambiguityResolution {
			errorReason = TransactionErrorReasonTransactionCommitAmbiguous
		}

		switch cerr.Class {
		case TransactionErrorClassFailAmbiguous:
			time.AfterFunc(3*time.Millisecond, func() {
				ambiguityResolution = true
				t.setATRCommittedLocked(ambiguityResolution, cb)
			})
			return
		case TransactionErrorClassFailTransient:
			if ambiguityResolution {
				time.AfterFunc(3*time.Millisecond, func() {
					t.setATRCommittedLocked(ambiguityResolution, cb)
				})
				return
			}

			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    false,
				ShouldNotRollback: false,
				Reason:            errorReason,
			}))
		case TransactionErrorClassFailPathAlreadyExists:
			t.resolveATRCommitConflictLocked(cb)
			return
		case TransactionErrorClassFailDocNotFound:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrNotFound),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            errorReason,
			}))
		case TransactionErrorClassFailPathNotFound:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrEntryNotFound),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            errorReason,
			}))
		case TransactionErrorClassFailOutOfSpace:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrFull),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            errorReason,
			}))
		case TransactionErrorClassFailExpiry:
			if errorReason == TransactionErrorReasonTransactionFailed {
				errorReason = TransactionErrorReasonTransactionExpired
			}

			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            errorReason,
			}))
		case TransactionErrorClassFailHard:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            errorReason,
			}))
		default:
			if ambiguityResolution {
				cb(t.operationFailed(operationFailedDef{
					Cerr:              cerr,
					ShouldNotRetry:    true,
					ShouldNotRollback: true,
					Reason:            errorReason,
				}))
				return
			}

			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: false,
				Reason:            errorReason,
			}))
		}
	}

	atrAgent := t.atrAgent
	atrOboUser := t.atrOboUser
	atrScopeName := t.atrScopeName
	atrKey := t.atrKey
	atrCollectionName := t.atrCollectionName

	insMutations := []jsonAtrMutation{}
	repMutations := []jsonAtrMutation{}
	remMutations := []jsonAtrMutation{}

	for _, mutation := range t.stagedMutations {
		jsonMutation := jsonAtrMutation{
			BucketName:     mutation.Agent.BucketName(),
			ScopeName:      mutation.ScopeName,
			CollectionName: mutation.CollectionName,
			DocID:          string(mutation.Key),
		}

		if mutation.OpType == TransactionStagedMutationInsert {
			insMutations = append(insMutations, jsonMutation)
		} else if mutation.OpType == TransactionStagedMutationReplace {
			repMutations = append(repMutations, jsonMutation)
		} else if mutation.OpType == TransactionStagedMutationRemove {
			remMutations = append(remMutations, jsonMutation)
		} else {
			ecCb(classifyError(wrapError(ErrIllegalState, "unexpected staged mutation type")))
			return
		}
	}

	t.checkExpiredAtomic(hookATRCommit, []byte{}, false, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeATRCommit(func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			var marshalErr error
			atrFieldOp := func(fieldName string, data interface{}, flags memd.SubdocFlag, op memd.SubDocOpType) SubDocOp {
				bytes, err := json.Marshal(data)
				if err != nil {
					marshalErr = err
				}

				return SubDocOp{
					Op:    op,
					Flags: flags,
					Path:  "attempts." + t.id + "." + fieldName,
					Value: bytes,
				}
			}

			atrOps := []SubDocOp{
				atrFieldOp("st", jsonAtrStateCommitted, memd.SubdocFlagXattrPath, memd.SubDocOpDictSet),
				atrFieldOp("tsc", "${Mutation.CAS}", memd.SubdocFlagXattrPath|memd.SubdocFlagExpandMacros, memd.SubDocOpDictSet),
				atrFieldOp("p", 0, memd.SubdocFlagXattrPath, memd.SubDocOpDictAdd),
				atrFieldOp("ins", insMutations, memd.SubdocFlagXattrPath, memd.SubDocOpDictSet),
				atrFieldOp("rep", repMutations, memd.SubdocFlagXattrPath, memd.SubDocOpDictSet),
				atrFieldOp("rem", remMutations, memd.SubdocFlagXattrPath, memd.SubDocOpDictSet),
			}
			if marshalErr != nil {
				ecCb(classifyError(marshalErr))
				return
			}

			_, err = atrAgent.MutateIn(MutateInOptions{
				ScopeName:              atrScopeName,
				CollectionName:         atrCollectionName,
				Key:                    atrKey,
				Ops:                    atrOps,
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Flags:                  memd.SubdocDocFlagNone,
				Deadline:               deadline,
				User:                   atrOboUser,
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

				t.hooks.AfterATRCommit(func(err error) {
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

func (t *transactionAttempt) setATRCompletedLocked(
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
					wrapError(ErrAttemptExpired, "completed atr removal failed during overtime")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailedPostCommit,
			}))
			return
		}

		switch cerr.Class {
		case TransactionErrorClassFailDocNotFound:
			fallthrough
		case TransactionErrorClassFailPathNotFound:
			// This is technically a full success, but FIT expects unstagingCompleted=false...
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailedPostCommit,
			}))
		case TransactionErrorClassFailExpiry:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrAttemptExpired, "completed atr removal operation expired")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailedPostCommit,
			}))
		case TransactionErrorClassFailHard:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailedPostCommit,
			}))
		default:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr,
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
				Reason:            TransactionErrorReasonTransactionFailedPostCommit,
			}))
		}
	}

	atrAgent := t.atrAgent
	atrOboUser := t.atrOboUser
	atrScopeName := t.atrScopeName
	atrKey := t.atrKey
	atrCollectionName := t.atrCollectionName

	t.checkExpiredAtomic(hookATRComplete, []byte{}, true, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeATRComplete(func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			atrOps := []SubDocOp{
				{
					Op:    memd.SubDocOpDelete,
					Flags: memd.SubdocFlagXattrPath,
					Path:  "attempts." + t.id,
				},
			}

			_, err = atrAgent.MutateIn(MutateInOptions{
				ScopeName:              atrScopeName,
				CollectionName:         atrCollectionName,
				Key:                    atrKey,
				Ops:                    atrOps,
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Deadline:               deadline,
				Flags:                  memd.SubdocDocFlagNone,
				User:                   atrOboUser,
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

				t.hooks.AfterATRComplete(func(err error) {
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

func (t *transactionAttempt) setATRAbortedLocked(
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
					wrapError(ErrAttemptExpired, "atr abort failed during overtime")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
			return
		}

		switch cerr.Class {
		case TransactionErrorClassFailExpiry:
			t.setExpiryOvertimeAtomic()
			time.AfterFunc(3*time.Millisecond, func() {
				t.setATRAbortedLocked(cb)
			})
		case TransactionErrorClassFailDocNotFound:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrNotFound),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		case TransactionErrorClassFailPathNotFound:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrEntryNotFound),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
		case TransactionErrorClassFailOutOfSpace:
			cb(t.operationFailed(operationFailedDef{
				Cerr:              cerr.Wrap(ErrAtrFull),
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
				t.setATRAbortedLocked(cb)
			})
		}
	}

	atrAgent := t.atrAgent
	atrOboUser := t.atrOboUser
	atrScopeName := t.atrScopeName
	atrKey := t.atrKey
	atrCollectionName := t.atrCollectionName

	insMutations := []jsonAtrMutation{}
	repMutations := []jsonAtrMutation{}
	remMutations := []jsonAtrMutation{}

	for _, mutation := range t.stagedMutations {
		jsonMutation := jsonAtrMutation{
			BucketName:     mutation.Agent.BucketName(),
			ScopeName:      mutation.ScopeName,
			CollectionName: mutation.CollectionName,
			DocID:          string(mutation.Key),
		}

		if mutation.OpType == TransactionStagedMutationInsert {
			insMutations = append(insMutations, jsonMutation)
		} else if mutation.OpType == TransactionStagedMutationReplace {
			repMutations = append(repMutations, jsonMutation)
		} else if mutation.OpType == TransactionStagedMutationRemove {
			remMutations = append(remMutations, jsonMutation)
		} else {
			ecCb(classifyError(wrapError(ErrIllegalState, "unexpected staged mutation type")))
			return
		}
	}

	t.checkExpiredAtomic(hookATRAbort, []byte{}, true, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeATRAborted(func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			var marshalErr error
			atrFieldOp := func(fieldName string, data interface{}, flags memd.SubdocFlag) SubDocOp {
				bytes, err := json.Marshal(data)
				if err != nil {
					marshalErr = err
				}

				return SubDocOp{
					Op:    memd.SubDocOpDictSet,
					Flags: flags,
					Path:  "attempts." + t.id + "." + fieldName,
					Value: bytes,
				}
			}

			atrOps := []SubDocOp{
				atrFieldOp("st", jsonAtrStateAborted, memd.SubdocFlagXattrPath),
				atrFieldOp("tsrs", "${Mutation.CAS}", memd.SubdocFlagXattrPath|memd.SubdocFlagExpandMacros),
				atrFieldOp("ins", insMutations, memd.SubdocFlagXattrPath),
				atrFieldOp("rep", repMutations, memd.SubdocFlagXattrPath),
				atrFieldOp("rem", remMutations, memd.SubdocFlagXattrPath),
			}
			if marshalErr != nil {
				ecCb(classifyError(marshalErr))
				return
			}

			_, err = atrAgent.MutateIn(MutateInOptions{
				ScopeName:              atrScopeName,
				CollectionName:         atrCollectionName,
				Key:                    atrKey,
				Ops:                    atrOps,
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Flags:                  memd.SubdocDocFlagNone,
				Deadline:               deadline,
				User:                   atrOboUser,
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

				t.hooks.AfterATRAborted(func(err error) {
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

func (t *transactionAttempt) setATRRolledBackLocked(
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
					wrapError(ErrAttemptExpired, "rolled back atr removal failed during overtime")),
				ShouldNotRetry:    true,
				ShouldNotRollback: true,
			}))
			return
		}

		switch cerr.Class {
		case TransactionErrorClassFailDocNotFound:
			fallthrough
		case TransactionErrorClassFailPathNotFound:
			cb(nil)
			return
		case TransactionErrorClassFailExpiry:
			cb(t.operationFailed(operationFailedDef{
				Cerr: classifyError(
					wrapError(ErrAttemptExpired, "rolled back atr removal operation expired")),
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
				t.setATRRolledBackLocked(cb)
			})
		}
	}

	atrAgent := t.atrAgent
	atrOboUser := t.atrOboUser
	atrScopeName := t.atrScopeName
	atrKey := t.atrKey
	atrCollectionName := t.atrCollectionName

	t.checkExpiredAtomic(hookATRRollback, []byte{}, true, func(cerr *classifiedError) {
		if cerr != nil {
			ecCb(cerr)
			return
		}

		t.hooks.BeforeATRRolledBack(func(err error) {
			if err != nil {
				ecCb(classifyHookError(err))
				return
			}

			deadline, duraTimeout := transactionsMutationTimeouts(t.keyValueTimeout, t.durabilityLevel)

			atrOps := []SubDocOp{
				{
					Op:    memd.SubDocOpDelete,
					Flags: memd.SubdocFlagXattrPath,
					Path:  "attempts." + t.id,
				},
			}

			_, err = atrAgent.MutateIn(MutateInOptions{
				ScopeName:              atrScopeName,
				CollectionName:         atrCollectionName,
				Key:                    atrKey,
				Ops:                    atrOps,
				DurabilityLevel:        transactionsDurabilityLevelToMemd(t.durabilityLevel),
				DurabilityLevelTimeout: duraTimeout,
				Deadline:               deadline,
				Flags:                  memd.SubdocDocFlagNone,
				User:                   atrOboUser,
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

				t.hooks.AfterATRRolledBack(func(err error) {
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
