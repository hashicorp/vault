// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package auth

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

// TestAgentBackoff verifies that auto-auth backoff grows within its jitter
// bounds, respects its maximum, and returns to the default minimum after reset.
func TestAgentBackoff(t *testing.T) {
	t.Parallel()

	max := 1024 * time.Second
	backoff := newAutoAuthBackoff(consts.DefaultMinBackoff, max, false)

	// Test initial value
	if backoff.backoff.Current() > consts.DefaultMinBackoff || backoff.backoff.Current() < consts.DefaultMinBackoff*3/4 {
		t.Fatalf("expected 1s initial backoff, got: %v", backoff.backoff.Current())
	}

	// Test that backoffSleep values are in expected range (75-100% of 2*previous)
	next, _ := backoff.backoff.Next()
	for i := 0; i < 9; i++ {
		old := next
		next, _ = backoff.backoff.Next()

		expMax := 2 * old
		expMin := 3 * expMax / 4

		if next < expMin || next > expMax {
			t.Fatalf("expected backoffSleep in range %v to %v, got: %v", expMin, expMax, backoff)
		}
	}

	// Test that backoffSleep is capped
	for i := 0; i < 100; i++ {
		_, _ = backoff.backoff.Next()
		if backoff.backoff.Current() > max {
			t.Fatalf("backoff exceeded max of 100s: %v", backoff)
		}
	}

	// Test reset
	backoff.backoff.Reset()
	if backoff.backoff.Current() > consts.DefaultMinBackoff || backoff.backoff.Current() < consts.DefaultMinBackoff*3/4 {
		t.Fatalf("expected 1s backoff after reset, got: %v", backoff.backoff.Current())
	}
}

// TestAgentMinBackoffCustom verifies that auto-auth honors custom minimum
// backoff values while defaulting non-positive values to one second.
func TestAgentMinBackoffCustom(t *testing.T) {
	t.Parallel()

	type test struct {
		minBackoff time.Duration
		want       time.Duration
	}

	tests := []test{
		{minBackoff: 0 * time.Second, want: 1 * time.Second},
		{minBackoff: 1 * time.Second, want: 1 * time.Second},
		{minBackoff: 5 * time.Second, want: 5 * time.Second},
		{minBackoff: 10 * time.Second, want: 10 * time.Second},
	}

	for _, test := range tests {
		max := 1024 * time.Second
		backoff := newAutoAuthBackoff(test.minBackoff, max, false)

		// Test initial value
		if backoff.backoff.Current() > test.want || backoff.backoff.Current() < test.want*3/4 {
			t.Fatalf("expected %d initial backoffSleep, got: %v", test.want, backoff.backoff.Current())
		}

		// Test that backoffSleep values are in expected range (75-100% of 2*previous)
		next, _ := backoff.backoff.Next()
		for i := 0; i < 5; i++ {
			old := next
			next, _ = backoff.backoff.Next()

			expMax := 2 * old
			expMin := 3 * expMax / 4

			if next < expMin || next > expMax {
				t.Fatalf("expected backoffSleep in range %v to %v, got: %v", expMin, expMax, backoff)
			}
		}

		// Test that backoffSleep is capped
		for i := 0; i < 100; i++ {
			next, _ = backoff.backoff.Next()
			if next > max {
				t.Fatalf("backoffSleep exceeded max of 100s: %v", backoff)
			}
		}

		// Test reset
		backoff.backoff.Reset()
		if backoff.backoff.Current() > test.want || backoff.backoff.Current() < test.want*3/4 {
			t.Fatalf("expected %d backoffSleep after reset, got: %v", test.want, backoff.backoff.Current())
		}
	}
}
