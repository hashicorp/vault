// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raftha

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	consulstorage "github.com/hashicorp/vault/helper/testhelpers/teststorage/consul"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestRaft_HA_NewCluster(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		t.Parallel()

		t.Run("no_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, teststorage.MakeFileBackend, false)
		})

		t.Run("with_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, teststorage.MakeFileBackend, true)
		})
	})

	t.Run("inmem", func(t *testing.T) {
		t.Parallel()

		t.Run("no_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, teststorage.MakeInmemBackend, false)
		})

		t.Run("with_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, teststorage.MakeInmemBackend, true)
		})
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()

		t.Run("no_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, consulstorage.MakeConsulBackend, false)
		})

		t.Run("with_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, consulstorage.MakeConsulBackend, true)
		})
	})
}

// TestRaftHA_Recover_Cluster test that we can recover data and re-boostrap a cluster
// that was created with raft HA enabled but is not using raft as the storage backend.
func TestRaftHA_Recover_Cluster(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())
	t.Run("file", func(t *testing.T) {
		physBundle := teststorage.MakeFileBackend(t, logger)
		testRaftHARecoverCluster(t, physBundle, logger)
	})
	t.Run("inmem", func(t *testing.T) {
		physBundle := teststorage.MakeInmemBackend(t, logger)
		testRaftHARecoverCluster(t, physBundle, logger)
	})
}

