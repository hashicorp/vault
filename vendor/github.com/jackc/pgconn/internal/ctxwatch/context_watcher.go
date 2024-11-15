package ctxwatch

import (
	"context"
	"sync"
)

// ContextWatcher watches a context and performs an action when the context is canceled. It can watch one context at a
// time.
type ContextWatcher struct {
	onCancel             func()
	onUnwatchAfterCancel func()
	unwatchChan          chan struct{}

	lock              sync.Mutex
	watchInProgress   bool
	onCancelWasCalled bool
}

// NewContextWatcher returns a ContextWatcher. onCancel will be called when a watched context is canceled.
// OnUnwatchAfterCancel will be called when Unwatch is called and the watched context had already been canceled and
// onCancel called.
func NewContextWatcher(onCancel func(), onUnwatchAfterCancel func()) *ContextWatcher {
	cw := &ContextWatcher{
		onCancel:             onCancel,
		onUnwatchAfterCancel: onUnwatchAfterCancel,
		unwatchChan:          make(chan struct{}),
	}

	return cw
}

// Watch starts watching ctx. If ctx is canceled then the onCancel function passed to NewContextWatcher will be called.
func (cw *ContextWatcher) Watch(ctx context.Context) {
	cw.lock.Lock()
	defer cw.lock.Unlock()

	if cw.watchInProgress {
		panic("Watch already in progress")
	}

	cw.onCancelWasCalled = false

	if ctx.Done() != nil {
		cw.watchInProgress = true
		go func() {
			select {
			case <-ctx.Done():
				cw.onCancel()
				cw.onCancelWasCalled = true
				<-cw.unwatchChan
			case <-cw.unwatchChan:
			}
		}()
	} else {
		cw.watchInProgress = false
	}
}

// Unwatch stops watching the previously watched context. If the onCancel function passed to NewContextWatcher was
// called then onUnwatchAfterCancel will also be called.
func (cw *ContextWatcher) Unwatch() {
	cw.lock.Lock()
	defer cw.lock.Unlock()

	if cw.watchInProgress {
		cw.unwatchChan <- struct{}{}
		if cw.onCancelWasCalled {
			cw.onUnwatchAfterCancel()
		}
		cw.watchInProgress = false
	}
}
