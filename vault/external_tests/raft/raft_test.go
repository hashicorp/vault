package rafttests

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func raftCluster(t testing.TB) *vault.TestCluster {
	conf := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}
	var opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	teststorage.RaftBackendSetup(conf, &opts)
	cluster := vault.NewTestCluster(t, conf, &opts)
	cluster.Start()
	vault.TestWaitActive(t, cluster.Cores[0].Core)
	return cluster
}

func TestRaft_Autopilot(t *testing.T) {
	t.Parallel()
	cluster := raftCluster(t)
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})
}

func TestRaft_RetryAutoJoin(t *testing.T) {
	t.Parallel()

	var (
		conf vault.CoreConfig

		opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	)

	teststorage.RaftBackendSetup(&conf, &opts)

	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, &conf, &opts)

	cluster.Start()
	defer cluster.Cleanup()

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}
	leaderCore := cluster.Cores[0]
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	{
		testhelpers.EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	leaderInfos := []*raft.LeaderJoinInfo{
		{
			AutoJoin:  "provider=aws region=eu-west-1 tag_key=consul tag_value=tag access_key_id=a secret_access_key=a",
			TLSConfig: leaderCore.TLSConfig,
			Retry:     true,
		},
	}

	{
		// expected to pass but not join as we're not actually discovering leader addresses
		core := cluster.Cores[1]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)

		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		require.NoError(t, err)
	}

	testhelpers.VerifyRaftPeers(t, cluster.Cores[0].Client, map[string]bool{
		"core-0": true,
	})
}

func TestRaft_Retry_Join(t *testing.T) {
	t.Parallel()
	var conf vault.CoreConfig
	var opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	teststorage.RaftBackendSetup(&conf, &opts)
	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer cluster.Cleanup()

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	{
		testhelpers.EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	leaderInfos := []*raft.LeaderJoinInfo{
		&raft.LeaderJoinInfo{
			LeaderAPIAddr: leaderAPI,
			TLSConfig:     leaderCore.TLSConfig,
			Retry:         true,
		},
	}

	{
		core := cluster.Cores[1]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(2 * time.Second)

		cluster.UnsealCore(t, core)
	}

	{
		core := cluster.Cores[2]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(2 * time.Second)

		cluster.UnsealCore(t, core)
	}

	testhelpers.VerifyRaftPeers(t, cluster.Cores[0].Client, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})
}

func TestRaft_Join(t *testing.T) {
	t.Parallel()
	var conf vault.CoreConfig
	var opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	teststorage.RaftBackendSetup(&conf, &opts)
	opts.SetupFunc = nil
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
		leaderCore.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
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

func TestRaft_RemovePeer(t *testing.T) {
	t.Parallel()
	cluster := raftCluster(t)
	defer cluster.Cleanup()

	for i, c := range cluster.Cores {
		if c.Core.Sealed() {
			t.Fatalf("failed to unseal core %d", i)
		}
	}

	client := cluster.Cores[0].Client

	testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})

	_, err := client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-2",
	})
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
		"core-1": true,
	})

	_, err = client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
	})
}

func TestRaft_Configuration(t *testing.T) {
	t.Parallel()
	cluster := raftCluster(t)
	defer cluster.Cleanup()

	for i, c := range cluster.Cores {
		if c.Core.Sealed() {
			t.Fatalf("failed to unseal core %d", i)
		}
	}

	client := cluster.Cores[0].Client
	secret, err := client.Logical().Read("sys/storage/raft/configuration")
	if err != nil {
		t.Fatal(err)
	}
	servers := secret.Data["config"].(map[string]interface{})["servers"].([]interface{})
	expected := map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	}
	if len(servers) != 3 {
		t.Fatalf("incorrect number of servers in the configuration")
	}
	for _, s := range servers {
		server := s.(map[string]interface{})
		nodeID := server["node_id"].(string)
		leader := server["leader"].(bool)
		switch nodeID {
		case "core-0":
			if !leader {
				t.Fatalf("expected server to be leader: %#v", server)
			}
		default:
			if leader {
				t.Fatalf("expected server to not be leader: %#v", server)
			}
		}

		delete(expected, nodeID)
	}
	if len(expected) != 0 {
		t.Fatalf("failed to read configuration successfully")
	}
}

func TestRaft_ShamirUnseal(t *testing.T) {
	t.Parallel()
	cluster := raftCluster(t)
	defer cluster.Cleanup()

	for i, c := range cluster.Cores {
		if c.Core.Sealed() {
			t.Fatalf("failed to unseal core %d", i)
		}
	}
}

