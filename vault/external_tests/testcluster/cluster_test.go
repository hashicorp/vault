// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package testcluster

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// These tests are focused on testing the test infrastructure itself.

func SyncTestNow(t *testing.T, f func(t *testing.T)) {
	toSleep := time.Now().Sub(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	synctest.Test(t, func(t *testing.T) {
		time.Sleep(toSleep)
		f(t)
	})
}

func TestNewTestClusterInmemNetworkListener(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(nil, nil, nil)
	opts.SyncTest = true
	// We'd need https://github.com/hashicorp/vault-plugin-secrets-kv/pull/243 to eliminate this.
	opts.SkipKVMount = true
	SyncTestNow(t, func(t *testing.T) {
		cluster := vault.NewTestCluster(t, conf, opts)
		defer func() {
			synctest.Wait()
			cluster.CleanupSyncTest()
		}()
		testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
		for i := range cluster.Cores {
			require.NoError(t, testcluster.NodeHealthy(t.Context(), cluster, i))
		}
	})
}
