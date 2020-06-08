package rafttests

import (
	"sync/atomic"
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
)

func TestRaft_HA_NewCluster(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		t.Run("no_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, teststorage.MakeFileBackend, false)
		})

		t.Run("with_client_certs", func(t *testing.T) {
			testRaftHANewCluster(t, teststorage.MakeFileBackend, true)
		})
	})

}

func testRaftHANewCluster(t *testing.T, bundler teststorage.PhysicalBackendBundler, addClientCerts bool) {
	var conf vault.CoreConfig
	var opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}

	teststorage.RaftHASetup(&conf, &opts, bundler)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer cluster.Cleanup()

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	// leaderAPI := leaderCore.Client.Address()
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
			// LeaderAPIAddr: leaderAPI,
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
	checkConfigFunc(t, leaderClient, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})

	// Test remove peers
	_, err := leaderClient.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
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
	checkConfigFunc(t, leaderClient, map[string]bool{
		"core-0": true,
	})
}

func TestRaft_HA_ExistingCluster(t *testing.T) {
	t.Skipf("not implemented")
}
