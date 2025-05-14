// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/framework"
)

// sendCurrentFragment is a no-op on OSS
func (a *ActivityLog) sendCurrentFragment(ctx context.Context) error {
	return nil
}

// setupClientIDsUsageInfo is a no-op on OSS
func (c *Core) setupClientIDsUsageInfo(ctx context.Context) {
}

// handleClientIDsInMemoryEndOfMonth is a no-op on OSS
func (a *ActivityLog) handleClientIDsInMemoryEndOfMonth(ctx context.Context, currentTime time.Time) {
}

// getStartEndTime parses input for start and end times
// If the end time is after the end of last month, it is adjusted to the last month
func getStartEndTime(d *framework.FieldData, now time.Time, billingStartTime time.Time) (time.Time, time.Time, StartEndTimesWarnings, error) {
	warnings := StartEndTimesWarnings{}
	startTime, endTime, err := parseStartEndTimes(d, billingStartTime)
	if err != nil {
		return startTime, endTime, warnings, err
	}
	// ensure end time is adjusted to the past month if it falls within the current month
	// or is in a future month
	if !endTime.Before(timeutil.StartOfMonth(now.UTC())) {
		endTime = timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, timeutil.StartOfMonth(now.UTC())))
		warnings.EndTimeAdjustedToPastMonth = true
	}

	return startTime, endTime, warnings, nil
}

// clientsGaugeCollectorCurrentBillingPeriod is no-op on CE
func (c *Core) clientsGaugeCollectorCurrentBillingPeriod(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	return []metricsutil.GaugeLabelValues{}, nil
}
