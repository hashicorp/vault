package backoff

import (
	"math"
	"time"
)

/*
ExponentialBackoff implements the Backoff interface. It represents an
instance that keeps track of retries, delays, and intervals for the
fibonacci backoff algorithm. This struct is instantiated by
the Exponential() function.
*/
type ExponentialBackoff struct {
	Retries    int
	MaxRetries int
	Delay      time.Duration
	Interval   time.Duration // time.Second, time.Millisecond, etc.
}

// Exponential creates a new instance of ExponentialBackoff.
func Exponential() *ExponentialBackoff {
	return &ExponentialBackoff{
		Retries:    0,
		MaxRetries: 5,
		Delay:      time.Duration(0),
		Interval:   time.Duration(1 * time.Second),
	}
}

/*
Next gets the next backoff delay. This method will increment the retries and check
if the maximum number of retries has been met. If this condition is satisfied, then
the function will return. Otherwise, the next backoff delay will be computed.

The exponential backoff delay is computed as follows:
`n = 2^c - 1` where `n` is the backoff delay and `c` is the number of retries.

Example, given a 1 second interval:

  Retry #        Backoff delay (in seconds)
    0                   0
    1                   1
    2                   3
    3                   7
    4                   15
    5                   31
*/
func (e *ExponentialBackoff) Next() bool {
	if e.Retries >= e.MaxRetries {
		return false
	}

	e.Retries++

	e.Delay = time.Duration(math.Pow(2, float64(e.Retries))-1) * e.Interval

	return true
}

/*
Retry will retry a function until the maximum number of retries is met. This method expects
the function `f` to return an error. If the failure condition is met, this method
will surface the error outputted from `f`, otherwise nil will be returned as normal.
*/
func (e *ExponentialBackoff) Retry(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for e.Next() {
		if err = f(); err == nil {
			return nil
		}

		time.Sleep(e.Delay)
	}

	return err
}

// Reset will reset the retry count and the backoff delay back to its initial state.
func (e *ExponentialBackoff) Reset() {
	e.Retries = 0
	e.Delay = time.Duration(0 * time.Second)
}
