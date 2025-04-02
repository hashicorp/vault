// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/timeutil"
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
// If the end time corresponds to the current month, it is adjusted to the last month
func getStartEndTime(startTime, endTime, billingStartTime time.Time) (time.Time, time.Time, StartEndTimesWarnings, error) {
	warnings := StartEndTimesWarnings{}

	// If a specific endTime is used, then respect that
	// otherwise we want to query up until the end of the current month.
	//
	// Also convert any user inputs to UTC to avoid
	// problems later.
	if endTime.IsZero() {
		endTime = timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, timeutil.StartOfMonth(time.Now().UTC())))
	} else {
		endTime = endTime.UTC()
		if timeutil.IsCurrentMonth(endTime, time.Now().UTC()) {
			endTime = timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, timeutil.StartOfMonth(endTime)))
			warnings.CurrentMonthAsEndTimeIgnored = true
		}
	}

	// If startTime is not specified, we would like to query
	// from the beginning of the billing period
	if startTime.IsZero() {
		startTime = billingStartTime
	} else {
		startTime = startTime.UTC()
	}
	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, StartEndTimesWarnings{}, fmt.Errorf("start_time is later than end_time")
	}

	return startTime, endTime, warnings, nil
}
