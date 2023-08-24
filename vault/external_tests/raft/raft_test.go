// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rafttests

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/benchhelpers"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

type RaftClusterOpts struct {
	DisableFollowerJoins           bool
	InmemCluster                   bool
	EnableAutopilot                bool
	PhysicalFactoryConfig          map[string]interface{}
	DisablePerfStandby             bool
	EnableResponseHeaderRaftNodeID bool
	NumCores                       int
	Seal                           vault.Seal
	VersionMap                     map[int]string
	RedundancyZoneMap              map[int]string
	EffectiveSDKVersionMap         map[int]string
}

func raftCluster(t testing.TB, ropts *RaftClusterOpts) (*vault.TestCluster, *vault.TestClusterOptions) {
	if ropts == nil {
		ropts = &RaftClusterOpts{}
	}

	conf := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		DisableAutopilot:               !ropts.EnableAutopilot,
		EnableResponseHeaderRaftNodeID: ropts.EnableResponseHeaderRaftNodeID,
		Seal:                           ropts.Seal,
	}

	opts := vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	}
	opts.InmemClusterLayers = ropts.InmemCluster
	opts.PhysicalFactoryConfig = ropts.PhysicalFactoryConfig
	conf.DisablePerformanceStandby = ropts.DisablePerfStandby
	opts.NumCores = ropts.NumCores
	opts.VersionMap = ropts.VersionMap
	opts.RedundancyZoneMap = ropts.RedundancyZoneMap
	opts.EffectiveSDKVersionMap = ropts.EffectiveSDKVersionMap

	teststorage.RaftBackendSetup(conf, &opts)

	if ropts.DisableFollowerJoins {
		opts.SetupFunc = nil
	}

	cluster := vault.NewTestCluster(benchhelpers.TBtoT(t), conf, &opts)
	cluster.Start()
	vault.TestWaitActive(benchhelpers.TBtoT(t), cluster.Cores[0].Core)
	return cluster, &opts
}

func TestRaft_BoltDBMetrics(t *testing.T) {
	t.Parallel()
	conf := vault.CoreConfig{}
	opts := vault.TestClusterOptions{
		HandlerFunc:            vaulthttp.Handler,
		NumCores:               1,
		CoreMetricSinkProvider: testhelpers.TestMetricSinkProvider(time.Minute),
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{
				Telemetry: configutil.ListenerTelemetry{
					UnauthenticatedMetricsAccess: true,
				},
			},
		},
	}

	teststorage.RaftBackendSetup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	leaderClient := cluster.Cores[0].Client

	// Write a few keys
	for i := 0; i < 50; i++ {
		_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			fmt.Sprintf("foo%d", i): fmt.Sprintf("bar%d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Even though there is a long delay between when we start the node and when we check for these metrics,
	// the core metrics loop isn't started until postUnseal, which happens after said delay. This means we
	// need a small artificial delay here as well, otherwise we won't see any metrics emitted.
	time.Sleep(5 * time.Second)
	data, err := testhelpers.SysMetricsReq(leaderClient, cluster, true)
	if err != nil {
		t.Fatal(err)
	}

	noBoltDBMetrics := true
	for _, g := range data.Gauges {
		if strings.HasPrefix(g.Name, "raft_storage.bolt.") {
			noBoltDBMetrics = false
			break
		}
	}

	if noBoltDBMetrics {
		t.Fatal("expected to find boltdb metrics being emitted from the raft backend, but there were none")
	}
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
			TLSConfig: leaderCore.TLSConfig(),
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

	err := testhelpers.VerifyRaftPeers(t, cluster.Cores[0].Client, map[string]bool{
		"core-0": true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRaft_Retry_Join(t *testing.T) {
	t.Parallel()
	var conf vault.CoreConfig
	opts := vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
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
	}

	leaderInfos := []*raft.LeaderJoinInfo{
		{
			LeaderAPIAddr: leaderAPI,
			TLSConfig:     leaderCore.TLSConfig(),
			Retry:         true,
		},
	}

	var wg sync.WaitGroup
	for _, clusterCore := range cluster.Cores[1:] {
		wg.Add(1)
		go func(t *testing.T, core *vault.TestClusterCore) {
			t.Helper()
			defer wg.Done()
			core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
			_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
			if err != nil {
				t.Error(err)
			}

			// Handle potential racy behavior with unseals. Retry the unseal until it succeeds.
			corehelpers.RetryUntil(t, 10*time.Second, func() error {
				return cluster.AttemptUnsealCore(core)
			})
		}(t, clusterCore)
	}

	// Unseal the leader and wait for the other cores to unseal
	cluster.UnsealCore(t, leaderCore)
	wg.Wait()

	vault.TestWaitActive(t, leaderCore.Core)

	corehelpers.RetryUntil(t, 10*time.Second, func() error {
		return testhelpers.VerifyRaftPeers(t, cluster.Cores[0].Client, map[string]bool{
			"core-0": true,
			"core-1": true,
			"core-2": true,
		})
	})
}

func TestRaft_Join(t *testing.T) {
	t.Parallel()
	var conf vault.CoreConfig
	opts := vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
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
	cluster, _ := raftCluster(t, nil)
	defer cluster.Cleanup()

	for i, c := range cluster.Cores {
		if c.Core.Sealed() {
			t.Fatalf("failed to unseal core %d", i)
		}
	}

	client := cluster.Cores[0].Client

	err := testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-2",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
		"core-1": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": "core-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRaft_NodeIDHeader(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description   string
		ropts         *RaftClusterOpts
		headerPresent bool
	}{
		{
			description:   "with no header configured",
			ropts:         nil,
			headerPresent: false,
		},
		{
			description: "with header configured",
			ropts: &RaftClusterOpts{
				EnableResponseHeaderRaftNodeID: true,
			},
			headerPresent: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cluster, _ := raftCluster(t, tc.ropts)
			defer cluster.Cleanup()

			for i, c := range cluster.Cores {
				if c.Core.Sealed() {
					t.Fatalf("failed to unseal core %d", i)
				}

				client := c.Client
				req := client.NewRequest("GET", "/v1/sys/seal-status")
				resp, err := client.RawRequest(req)
				if err != nil {
					t.Fatalf("err: %s", err)
				}
				if resp == nil {
					t.Fatalf("nil response")
				}

				rniHeader := resp.Header.Get("X-Vault-Raft-Node-ID")
				nodeID := c.Core.GetRaftNodeID()

				if tc.headerPresent && rniHeader == "" {
					t.Fatal("missing 'X-Vault-Raft-Node-ID' header entry in response")
				}
				if tc.headerPresent && rniHeader != nodeID {
					t.Fatalf("got the wrong raft node id. expected %s to equal %s", rniHeader, nodeID)
				}
				if !tc.headerPresent && rniHeader != "" {
					t.Fatal("didn't expect 'X-Vault-Raft-Node-ID' header but it was present anyway")
				}
			}
		})
	}
}

func TestRaft_Configuration(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, nil)
	defer cluster.Cleanup()
	Raft_Configuration_Test(t, cluster)
}

func TestRaft_ShamirUnseal(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, nil)
	defer cluster.Cleanup()

	for i, c := range cluster.Cores {
		if c.Core.Sealed() {
			t.Fatalf("failed to unseal core %d", i)
		}
	}
}

