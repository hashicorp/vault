// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly && !enterprise

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/stretchr/testify/require"
)

// TestActivityLog_setupClientIDsUsageInfo_CE verifies that upon startup, the client IDs are not loaded in CE
func TestActivityLog_setupClientIDsUsageInfo_CE(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	core.setupClientIDsUsageInfo(context.Background())

	// wait for clientIDs to be loaded into memory
	verifyClientsLoadedInMemory := func() {
		corehelpers.RetryUntil(t, 60*time.Second, func() error {
			if a.GetClientIDsUsageInfoLoaded() {
				return fmt.Errorf("loaded clientIDs to memory")
			}
			return nil
		})
	}
	verifyClientsLoadedInMemory()

	require.Len(t, a.GetClientIDsUsageInfo(), 0)
}

// TestGetStartEndTime_EndTimeAdjustedToPastMonth tests getStartEndTime for proper adjustment of given end time to past month
func TestGetStartEndTime_EndTimeAdjustedToPastMonth(t *testing.T) {
	now := time.Now().UTC()
	currentMonthStart := timeutil.StartOfMonth(now)
	previousMonthEnd := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))

	// billing start time is zero for CE
	billingStartTime := time.Time{}
	sixMonthsAgo := now.AddDate(0, -6, 0)

	tests := []struct {
		name           string
		givenStartTime time.Time
		givenEndTime   time.Time
		expectedStart  time.Time
		expectedEnd    time.Time
		expectWarning  bool
		expectErr      bool
	}{
		{
			name:           "End time in the past is unchanged",
			givenStartTime: sixMonthsAgo,
			givenEndTime:   now.AddDate(0, -1, -1),
			expectedStart:  sixMonthsAgo,
			expectedEnd:    now.AddDate(0, -1, -1).UTC(),
			expectWarning:  false,
			expectErr:      false,
		},
		{
			name:           "End time in the current month is clamped to previous month",
			givenStartTime: sixMonthsAgo,
			givenEndTime:   currentMonthStart.AddDate(0, 0, 5).Add(2 * time.Hour),
			expectedStart:  sixMonthsAgo,
			expectedEnd:    previousMonthEnd,
			expectWarning:  true,
			expectErr:      false,
		},
		{
			name:           "End time in the future is clamped to previous month",
			givenStartTime: sixMonthsAgo,
			givenEndTime:   now.AddDate(0, 1, 0),
			expectedStart:  sixMonthsAgo,
			expectedEnd:    previousMonthEnd,
			expectWarning:  true,
			expectErr:      false,
		},
		{
			name:           "End time is zero and gets clamped to previous month",
			givenStartTime: sixMonthsAgo,
			givenEndTime:   time.Time{},
			expectedStart:  sixMonthsAgo,
			expectedEnd:    previousMonthEnd,
			expectWarning:  true,
			expectErr:      false,
		},
		{
			name:           "Start time after end time causes error",
			givenStartTime: now.AddDate(0, 2, 0),
			givenEndTime:   now.AddDate(0, 1, 0),
			expectWarning:  false,
			expectErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &framework.FieldData{
				Schema: map[string]*framework.FieldSchema{
					"start_time": {
						Type: framework.TypeTime,
					},
					"end_time": {
						Type: framework.TypeTime,
					},
				},
				Raw: map[string]any{
					"start_time": tc.givenStartTime.Format(time.RFC3339Nano),
					"end_time":   tc.givenEndTime.Format(time.RFC3339Nano),
				},
			}

			start, end, warnings, err := getStartEndTime(d, time.Now(), billingStartTime)

			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.WithinDuration(t, tc.expectedStart, start, time.Second, "Expected start time did not match")
			require.WithinDuration(t, tc.expectedEnd, end, time.Second, "Expected end time did not match")
			require.Equal(t, tc.expectWarning, warnings.EndTimeAdjustedToPastMonth)
		})
	}
}

// TestActivityLog_GetBillingPeriodActivityTelemetryMetric_CE verifies that vault.client.billing_period.activity is not emitted in ce
func TestActivityLog_GetBillingPeriodActivityTelemetryMetric_CE(t *testing.T) {
	inMemSink := metrics.NewInmemSink(1*time.Second, 10000*time.Hour)
	sink := metricsutil.NewClusterMetricSink("test", inMemSink)
	testClusterName := "test-cluster"
	config := &CoreConfig{
		MetricSink: sink,
		ActivityLogConfig: ActivityLogCoreConfig{
			ForceEnable:   true,
			DisableTimers: true,
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, config)
	a := core.activityLog
	a.SetEnable(true)

	core.setupClientIDsUsageInfo(context.Background())

	// wait for clientIDs to be loaded into memory
	verifyClientsLoadedInMemory := func() {
		corehelpers.RetryUntil(t, 60*time.Second, func() error {
			if a.GetClientIDsUsageInfoLoaded() {
				return fmt.Errorf("loaded clientIDs to memory")
			}
			return nil
		})
	}
	verifyClientsLoadedInMemory()

	require.Len(t, a.GetClientIDsUsageInfo(), 0)

	// verify telemetry metric value client.billing_period.activity
	intervals := inMemSink.Data()

	// Test crossed an interval boundary, don't try to deal with it.
	// If we start close to the end of an interval, metrics will be split across two buckets
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}

	// vault.client.billing_period.activity should not be present in ce
	fullMetricName := fmt.Sprintf("client.billing_period.activity;cluster=%s", testClusterName)
	_, ok := intervals[0].Gauges[fullMetricName]
	require.False(t, ok)
}
