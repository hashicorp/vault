// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly && !enterprise

package activity_testonly

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/stretchr/testify/require"
)

// Test_ActivityLog_Disable writes data for a past month and a current month and
// then disables the activity log. The test then queries for a timeframe that
// includes both the disabled and enabled dates. The test verifies that the past
// month's data is returned, but there is no current month data.
func Test_ActivityLog_Disable(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(1).
		NewClientsSeen(5).
		NewCurrentMonthData().
		NewClientsSeen(5).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	_, err = client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "disable",
	})

	require.NoError(t, err)

	now := time.Now().UTC()
	// query from the beginning of the previous month to the end of this month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now)).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	monthsResponse := getMonthsData(t, resp)

	// we only expect data for the previous month
	require.Len(t, monthsResponse, 1)
	lastMonthResp := monthsResponse[0]
	ts, err := time.Parse(time.RFC3339, lastMonthResp.Timestamp)
	require.NoError(t, err)
	require.Equal(t, ts.UTC(), timeutil.StartOfPreviousMonth(now.UTC()))
}
