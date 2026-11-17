// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import "sync"

// keyedLock is a set of mutexes, one per key. A key's mutex is created on
// first use and deleted once nothing holds or waits on it, so memory is
// bounded by the number of keys in use. Unlike striped locks, distinct keys
// never contend.
type keyedLock struct {
	mu    sync.Mutex
	locks map[string]*keyedLockEntry
}

type keyedLockEntry struct {
	mu sync.Mutex

	// refs counts the goroutines holding or waiting on mu. It is only
	// accessed while holding keyedLock.mu.
	refs int
}

func newKeyedLock() *keyedLock {
	return &keyedLock{locks: make(map[string]*keyedLockEntry)}
}

// Lock locks the mutex for key, blocking until it is available.
func (l *keyedLock) Lock(key string) {
	l.mu.Lock()
	e, ok := l.locks[key]
	if !ok {
		e = &keyedLockEntry{}
		l.locks[key] = e
	}
	e.refs++
	l.mu.Unlock()

	e.mu.Lock()
}

// Unlock unlocks the mutex for key. It panics if key is not locked.
func (l *keyedLock) Unlock(key string) {
	l.mu.Lock()
	e, ok := l.locks[key]
	if !ok {
		l.mu.Unlock()
		panic("keyedLock: unlock of unlocked key")
	}
	e.refs--
	if e.refs == 0 {
		delete(l.locks, key)
	}
	l.mu.Unlock()

	e.mu.Unlock()
}
