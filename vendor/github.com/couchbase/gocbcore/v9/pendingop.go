package gocbcore

import "sync/atomic"

// PendingOp represents an outstanding operation within the client.
// This can be used to cancel an operation before it completes.
// This can also be used to Get information about the operation once
// it has completed (cancelled or successful).
type PendingOp interface {
	Cancel()
}

type multiPendingOp struct {
	ops          []PendingOp
	completedOps uint32
	isIdempotent bool
}

func (mp *multiPendingOp) Cancel() {
	for _, op := range mp.ops {
		op.Cancel()
	}
}

func (mp *multiPendingOp) CompletedOps() uint32 {
	return atomic.LoadUint32(&mp.completedOps)
}

func (mp *multiPendingOp) IncrementCompletedOps() uint32 {
	return atomic.AddUint32(&mp.completedOps, 1)
}
