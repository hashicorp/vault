// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestGrabLockOrStopped is a non-deterministic test to detect deadlocks in the
// grabLockOrStopped function. This test starts a bunch of workers which
// continually lock/unlock and rlock/runlock the same RWMutex. Each worker also
// starts a goroutine which closes the stop channel 1/2 the time, which races
// with acquisition of the lock.
func TestGrabLockOrStop(t *testing.T) {
	// Stop the test early if we deadlock.
	const (
		workers      = 100
		testDuration = time.Second
		testTimeout  = 10 * testDuration
	)
	done := make(chan struct{})
	defer close(done)
	var lockCount int64
	go func() {
		select {
		case <-done:
		case <-time.After(testTimeout):
			panic(fmt.Sprintf("deadlock after %d lock count",
				atomic.LoadInt64(&lockCount)))
		}
	}()

	// lock is locked/unlocked and rlocked/runlocked concurrently.
	var lock sync.RWMutex
	start := time.Now()

	// workerWg is used to wait until all workers exit.
	var workerWg sync.WaitGroup
	workerWg.Add(workers)

	// Start a bunch of worker goroutines.
	for g := 0; g < workers; g++ {
		g := g
		go func() {
			defer workerWg.Done()
			for time.Now().Sub(start) < testDuration {
				stop := make(chan struct{})

				// closerWg waits until the closer goroutine exits before we do
				// another iteration. This makes sure goroutines don't pile up.
				var closerWg sync.WaitGroup
				closerWg.Add(1)
				go func() {
					defer closerWg.Done()
					// Close the stop channel half the time.
					if rand.Int()%2 == 0 {
						close(stop)
					}
				}()

				// Half the goroutines lock/unlock and the other half rlock/runlock.
				if g%2 == 0 {
					if !grabLockOrStop(lock.Lock, lock.Unlock, stop) {
						lock.Unlock()
					}
				} else {
					if !grabLockOrStop(lock.RLock, lock.RUnlock, stop) {
						lock.RUnlock()
					}
				}

				closerWg.Wait()

				// This lets us know how many lock/unlock and rlock/runlock have
				// happened if there's a deadlock.
				atomic.AddInt64(&lockCount, 1)
			}
		}()
	}
	workerWg.Wait()
}
