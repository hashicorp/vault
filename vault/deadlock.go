// +build deadlock

package vault

import (
	"github.com/sasha-s/go-deadlock"
)

type DeadlockMutexW struct {
	deadlock.Mutex
}
type DeadlockMutexRW struct {
	deadlock.RWMutex
}
