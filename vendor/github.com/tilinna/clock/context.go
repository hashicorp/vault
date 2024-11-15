package clock

import (
	"context"
	"time"
)

type clockKey struct{}

// Context returns a copy of parent in which the Clock is associated with.
func Context(parent context.Context, c Clock) context.Context {
	return context.WithValue(parent, clockKey{}, c)
}

// FromContext returns the Clock associated with the context, or Realtime().
func FromContext(ctx context.Context) Clock {
	if c, ok := ctx.Value(clockKey{}).(Clock); ok {
		return c
	}
	return Realtime()
}

// After is a convenience wrapper for FromContext(ctx).After.
func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return FromContext(ctx).After(d)
}

// AfterFunc is a convenience wrapper for FromContext(ctx).AfterFunc.
func AfterFunc(ctx context.Context, d time.Duration, f func()) *Timer {
	return FromContext(ctx).AfterFunc(d, f)
}

// NewTicker is a convenience wrapper for FromContext(ctx).NewTicker.
func NewTicker(ctx context.Context, d time.Duration) *Ticker {
	return FromContext(ctx).NewTicker(d)
}

// NewTimer is a convenience wrapper for FromContext(ctx).NewTimer.
func NewTimer(ctx context.Context, d time.Duration) *Timer {
	return FromContext(ctx).NewTimer(d)
}

// Now is a convenience wrapper for FromContext(ctx).Now.
func Now(ctx context.Context) time.Time {
	return FromContext(ctx).Now()
}

// Since is a convenience wrapper for FromContext(ctx).Since.
func Since(ctx context.Context, t time.Time) time.Duration {
	return FromContext(ctx).Since(t)
}

// Sleep is a convenience wrapper for FromContext(ctx).Sleep.
func Sleep(ctx context.Context, d time.Duration) {
	FromContext(ctx).Sleep(d)
}

// Tick is a convenience wrapper for FromContext(ctx).Tick.
func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return FromContext(ctx).Tick(d)
}

// Until is a convenience wrapper for FromContext(ctx).Until.
func Until(ctx context.Context, t time.Time) time.Duration {
	return FromContext(ctx).Until(t)
}

// DeadlineContext is a convenience wrapper for FromContext(ctx).DeadlineContext.
func DeadlineContext(ctx context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return FromContext(ctx).DeadlineContext(ctx, d)
}

// TimeoutContext is a convenience wrapper for FromContext(ctx).TimeoutContext.
func TimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return FromContext(ctx).TimeoutContext(ctx, timeout)
}
