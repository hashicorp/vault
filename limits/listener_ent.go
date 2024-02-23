// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1
package limits

import (
	"sync/atomic"

	"github.com/armon/go-metrics"
	"github.com/platinummonkey/go-concurrency-limits/limiter"
)

// RequestListener is a thin wrapper for limiter.DefaultLimiter to handle the
// case where request limiting is turned off.
type RequestListener struct {
	*limiter.DefaultListener
	released *atomic.Bool
}

// OnSuccess is called as a notification that the operation succeeded and
// internally measured latency should be used as an RTT sample.
func (l *RequestListener) OnSuccess() {
	if l.DefaultListener != nil {
		metrics.IncrCounter(([]string{"limits", "concurrency", "success"}), 1)
		l.DefaultListener.OnSuccess()
		l.released.Store(true)
	}
}

// OnDropped is called to indicate the request failed and was dropped due to an
// internal server error. Note that this does not include ErrCapacity.
func (l *RequestListener) OnDropped() {
	if l.DefaultListener != nil {
		metrics.IncrCounter(([]string{"limits", "concurrency", "dropped"}), 1)
		l.DefaultListener.OnDropped()
		l.released.Store(true)
	}
}

// OnIgnore is called to indicate the operation failed before any meaningful RTT
// measurement could be made and should be ignored to not introduce an
// artificially low RTT. It also provides an extra layer of protection against
// leaks of the underlying StrategyToken during recoverable panics in the
// request handler. We treat these as Ignored, discard the measurement, and mark
// the listener as released.
func (l *RequestListener) OnIgnore() {
	if l.DefaultListener != nil && l.released.Load() != true {
		metrics.IncrCounter(([]string{"limits", "concurrency", "ignored"}), 1)
		l.DefaultListener.OnIgnore()
		l.released.Store(true)
	}
}
