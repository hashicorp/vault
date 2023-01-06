package raftha

import (
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	consulstorage "github.com/hashicorp/vault/helper/testhelpers/teststorage/consul"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
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

func Test_RaftHA_Recover_Cluster(t *testing.T) {
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
	cluster.Start()
	defer cluster.Cleanup()

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	// Seal the leader so we can install an address provider
	{
		testhelpers.EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	// Now unseal core for join commands to work
	testhelpers.EnsureCoresUnsealed(t, cluster)

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

// In this test, we're going to create a raft HA cluster and store a test secret in a KVv2
// We're going to simulate an outage and destroy the cluster but we'll keep the storage backend.
// We'll recreate a new cluster with the same storage backend and ensure that we can recover using
// sys/storage/raft/bootstrap with `reset_tls_keyring` set to true. We'll check that the new cluster
// is functional and no data was lost: we can get the test secret from the KVv2.
func testRaftHARecoverCluster(t *testing.T, physBundle *vault.PhysicalBackendBundle, logger hclog.Logger) {
	t.Log("Simulating cluster recovery with raft as HABackend but not storage")
	var conf vault.CoreConfig
	opts := vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		// We're not testing the HA, only that it can be recovered. No need for multiple cores.
		NumCores: 1,
	}

	haStorage, haCleanup := teststorage.MakeReusableRaftHAStorage(t, logger, opts.NumCores, physBundle)
	defer haCleanup()
	haStorage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer cluster.Cleanup()

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	// Seal the leader so we can install an address provider
	{
		testhelpers.EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	var (
		clusterBarrierKeys [][]byte
		clusterRootToken   string
	)
	clusterBarrierKeys = cluster.BarrierKeys
	clusterRootToken = cluster.RootToken
	// Now unseal core
	testhelpers.EnsureCoresUnsealed(t, cluster)

	// Ensure peers are added
	leaderClient := cluster.Cores[0].Client

	// Mount a KVv2 backend to store a test data
	err := leaderClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	kvData := map[string]interface{}{
		"data": map[string]interface{}{
			"kittens": "awesome",
		},
	}

	// Store the test data in the KVv2 backend
	secretRaw, err := leaderClient.Logical().Write("kv/test_known_data", kvData)
	if err != nil {
		t.Fatalf("write secret failed - err :%#v, resp: %#v\n", err, secretRaw)
	}

	// We delete the current cluster. We keep the storage backend so we can recover the cluster
	cluster.Cleanup()

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
	haStorage.Setup(&conf, &opts)
	cluster_restored := vault.NewTestCluster(t, &conf, &opts)
	cluster_restored.Start()
	defer cluster_restored.Cleanup()

	addressProviderRestored := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster_restored}

	cluster_restored.BarrierKeys = clusterBarrierKeys
	cluster_restored.RootToken = clusterRootToken
	leaderCoreRestored := cluster_restored.Cores[0]
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	// Seal the leader so we can install an address provider
	{
		testhelpers.EnsureCoreSealed(t, leaderCoreRestored)
		leaderCore.UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProviderRestored)
		cluster_restored.UnsealCore(t, leaderCoreRestored)
	}

	testhelpers.EnsureCoresUnsealed(t, cluster_restored)

	leaderClientRestored := cluster_restored.Cores[0].Client
	// Trying to bootstrap the cluster, it should fail as it already exists.
	result, err := leaderClientRestored.Logical().Write("sys/storage/raft/bootstrap", nil)
	if err == nil || !strings.Contains(err.Error(), "could not generate TLS keyring during bootstrap: TLS keyring already present") {
		t.Fatalf("re-bootstraping backend should error as TLS keyring already exists.")
	}

	// We now reset the TLS keyring and bootstrap the cluster again.
	result, err = leaderClientRestored.Logical().Write("sys/storage/raft/bootstrap", map[string]interface{}{
		"reset_tls_keyring": true,
	})
	if err != nil {
		t.Fatalf("failed re-bootstraping backend - err :%#v, resp: %#v\n", err, result)
	}

	vault.TestWaitActive(t, leaderCoreRestored.Core)
	// Core should be active and cluster in a working state. We should be able to
	// read the data from the KVv2 backend.
	leaderClientRestored.SetToken(clusterRootToken)
	secretRaw, err = leaderClientRestored.Logical().Read("kv/test_known_data")
	if err != nil {
		t.Fatalf("read secret failed - err :%#v, resp: %#v\n", err, secretRaw)
	}

	kittens := secretRaw.Data["data"].(map[string]interface{})["kittens"]
	if kittens != "awesome" {
		t.Fatalf("expected kittens secret to be awesome but it was %q", kittens)
	}
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
		cluster.Start()
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
		cluster.Start()
		defer func() {
			cluster.Cleanup()
			haStorage.Cleanup(t, cluster)
		}()

		// Set cluster values
		cluster.BarrierKeys = clusterBarrierKeys
		cluster.RootToken = clusterRootToken

		addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}
		atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

		// Seal the leader so we can install an address provider
		leaderCore := cluster.Cores[0]
		{
			testhelpers.EnsureCoreSealed(t, leaderCore)
			leaderCore.UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
			testhelpers.EnsureCoreUnsealed(t, cluster, leaderCore)
		}

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

		// Set address provider
		cluster.Cores[1].UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.Cores[2].UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)

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
