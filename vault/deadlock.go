// +build deadlock

package vault

import (
	"github.com/sasha-s/go-deadlock"
)

// DeadlockMutex, when the build tag `deadlock` is present, behaves like a
// sync.Mutex but does periodic checking to see if outstanding locks and requests
// look like a deadlock.  If it finds a deadlock candidate it will output it
// prefixed with "POTENTIAL DEADLOCK", as described at
// https://github.com/sasha-s/go-deadlock
type DeadlockMutex struct {
	deadlock.Mutex
}

// DeadlockRWMutex is the RW version of DeadlockMutex.
type DeadlockRWMutex struct {
	deadlock.RWMutex
}
