package gocb

import (
	"github.com/couchbase/gocbcore/v10"
)

// TransactionHooks provides a number of internal hooks used for testing.
// Internal: This should never be used and is not supported.
type TransactionHooks interface {
	BeforeATRCommit(ctx TransactionAttemptContext) error
	AfterATRCommit(ctx TransactionAttemptContext) error
	BeforeDocCommitted(ctx TransactionAttemptContext, docID string) error
	BeforeRemovingDocDuringStagedInsert(ctx TransactionAttemptContext, docID string) error
	BeforeRollbackDeleteInserted(ctx TransactionAttemptContext, docID string) error
	AfterDocCommittedBeforeSavingCAS(ctx TransactionAttemptContext, docID string) error
	AfterDocCommitted(ctx TransactionAttemptContext, docID string) error
	BeforeStagedInsert(ctx TransactionAttemptContext, docID string) error
	BeforeStagedRemove(ctx TransactionAttemptContext, docID string) error
	BeforeStagedReplace(ctx TransactionAttemptContext, docID string) error
	BeforeDocRemoved(ctx TransactionAttemptContext, docID string) error
	BeforeDocRolledBack(ctx TransactionAttemptContext, docID string) error
	AfterDocRemovedPreRetry(ctx TransactionAttemptContext, docID string) error
	AfterDocRemovedPostRetry(ctx TransactionAttemptContext, docID string) error
	AfterGetComplete(ctx TransactionAttemptContext, docID string) error
	AfterStagedReplaceComplete(ctx TransactionAttemptContext, docID string) error
	AfterStagedRemoveComplete(ctx TransactionAttemptContext, docID string) error
	AfterStagedInsertComplete(ctx TransactionAttemptContext, docID string) error
	AfterRollbackReplaceOrRemove(ctx TransactionAttemptContext, docID string) error
	AfterRollbackDeleteInserted(ctx TransactionAttemptContext, docID string) error
	BeforeCheckATREntryForBlockingDoc(ctx TransactionAttemptContext, docID string) error
	BeforeDocGet(ctx TransactionAttemptContext, docID string) error
	BeforeGetDocInExistsDuringStagedInsert(ctx TransactionAttemptContext, docID string) error
	BeforeRemoveStagedInsert(ctx TransactionAttemptContext, docID string) error
	AfterRemoveStagedInsert(ctx TransactionAttemptContext, docID string) error
	AfterDocsCommitted(ctx TransactionAttemptContext) error
	AfterDocsRemoved(ctx TransactionAttemptContext) error
	AfterATRPending(ctx TransactionAttemptContext) error
	BeforeATRPending(ctx TransactionAttemptContext) error
	BeforeATRComplete(ctx TransactionAttemptContext) error
	BeforeATRRolledBack(ctx TransactionAttemptContext) error
	AfterATRComplete(ctx TransactionAttemptContext) error
	BeforeATRAborted(ctx TransactionAttemptContext) error
	AfterATRAborted(ctx TransactionAttemptContext) error
	AfterATRRolledBack(ctx TransactionAttemptContext) error
	BeforeATRCommitAmbiguityResolution(ctx TransactionAttemptContext) error
	RandomATRIDForVbucket(ctx TransactionAttemptContext) (string, error)
	HasExpiredClientSideHook(ctx TransactionAttemptContext, stage string, vbID string) (bool, error)
	BeforeQuery(ctx TransactionAttemptContext, statement string) error
	AfterQuery(ctx TransactionAttemptContext, statement string) error
}

// TransactionCleanupHooks provides a number of internal hooks used for testing.
// Internal: This should never be used and is not supported.
type TransactionCleanupHooks interface {
	BeforeATRGet(id string) error
	BeforeDocGet(id string) error
	BeforeRemoveLinks(id string) error
	BeforeCommitDoc(id string) error
	BeforeRemoveDocStagedForRemoval(id string) error
	BeforeRemoveDoc(id string) error
	BeforeATRRemove(id string) error
}

