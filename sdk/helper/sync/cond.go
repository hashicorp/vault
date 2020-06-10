package sync

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Conditional variable implementation that uses channels for notifications.
// Only supports .Broadcast() method, however supports timeout based Wait() calls
// unlike regular sync.TimeoutCond.
//
// Credit to Zviad Metreveli
type TimeoutCond struct {
	L sync.Locker
	n unsafe.Pointer
}

func NewTimeoutCond(l sync.Locker) *TimeoutCond {
	c := &TimeoutCond{L: l}
	n := make(chan struct{})
	c.n = unsafe.Pointer(&n)
	return c
}

// Waits for Broadcast calls. Similar to regular sync.TimeoutCond, this unlocks the underlying
// locker first, waits on changes and re-locks it before returning.
func (c *TimeoutCond) Wait() {
	n := c.NotifyChan()
	c.L.Unlock()
	<-n
	c.L.Lock()
}

// Same as Wait() call, but will only wait up to a given timeout.
func (c *TimeoutCond) WaitWithTimeout(t time.Duration) {
	n := c.NotifyChan()
	c.L.Unlock()
	select {
	case <-n:
	case <-time.After(t):
	}
	c.L.Lock()
}

// Returns a channel that can be used to wait for next Broadcast() call.
func (c *TimeoutCond) NotifyChan() <-chan struct{} {
	ptr := atomic.LoadPointer(&c.n)
	return *((*chan struct{})(ptr))
}

// Broadcast call notifies everyone that something has changed.
func (c *TimeoutCond) Broadcast() {
	n := make(chan struct{})
	ptrOld := atomic.SwapPointer(&c.n, unsafe.Pointer(&n))
	close(*(*chan struct{})(ptrOld))
}
