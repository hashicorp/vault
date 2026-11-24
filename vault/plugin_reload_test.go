// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/locking"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/stretchr/testify/require"
)

// TestRetryPersistOnOverload_SucceedsImmediately verifies that when persist
// returns nil on the first call, retryPersistOnOverload does not retry.
func TestRetryPersistOnOverload_SucceedsImmediately(t *testing.T) {
	calls := 0
	retryPersistOnOverload(context.Background(), log.NewNullLogger(), func() error {
		calls++
		return nil
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

// TestRetryPersistOnOverload_RetriesOnOverloaded verifies that ErrOverloaded
// causes retries and that a subsequent nil return is treated as success.
func TestRetryPersistOnOverload_RetriesOnOverloaded(t *testing.T) {
	// Use a tiny backoff so the test doesn't sleep for real durations.
	origMin := reloadPersistRetryMin
	origMax := reloadPersistRetryMax
	reloadPersistRetryMin = time.Millisecond
	reloadPersistRetryMax = 5 * time.Millisecond
	t.Cleanup(func() {
		reloadPersistRetryMin = origMin
		reloadPersistRetryMax = origMax
	})

	const failTimes = 3
	calls := 0
	retryPersistOnOverload(context.Background(), log.NewNullLogger(), func() error {
		calls++
		if calls <= failTimes {
			return consts.ErrOverloaded
		}
		return nil
	})
	if calls != failTimes+1 {
		t.Fatalf("expected %d calls, got %d", failTimes+1, calls)
	}
}

// TestRetryPersistOnOverload_MaxRetriesExhausted verifies that after
// reloadPersistMaxRetries+1 attempts all returning ErrOverloaded, the function
// returns without panicking (error is logged and dropped).
func TestRetryPersistOnOverload_MaxRetriesExhausted(t *testing.T) {
	origMin := reloadPersistRetryMin
	origMax := reloadPersistRetryMax
	reloadPersistRetryMin = time.Millisecond
	reloadPersistRetryMax = 5 * time.Millisecond
	t.Cleanup(func() {
		reloadPersistRetryMin = origMin
		reloadPersistRetryMax = origMax
	})

	calls := 0
	retryPersistOnOverload(context.Background(), log.NewNullLogger(), func() error {
		calls++
		return consts.ErrOverloaded
	})
	// NewBackoff with maxRetries=N calls Next() N times before returning
	// ErrMaxRetry, so persist is called N+1 times total.
	expected := reloadPersistMaxRetries + 1
	if calls != expected {
		t.Fatalf("expected %d calls, got %d", expected, calls)
	}
}

// TestRetryPersistOnOverload_NonOverloadErrorImmediate verifies that a
// non-ErrOverloaded error is not retried — the function returns after exactly
// one invocation.
func TestRetryPersistOnOverload_NonOverloadErrorImmediate(t *testing.T) {
	calls := 0
	sentinel := errors.New("leader stepdown")
	retryPersistOnOverload(context.Background(), log.NewNullLogger(), func() error {
		calls++
		return sentinel
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

// TestRetryPersistOnOverload_ContextCancelled verifies that a cancelled context
// aborts the retry loop promptly after the first ErrOverloaded.
func TestRetryPersistOnOverload_ContextCancelled(t *testing.T) {
	origMin := reloadPersistRetryMin
	origMax := reloadPersistRetryMax
	reloadPersistRetryMin = 50 * time.Millisecond
	reloadPersistRetryMax = 500 * time.Millisecond
	t.Cleanup(func() {
		reloadPersistRetryMin = origMin
		reloadPersistRetryMax = origMax
	})

	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	// Cancel the context immediately after the first overload is returned so
	// that the select in retryPersistOnOverload picks ctx.Done().
	retryPersistOnOverload(ctx, log.NewNullLogger(), func() error {
		calls++
		cancel() // cancel before sleeping
		return consts.ErrOverloaded
	})

	if calls != 1 {
		t.Fatalf("expected 1 call before context cancel, got %d", calls)
	}
}

// TestDispatchReloads_AllComplete verifies that all entries are processed and
// results are populated when there are fewer entries than the semaphore limit.
func TestDispatchReloads_AllComplete(t *testing.T) {
	t.Parallel()
	entries := []int{1, 2, 3}
	results, skippedFrom := dispatchReloads(context.Background(), entries,
		func(_ context.Context, n int) int { return n * 2 },
	)
	if skippedFrom != len(entries) {
		t.Fatalf("expected skippedFrom=%d, got %d", len(entries), skippedFrom)
	}
	for i, want := range []int{2, 4, 6} {
		if results[i] != want {
			t.Errorf("results[%d]: want %d, got %d", i, want, results[i])
		}
	}
}

// TestDispatchReloads_ContextCancelledWhileWaitingForSlot is the core
// regression test. It fills the semaphore with concurrentReloadWorkers parked
// goroutines, then cancels ctx while the coordinator is waiting for a free
// slot, and asserts the call returns promptly rather than blocking.
func TestDispatchReloads_ContextCancelledWhileWaitingForSlot(t *testing.T) {
	t.Parallel()

	// gate keeps every dispatched goroutine parked until we release it.
	gate := make(chan struct{})
	// ready is closed once all concurrentReloadWorkers slots are confirmed occupied.
	ready := make(chan struct{})
	var readyOnce sync.Once
	var occupiedCount sync.WaitGroup
	occupiedCount.Add(concurrentReloadWorkers)

	// concurrentReloadWorkers+1 entries: the first N fill the semaphore, the
	// (N+1)-th entry is the one that must not block the coordinator.
	n := concurrentReloadWorkers + 1
	entries := make([]int, n)

	ctx, cancel := context.WithCancel(context.Background())

	type result struct {
		val int
		err error
	}

	// watcher closes ready once all N slots are occupied.
	go func() {
		occupiedCount.Wait()
		readyOnce.Do(func() { close(ready) })
	}()

	work := func(ctx context.Context, v int) result {
		// Signal that this slot is now occupied, then park.
		occupiedCount.Done()
		select {
		case <-gate:
			return result{val: v}
		case <-ctx.Done():
			return result{err: ctx.Err()}
		}
	}

	done := make(chan struct{})
	var results []result
	var skippedFrom int
	go func() {
		defer close(done)
		results, skippedFrom = dispatchReloads(ctx, entries, work)
	}()

	// Wait until all slots are occupied, then cancel.
	select {
	case <-ready:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for workers to fill semaphore")
	}
	cancel()

	// dispatchReloads must return promptly — before the gate is ever opened.
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("dispatchReloads did not return after context cancellation — coordinator goroutine is blocked")
	}

	// Unblock parked goroutines so they can exit cleanly.
	close(gate)

	// At least the last entry was skipped.
	if skippedFrom == len(entries) {
		t.Errorf("expected at least one skipped entry, got skippedFrom=%d (==len)", skippedFrom)
	}
	// Skipped slots hold zero values (nil err field).
	for i := skippedFrom; i < len(results); i++ {
		if results[i].err != nil {
			t.Errorf("results[%d].err = %v; want nil (zero value — skipped slot)", i, results[i].err)
		}
	}
}

// TestDispatchReloads_EmptyEntries verifies that an empty slice is a no-op.
func TestDispatchReloads_EmptyEntries(t *testing.T) {
	t.Parallel()
	results, skippedFrom := dispatchReloads(context.Background(), []int{},
		func(_ context.Context, n int) int { return n },
	)
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %v", results)
	}
	if skippedFrom != 0 {
		t.Fatalf("expected skippedFrom=0, got %d", skippedFrom)
	}
}

// TestDispatchReloads_AlreadyCancelledContext verifies that if ctx is already
// cancelled before dispatch begins, the coordinator stops dispatching as soon
// as the select picks ctx.Done(). Because Go's select is non-deterministic when
// multiple cases are ready, any number of these entries may be dispatched.
func TestDispatchReloads_AlreadyCancelledContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before any dispatch

	entries := []int{1, 2, 3}
	results, skippedFrom := dispatchReloads(ctx, entries,
		func(_ context.Context, n int) int { return n * 2 },
	)
	if skippedFrom > len(entries) {
		t.Errorf("skippedFrom = %d, exceeds %d entries", skippedFrom, len(entries))
	}
	for i := 0; i < skippedFrom; i++ {
		if results[i] != entries[i]*2 {
			t.Errorf("results[%d]: expected dispatched value %d, got %d", i, entries[i]*2, results[i])
		}
	}
	// Entries at and after skippedFrom must hold the zero value (not dispatched).
	for i := skippedFrom; i < len(results); i++ {
		if results[i] != 0 {
			t.Errorf("results[%d]: expected zero value (skipped), got %d", i, results[i])
		}
	}
}

// TestDispatchReloads_SkippedFromFillPattern verifies the caller convention:
// after dispatchReloads returns, results[skippedFrom:] are zero values, and the
// caller stamps ctx.Err() into them so they surface in the multierror.
func TestDispatchReloads_SkippedFromFillPattern(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	type res struct {
		val int
		err error
	}
	entries := []int{10, 20, 30}
	results, skippedFrom := dispatchReloads(ctx, entries,
		func(_ context.Context, n int) res { return res{val: n} },
	)
	// Simulate the caller's fill loop (as written at each call site).
	for i := skippedFrom; i < len(results); i++ {
		if results[i].err == nil {
			results[i].err = ctx.Err()
		}
	}
	for i := skippedFrom; i < len(results); i++ {
		if !errors.Is(results[i].err, context.Canceled) {
			t.Errorf("results[%d].err = %v; want context.Canceled", i, results[i].err)
		}
	}
}

// TestReloadMatchingPlugin_Cancelled preserves the mount path in errors for
// skipped secret and auth reloads when cancellation interrupts dispatch.
func TestReloadMatchingPlugin_Cancelled(t *testing.T) {
	t.Parallel()

	for _, pluginType := range []consts.PluginType{consts.PluginTypeSecrets, consts.PluginTypeCredential} {
		t.Run(pluginType.String(), func(t *testing.T) {
			t.Parallel()

			entries := make([]*MountEntry, concurrentReloadWorkers+1)
			for i := range entries {
				entries[i] = &MountEntry{
					Path:      fmt.Sprintf("mount-%d/", i),
					Type:      "mock-plugin",
					namespace: namespace.RootNamespace,
				}
			}
			core := &Core{
				router:     NewRouter(),
				mountsLock: locking.CreateConfigurableRWMutex(nil, "mountsLock"),
				authLock:   locking.CreateConfigurableRWMutex(nil, "authLock"),
			}
			if pluginType == consts.PluginTypeSecrets {
				core.mounts = &MountTable{Entries: entries}
			} else {
				core.auth = &MountTable{Entries: entries}
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := core.reloadMatchingPlugin(ctx, nil, pluginType, "mock-plugin")
			require.Error(t, err)
			require.ErrorIs(t, err, context.Canceled)
			for _, entry := range entries {
				require.ErrorContains(t, err, entry.Path)
			}
		})
	}
}
