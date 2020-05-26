package gocbcore

import (
	"sync"
	"time"
)

var globalTimerPool sync.Pool

// AcquireTimer acquires a time from a global pool of timers maintained by the library.
func AcquireTimer(d time.Duration) *time.Timer {
	tmr, isTmr := globalTimerPool.Get().(*time.Timer)
	if tmr == nil || !isTmr {
		if !isTmr && tmr != nil {
			logErrorf("Encountered non-timer in timer pool")
		}

		return time.NewTimer(d)
	}
	tmr.Reset(d)
	return tmr
}

// ReleaseTimer returns a timer to the global pool of timers maintained by the library.
func ReleaseTimer(t *time.Timer, wasRead bool) {
	stopped := t.Stop()
	if !wasRead && !stopped {
		<-t.C
	}
	globalTimerPool.Put(t)
}
