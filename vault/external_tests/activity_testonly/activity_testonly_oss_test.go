// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly && !enterprise

package activity_testonly

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/timeutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
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

// Test_ActivityLog_LoseLeadership writes data for the second last month, then causes the
// active node to lose leadership. Once a new node becomes the leader, then the
// test queries for the second last month data and verifies that the data from
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
		NewPreviousMonthData(1).
		NewClientsSeen(10).
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES, generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES)
	require.NoError(t, err)
	now := time.Now().UTC()

	testhelpers.EnsureCoreSealed(t, active)
	newActive := testhelpers.WaitForActiveNode(t, cluster)
	standby := active
	testhelpers.WaitForStandbyNode(t, standby)
	testhelpers.EnsureCoreUnsealed(t, cluster, standby)

	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	startPastMonth := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now))
	resp, err := newActive.Client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {endPastMonth.Format(time.RFC3339)},
		"start_time": {startPastMonth.Format(time.RFC3339)},
	})
	require.NoError(t, err)
	// verify start and end times in the response
	require.Equal(t, resp.Data["start_time"], startPastMonth.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))
	monthResponse := getMonthsData(t, resp)
	require.Len(t, monthResponse, 1)
	require.Equal(t, 10, monthResponse[0].NewClients.Counts.Clients)
}

// Test_ActivityLog_ClientsOverlapping writes data for the second last month and
// the previous month. In the second last month, 7 new clients are seen. In the previous
// month, there are 5 repeated and 2 new clients. The test queries over the
// second last and previous months, and verifies that the repeated clients are not
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
		NewPreviousMonthData(2).
		NewClientsSeen(7).
		NewPreviousMonthData(1).
		RepeatedClientsSeen(5).
		NewClientsSeen(2).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()
	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	startTwoMonthsAgo := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(2, now))
	// query from the beginning of the second last month to the end of the previous month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {endPastMonth.Format(time.RFC3339)},
		"start_time": {startTwoMonthsAgo.Format(time.RFC3339)},
	})
	require.NoError(t, err)
	// verify start and end times in the response
	require.Equal(t, resp.Data["start_time"], startTwoMonthsAgo.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))
	monthsResponse := getMonthsData(t, resp)
	require.Len(t, monthsResponse, 2)
	for _, month := range monthsResponse {
		ts, err := time.Parse(time.RFC3339, month.Timestamp)
		require.NoError(t, err)
		// The previous month should have a total of 7 clients
		// 2 of those will be considered new
		if ts.UTC().Equal(timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now))) {
			require.Equal(t, month.Counts.Clients, 7)
			require.Equal(t, month.NewClients.Counts.Clients, 2)
		} else {
			// All clients will be considered new for the second last month
			require.Equal(t, month.Counts.Clients, 7)
			require.Equal(t, month.NewClients.Counts.Clients, 7)

		}
	}
}

// Test_ActivityLog_ClientsNewCurrentMonth writes data for the second last month and
// past month with 5 repeated clients and 2 new clients in the past month.
// The test then queries the activity log for only the past month, and
// verifies that all 7 clients seen the past month are considered new.
func Test_ActivityLog_ClientsNewCurrentMonth(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(2).
		NewClientsSeen(5).
		NewPreviousMonthData(1).
		RepeatedClientsSeen(5).
		NewClientsSeen(2).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()
	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	startPastMonth := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now))
	// query from the beginning of the second last month to the end of the previous month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {endPastMonth.Format(time.RFC3339)},
		"start_time": {startPastMonth.Format(time.RFC3339)},
	})
	// verify start and end times in the response
	require.Equal(t, resp.Data["start_time"], startPastMonth.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))
	require.NoError(t, err)
	monthsResponse := getMonthsData(t, resp)
	require.Len(t, monthsResponse, 1)
	require.Equal(t, 7, monthsResponse[0].NewClients.Counts.Clients)
}

// Test_ActivityLog_EmptyDataMonths writes data for only the past month,
// then queries a timeframe of several months in the past to now. The test
// verifies that empty months of data are returned for the past, and the past
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
		NewPreviousMonthData(1).
		NewClientsSeen(10).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()
	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	startFourMonthsAgo := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(4, now))
	// query from the beginning of 4 months ago to the end of past month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {endPastMonth.Format(time.RFC3339)},
		"start_time": {startFourMonthsAgo.Format(time.RFC3339)},
	})
	require.NoError(t, err)
	// verify start and end times in the response
	require.Equal(t, resp.Data["start_time"], startFourMonthsAgo.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))
	monthsResponse := getMonthsData(t, resp)

	require.Len(t, monthsResponse, 4)
	for _, month := range monthsResponse {
		ts, err := time.Parse(time.RFC3339, month.Timestamp)
		require.NoError(t, err)
		// past month should have data
		if ts.UTC().Equal(timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now))) {
			require.Equal(t, month.Counts.Clients, 10)
		} else {
			// other months should be empty
			require.Nil(t, month.Counts)
		}
	}
}

