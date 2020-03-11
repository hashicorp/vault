// +build !deadlock

package vault

import (
	"sync"
)

// DeadlockMutex is just a sync.Mutex when the build tag `deadlock` is absent.
// See its other definition in the corresponding deadlock-build-tag-constrained
// file for more details.
type DeadlockMutex struct {
	sync.Mutex
}

// DeadlockRWMutex is the RW version of DeadlockMutex.
type DeadlockRWMutex struct {
	sync.RWMutex
}
