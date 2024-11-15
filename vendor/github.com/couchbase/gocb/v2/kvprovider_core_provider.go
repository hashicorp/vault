package gocb

import (
	"context"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

// kvProviderCoreProvider provides us with a way to unit test what the kvProviderCore layer is doing.
type kvProviderCoreProvider interface {
	Add(opts gocbcore.AddOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error)
	Set(opts gocbcore.SetOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error)
	Replace(opts gocbcore.ReplaceOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error)
	Get(opts gocbcore.GetOptions, cb gocbcore.GetCallback) (gocbcore.PendingOp, error)
	GetOneReplica(opts gocbcore.GetOneReplicaOptions, cb gocbcore.GetReplicaCallback) (gocbcore.PendingOp, error)
	Observe(opts gocbcore.ObserveOptions, cb gocbcore.ObserveCallback) (gocbcore.PendingOp, error)
	ObserveVb(opts gocbcore.ObserveVbOptions, cb gocbcore.ObserveVbCallback) (gocbcore.PendingOp, error)
	GetMeta(opts gocbcore.GetMetaOptions, cb gocbcore.GetMetaCallback) (gocbcore.PendingOp, error)
	Delete(opts gocbcore.DeleteOptions, cb gocbcore.DeleteCallback) (gocbcore.PendingOp, error)
	LookupIn(opts gocbcore.LookupInOptions, cb gocbcore.LookupInCallback) (gocbcore.PendingOp, error)
	MutateIn(opts gocbcore.MutateInOptions, cb gocbcore.MutateInCallback) (gocbcore.PendingOp, error)
	GetAndTouch(opts gocbcore.GetAndTouchOptions, cb gocbcore.GetAndTouchCallback) (gocbcore.PendingOp, error)
	GetAndLock(opts gocbcore.GetAndLockOptions, cb gocbcore.GetAndLockCallback) (gocbcore.PendingOp, error)
	Unlock(opts gocbcore.UnlockOptions, cb gocbcore.UnlockCallback) (gocbcore.PendingOp, error)
	Touch(opts gocbcore.TouchOptions, cb gocbcore.TouchCallback) (gocbcore.PendingOp, error)
	Increment(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error)
	Decrement(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error)
	Append(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error)
	Prepend(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error)
	WaitForConfigSnapshot(deadline time.Time, opts gocbcore.WaitForConfigSnapshotOptions, cb gocbcore.WaitForConfigSnapshotCallback) (gocbcore.PendingOp, error)
	RangeScanCreate(vbID uint16, opts gocbcore.RangeScanCreateOptions, cb gocbcore.RangeScanCreateCallback) (gocbcore.PendingOp, error)
	GetCollectionID(scopeName string, collectionName string, opts gocbcore.GetCollectionIDOptions, cb gocbcore.GetCollectionIDCallback) (gocbcore.PendingOp, error)
}

type kvProviderConfigSnapshotProvider interface {
	WaitForConfigSnapshot(ctx context.Context, deadline time.Time) (coreConfigSnapshot, error)
}

type coreConfigSnapshot interface {
	RevID() int64
	NumVbuckets() (int, error)
	NumReplicas() (int, error)
	NumServers() (int, error)
	VbucketsOnServer(index int) ([]uint16, error)
	KeyToServersByServerGroup(key []byte) (map[string][]int, error)
}

type stdCoreConfigSnapshotProvider struct {
	agent kvProviderCoreProvider
}

func (p *stdCoreConfigSnapshotProvider) WaitForConfigSnapshot(ctx context.Context, deadline time.Time) (coreConfigSnapshot, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var snapOut coreConfigSnapshot
	var errOut error
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(p.agent.WaitForConfigSnapshot(deadline, gocbcore.WaitForConfigSnapshotOptions{}, func(result *gocbcore.WaitForConfigSnapshotResult, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		snapOut = result.Snapshot
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return snapOut, errOut
}
