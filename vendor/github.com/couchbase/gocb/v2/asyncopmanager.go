package gocb

import (
	gocbcore "github.com/couchbase/gocbcore/v9"
)

type asyncOpManager struct {
	signal chan struct{}

	wasResolved bool
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

	<-m.signal

	return nil
}

func newAsyncOpManager() *asyncOpManager {
	return &asyncOpManager{
		signal: make(chan struct{}, 1),
	}
}
