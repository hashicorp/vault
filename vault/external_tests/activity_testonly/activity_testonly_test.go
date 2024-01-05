// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package activity_testonly

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/timeutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

// Test_ActivityLog_LoseLeadership writes data for this month, then causes the
// active node to lose leadership. Once a new node becomes the leader, then the
// test queries for the current month data and verifies that the data from
// before the leadership transfer is returned
func Test_ActivityLog_LoseLeadership(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    2,
	})
	cluster.Start()
	defer cluster.Cleanup()

	active := testhelpers.DeriveStableActiveCore(t, cluster)
	client := active.Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)

	_, err = clientcountutil.NewActivityLogData(client).
		NewCurrentMonthData().
		NewClientsSeen(10).
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)
	now := time.Now().UTC()

	testhelpers.EnsureCoreSealed(t, active)
	newActive := testhelpers.WaitForActiveNode(t, cluster)
	standby := active
	testhelpers.WaitForStandbyNode(t, standby)
	testhelpers.EnsureCoreUnsealed(t, cluster, standby)

	resp, err := newActive.Client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"start_time": {timeutil.StartOfMonth(now).Format(time.RFC3339)},
	})
	monthResponse := getMonthsData(t, resp)
	require.Len(t, monthResponse, 1)
	require.Equal(t, 10, monthResponse[0].NewClients.Counts.Clients)
}

// Test_ActivityLog_ClientsOverlapping writes data for the previous month and
// current month. In the previous month, 7 new clients are seen. In the current
// month, there are 5 repeated and 2 new clients. The test queries over the
// previous and current months, and verifies that the repeated clients are not
// considered new
func Test_ActivityLog_ClientsOverlapping(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(1).
		NewClientsSeen(7).
		NewCurrentMonthData().
		RepeatedClientsSeen(5).
		NewClientsSeen(2).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()

	// query from the beginning of the previous month to the end of this month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now)).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	monthsResponse := getMonthsData(t, resp)
	require.Len(t, monthsResponse, 2)
	for _, month := range monthsResponse {
		ts, err := time.Parse(time.RFC3339, month.Timestamp)
		require.NoError(t, err)
		// This month should have a total of 7 clients
		// 2 of those will be considered new
		if ts.UTC().Equal(timeutil.StartOfMonth(now)) {
			require.Equal(t, month.Counts.Clients, 7)
			require.Equal(t, month.NewClients.Counts.Clients, 2)
		} else {
			// All clients will be considered new for the previous month
			require.Equal(t, month.Counts.Clients, 7)
			require.Equal(t, month.NewClients.Counts.Clients, 7)

		}
	}
}

// Test_ActivityLog_ClientsNewCurrentMonth writes data for the past month and
// current month with 5 repeated clients and 2 new clients in the current month.
// The test then queries the activity log for only the current month, and
// verifies that all 7 clients seen this month are considered new.
func Test_ActivityLog_ClientsNewCurrentMonth(t *testing.T) {
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
		RepeatedClientsSeen(5).
		NewClientsSeen(2).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()

	// query from the beginning of this month to the end of this month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"start_time": {timeutil.StartOfMonth(now).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	monthsResponse := getMonthsData(t, resp)
	require.Len(t, monthsResponse, 1)
	require.Equal(t, 7, monthsResponse[0].NewClients.Counts.Clients)
}

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

// Test_ActivityLog_EmptyDataMonths writes data for only the current month,
// then queries a timeframe of several months in the past to now. The test
// verifies that empty months of data are returned for the past, and the current
// month data is correct.
func Test_ActivityLog_EmptyDataMonths(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	_, err = clientcountutil.NewActivityLogData(client).
		NewCurrentMonthData().
		NewClientsSeen(10).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()
	// query from the beginning of 3 months ago to the end of this month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(3, now)).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	monthsResponse := getMonthsData(t, resp)

	require.Len(t, monthsResponse, 4)
	for _, month := range monthsResponse {
		ts, err := time.Parse(time.RFC3339, month.Timestamp)
		require.NoError(t, err)
		// current month should have data
		if ts.UTC().Equal(timeutil.StartOfMonth(now)) {
			require.Equal(t, month.Counts.Clients, 10)
		} else {
			// other months should be empty
			require.Nil(t, month.Counts)
		}
	}
}

func getMonthsData(t *testing.T, resp *api.Secret) []vault.ResponseMonth {
	t.Helper()
	monthsRaw, ok := resp.Data["months"]
	require.True(t, ok)
	monthsResponse := make([]vault.ResponseMonth, 0)
	err := mapstructure.Decode(monthsRaw, &monthsResponse)
	require.NoError(t, err)
	return monthsResponse
}