func TestRaft_SnapshotAPI(t *testing.T) {
	t.Parallel()
	cluster := raftCluster(t)
	defer cluster.Cleanup()

	leaderClient := cluster.Cores[0].Client

	// Write a few keys
	for i := 0; i < 10; i++ {
		_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: transport,
	}

	// Take a snapshot
	req := leaderClient.NewRequest("GET", "/v1/sys/storage/raft/snapshot")
	httpReq, err := req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	snap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(snap) == 0 {
		t.Fatal("no snapshot returned")
	}

	// Write a few more keys
	for i := 10; i < 20; i++ {
		_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Restore snapshot
	req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot")
	req.Body = bytes.NewBuffer(snap)
	httpReq, err = req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	// List kv to make sure we removed the extra keys
	secret, err := leaderClient.Logical().List("secret/")
	if err != nil {
		t.Fatal(err)
	}

	if len(secret.Data["keys"].([]interface{})) != 10 {
		t.Fatal("snapshot didn't apply correctly")
	}
}

func TestRaft_SnapshotAPI_RekeyRotate_Backward(t *testing.T) {
	tCases := []struct {
		Name   string
		Rekey  bool
		Rotate bool
	}{
		{
			Name:   "rekey",
			Rekey:  true,
			Rotate: false,
		},
		{
			Name:   "rotate",
			Rekey:  false,
			Rotate: true,
		},
		{
			Name:   "both",
			Rekey:  true,
			Rotate: true,
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.Name, func(t *testing.T) {
			// bind locally
			tCaseLocal := tCase
			t.Parallel()

			cluster := raftCluster(t)
			defer cluster.Cleanup()

			leaderClient := cluster.Cores[0].Client

			// Write a few keys
			for i := 0; i < 10; i++ {
				_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
					"test": "data",
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			transport := cleanhttp.DefaultPooledTransport()
			transport.TLSClientConfig = cluster.Cores[0].TLSConfig.Clone()
			if err := http2.ConfigureTransport(transport); err != nil {
				t.Fatal(err)
			}
			client := &http.Client{
				Transport: transport,
			}

			// Take a snapshot
			req := leaderClient.NewRequest("GET", "/v1/sys/storage/raft/snapshot")
			httpReq, err := req.ToHTTP()
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			snap, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(snap) == 0 {
				t.Fatal("no snapshot returned")
			}

			// cache the original barrier keys
			barrierKeys := cluster.BarrierKeys

			if tCaseLocal.Rotate {
				// Rotate
				err = leaderClient.Sys().Rotate()
				if err != nil {
					t.Fatal(err)
				}
			}

			if tCaseLocal.Rekey {
				// Rekey
				cluster.BarrierKeys = testhelpers.RekeyCluster(t, cluster, false)
			}

			if tCaseLocal.Rekey {
				// Restore snapshot, should fail.
				req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot")
				req.Body = bytes.NewBuffer(snap)
				httpReq, err = req.ToHTTP()
				if err != nil {
					t.Fatal(err)
				}
				resp, err = client.Do(httpReq)
				if err != nil {
					t.Fatal(err)
				}
				// Parse Response
				apiResp := api.Response{Response: resp}
				if !strings.Contains(apiResp.Error().Error(), "could not verify hash file, possibly the snapshot is using a different set of unseal keys") {
					t.Fatal(apiResp.Error())
				}
			}

			// Restore snapshot force
			req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot-force")
			req.Body = bytes.NewBuffer(snap)
			httpReq, err = req.ToHTTP()
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}

			testhelpers.EnsureStableActiveNode(t, cluster)

			// Write some data so we can make sure we can read it later. This is testing
			// that we correctly reload the keyring
			_, err = leaderClient.Logical().Write("secret/foo", map[string]interface{}{
				"test": "data",
			})
			if err != nil {
				t.Fatal(err)
			}

			testhelpers.EnsureCoresSealed(t, cluster)

			cluster.BarrierKeys = barrierKeys
			testhelpers.EnsureCoresUnsealed(t, cluster)
			testhelpers.WaitForActiveNode(t, cluster)
			activeCore := testhelpers.DeriveStableActiveCore(t, cluster)

			// Read the value.
			data, err := activeCore.Client.Logical().Read("secret/foo")
			if err != nil {
				t.Fatal(err)
			}
			if data.Data["test"].(string) != "data" {
				t.Fatal(data)
			}
		})
	}
}

func TestRaft_SnapshotAPI_RekeyRotate_Forward(t *testing.T) {
	tCases := []struct {
		Name       string
		Rekey      bool
		Rotate     bool
		ShouldSeal bool
	}{
		{
			Name:       "rekey",
			Rekey:      true,
			Rotate:     false,
			ShouldSeal: false,
		},
		{
			Name:   "rotate",
			Rekey:  false,
			Rotate: true,
			// Rotate writes a new master key upgrade using the new term, which
			// we can no longer decrypt. We must seal here.
			ShouldSeal: true,
		},
		{
			Name:   "both",
			Rekey:  true,
			Rotate: true,
			// If we are moving forward and we have rekeyed and rotated there
			// isn't any way to restore the latest keys so expect to seal.
			ShouldSeal: true,
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.Name, func(t *testing.T) {
			// bind locally
			tCaseLocal := tCase
			t.Parallel()

			cluster := raftCluster(t)
			defer cluster.Cleanup()

			leaderClient := cluster.Cores[0].Client

			// Write a few keys
			for i := 0; i < 10; i++ {
				_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
					"test": "data",
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			transport := cleanhttp.DefaultPooledTransport()
			transport.TLSClientConfig = cluster.Cores[0].TLSConfig.Clone()
			if err := http2.ConfigureTransport(transport); err != nil {
				t.Fatal(err)
			}
			client := &http.Client{
				Transport: transport,
			}

			// Take a snapshot
			req := leaderClient.NewRequest("GET", "/v1/sys/storage/raft/snapshot")
			httpReq, err := req.ToHTTP()
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}

			snap, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if len(snap) == 0 {
				t.Fatal("no snapshot returned")
			}

			if tCaseLocal.Rekey {
				// Rekey
				cluster.BarrierKeys = testhelpers.RekeyCluster(t, cluster, false)
			}
			if tCaseLocal.Rotate {
				// Set the key clean up to 0 so it's cleaned immediately. This
				// will simulate that there are no ways to upgrade to the latest
				// term.
				for _, c := range cluster.Cores {
					c.Core.SetKeyRotateGracePeriod(0)
				}

				// Rotate
				err = leaderClient.Sys().Rotate()
				if err != nil {
					t.Fatal(err)
				}
				// Let the key upgrade get deleted
				time.Sleep(1 * time.Second)
			}

			// cache the new barrier keys
			newBarrierKeys := cluster.BarrierKeys

			// Take another snapshot for later use in "jumping" forward
			req = leaderClient.NewRequest("GET", "/v1/sys/storage/raft/snapshot")
			httpReq, err = req.ToHTTP()
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}

			snap2, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if len(snap2) == 0 {
				t.Fatal("no snapshot returned")
			}

			// Restore snapshot to move us back in time so we can test going
			// forward
			req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot-force")
			req.Body = bytes.NewBuffer(snap)
			httpReq, err = req.ToHTTP()
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}

			testhelpers.EnsureStableActiveNode(t, cluster)
			if tCaseLocal.Rekey {
				// Restore snapshot, should fail.
				req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot")
				req.Body = bytes.NewBuffer(snap2)
				httpReq, err = req.ToHTTP()
				if err != nil {
					t.Fatal(err)
				}
				resp, err = client.Do(httpReq)
				if err != nil {
					t.Fatal(err)
				}
				// Parse Response
				apiResp := api.Response{Response: resp}
				if apiResp.Error() == nil || !strings.Contains(apiResp.Error().Error(), "could not verify hash file, possibly the snapshot is using a different set of unseal keys") {
					t.Fatalf("expected error verifying hash file, got %v", apiResp.Error())
				}
			}

			// Restore snapshot force
			req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot-force")
			req.Body = bytes.NewBuffer(snap2)
			httpReq, err = req.ToHTTP()
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}

			switch tCaseLocal.ShouldSeal {
			case true:
				testhelpers.WaitForNCoresSealed(t, cluster, 3)

			case false:
				testhelpers.EnsureStableActiveNode(t, cluster)

				// Write some data so we can make sure we can read it later. This is testing
				// that we correctly reload the keyring
				_, err = leaderClient.Logical().Write("secret/foo", map[string]interface{}{
					"test": "data",
				})
				if err != nil {
					t.Fatal(err)
				}

				testhelpers.EnsureCoresSealed(t, cluster)

				cluster.BarrierKeys = newBarrierKeys
				testhelpers.EnsureCoresUnsealed(t, cluster)
				testhelpers.WaitForActiveNode(t, cluster)
				activeCore := testhelpers.DeriveStableActiveCore(t, cluster)

				// Read the value.
				data, err := activeCore.Client.Logical().Read("secret/foo")
				if err != nil {
					t.Fatal(err)
				}
				if data.Data["test"].(string) != "data" {
					t.Fatal(data)
				}
			}
		})
	}
}