// TransactionClientRecordHooks provides a number of internal hooks used for testing.
// Internal: This should never be used and is not supported.
type TransactionClientRecordHooks interface {
	BeforeCreateRecord() error
	BeforeRemoveClient() error
	BeforeUpdateCAS() error
	BeforeGetRecord() error
	BeforeUpdateRecord() error
}

type transactionHooksWrapper interface {
	SetAttemptContext(ctx TransactionAttemptContext)
	gocbcore.TransactionHooks
	Hooks() TransactionHooks
}

type transactionCleanupHooksWrapper interface {
	gocbcore.TransactionCleanUpHooks
}

type coreTxnsHooksWrapper struct {
	ctx   TransactionAttemptContext
	hooks TransactionHooks
}

type clientRecordHooksWrapper interface {
	gocbcore.TransactionClientRecordHooks
}

func (cthw *coreTxnsHooksWrapper) SetAttemptContext(ctx TransactionAttemptContext) {
	cthw.ctx = ctx
}

func (cthw *coreTxnsHooksWrapper) Hooks() TransactionHooks {
	return cthw.hooks
}

func (cthw *coreTxnsHooksWrapper) BeforeATRCommit(cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeATRCommit(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterATRCommit(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterATRCommit(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeDocCommitted(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeDocCommitted(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeRemovingDocDuringStagedInsert(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeRemovingDocDuringStagedInsert(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeRollbackDeleteInserted(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeRollbackDeleteInserted(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterDocCommittedBeforeSavingCAS(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterDocCommittedBeforeSavingCAS(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterDocCommitted(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterDocCommitted(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeStagedInsert(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeStagedInsert(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeStagedRemove(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeStagedRemove(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeStagedReplace(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeStagedReplace(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeDocRemoved(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeDocRemoved(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeDocRolledBack(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeDocRolledBack(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterDocRemovedPreRetry(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterDocRemovedPreRetry(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterDocRemovedPostRetry(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterDocRemovedPostRetry(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterGetComplete(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterGetComplete(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterStagedReplaceComplete(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterStagedReplaceComplete(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterStagedRemoveComplete(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterStagedRemoveComplete(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterStagedInsertComplete(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterStagedInsertComplete(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterRollbackReplaceOrRemove(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterRollbackReplaceOrRemove(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterRollbackDeleteInserted(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterRollbackDeleteInserted(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeCheckATREntryForBlockingDoc(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeCheckATREntryForBlockingDoc(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeDocGet(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeDocGet(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeGetDocInExistsDuringStagedInsert(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeGetDocInExistsDuringStagedInsert(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeRemoveStagedInsert(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeRemoveStagedInsert(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterRemoveStagedInsert(docID []byte, cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterRemoveStagedInsert(cthw.ctx, string(docID)))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterDocsCommitted(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterDocsCommitted(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterDocsRemoved(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterDocsRemoved(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterATRPending(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterATRPending(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeATRPending(cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeATRPending(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeATRComplete(cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeATRComplete(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeATRRolledBack(cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeATRRolledBack(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterATRComplete(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterATRComplete(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeATRAborted(cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeATRAborted(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterATRAborted(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterATRAborted(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) AfterATRRolledBack(cb func(err error)) {
	go func() {
		cb(cthw.hooks.AfterATRRolledBack(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) BeforeATRCommitAmbiguityResolution(cb func(err error)) {
	go func() {
		cb(cthw.hooks.BeforeATRCommitAmbiguityResolution(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) RandomATRIDForVbucket(cb func(string, error)) {
	go func() {
		cb(cthw.hooks.RandomATRIDForVbucket(cthw.ctx))
	}()
}

func (cthw *coreTxnsHooksWrapper) HasExpiredClientSideHook(stage string, vbID []byte, cb func(bool, error)) {
	go func() {
		cb(cthw.hooks.HasExpiredClientSideHook(cthw.ctx, stage, string(vbID)))
	}()
}

type coreTxnsCleanupHooksWrapper struct {
	CleanupHooks TransactionCleanupHooks
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeATRGet(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeATRGet(string(id)))
	}()
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeDocGet(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeDocGet(string(id)))
	}()
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeRemoveLinks(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeRemoveLinks(string(id)))
	}()
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeCommitDoc(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeCommitDoc(string(id)))
	}()
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeRemoveDocStagedForRemoval(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeRemoveDocStagedForRemoval(string(id)))
	}()
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeRemoveDoc(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeRemoveDoc(string(id)))
	}()
}

func (cthw *coreTxnsCleanupHooksWrapper) BeforeATRRemove(id []byte, cb func(error)) {
	go func() {
		cb(cthw.CleanupHooks.BeforeATRRemove(string(id)))
	}()
}

type coreTxnsClientRecordHooksWrapper struct {
	coreTxnsCleanupHooksWrapper
	ClientRecordHooks TransactionClientRecordHooks
}

func (hw *coreTxnsClientRecordHooksWrapper) BeforeCreateRecord(cb func(error)) {
	go func() {
		cb(hw.ClientRecordHooks.BeforeCreateRecord())
	}()
}

func (hw *coreTxnsClientRecordHooksWrapper) BeforeRemoveClient(cb func(error)) {
	go func() {
		cb(hw.ClientRecordHooks.BeforeRemoveClient())
	}()
}

func (hw *coreTxnsClientRecordHooksWrapper) BeforeUpdateCAS(cb func(error)) {
	go func() {
		cb(hw.ClientRecordHooks.BeforeUpdateCAS())
	}()
}

func (hw *coreTxnsClientRecordHooksWrapper) BeforeGetRecord(cb func(error)) {
	go func() {
		cb(hw.ClientRecordHooks.BeforeGetRecord())
	}()
}

func (hw *coreTxnsClientRecordHooksWrapper) BeforeUpdateRecord(cb func(error)) {
	go func() {
		cb(hw.ClientRecordHooks.BeforeUpdateRecord())
	}()
}

type noopHooksWrapper struct {
	gocbcore.TransactionDefaultHooks
	hooks transactionsDefaultHooks
}

func (nhw *noopHooksWrapper) SetAttemptContext(ctx TransactionAttemptContext) {
}

func (nhw *noopHooksWrapper) Hooks() TransactionHooks {
	return nhw.hooks
}

type noopCleanupHooksWrapper struct {
	gocbcore.TransactionDefaultCleanupHooks
}

type noopClientRecordHooksWrapper struct {
	gocbcore.TransactionDefaultCleanupHooks
	gocbcore.TransactionDefaultClientRecordHooks
}

type transactionsDefaultHooks struct {
}

func (d transactionsDefaultHooks) BeforeATRCommit(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) AfterATRCommit(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeDocCommitted(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeRemovingDocDuringStagedInsert(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeRollbackDeleteInserted(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterDocCommittedBeforeSavingCAS(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterDocCommitted(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeStagedInsert(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeStagedRemove(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeStagedReplace(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeDocRemoved(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeDocRolledBack(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterDocRemovedPreRetry(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterDocRemovedPostRetry(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterGetComplete(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterStagedReplaceComplete(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterStagedRemoveComplete(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterStagedInsertComplete(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterRollbackReplaceOrRemove(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterRollbackDeleteInserted(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeCheckATREntryForBlockingDoc(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeDocGet(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeGetDocInExistsDuringStagedInsert(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeRemoveStagedInsert(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterRemoveStagedInsert(ctx TransactionAttemptContext, docID string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterDocsCommitted(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) AfterDocsRemoved(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) AfterATRPending(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeATRPending(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeATRComplete(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeATRRolledBack(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) AfterATRComplete(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeATRAborted(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) AfterATRAborted(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) AfterATRRolledBack(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) BeforeATRCommitAmbiguityResolution(ctx TransactionAttemptContext) error {
	return nil
}

func (d transactionsDefaultHooks) RandomATRIDForVbucket(ctx TransactionAttemptContext) (string, error) {
	return "", nil
}

func (d transactionsDefaultHooks) HasExpiredClientSideHook(ctx TransactionAttemptContext, stage string, vbID string) (bool, error) {
	return false, nil
}

func (d transactionsDefaultHooks) BeforeQuery(ctx TransactionAttemptContext, statement string) error {
	return nil
}

func (d transactionsDefaultHooks) AfterQuery(ctx TransactionAttemptContext, statement string) error {
	return nil
}
