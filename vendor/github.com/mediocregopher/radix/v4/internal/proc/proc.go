// Package proc implements a simple framework for initializing and cleanly
// shutting down components.
package proc

import (
	"context"
	"errors"
	"sync"
)

// ErrClosed indicates that an operation could not be completed because the
// component has been closed.
var ErrClosed = errors.New("closed")

// Proc implements a lightweight pattern for setting up and tearing down
// components cleanly and consistently.
type Proc struct {
	ctx         context.Context
	ctxCancelFn context.CancelFunc
	ctxDoneCh   <-chan struct{}

	closeOnce sync.Once
	closed    bool
	wg        sync.WaitGroup

	lock sync.RWMutex
}

// New initializes and returns a clean Proc.
func New() *Proc {
	ctx, cancel := context.WithCancel(context.Background())
	return &Proc{
		ctx:         ctx,
		ctxCancelFn: cancel,
		ctxDoneCh:   ctx.Done(),
	}
}

// Run spawns a new go-routine which will run with the given callback. The
// callback's context will be closed when Close is called on Proc, and the
// go-routine must return for Close to return.
func (p *Proc) Run(fn func(ctx context.Context)) {
	p.wg.Add(1)
	go func() {
		fn(p.ctx)
		p.wg.Done()
	}()
}

// Close marks this Proc as having been closed (all methods will return
// ErrClosed after this), waits for all go-routines spawned with Run to return,
// and then calls the given callback. If Close is called multiple times it will
// return ErrClosed the subsequent times without taking any other action.
func (p *Proc) Close(fn func() error) error {
	return p.PrefixedClose(func() error { return nil }, fn)
}

// PrefixedClose is like Close but it will additionally perform a callback prior
// to marking the Proc as closed.
func (p *Proc) PrefixedClose(prefixFn, fn func() error) error {
	err := ErrClosed
	p.closeOnce.Do(func() {
		err = prefixFn()
		p.ctxCancelFn()

		p.lock.Lock()
		p.closed = true
		p.lock.Unlock()

		p.wg.Wait()
		if fn != nil {
			if fnErr := fn(); err == nil {
				err = fnErr
			}
		}
	})
	return err
}

// ClosedCh returns a channel which will be closed when the Proc is closed.
func (p *Proc) ClosedCh() <-chan struct{} {
	return p.ctxDoneCh
}

// IsClosed returns true if Close has been called.
func (p *Proc) IsClosed() bool {
	select {
	case <-p.ctxDoneCh:
		return true
	default:
		return false
	}
}

// WithRLock performs the given callback while holding a read lock on an
// internal RWMutex. This will return ErrClosed without calling the callback if
// Close has already been called.
func (p *Proc) WithRLock(fn func() error) error {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.closed {
		return ErrClosed
	}
	return fn()
}

// WithLock performs the given callback while holding a write lock on an
// internal RWMutex. This will return ErrClosed without calling the callback if
// Close has already been called.
func (p *Proc) WithLock(fn func() error) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed {
		return ErrClosed
	}
	return fn()
}
