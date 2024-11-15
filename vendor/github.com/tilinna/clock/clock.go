// Package clock implements a library for mocking time.
//
// Usage
//
// Include a Clock variable on your application and initialize it with
// a Realtime() by default. Then use the Clock for all time-related API
// calls. So instead of time.NewTimer(), say myClock.NewTimer().
//
// On a test setup, override or inject the variable with a Mock instance
// and use it to control how the time behaves during each test phase.
//
// To mock context.WithTimeout and context.WithDeadline, use the included
// Context, TimeoutContext and DeadlineContext methods.
//
// The Context method is also useful in cases where you need to pass a
// Clock via an 'func(ctx Context, ..)' API you can't change yourself.
// The FromContext method will then return the associated Clock instance.
// Alternatively, use the context'ed methods like Sleep(ctx) directly.
//
// All methods are safe for concurrent use.
package clock

import (
	"context"
	"time"
)

// Clock represents an interface to the functions in the standard time and context packages.
type Clock interface {
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	NewTicker(d time.Duration) *Ticker
	NewTimer(d time.Duration) *Timer
	Now() time.Time
	Since(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Until(t time.Time) time.Duration

	// DeadlineContext returns a copy of the parent context with the associated
	// Clock deadline adjusted to be no later than d.
	DeadlineContext(parent context.Context, d time.Time) (context.Context, context.CancelFunc)

	// TimeoutContext returns DeadlineContext(parent, Now(parent).Add(timeout)).
	TimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
}

type clock struct{}

var realtime = clock{}

// Realtime returns the standard real-time Clock.
func Realtime() Clock {
	return realtime
}

func (clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (clock) AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{timer: time.AfterFunc(d, f)}
}

func (clock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{
		C:      t.C,
		ticker: t,
	}
}

func (clock) NewTimer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{
		C:     t.C,
		timer: t,
	}
}

func (clock) Now() time.Time {
	return time.Now()
}

func (clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (clock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (clock) Tick(d time.Duration) <-chan time.Time {
	// Using time.Tick would trigger a vet tool warning.
	if d <= 0 {
		return nil
	}
	return time.NewTicker(d).C
}

func (clock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (clock) DeadlineContext(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}

func (clock) TimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
