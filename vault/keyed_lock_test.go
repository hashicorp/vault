// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestKeyedLock_SameKeyExcludes verifies that a second Lock on a held key
// blocks until the key is unlocked.
func TestKeyedLock_SameKeyExcludes(t *testing.T) {
	t.Parallel()

	l := newKeyedLock()
	l.Lock("a")

	acquired := make(chan struct{})
	go func() {
		l.Lock("a")
		close(acquired)
	}()

	select {
	case <-acquired:
		t.Fatal("second Lock acquired a held key")
	case <-time.After(100 * time.Millisecond):
	}

	l.Unlock("a")
	select {
	case <-acquired:
	case <-time.After(10 * time.Second):
		t.Fatal("second Lock did not acquire the key after Unlock")
	}
	l.Unlock("a")
}

// TestKeyedLock_DistinctKeysDoNotBlock verifies that holding one key does not
// block locking another.
func TestKeyedLock_DistinctKeysDoNotBlock(t *testing.T) {
	t.Parallel()

	l := newKeyedLock()
	l.Lock("a")
	defer l.Unlock("a")

	acquired := make(chan struct{})
	go func() {
		l.Lock("b")
		l.Unlock("b")
		close(acquired)
	}()

	select {
	case <-acquired:
	case <-time.After(10 * time.Second):
		t.Fatal("locking a distinct key blocked")
	}
}

// TestKeyedLock_Concurrent verifies mutual exclusion per key under contention
// and that every entry is released once all keys are unlocked.
func TestKeyedLock_Concurrent(t *testing.T) {
	t.Parallel()

	const (
		keys       = 8
		goroutines = 64
		iterations = 500
	)

	l := newKeyedLock()
	counters := make([]int, keys)
	var wg sync.WaitGroup
	for g := range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range iterations {
				k := (g + i) % keys
				key := strconv.Itoa(k)
				l.Lock(key)
				counters[k]++
				l.Unlock(key)
			}
		}()
	}
	wg.Wait()

	total := 0
	for _, c := range counters {
		total += c
	}
	require.Equal(t, goroutines*iterations, total)
	require.Empty(t, l.locks)
}

// TestKeyedLock_UnlockUnlockedPanics verifies that unlocking a key that is not
// locked panics, matching sync.Mutex.
func TestKeyedLock_UnlockUnlockedPanics(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() { newKeyedLock().Unlock("a") })
}
