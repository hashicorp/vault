package gocb

import (
	"context"
	gocbcore "github.com/couchbase/gocbcore/v10"
)

type asyncOpManager struct {
	signal chan struct{}

	cancelCh    chan struct{}
	wasResolved bool
	ctx         context.Context
}

func (m *asyncOpManager) SetCancelCh(cancelCh chan struct{}) {
	m.cancelCh = cancelCh
}

func (m *asyncOpManager) Reject() {
	m.signal <- struct{}{}
}

func (m *asyncOpManager) Resolve() {
	m.wasResolved = true
	m.signal <- struct{}{}
}

func (m *asyncOpManager) Wait(op gocbcore.PendingOp, err error) error {
	if err != nil {
		return err
	}

	select {
	case <-m.signal:
		// Good to go
	case <-m.ctx.Done():
		op.Cancel()
		<-m.signal
	case <-m.cancelCh:
		op.Cancel()
		<-m.signal
	}

	return nil
}

func newAsyncOpManager(ctx context.Context) *asyncOpManager {
	if ctx == nil {
		ctx = context.Background()
	}
	return &asyncOpManager{
		signal: make(chan struct{}, 1),
		ctx:    ctx,
	}
}
