package clock

import (
	"context"
	"sync"
	"time"
)

type mockTimers interface {
	start(t *mockTimer)
	stop(t *mockTimer)
	reset(t *mockTimer)
	next() *mockTimer
	len() int
}

// Mock implements a Clock that only moves with Add, AddNext and Set.
//
// The clock can be suspended with Lock and resumed with Unlock.
// While suspended, all attempts to use the API will block.
//
// To increase predictability, all Mock methods acquire
// and release the Mutex only once during their execution.
type Mock struct {
	sync.Mutex
	now time.Time
	mockTimers
}

// NewMock returns a new Mock with current time set to now.
//
// Use Realtime to get the real-time Clock.
func NewMock(now time.Time) *Mock {
	return &Mock{
		now:        now,
		mockTimers: &timerHeap{},
	}
}

// Add advances the current time by duration d and fires all expired timers.
//
// Returns the new current time.
// To increase predictability and speed, Tickers are ticked only once per call.
func (m *Mock) Add(d time.Duration) time.Time {
	m.Lock()
	defer m.Unlock()
	now, _ := m.set(m.now.Add(d))
	return now
}

// AddNext advances the current time to the next available timer deadline
// and fires all expired timers.
//
// Returns the new current time and the advanced duration.
func (m *Mock) AddNext() (time.Time, time.Duration) {
	m.Lock()
	defer m.Unlock()
	t := m.next()
	if t == nil {
		return m.now, 0
	}
	return m.set(t.deadline)
}

// Set advances the current time to t and fires all expired timers.
//
// Returns the advanced duration.
// To increase predictability and speed, Tickers are ticked only once per call.
func (m *Mock) Set(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	_, d := m.set(t)
	return d
}

func (m *Mock) set(now time.Time) (time.Time, time.Duration) {
	cur := m.now
	for {
		t := m.next()
		if t == nil || t.deadline.After(now) {
			m.now = now
			return m.now, m.now.Sub(cur)
		}
		m.now = t.deadline
		if d := t.fire(); d == 0 {
			// Timers are always stopped.
			m.stop(t)
		} else {
			// Ticker's next deadline is set to the first tick after the new now.
			dd := (now.Sub(m.now)/d + 1) * d
			t.deadline = m.now.Add(dd)
			m.reset(t)
		}
	}
}

// Len returns the number of active timers.
func (m *Mock) Len() int {
	m.Lock()
	defer m.Unlock()
	return m.len()
}

// Now returns the current mocked time.
func (m *Mock) Now() time.Time {
	m.Lock()
	defer m.Unlock()
	return m.now
}

// Since returns the time elapsed since t.
func (m *Mock) Since(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return m.now.Sub(t)
}

// Until returns the duration until t.
func (m *Mock) Until(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return t.Sub(m.now)
}

// DeadlineContext implements Clock.
func (m *Mock) DeadlineContext(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.deadlineContext(parent, d)
}

// TimeoutContext implements Clock.
func (m *Mock) TimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.deadlineContext(parent, m.now.Add(timeout))
}

func (m *Mock) deadlineContext(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	cancelCtx, cancel := context.WithCancel(Context(parent, m))
	if pd, ok := parent.Deadline(); ok && !pd.After(deadline) {
		return cancelCtx, cancel
	}
	ctx := &mockCtx{
		Context:  cancelCtx,
		done:     make(chan struct{}),
		deadline: deadline,
	}
	t := m.newTimerFunc(deadline, nil)
	go func() {
		select {
		case <-t.C:
			ctx.err = context.DeadlineExceeded
		case <-cancelCtx.Done():
			ctx.err = cancelCtx.Err()
			defer t.Stop()
		}
		close(ctx.done)
	}()
	return ctx, cancel
}

type mockCtx struct {
	context.Context
	deadline time.Time
	done     chan struct{}
	err      error
}

func (ctx *mockCtx) Deadline() (time.Time, bool) {
	return ctx.deadline, true
}

func (ctx *mockCtx) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *mockCtx) Err() error {
	select {
	case <-ctx.done:
		return ctx.err
	default:
		return nil
	}
}