func TestRaft_SnapshotAPI_DifferentCluster(t *testing.T) {
	t.Parallel()
	cluster := raftCluster(t)
	defer cluster.Cleanup()

	leaderClient := cluster.Cores[0].Client

	// Write a few keys
	for i := 0; i < 10; i++ {
		_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: transport,
	}

	// Take a snapshot
	req := leaderClient.NewRequest("GET", "/v1/sys/storage/raft/snapshot")
	httpReq, err := req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	snap, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if len(snap) == 0 {
		t.Fatal("no snapshot returned")
	}

	// Cluster 2
	{
		cluster2 := raftCluster(t)
		defer cluster2.Cleanup()

		leaderClient := cluster2.Cores[0].Client

		transport := cleanhttp.DefaultPooledTransport()
		transport.TLSClientConfig = cluster2.Cores[0].TLSConfig.Clone()
		if err := http2.ConfigureTransport(transport); err != nil {
			t.Fatal(err)
		}
		client := &http.Client{
			Transport: transport,
		}
		// Restore snapshot, should fail.
		req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot")
		req.Body = bytes.NewBuffer(snap)
		httpReq, err = req.ToHTTP()
		if err != nil {
			t.Fatal(err)
		}
		resp, err = client.Do(httpReq)
		if err != nil {
			t.Fatal(err)
		}
		// Parse Response
		apiResp := api.Response{Response: resp}
		if !strings.Contains(apiResp.Error().Error(), "could not verify hash file, possibly the snapshot is using a different set of unseal keys") {
			t.Fatal(apiResp.Error())
		}

		// Restore snapshot force
		req = leaderClient.NewRequest("POST", "/v1/sys/storage/raft/snapshot-force")
		req.Body = bytes.NewBuffer(snap)
		httpReq, err = req.ToHTTP()
		if err != nil {
			t.Fatal(err)
		}
		resp, err = client.Do(httpReq)
		if err != nil {
			t.Fatal(err)
		}

		testhelpers.WaitForNCoresSealed(t, cluster2, 3)
	}
}

