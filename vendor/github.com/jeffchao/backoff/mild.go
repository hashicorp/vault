package backoff

import (
	"time"
)

/*
MILDBackoff implements the Backoff interface. It represents an
instance that keeps track of retries, delays, and intervals for the
fibonacci backoff algorithm. This struct is instantiated by
the MILD() function.
*/
type MILDBackoff struct {
	Retries    int
	MaxRetries int
	Delay      time.Duration
	Interval   time.Duration // time.Second, time.Millisecond, etc.
	Slots      []time.Duration
}

// MILD creates a new instance of MILDBackoff.
func MILD() *MILDBackoff {
	return &MILDBackoff{
		Retries:    0,
		MaxRetries: 5,
		Delay:      time.Duration(0),
	}
}

/*
Next gets the next backoff delay. This method will increment the retries and check
if the maximum number of retries has been met. If this condition is satisfied, then
the function will return. Otherwise, the next backoff delay will be computed.

The MILD backoff delay is computed as follows:
`n = min(1.5 * n, len(slots)) upon failure; n = max(slots(c) - 1, 0) upon success;
n(0) = 0, n(1) = 1`
where `n` is the backoff delay, `c` is the retry slot, and `slots` is an array of retry delays.

This means a method must repeatedly succeed until `slots` is empty for the overall
backoff mechanism to terminate. Conversely, a repeated number of failures until the
maximum number of retries will result in a failure.

Example, given a 1 second interval, with max retries of 5:

  Retry #        Backoff delay (in seconds)       success/fail
    1                   1                             fail
    2                   1.5                           fail
    3                   1                             success
    4                   1.5                           fail
    5                   2.25                          fail

  Retry #        Backoff delay (in seconds)       success/fail
    1                   1                             fail
    2                   1.5                           fail
    3                   1                             success
    4                   0                             success
    5                   -                             success
*/
func (m *MILDBackoff) Next() bool {
	if m.Retries >= m.MaxRetries {
		return false
	}

	m.increment()

	return true
}

/*
Retry will retry a function until the maximum number of retries is met. This method expects
the function `f` to return an error. If the failure condition is met, this method
will surface the error outputted from `f`, otherwise nil will be returned as normal.
*/
func (m *MILDBackoff) Retry(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for m.Next() {
		if err := f(); err == nil {
			if len(m.Slots) == 0 {
				return nil
			}
			m.decrement()
		}

		time.Sleep(m.Delay)
	}

	return err
}

func (m *MILDBackoff) increment() {
	m.Retries++

	if m.Delay == 0 {
		m.Delay = time.Duration(1 * m.Interval)
	} else {
		m.Delay = m.Delay + (m.Delay / 2)
	}

	m.Slots = append(m.Slots, m.Delay)
}

func (m *MILDBackoff) decrement() {
	copy(m.Slots[len(m.Slots)-1:], m.Slots[len(m.Slots):])
	m.Slots[len(m.Slots)-1] = time.Duration(0 * m.Interval)
	m.Slots = m.Slots[:len(m.Slots)-1]
	m.Retries--
	if len(m.Slots) == 0 {
		m.Delay = time.Duration(0 * m.Interval)
	} else {
		m.Delay = m.Slots[len(m.Slots)-1]
	}
}

// Reset will reset the retry count, the backoff delay, and backoff slots back to its initial state.
func (m *MILDBackoff) Reset() {
	m.Retries = 0
	m.Delay = time.Duration(0 * m.Interval)
	m.Slots = nil
	m.Slots = make([]time.Duration, 0, m.MaxRetries)
}
