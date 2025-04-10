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
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestACMERegeneration_RegenerateWithCurrentMonth writes segments for previous
// months and the current month. The test regenerates the precomputed queries,
// and verifies that the counts are correct when querying both with and without
// the current month
func TestACMERegeneration_RegenerateWithCurrentMonth(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{EnableRaw: true})
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)
	now := time.Now().UTC()
	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(3).
		// 3 months ago, 15 non-entity clients and 10 ACME clients
		NewClientsSeen(15, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(10, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewPreviousMonthData(2).
		// 2 months ago, 7 new non-entity clients and 5 new ACME clients
		RepeatedClientsSeen(2, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(7, clientcountutil.WithClientType("non-entity-token")).
		RepeatedClientsSeen(5, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewClientsSeen(5, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewPreviousMonthData(1).
		// 1 months ago, 4 new non-entity clients and 2 new ACME clients
		RepeatedClientsSeen(3, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(4, clientcountutil.WithClientType("non-entity-token")).
		RepeatedClientsSeen(1, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewClientsSeen(2, clientcountutil.WithClientType(vault.ACMEActivityType)).

		// current month, 10 new non-entity clients and 20 new ACME clients
		NewCurrentMonthData().
		NewClientsSeen(10, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(20, clientcountutil.WithClientType(vault.ACMEActivityType)).
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)

	require.NoError(t, err)

	forceRegeneration(t, cluster)

	startFiveMonthsAgo := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(5, now))
	endPastMonth := timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now))
	// current month isn't included in this query
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"start_time": {startFiveMonthsAgo.Format(time.RFC3339)},
		"end_time":   {endPastMonth.Format(time.RFC3339)},
	})
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityClients: 26,
		Clients:          43,
		ACMEClients:      17,
	}, getTotals(t, resp))
	// verify start and end times in the response
	require.Equal(t, resp.Data["start_time"], startFiveMonthsAgo.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))

	// explicitly include the current month in the request
	// the given end time is adjusted to the last month, excluding the current month at the API
	respWithCurrent, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"start_time": {startFiveMonthsAgo.Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
	})
	// verify start and end times in the response
	// end time is expected to be adjusted to the past month, excluding the current month
	require.Equal(t, resp.Data["start_time"], startFiveMonthsAgo.UTC().Format(time.RFC3339))
	require.Equal(t, resp.Data["end_time"], endPastMonth.UTC().Format(time.RFC3339))
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityClients: 26,
		Clients:          43,
		ACMEClients:      17,
	}, getTotals(t, respWithCurrent))
}
