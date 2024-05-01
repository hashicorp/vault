// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package activity_testonly

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestACMERegeneration_RegeneratePreviousMonths(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

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

	testhelpers.EnsureCoresSealed(t, cluster)
	testhelpers.EnsureCoresUnsealed(t, cluster)
	testhelpers.WaitForActiveNode(t, cluster)

	testhelpers.RetryUntil(t, 10*time.Second, func() error {
		_, err := client.Logical().Read("sys/raw/sys/counters/activity/acme-regeneration")
		if err != nil {
			return err
		}
		return nil
	})
	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{})
	require.NoError(t, err)
	require.Equal(t, vault.ResponseCounts{
		NonEntityTokens:  26,
		NonEntityClients: 26,
		Clients:          43,
		ACMEClients:      17,
	}, getTotals(t, resp))
}