func TestRaft_SnapshotAPI(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, nil)
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

	// Take a snapshot
	buf := new(bytes.Buffer)
	err := leaderClient.Sys().RaftSnapshot(buf)
	if err != nil {
		t.Fatal(err)
	}
	snap, err := io.ReadAll(buf)
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
	err = leaderClient.Sys().RaftSnapshotRestore(bytes.NewReader(snap), false)
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

func TestRaft_SnapshotAPI_MidstreamFailure(t *testing.T) {
	// defer goleak.VerifyNone(t)
	t.Parallel()

	seal, setErr := vaultseal.NewToggleableTestSeal(nil)
	autoSeal := vault.NewAutoSeal(seal)
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		NumCores: 1,
		Seal:     autoSeal,
	})
	defer cluster.Cleanup()

	leaderClient := cluster.Cores[0].Client

	// Write a bunch of keys; if too few, the detection code in api.RaftSnapshot
	// will never make it into the tar part, it'll fail merely when trying to
	// decompress the stream.
	for i := 0; i < 1000; i++ {
		_, err := leaderClient.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	r, w := io.Pipe()
	var snap []byte
	var wg sync.WaitGroup
	wg.Add(1)

	var readErr error
	go func() {
		snap, readErr = ioutil.ReadAll(r)
		wg.Done()
	}()

	setErr[0](errors.New("seal failure"))
	// Take a snapshot
	err := leaderClient.Sys().RaftSnapshot(w)
	w.Close()
	if err == nil || err != api.ErrIncompleteSnapshot {
		t.Fatalf("expected err=%v, got: %v", api.ErrIncompleteSnapshot, err)
	}
	wg.Wait()
	if len(snap) == 0 && readErr == nil {
		readErr = errors.New("no bytes read")
	}
	if readErr != nil {
		t.Fatal(readErr)
	}
}

