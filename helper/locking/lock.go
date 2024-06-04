// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package locking

import (
	"sync"

	"github.com/sasha-s/go-deadlock"
)

// Common mutex interface to allow either built-in or imported deadlock use
type Mutex interface {
	Lock()
	Unlock()
}

// Common r/w mutex interface to allow either built-in or imported deadlock use
type RWMutex interface {
	Lock()
	RLock()
	RLocker() sync.Locker
	RUnlock()
	Unlock()
}

// DeadlockMutex (used when requested via config option `detact_deadlocks`),
// behaves like a sync.Mutex but does periodic checking to see if outstanding
// locks and requests look like a deadlock.  If it finds a deadlock candidate it
// will output it prefixed with "POTENTIAL DEADLOCK", as described at
// https://github.com/sasha-s/go-deadlock
type DeadlockMutex struct {
	deadlock.Mutex
}

// DeadlockRWMutex is the RW version of DeadlockMutex.
type DeadlockRWMutex struct {
	deadlock.RWMutex
}

// Regular sync/mutex.
type SyncMutex struct {
	sync.Mutex
}

// DeadlockRWMutex is the RW version of SyncMutex.
type SyncRWMutex struct {
	sync.RWMutex
}
