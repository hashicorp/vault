// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package activity_testonly

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func forceRegeneration(t *testing.T, cluster *vault.TestCluster) {
	t.Helper()
	client := cluster.Cores[0].Client
	_, err := client.Logical().Delete("sys/raw/sys/counters/activity/acme-regeneration")
	require.NoError(t, err)
	testhelpers.EnsureCoresSealed(t, cluster)
	testhelpers.EnsureCoresUnsealed(t, cluster)
	testhelpers.WaitForActiveNode(t, cluster)

	testhelpers.RetryUntil(t, 10*time.Second, func() error {
		r, err := client.Logical().Read("sys/raw/sys/counters/activity/acme-regeneration")
		if err != nil {
			return err
		}
		if r == nil {
			return errors.New("no response")
		}
		return nil
	})
}

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

	// current month isn't included in this query
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(5, now)).Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, now)).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityClients: 26,
		Clients:          43,
		ACMEClients:      17,
	}, getTotals(t, resp))

	// explicitly include the current month
	respWithCurrent, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(5, now)).Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityClients: 36,
		Clients:          73,
		ACMEClients:      37,
	}, getTotals(t, respWithCurrent))
}

// TestACMERegeneration_RegenerateMuchOlder creates segments 5 months ago, 4
// months ago, and 3 months ago. The test regenerates the precomputed queries
// and then verifies that this older data is included in the generated results.
func TestACMERegeneration_RegenerateMuchOlder(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{EnableRaw: true})
	client := cluster.Cores[0].Client

	now := time.Now().UTC()
	_, err := clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(5).
		// 5 months ago, 15 non-entity clients and 10 ACME clients
		NewClientsSeen(15, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(10, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewPreviousMonthData(4).
		// 4 months ago, 7 new non-entity clients and 5 new ACME clients
		RepeatedClientsSeen(2, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(7, clientcountutil.WithClientType("non-entity-token")).
		RepeatedClientsSeen(5, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewClientsSeen(5, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewPreviousMonthData(3).
		// 3 months ago, 4 new non-entity clients and 2 new ACME clients
		RepeatedClientsSeen(3, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(4, clientcountutil.WithClientType("non-entity-token")).
		RepeatedClientsSeen(1, clientcountutil.WithClientType(vault.ACMEActivityType)).
		NewClientsSeen(2, clientcountutil.WithClientType(vault.ACMEActivityType)).
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)

	require.NoError(t, err)

	forceRegeneration(t, cluster)
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(5, now)).Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityClients: 26,
		Clients:          43,
		ACMEClients:      17,
	}, getTotals(t, resp))
}

// TestACMERegeneration_RegeneratePreviousMonths creates segments for the
// previous 3 months, and no segments for the current month. The test verifies
// that the older data gets regenerated
func TestACMERegeneration_RegeneratePreviousMonths(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{EnableRaw: true})
	client := cluster.Cores[0].Client

	now := time.Now().UTC()
	_, err := clientcountutil.NewActivityLogData(client).
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
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)

	require.NoError(t, err)

	forceRegeneration(t, cluster)

	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
		"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(5, now)).Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
	})
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityClients: 26,
		Clients:          43,
		ACMEClients:      17,
	}, getTotals(t, resp))
}
