// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

var ErrMaxRetry = errors.New("exceeded maximum number of retries")

const maxJitter = 0.25

// Backoff is used to do capped exponential backoff with jitter, with a maximum number of retries.
// Generally, use this struct by calling Next() or NextSleep() after a failure.
// If configured for N max retries, Next() and NextSleep() will return an error on the call N+1.
// The jitter is set to 25%, so values returned will have up to 25% less than twice the previous value.
// The min value will also include jitter, so the first call will almost always be less than the requested minimum value.
// Backoff is not thread-safe.
type Backoff struct {
	currentAttempt int
	maxRetries     int
	min            time.Duration
	max            time.Duration
	current        time.Duration
}

// NewBackoff creates a new exponential backoff with the given number of maximum retries and min/max durations.
func NewBackoff(maxRetries int, min, max time.Duration) *Backoff {
	b := &Backoff{
		maxRetries: maxRetries,
		max:        max,
		min:        min,
	}
	b.Reset()
	return b
}

// Current returns the next time that will be returned by Next() (or slept in NextSleep()).
func (b *Backoff) Current() time.Duration {
	return b.current
}

// Next determines the next backoff duration that is roughly twice
// the current value, capped to a max value, with a measure of randomness.
// It returns an error if there are no more retries left.
func (b *Backoff) Next() (time.Duration, error) {
	if b.currentAttempt >= b.maxRetries {
		return time.Duration(-1), ErrMaxRetry
	}
	defer func() {
		b.currentAttempt += 1
	}()
	if b.currentAttempt == 0 {
		return b.current, nil
	}
	next := 2 * b.current
	if next > b.max {
		next = b.max
	}
	next = jitter(next)
	b.current = next
	return next, nil
}

// NextSleep will synchronously sleep the next backoff amount (see Next()).
// It returns an error if there are no more retries left.
func (b *Backoff) NextSleep() error {
	next, err := b.Next()
	if err != nil {
		return err
	}
	time.Sleep(next)
	return nil
}

// Reset resets the state to the initial backoff amount and 0 retries.
func (b *Backoff) Reset() {
	b.current = b.min
	b.current = jitter(b.current)
	b.currentAttempt = 0
}

func jitter(t time.Duration) time.Duration {
	f := float64(t) * (1.0 - maxJitter*rand.Float64())
	return time.Duration(math.Floor(f))
}

// Retry calls the given function until it does not return an error, at least once and up to max_retries + 1 times.
// If the number of retries is exceeded, Retry() will return the last error seen joined with ErrMaxRetry.
func (b *Backoff) Retry(f func() error) error {
	for {
		err := f()
		if err == nil {
			return nil
		}

		maxRetryErr := b.NextSleep()
		if maxRetryErr != nil {
			return errors.Join(maxRetryErr, err)
		}
	}
	return nil // unreachable
}
