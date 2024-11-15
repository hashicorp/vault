// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package csot

import (
	"context"
	"time"
)

type timeoutKey struct{}

// MakeTimeoutContext returns a new context with Client-Side Operation Timeout (CSOT) feature-gated behavior
// and a Timeout set to the passed in Duration. Setting a Timeout on a single operation is not supported in
// public API.
//
// TODO(GODRIVER-2348) We may be able to remove this function once CSOT feature-gated behavior becomes the
// TODO default behavior.
func MakeTimeoutContext(ctx context.Context, to time.Duration) (context.Context, context.CancelFunc) {
	// Only use the passed in Duration as a timeout on the Context if it
	// is non-zero and if the Context doesn't already have a timeout.
	cancelFunc := func() {}
	if _, deadlineSet := ctx.Deadline(); to != 0 && !deadlineSet {
		ctx, cancelFunc = context.WithTimeout(ctx, to)
	}

	// Add timeoutKey either way to indicate CSOT is enabled.
	return context.WithValue(ctx, timeoutKey{}, true), cancelFunc
}

func IsTimeoutContext(ctx context.Context) bool {
	return ctx.Value(timeoutKey{}) != nil
}

// ZeroRTTMonitor implements the RTTMonitor interface and is used internally for testing. It returns 0 for all
// RTT calculations and an empty string for RTT statistics.
type ZeroRTTMonitor struct{}

// EWMA implements the RTT monitor interface.
func (zrm *ZeroRTTMonitor) EWMA() time.Duration {
	return 0
}

// Min implements the RTT monitor interface.
func (zrm *ZeroRTTMonitor) Min() time.Duration {
	return 0
}

// P90 implements the RTT monitor interface.
func (zrm *ZeroRTTMonitor) P90() time.Duration {
	return 0
}

// Stats implements the RTT monitor interface.
func (zrm *ZeroRTTMonitor) Stats() string {
	return ""
}