func TestRaft_SnapshotAPI_RekeyRotate_Backward(t *testing.T) {
	t.Parallel()
	type testCase struct {
		Name               string
		Rekey              bool
		Rotate             bool
		DisablePerfStandby bool
	}

	tCases := []testCase{
		{
			Name:               "rekey",
			Rekey:              true,
			Rotate:             false,
			DisablePerfStandby: true,
		},
		{
			Name:               "rotate",
			Rekey:              false,
			Rotate:             true,
			DisablePerfStandby: true,
		},
		{
			Name:               "both",
			Rekey:              true,
			Rotate:             true,
			DisablePerfStandby: true,
		},
	}

	if constants.IsEnterprise {
		tCases = append(tCases, []testCase{
			{
				Name:               "rekey-with-perf-standby",
				Rekey:              true,
				Rotate:             false,
				DisablePerfStandby: false,
			},
			{
				Name:               "rotate-with-perf-standby",
				Rekey:              false,
				Rotate:             true,
				DisablePerfStandby: false,
			},
			{
				Name:               "both-with-perf-standby",
				Rekey:              true,
				Rotate:             true,
				DisablePerfStandby: false,
			},
		}...)
	}

	for _, tCase := range tCases {
		t.Run(tCase.Name, func(t *testing.T) {
			// bind locally
			tCaseLocal := tCase
			t.Parallel()

			cluster, _ := raftCluster(t, &RaftClusterOpts{DisablePerfStandby: tCaseLocal.DisablePerfStandby})
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
			transport.TLSClientConfig = cluster.Cores[0].TLSConfig()
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

				testhelpers.EnsureStableActiveNode(t, cluster)
				testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
			}

			if tCaseLocal.Rekey {
				// Rekey
				cluster.BarrierKeys = testhelpers.RekeyCluster(t, cluster, false)

				testhelpers.EnsureStableActiveNode(t, cluster)
				testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
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
			testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

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
	t.Parallel()
	type testCase struct {
		Name               string
		Rekey              bool
		Rotate             bool
		ShouldSeal         bool
		DisablePerfStandby bool
	}

	tCases := []testCase{
		{
			Name:               "rekey",
			Rekey:              true,
			Rotate:             false,
			ShouldSeal:         false,
			DisablePerfStandby: true,
		},
		{
			Name:   "rotate",
			Rekey:  false,
			Rotate: true,
			// Rotate writes a new master key upgrade using the new term, which
			// we can no longer decrypt. We must seal here.
			ShouldSeal:         true,
			DisablePerfStandby: true,
		},
		{
			Name:   "both",
			Rekey:  true,
			Rotate: true,
			// If we are moving forward and we have rekeyed and rotated there
			// isn't any way to restore the latest keys so expect to seal.
			ShouldSeal:         true,
			DisablePerfStandby: true,
		},
	}

	if constants.IsEnterprise {
		tCases = append(tCases, []testCase{
			{
				Name:               "rekey-with-perf-standby",
				Rekey:              true,
				Rotate:             false,
				ShouldSeal:         false,
				DisablePerfStandby: false,
			},
			{
				Name:   "rotate-with-perf-standby",
				Rekey:  false,
				Rotate: true,
				// Rotate writes a new master key upgrade using the new term, which
				// we can no longer decrypt. We must seal here.
				ShouldSeal:         true,
				DisablePerfStandby: false,
			},
			{
				Name:   "both-with-perf-standby",
				Rekey:  true,
				Rotate: true,
				// If we are moving forward and we have rekeyed and rotated there
				// isn't any way to restore the latest keys so expect to seal.
				ShouldSeal:         true,
				DisablePerfStandby: false,
			},
		}...)
	}

	for _, tCase := range tCases {
		t.Run(tCase.Name, func(t *testing.T) {
			// bind locally
			tCaseLocal := tCase
			t.Parallel()

			cluster, _ := raftCluster(t, &RaftClusterOpts{DisablePerfStandby: tCaseLocal.DisablePerfStandby})
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
			transport.TLSClientConfig = cluster.Cores[0].TLSConfig()
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

				testhelpers.EnsureStableActiveNode(t, cluster)
				testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
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

				if !tCaseLocal.DisablePerfStandby {
					// Without the key upgrade the perf standby nodes will seal and
					// raft will get into a failure state. Make sure we get the
					// cluster back into a healthy state before moving forward.
					testhelpers.WaitForNCoresSealed(t, cluster, 2)
					testhelpers.EnsureCoresUnsealed(t, cluster)
					testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

					active := testhelpers.DeriveActiveCore(t, cluster)
					leaderClient = active.Client
				}
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
			testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
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
				testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

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
	cluster, _ := raftCluster(t, nil)
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
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig()
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
		cluster2, _ := raftCluster(t, nil)
		defer cluster2.Cleanup()

		leaderClient := cluster2.Cores[0].Client

		transport := cleanhttp.DefaultPooledTransport()
		transport.TLSClientConfig = cluster2.Cores[0].TLSConfig()
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
	cluster, _ := raftCluster(b, nil)
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
	opts := vault.TestClusterOptions{HandlerFunc: vaulthttp.Handler}
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