func testRaftHANewCluster(t *testing.T, bundler teststorage.PhysicalBackendBundler, addClientCerts bool) {
	var conf vault.CoreConfig
	opts := vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}

	teststorage.RaftHASetup(&conf, &opts, bundler)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	defer cluster.Cleanup()

	joinFunc := func(client *api.Client, addClientCerts bool) {
		req := &api.RaftJoinRequest{
			LeaderCACert: string(cluster.CACertPEM),
		}
		if addClientCerts {
			req.LeaderClientCert = string(cluster.CACertPEM)
			req.LeaderClientKey = string(cluster.CAKeyPEM)
		}
		resp, err := client.Sys().RaftJoin(req)
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Joined {
			t.Fatalf("failed to join raft cluster")
		}
	}

	joinFunc(cluster.Cores[1].Client, addClientCerts)
	joinFunc(cluster.Cores[2].Client, addClientCerts)

	// Ensure peers are added
	leaderClient := cluster.Cores[0].Client
	err := testhelpers.VerifyRaftPeers(t, leaderClient, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test remove peers
	_, err = leaderClient.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = leaderClient.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure peers are removed
	err = testhelpers.VerifyRaftPeers(t, leaderClient, map[string]bool{
		"core-0": true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

// testRaftHARecoverCluster : in this test, we're going to create a raft HA cluster and store a test secret in a KVv2
// We're going to simulate an outage and destroy the cluster but we'll keep the storage backend.
// We'll recreate a new cluster with the same storage backend and ensure that we can recover using
// sys/storage/raft/bootstrap. We'll check that the new cluster
// is functional and no data was lost: we can get the test secret from the KVv2.
func testRaftHARecoverCluster(t *testing.T, physBundle *vault.PhysicalBackendBundle, logger hclog.Logger) {
	opts := vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		// We're not testing the HA, only that it can be recovered. No need for multiple cores.
		NumCores: 1,
	}

	haStorage, haCleanup := teststorage.MakeReusableRaftHAStorage(t, logger, opts.NumCores, physBundle)
	defer haCleanup()
	haStorage.Setup(nil, &opts)
	cluster := vault.NewTestCluster(t, nil, &opts)

	var (
		clusterBarrierKeys [][]byte
		clusterRootToken   string
	)
	clusterBarrierKeys = cluster.BarrierKeys
	clusterRootToken = cluster.RootToken
	leaderCore := cluster.Cores[0]
	testhelpers.EnsureCoreUnsealed(t, cluster, leaderCore)

	leaderClient := cluster.Cores[0].Client
	leaderClient.SetToken(clusterRootToken)

	// Mount a KVv2 backend to store a test data
	err := leaderClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	kvData := map[string]interface{}{
		"data": map[string]interface{}{
			"kittens": "awesome",
		},
	}

	// Store the test data in the KVv2 backend
	_, err = leaderClient.Logical().Write("kv/data/test_known_data", kvData)
	require.NoError(t, err)

	// We now have a raft HA cluster with a KVv2 backend enabled and a test data.
	// We're now going to delete the cluster and create a new raft HA cluster with the same backend storage
	// and ensure we can recover to a working vault cluster and don't lose the data from the backend storage.

	opts = vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		// We're not testing the HA, only that it can be recovered. No need for multiple cores.
		NumCores: 1,
		// It's already initialized as we keep the same storage backend.
		SkipInit: true,
	}
	haStorage, haCleanup = teststorage.MakeReusableRaftHAStorage(t, logger, opts.NumCores, physBundle)
	defer haCleanup()
	haStorage.Setup(nil, &opts)
	clusterRestored := vault.NewTestCluster(t, nil, &opts)

	clusterRestored.BarrierKeys = clusterBarrierKeys
	clusterRestored.RootToken = clusterRootToken
	leaderCoreRestored := clusterRestored.Cores[0]

	testhelpers.EnsureCoresUnsealed(t, clusterRestored)

	leaderClientRestored := clusterRestored.Cores[0].Client

	// We now reset the TLS keyring and bootstrap the cluster again.
	_, err = leaderClientRestored.Logical().Write("sys/storage/raft/bootstrap", nil)
	require.NoError(t, err)

	vault.TestWaitActive(t, leaderCoreRestored.Core)
	// Core should be active and cluster in a working state. We should be able to
	// read the data from the KVv2 backend.
	leaderClientRestored.SetToken(clusterRootToken)
	secretRaw, err := leaderClientRestored.Logical().Read("kv/data/test_known_data")
	require.NoError(t, err)

	data := secretRaw.Data["data"]
	dataAsMap := data.(map[string]interface{})
	require.NotNil(t, dataAsMap)
	require.Equal(t, "awesome", dataAsMap["kittens"])

	// Ensure no writes are happening before we try to clean it up, to prevent
	// issues deleting the files.
	clusterRestored.EnsureCoresSealed(t)
	clusterRestored.Cleanup()
}

func TestRaft_HA_ExistingCluster(t *testing.T) {
	t.Parallel()
	conf := vault.CoreConfig{
		DisablePerformanceStandby: true,
	}
	opts := vault.TestClusterOptions{
		HandlerFunc:        vaulthttp.Handler,
		NumCores:           vault.DefaultNumCores,
		KeepStandbysSealed: true,
	}
	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())

	physBundle := teststorage.MakeInmemBackend(t, logger)
	physBundle.HABackend = nil

	storage, cleanup := teststorage.MakeReusableStorage(t, logger, physBundle)
	defer cleanup()

	var (
		clusterBarrierKeys [][]byte
		clusterRootToken   string
	)
	createCluster := func(t *testing.T) {
		t.Log("simulating cluster creation without raft as HABackend")

		storage.Setup(&conf, &opts)

		cluster := vault.NewTestCluster(t, &conf, &opts)
		defer func() {
			cluster.Cleanup()
			storage.Cleanup(t, cluster)
		}()

		clusterBarrierKeys = cluster.BarrierKeys
		clusterRootToken = cluster.RootToken
	}

	createCluster(t)

	haStorage, haCleanup := teststorage.MakeReusableRaftHAStorage(t, logger, opts.NumCores, physBundle)
	defer haCleanup()

	updateCluster := func(t *testing.T) {
		t.Log("simulating cluster update with raft as HABackend")

		opts.SkipInit = true
		haStorage.Setup(&conf, &opts)

		cluster := vault.NewTestCluster(t, &conf, &opts)
		defer func() {
			cluster.Cleanup()
			haStorage.Cleanup(t, cluster)
		}()

		// Set cluster values
		cluster.BarrierKeys = clusterBarrierKeys
		cluster.RootToken = clusterRootToken

		leaderCore := cluster.Cores[0]
		testhelpers.EnsureCoreUnsealed(t, cluster, leaderCore)

		// Call the bootstrap on the leader and then ensure that it becomes active
		leaderClient := cluster.Cores[0].Client
		leaderClient.SetToken(clusterRootToken)
		{
			_, err := leaderClient.Logical().Write("sys/storage/raft/bootstrap", nil)
			if err != nil {
				t.Fatal(err)
			}
			vault.TestWaitActive(t, leaderCore.Core)
		}

		// Now unseal core for join commands to work
		testhelpers.EnsureCoresUnsealed(t, cluster)

		joinFunc := func(client *api.Client) {
			req := &api.RaftJoinRequest{
				LeaderCACert: string(cluster.CACertPEM),
			}
			resp, err := client.Sys().RaftJoin(req)
			if err != nil {
				t.Fatal(err)
			}
			if !resp.Joined {
				t.Fatalf("failed to join raft cluster")
			}
		}

		joinFunc(cluster.Cores[1].Client)
		joinFunc(cluster.Cores[2].Client)

		// Ensure peers are added
		err := testhelpers.VerifyRaftPeers(t, leaderClient, map[string]bool{
			"core-0": true,
			"core-1": true,
			"core-2": true,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	updateCluster(t)
}

// TestRaftHACluster_Removed_ReAdd creates a raft HA cluster with a file
// backend. The test adds two standbys to the cluster and then removes one of
// them. The removed follower tries to re-join, and the test verifies that it
// errors and cannot join.
func TestRaftHACluster_Removed_ReAdd(t *testing.T) {
	t.Parallel()
	var conf vault.CoreConfig
	opts := vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	teststorage.RaftHASetup(&conf, &opts, teststorage.MakeFileBackend)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	defer cluster.Cleanup()
	vault.TestWaitActive(t, cluster.Cores[0].Core)

	leader := cluster.Cores[0]
	follower := cluster.Cores[2]
	joinReq := &api.RaftJoinRequest{LeaderCACert: string(cluster.CACertPEM)}
	_, err := follower.Client.Sys().RaftJoin(joinReq)
	require.NoError(t, err)
	_, err = cluster.Cores[1].Client.Sys().RaftJoin(joinReq)
	require.NoError(t, err)

	testhelpers.RetryUntil(t, 3*time.Second, func() error {
		resp, err := leader.Client.Sys().RaftAutopilotState()
		if err != nil {
			return err
		}
		if len(resp.Servers) != 3 {
			return errors.New("need 3 servers")
		}
		for serverID, server := range resp.Servers {
			if !server.Healthy {
				return fmt.Errorf("server %s is unhealthy", serverID)
			}
			if server.NodeType != "voter" {
				return fmt.Errorf("server %s has type %s", serverID, server.NodeType)
			}
		}
		return nil
	})
	_, err = leader.Client.Logical().Write("/sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": follower.NodeID,
	})
	require.NoError(t, err)
	require.Eventually(t, follower.Sealed, 10*time.Second, 250*time.Millisecond)

	_, err = follower.Client.Sys().RaftJoin(joinReq)
	require.Error(t, err)
}
