// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestShouldSkipUnifiedTransfer verifies the throttle predicate used by runUnifiedTransfer.
func TestShouldSkipUnifiedTransfer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		forceRerun bool
		lastRun    time.Time
		wantSkip   bool
	}{
		{
			name:       "first run is never throttled",
			forceRerun: false,
			lastRun:    time.Time{},
			wantSkip:   false,
		},
		{
			name:       "force flag bypasses throttle within min delay",
			forceRerun: true,
			lastRun:    time.Now(),
			wantSkip:   false,
		},
		{
			name:       "no force flag within min delay is throttled",
			forceRerun: false,
			lastRun:    time.Now(),
			wantSkip:   true,
		},
		{
			name:       "no force flag after min delay elapses is not throttled",
			forceRerun: false,
			lastRun:    time.Now().Add(-2 * minUnifiedTransferDelay),
			wantSkip:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := shouldSkipUnifiedTransfer(tc.forceRerun, tc.lastRun)
			require.Equal(t, tc.wantSkip, got,
				"throttle check mismatch for case %q (forceRerun=%v, lastRun=%v)",
				tc.name, tc.forceRerun, tc.lastRun)
		})
	}
}
