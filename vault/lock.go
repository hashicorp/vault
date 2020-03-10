// +build !deadlock

package vault

import (
	"sync"
)

type DeadlockMutexW struct {
	sync.Mutex
}
type DeadlockMutexRW struct {
	sync.RWMutex
}
