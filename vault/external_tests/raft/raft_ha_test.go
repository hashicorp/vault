package rafttests

import (
	"sync/atomic"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
)

func TestRaft_HA_NewCluster(t *testing.T) {
	t.Parallel()
	t.Run("inmem", func(t *testing.T) {
		t.Parallel()
		testRaftHANewCluster(t, teststorage.InmemBackendOnHARaftSetup)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		testRaftHANewCluster(t, teststorage.FileBackendOnHARaftSetup)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()
		testRaftHANewCluster(t, teststorage.ConsulBackendOnHARaftSetup)
	})
}

func testRaftHANewCluster(t *testing.T, setup teststorage.ClusterSetupMutator) {
	var conf vault.CoreConfig
	var opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	// teststorage.FileBackendOnHARaftSetup(&conf, &opts)
	setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer cluster.Cleanup()

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	// Seal the leader so we can install an address provider
	{
		testhelpers.EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingHAStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	joinFunc := func(client *api.Client, addClientCerts bool) {
		req := &api.RaftJoinRequest{
			LeaderAPIAddr: leaderAPI,
			LeaderCACert:  string(cluster.CACertPEM),
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

	joinFunc(cluster.Cores[1].Client, false)
	joinFunc(cluster.Cores[2].Client, false)

	_, err := cluster.Cores[0].Client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = cluster.Cores[0].Client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-2",
	})
	if err != nil {
		t.Fatal(err)
	}

	joinFunc(cluster.Cores[1].Client, true)
	joinFunc(cluster.Cores[2].Client, true)
}

func TestRaft_HA_ExistingCluster(t *testing.T) {
	t.Skipf("not implemented")
}