// Test_ActivityLog_FutureEndDate queries a start time from the past
// and an end date in the future. The test
// verifies that the current month is returned in the response.
func Test_ActivityLog_FutureEndDate(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(1).
		NewClientsSeen(10).
		NewCurrentMonthData().
		NewClientsSeen(10).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	now := time.Now().UTC()
	startThreeMonthsAgo := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(3, now))
	// query from the beginning of 3 months ago to beginning of next month
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.StartOfNextMonth(now).Format(time.RFC3339)},
		"start_time": {startThreeMonthsAgo.Format(time.RFC3339)},
	})
	require.NoError(t, err)
	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	// verify start and end times in the response
	// end time must be adjusted to past month if within current month or in future date
	require.Equal(t, resp.Data["start_time"], startThreeMonthsAgo.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))
	monthsResponse := getMonthsData(t, resp)

	require.Len(t, monthsResponse, 3)

	// Get the last month of data in the slice
	expectedCurrentMonthData := monthsResponse[2]
	expectedTime, err := time.Parse(time.RFC3339, expectedCurrentMonthData.Timestamp)
	require.NoError(t, err)
	if !timeutil.IsCurrentMonth(expectedTime, endPastMonth) {
		t.Fatalf("final month data is not past month")
	}
}

// Test_ActivityLog_ClientTypeResponse runs for each client type. In the
// subtests, 10 clients of the type are created and the test verifies that the
// activity log query response returns 10 clients of that type at every level of
// the response hierarchy
func Test_ActivityLog_ClientTypeResponse(t *testing.T) {
	t.Parallel()
	for _, tc := range allClientTypeTestCases {
		tc := tc
		t.Run(tc.clientType, func(t *testing.T) {
			t.Parallel()
			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
				"enabled": "enable",
			})
			_, err = clientcountutil.NewActivityLogData(client).
				NewPreviousMonthData(1).
				NewClientsSeen(10, clientcountutil.WithClientType(tc.clientType)).
				Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES, generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES)
			require.NoError(t, err)

			now := time.Now().UTC()
			endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
			startPastMonth := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now))
			resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
				"end_time":   {endPastMonth.Format(time.RFC3339)},
				"start_time": {startPastMonth.Format(time.RFC3339)},
			})
			require.NoError(t, err)
			// verify start and end times in the response
			require.Equal(t, resp.Data["start_time"], startPastMonth.UTC().Format(time.RFC3339))
			require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))

			total := getTotals(t, resp)
			require.Equal(t, 10, tc.responseCountsFn(total))
			require.Equal(t, 10, total.Clients)

			byNamespace := getNamespaceData(t, resp)
			require.Equal(t, 10, tc.responseCountsFn(byNamespace[0].Counts))
			require.Equal(t, 10, tc.responseCountsFn(*byNamespace[0].Mounts[0].Counts))
			require.Equal(t, 10, byNamespace[0].Counts.Clients)
			require.Equal(t, 10, byNamespace[0].Mounts[0].Counts.Clients)

			byMonth := getMonthsData(t, resp)
			require.Equal(t, 10, tc.responseCountsFn(*byMonth[0].NewClients.Counts))
			require.Equal(t, 10, tc.responseCountsFn(*byMonth[0].Counts))
			require.Equal(t, 10, tc.responseCountsFn(byMonth[0].Namespaces[0].Counts))
			require.Equal(t, 10, tc.responseCountsFn(*byMonth[0].Namespaces[0].Mounts[0].Counts))
			require.Equal(t, 10, byMonth[0].NewClients.Counts.Clients)
			require.Equal(t, 10, byMonth[0].Counts.Clients)
			require.Equal(t, 10, byMonth[0].Namespaces[0].Counts.Clients)
			require.Equal(t, 10, byMonth[0].Namespaces[0].Mounts[0].Counts.Clients)
		})

	}
}

