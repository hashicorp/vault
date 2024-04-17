// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package activity_testonly

import (
	"context"
	"encoding/json"
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

var allClientTypeTestCases = []struct {
	clientType       string
	topLevelJSONKey  string
	responseCountsFn func(r vault.ResponseCounts) int
}{
	{
		clientType:      vault.ACMEActivityType,
		topLevelJSONKey: "acme_clients",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.ACMEClients
		},
	},
	{
		clientType:      "secret-sync",
		topLevelJSONKey: "secret_syncs",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.SecretSyncs
		},
	},
	{
		clientType:      "entity",
		topLevelJSONKey: "entity_clients",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.EntityClients
		},
	},
	{
		clientType:      "non-entity-token",
		topLevelJSONKey: "non_entity_clients",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.NonEntityClients
		},
	},
}

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

func getNamespaceData(t *testing.T, resp *api.Secret) []vault.ResponseNamespace {
	t.Helper()
	nsRaw, ok := resp.Data["by_namespace"]
	require.True(t, ok)
	nsResponse := make([]vault.ResponseNamespace, 0)
	err := mapstructure.Decode(nsRaw, &nsResponse)
	require.NoError(t, err)
	return nsResponse
}

func getTotals(t *testing.T, resp *api.Secret) vault.ResponseCounts {
	t.Helper()
	totalRaw, ok := resp.Data["total"]
	require.True(t, ok)
	total := vault.ResponseCounts{}
	err := mapstructure.Decode(totalRaw, &total)
	require.NoError(t, err)
	return total
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
				NewCurrentMonthData().
				NewClientsSeen(10, clientcountutil.WithClientType(tc.clientType)).
				Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)
			require.NoError(t, err)

			now := time.Now().UTC()
			resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
				"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
				"start_time": {timeutil.StartOfMonth(now).Format(time.RFC3339)},
			})
			require.NoError(t, err)

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

// Test_ActivityLogCurrentMonth_Response runs for each client type. The subtest
// creates 10 clients of the type and verifies that the activity log partial
// month response returns 10 clients of that type at every level of the response
// hierarchy
func Test_ActivityLogCurrentMonth_Response(t *testing.T) {
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
				NewCurrentMonthData().
				NewClientsSeen(10, clientcountutil.WithClientType(tc.clientType)).
				Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)
			require.NoError(t, err)

			resp, err := client.Logical().Read("sys/internal/counters/activity/monthly")
			require.NoError(t, err)

			clientsOfType, ok := resp.Data[tc.topLevelJSONKey]
			require.True(t, ok)
			require.Equal(t, json.Number("10"), clientsOfType)
			clients, ok := resp.Data["clients"]
			require.True(t, ok)
			require.Equal(t, json.Number("10"), clients)

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

// Test_ActivityLog_Deduplication runs for all client types. The subtest
// verifies that the clients of that type are deduplicated across months. The
// test creates 10 clients and repeats those clients in later months, then also
// registers 3 and then 2 new clients. The test verifies that the total number
// of clients is 15 (10 + 2 + 3), ensuring that the duplicates are not included
func Test_ActivityLog_Deduplication(t *testing.T) {
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
				NewPreviousMonthData(3).
				NewClientsSeen(10, clientcountutil.WithClientType(tc.clientType)).
				NewPreviousMonthData(2).
				RepeatedClientsSeen(4, clientcountutil.WithClientType(tc.clientType)).
				NewClientsSeen(3, clientcountutil.WithClientType(tc.clientType)).
				NewPreviousMonthData(1).
				RepeatedClientsSeen(5, clientcountutil.WithClientType(tc.clientType)).
				NewClientsSeen(2, clientcountutil.WithClientType(tc.clientType)).
				Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES)
			require.NoError(t, err)

			now := time.Now().UTC()
			resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
				"end_time":   {timeutil.StartOfMonth(now).Format(time.RFC3339)},
				"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(4, now)).Format(time.RFC3339)},
			},
			)
			require.NoError(t, err)

			total := getTotals(t, resp)
			require.Equal(t, 15, tc.responseCountsFn(total))
			require.Equal(t, 15, total.Clients)
		})
	}
}

// Test_ActivityLog_MountDeduplication writes data for the previous
// month across 4 mounts. The cubbyhole and sys mounts have clients in the
// current month as well. The test verifies that the mount counts are correctly
// summed in the results when the previous and current month are queried.
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
		NewPreviousMonthData(1).
		NewClientSeen(clientcountutil.WithClientMount("sys")).
		NewClientSeen(clientcountutil.WithClientMount("secret")).
		NewClientSeen(clientcountutil.WithClientMount("cubbyhole")).
		NewClientSeen(clientcountutil.WithClientMount("identity")).
		NewCurrentMonthData().
		NewClientSeen(clientcountutil.WithClientMount("cubbyhole")).
		NewClientSeen(clientcountutil.WithClientMount("sys")).
		Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES, generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now)).Format(time.RFC3339)},
	})

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
