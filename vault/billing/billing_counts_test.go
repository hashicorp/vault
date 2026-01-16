// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestGetMonthlyBillingPath verifies the GetMonthlyBillingPath function
// returns the correct billing path for the given product area and month
func TestGetMonthlyBillingPath(t *testing.T) {
	ts := time.Date(2026, time.January, 5, 12, 0, 0, 0, time.UTC)

	got := GetMonthlyBillingPath(ReplicatedPrefix, ts, KvHWMCountsHWM)
	want := "replicated/2026/01/maxKvCounts/"
	require.Equal(t, got, want)
}