// Test_ActivityLog_MountDeduplication writes data for the second last
// month across 4 mounts. The cubbyhole and sys mounts have clients in the
// past month as well. The test verifies that the mount counts are correctly
// summed in the results when the second last month and past month are queried.
func Test_ActivityLog_MountDeduplication(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	now := time.Now().UTC()

	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(2).
		NewClientSeen(clientcountutil.WithClientMount("sys")).
		NewClientSeen(clientcountutil.WithClientMount("secret")).
		NewClientSeen(clientcountutil.WithClientMount("cubbyhole")).
		NewClientSeen(clientcountutil.WithClientMount("identity")).
		NewPreviousMonthData(1).
		NewClientSeen(clientcountutil.WithClientMount("cubbyhole")).
		NewClientSeen(clientcountutil.WithClientMount("sys")).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	startTwoMonthsAgo := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(2, now))
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {endPastMonth.Format(time.RFC3339)},
		"start_time": {startTwoMonthsAgo.Format(time.RFC3339)},
	})
	// verify start and end times in the response
	require.Equal(t, resp.Data["start_time"], startTwoMonthsAgo.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))

	require.NoError(t, err)
	byNamespace := getNamespaceData(t, resp)
	require.Len(t, byNamespace, 1)
	require.Len(t, byNamespace[0].Mounts, 4)
	mountSet := make(map[string]int, 4)
	for _, mount := range byNamespace[0].Mounts {
		mountSet[mount.MountPath] = mount.Counts.Clients
	}
	require.Equal(t, map[string]int{
		"identity/":  1,
		"sys/":       2,
		"cubbyhole/": 2,
		"secret/":    1,
	}, mountSet)
}

// TestHandleQuery_MultipleMounts creates a cluster with
// two userpass mounts. It then tests verifies that
// the total new counts are calculated within a reasonably level of accuracy for
// various numbers of clients in each mount.
func TestHandleQuery_MultipleMounts(t *testing.T) {
	tests := map[string]struct {
		twoMonthsAgo          [][]int
		oneMonthAgo           [][]int
		currentMonth          [][]int
		expectedNewClients    int
		expectedTotalAccuracy float64
	}{
		"low volume, all mounts": {
			twoMonthsAgo: [][]int{
				{20, 20},
			},
			oneMonthAgo: [][]int{
				{30, 30},
			},
			currentMonth: [][]int{
				{40, 40},
			},
			expectedNewClients:    80,
			expectedTotalAccuracy: 1,
		},
		"medium volume, all mounts": {
			twoMonthsAgo: [][]int{
				{200, 200},
			},
			oneMonthAgo: [][]int{
				{300, 300},
			},
			currentMonth: [][]int{
				{400, 400},
			},
			expectedNewClients:    800,
			expectedTotalAccuracy: 0.98,
		},
		"higher volume, all mounts": {
			twoMonthsAgo: [][]int{
				{200, 200},
			},
			oneMonthAgo: [][]int{
				{300, 300},
			},
			currentMonth: [][]int{
				{2000, 5000},
			},
			expectedNewClients:    7000,
			expectedTotalAccuracy: 0.95,
		},
		"higher volume, no repeats": {
			twoMonthsAgo: [][]int{
				{200, 200},
			},
			oneMonthAgo: [][]int{
				{300, 300},
			},
			currentMonth: [][]int{
				{4000, 6000},
			},
			expectedNewClients:    10000,
			expectedTotalAccuracy: 0.98,
		},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%s", i)
		t.Run(testname, func(t *testing.T) {
			var err error
			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			_, err = client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
				"enabled": "enable",
			})
			require.NoError(t, err)

			// Create two namespaces
			namespaces := []string{namespace.RootNamespaceID}
			mounts := make(map[string][]string)

			// Add two userpass mounts to each namespace
			for _, ns := range namespaces {
				err = client.WithNamespace(ns).Sys().EnableAuthWithOptions("userpass1", &api.EnableAuthOptions{
					Type: "userpass",
				})
				require.NoError(t, err)
				err = client.WithNamespace(ns).Sys().EnableAuthWithOptions("userpass2", &api.EnableAuthOptions{
					Type: "userpass",
				})
				require.NoError(t, err)
				mounts[ns] = []string{"auth/userpass1", "auth/userpass2"}
			}

			activityLogGenerator := clientcountutil.NewActivityLogData(client)

			// Write three months ago data
			activityLogGenerator = activityLogGenerator.NewPreviousMonthData(3)
			for nsIndex, nsId := range namespaces {
				for mountIndex, mount := range mounts[nsId] {
					activityLogGenerator = activityLogGenerator.
						NewClientsSeen(tt.twoMonthsAgo[nsIndex][mountIndex], clientcountutil.WithClientNamespace(nsId), clientcountutil.WithClientMount(mount))
				}
			}

			// Write two months ago data
			activityLogGenerator = activityLogGenerator.NewPreviousMonthData(2)
			for nsIndex, nsId := range namespaces {
				for mountIndex, mount := range mounts[nsId] {
					activityLogGenerator = activityLogGenerator.
						NewClientsSeen(tt.oneMonthAgo[nsIndex][mountIndex], clientcountutil.WithClientNamespace(nsId), clientcountutil.WithClientMount(mount))
				}
			}

			// Write previous month data
			activityLogGenerator = activityLogGenerator.NewPreviousMonthData(1)
			for nsIndex, nsPath := range namespaces {
				for mountIndex, mount := range mounts[nsPath] {
					activityLogGenerator = activityLogGenerator.
						RepeatedClientSeen(clientcountutil.WithClientNamespace(nsPath), clientcountutil.WithClientMount(mount)).
						NewClientsSeen(tt.currentMonth[nsIndex][mountIndex], clientcountutil.WithClientNamespace(nsPath), clientcountutil.WithClientMount(mount))
				}
			}

			// Write all the client count data
			_, err = activityLogGenerator.Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
			require.NoError(t, err)

			endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, time.Now()).UTC())
			startThreeMonthsAgo := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(3, time.Now().UTC()))

			// query activity log
			resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
				"end_time":   {endPastMonth.Format(time.RFC3339)},
				"start_time": {startThreeMonthsAgo.Format(time.RFC3339)},
			})
			require.NoError(t, err)
			// verify start and end times in the response
			require.Equal(t, resp.Data["start_time"], startThreeMonthsAgo.UTC().Format(time.RFC3339))
			require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))

			// Ensure that the month response is the same as the totals, because all clients
			// are new clients and there will be no approximation in the single month partial
			// case
			monthsRaw, ok := resp.Data["months"]
			if !ok {
				t.Fatalf("malformed results. got %v", resp.Data)
			}
			monthsResponse := make([]*vault.ResponseMonth, 0)
			err = mapstructure.Decode(monthsRaw, &monthsResponse)

			currentMonthClients := monthsResponse[len(monthsResponse)-1]

			// Now verify that the new client totals for ALL namespaces are approximately accurate (there are no namespaces in CE)
			newClientsError := math.Abs((float64)(currentMonthClients.NewClients.Counts.Clients - tt.expectedNewClients))
			newClientsErrorMargin := newClientsError / (float64)(tt.expectedNewClients)
			expectedAccuracyCalc := (1 - tt.expectedTotalAccuracy) * 100 / 100
			if newClientsErrorMargin > expectedAccuracyCalc {
				t.Fatalf("bad accuracy: expected %+v, found %+v", expectedAccuracyCalc, newClientsErrorMargin)
			}

			// Verify that the totals for the clients are visibly sensible (that is the total of all the individual new clients per namespace)
			total := 0
			for _, newClientCounts := range currentMonthClients.NewClients.Namespaces {
				total += newClientCounts.Counts.Clients
			}
			if diff := math.Abs(float64(currentMonthClients.NewClients.Counts.Clients - total)); diff >= 1 {
				t.Fatalf("total expected was %d but got %d", currentMonthClients.NewClients.Counts.Clients, total)
			}
		})
	}
}

