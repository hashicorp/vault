/*
This code was taken from the same implementation in a branch from Consul and then
had the package updated and the mutex type unexported.
*/
package autopilot

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type mutex semaphore.Weighted

// New returns a Mutex that is ready for use.
func newMutex() *mutex {
	return (*mutex)(semaphore.NewWeighted(1))
}

func (m *mutex) Lock() {
	_ = (*semaphore.Weighted)(m).Acquire(context.Background(), 1)
}

func (m *mutex) Unlock() {
	(*semaphore.Weighted)(m).Release(1)
}

// TryLock acquires the mutex, blocking until resources are available or ctx is
// done. On success, returns nil. On failure, returns ctx.Err() and leaves the
// semaphore unchanged.
//
// If ctx is already done, Acquire may still succeed without blocking.
func (m *mutex) TryLock(ctx context.Context) error {
	return (*semaphore.Weighted)(m).Acquire(ctx, 1)
}
