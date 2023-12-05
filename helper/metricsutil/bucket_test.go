// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package metricsutil

import (
	"testing"
	"time"
)

func TestTTLBucket_Lookup(t *testing.T) {
	testCases := []struct {
		Input    time.Duration
		Expected string
	}{
		{30 * time.Second, "1m"},
		{0 * time.Second, "1m"},
		{2 * time.Hour, "2h"},
		{2*time.Hour - time.Second, "2h"},
		{2*time.Hour + time.Second, "1d"},
		{30 * 24 * time.Hour, "30d"},
		{31 * 24 * time.Hour, "+Inf"},
	}

	for _, tc := range testCases {
		bucket := TTLBucket(tc.Input)
		if bucket != tc.Expected {
			t.Errorf("Expected %q, got %q for duration %v.", tc.Expected, bucket, tc.Input)
		}
	}
}