func BenchmarkRaft_SingleNode(b *testing.B) {
	cluster := raftCluster(b)
	defer cluster.Cleanup()

	leaderClient := cluster.Cores[0].Client

	bench := func(b *testing.B, dataSize int) {
		data, err := uuid.GenerateRandomBytes(dataSize)
		if err != nil {
			b.Fatal(err)
		}

		testName := b.Name()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("secret/%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			_, err := leaderClient.Logical().Write(key, map[string]interface{}{
				"test": data,
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("256b", func(b *testing.B) { bench(b, 25) })
}

func TestRaft_Join_InitStatus(t *testing.T) {
	t.Parallel()
	var conf vault.CoreConfig
	var opts = vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
	teststorage.RaftBackendSetup(&conf, &opts)
	opts.SetupFunc = nil
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
		leaderCore.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	joinFunc := func(client *api.Client) {
		req := &api.RaftJoinRequest{
			LeaderAPIAddr: leaderAPI,
			LeaderCACert:  string(cluster.CACertPEM),
		}
		resp, err := client.Sys().RaftJoin(req)
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Joined {
			t.Fatalf("failed to join raft cluster")
		}
	}

	verifyInitStatus := func(coreIdx int, expected bool) {
		t.Helper()
		client := cluster.Cores[coreIdx].Client

		initialized, err := client.Sys().InitStatus()
		if err != nil {
			t.Fatal(err)
		}

		if initialized != expected {
			t.Errorf("core %d: expected init=%v, sys/init returned %v", coreIdx, expected, initialized)
		}

		status, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}

		if status.Initialized != expected {
			t.Errorf("core %d: expected init=%v, sys/seal-status returned %v", coreIdx, expected, status.Initialized)
		}

		health, err := client.Sys().Health()
		if err != nil {
			t.Fatal(err)
		}
		if health.Initialized != expected {
			t.Errorf("core %d: expected init=%v, sys/health returned %v", coreIdx, expected, health.Initialized)
		}
	}

	for i := range cluster.Cores {
		verifyInitStatus(i, i < 1)
	}

	joinFunc(cluster.Cores[1].Client)
	for i, core := range cluster.Cores {
		verifyInitStatus(i, i < 2)
		if i == 1 {
			cluster.UnsealCore(t, core)
			verifyInitStatus(i, true)
		}
	}

	joinFunc(cluster.Cores[2].Client)
	for i, core := range cluster.Cores {
		verifyInitStatus(i, true)
		if i == 2 {
			cluster.UnsealCore(t, core)
			verifyInitStatus(i, true)
		}
	}

	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
	for i := range cluster.Cores {
		verifyInitStatus(i, true)
	}
}
