// +build !deadlock

package vault

import (
	"sync"
)

// DeadlockMutexW is just a sync.Mutex when the build tag `deadlock` is absent.
// See its other definition in the corresponding deadlock-build-tag-constrained
// file for more details.
type DeadlockMutexW struct {
	sync.Mutex
}

// DeadlockMutexRW is the RW version of DeadlockMutexW.
type DeadlockMutexRW struct {
	sync.RWMutex
}