// TestActivityLog_CountersAPI_NoErrorOnLoadingClientIDsToMemoryFlag_CE verifies that default counters api is not blocked by clientIDsUsageInfoLoaded flag as it always remains false on CE
func TestActivityLog_CountersAPI_NoErrorOnLoadingClientIDsToMemoryFlag_CE(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	core := cluster.Cores[0]
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	a := core.GetActivityLog()
	a.SetEnable(true)
	now := time.Now().UTC()

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

	// add some data to previous months
	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(2).
		NewClientSeen(clientcountutil.WithClientMount("sys")).
		NewClientSeen(clientcountutil.WithClientMount("secret")).
		NewClientSeen(clientcountutil.WithClientMount("cubbyhole")).
		NewClientSeen(clientcountutil.WithClientMount("identity")).
		NewPreviousMonthData(1).
		NewClientSeen(clientcountutil.WithClientMount("cubbyhole")).
		NewClientSeen(clientcountutil.WithClientMount("sys")).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	// clientIDs in memory should be 0 as they are not updated in CE
	require.Len(t, a.GetClientIDsUsageInfo(), 0)

	// default counters api query
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{})
	require.NoError(t, err)

	// verify query response response
	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))

	// start time will be default time for time.Time{} as no start time is specified in input params
	require.Equal(t, resp.Data["start_time"], time.Time{}.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))

	byNamespace := getNamespaceData(t, resp)
	require.Len(t, byNamespace, 1)
	require.Len(t, byNamespace[0].Mounts, 4)
	mountSet := make(map[string]int, 4)
	for _, mount := range byNamespace[0].Mounts {
		mountSet[mount.MountPath] = mount.Counts.Clients
	}
	require.Equal(t, map[string]int{
		"identity/":  1,
		"sys/":       2,
		"cubbyhole/": 2,
		"secret/":    1,
	}, mountSet)
}
