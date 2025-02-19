// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rafttests

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cleanhttp"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/cluster"
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
	PerNodePhysicalFactoryConfig   map[int]map[string]interface{}
}

func raftClusterBuilder(t testing.TB, ropts *RaftClusterOpts) (*vault.CoreConfig, vault.TestClusterOptions) {
	if ropts == nil {
		ropts = &RaftClusterOpts{
			InmemCluster: true,
		}
	}

	conf := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		DisableAutopilot:               !ropts.EnableAutopilot,
		EnableResponseHeaderRaftNodeID: ropts.EnableResponseHeaderRaftNodeID,
		Seal:                           ropts.Seal,
		EnableRaw:                      true,
	}

	opts := vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	}
	opts.InmemClusterLayers = ropts.InmemCluster
	opts.PhysicalFactoryConfig = ropts.PhysicalFactoryConfig
	conf.DisablePerformanceStandby = ropts.DisablePerfStandby
	opts.NumCores = ropts.NumCores
	opts.EffectiveSDKVersionMap = ropts.EffectiveSDKVersionMap
	opts.PerNodePhysicalFactoryConfig = ropts.PerNodePhysicalFactoryConfig
	if len(ropts.VersionMap) > 0 || len(ropts.RedundancyZoneMap) > 0 {
		if opts.PerNodePhysicalFactoryConfig == nil {
			opts.PerNodePhysicalFactoryConfig = map[int]map[string]interface{}{}
		}
		for idx, ver := range ropts.VersionMap {
			if opts.PerNodePhysicalFactoryConfig[idx] == nil {
				opts.PerNodePhysicalFactoryConfig[idx] = map[string]interface{}{}
			}
			opts.PerNodePhysicalFactoryConfig[idx]["autopilot_upgrade_version"] = ver
		}
		for idx, zone := range ropts.RedundancyZoneMap {
			if opts.PerNodePhysicalFactoryConfig[idx] == nil {
				opts.PerNodePhysicalFactoryConfig[idx] = map[string]interface{}{}
			}
			opts.PerNodePhysicalFactoryConfig[idx]["autopilot_redundancy_zone"] = zone
		}
	}

	teststorage.RaftBackendSetup(conf, &opts)

	if ropts.DisableFollowerJoins {
		opts.SetupFunc = nil
	}
	return conf, opts
}

func raftCluster(t testing.TB, ropts *RaftClusterOpts) (*vault.TestCluster, *vault.TestClusterOptions) {
	conf, opts := raftClusterBuilder(t, ropts)
	cluster := vault.NewTestCluster(t, conf, &opts)
	vault.TestWaitActive(t, cluster.Cores[0].Core)
	return cluster, &opts
}

func TestRaft_BoltDBMetrics(t *testing.T) {
	t.Parallel()
	conf, opts := raftClusterBuilder(t, &RaftClusterOpts{
		InmemCluster: true,
		NumCores:     1,
	})
	opts.CoreMetricSinkProvider = testhelpers.TestMetricSinkProvider(time.Minute)
	opts.DefaultHandlerProperties = vault.HandlerProperties{
		ListenerConfig: &configutil.Listener{
			Telemetry: configutil.ListenerTelemetry{
				UnauthenticatedMetricsAccess: true,
			},
		},
	}
	cluster := vault.NewTestCluster(t, conf, &opts)
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

	cluster, _ := raftCluster(t, &RaftClusterOpts{
		InmemCluster:         true,
		DisableFollowerJoins: true,
	})
	defer cluster.Cleanup()

	leaderCore := cluster.Cores[0]

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
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		InmemCluster:         true,
		DisableFollowerJoins: true,
	})
	defer cluster.Cleanup()

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()

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

// TestRaftChallenge_sameAnswerSameID_concurrent verifies that 10 goroutines
// all requesting a raft challenge with the same ID all return the same answer.
// This is a regression test for a TOCTTOU race found during testing.
func TestRaftChallenge_sameAnswerSameID_concurrent(t *testing.T) {
	t.Parallel()

	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		NumCores:             1,
	})
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	challenges := make(chan string, 15)
	wg := sync.WaitGroup{}
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
				"server_id": "node1",
			})
			require.NoError(t, err)
			challenges <- res.Data["challenge"].(string)
		}()
	}

	wg.Wait()
	challengeSet := make(map[string]struct{})
	close(challenges)
	for challenge := range challenges {
		challengeSet[challenge] = struct{}{}
	}

	require.Len(t, challengeSet, 1)
}

// TestRaftChallenge_sameAnswerSameID verifies that repeated bootstrap requests
// with the same node ID return the same challenge, but that a different node ID
// returns a different challenge
func TestRaftChallenge_sameAnswerSameID(t *testing.T) {
	t.Parallel()

	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		NumCores:             1,
	})
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client
	res, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
		"server_id": "node1",
	})
	require.NoError(t, err)

	// querying the same ID returns the same challenge
	challenge := res.Data["challenge"]
	resSameID, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
		"server_id": "node1",
	})
	require.NoError(t, err)
	require.Equal(t, challenge, resSameID.Data["challenge"])

	// querying a different ID returns a new challenge
	resDiffID, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
		"server_id": "node2",
	})
	require.NoError(t, err)
	require.NotEqual(t, challenge, resDiffID.Data["challenge"])
}

// TestRaftChallenge_evicted verifies that a valid answer errors if there have
// been more than 20 challenge requests after it, because our cache of pending
// bootstraps is limited to 20
func TestRaftChallenge_evicted(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		NumCores:             1,
	})
	defer cluster.Cleanup()
	firstResponse := map[string]interface{}{}
	client := cluster.Cores[0].Client
	for i := 0; i < vault.RaftInitialChallengeLimit+1; i++ {
		if i == vault.RaftInitialChallengeLimit {
			// wait before sending the last request, so we don't get rate
			// limited
			time.Sleep(2 * time.Second)
		}
		res, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
			"server_id": fmt.Sprintf("node-%d", i),
		})
		require.NoError(t, err)

		// save the response from the first challenge
		if i == 0 {
			firstResponse = res.Data
		}
	}

	// get the answer to the challenge
	challengeRaw, err := base64.StdEncoding.DecodeString(firstResponse["challenge"].(string))
	require.NoError(t, err)
	eBlob := &wrapping.BlobInfo{}
	err = proto.Unmarshal(challengeRaw, eBlob)
	require.NoError(t, err)
	access := cluster.Cores[0].SealAccess().GetAccess()
	multiWrapValue := &vaultseal.MultiWrapValue{
		Generation: access.Generation(),
		Slots:      []*wrapping.BlobInfo{eBlob},
	}
	plaintext, _, err := access.Decrypt(context.Background(), multiWrapValue)
	require.NoError(t, err)

	// send the answer
	_, err = client.Logical().Write("sys/storage/raft/bootstrap/answer", map[string]interface{}{
		"answer":          base64.StdEncoding.EncodeToString(plaintext),
		"server_id":       "node-0",
		"cluster_addr":    "127.0.0.1:8200",
		"sdk_version":     "1.1.1",
		"upgrade_version": "1.2.3",
		"non_voter":       false,
	})

	require.ErrorContains(t, err, "no expected answer for the server id provided")
}

// TestRaft_ChallengeSpam creates 40 raft bootstrap challenges. The first 20
// should succeed. After 20 challenges have been created, slow down the requests
// so that there are 2.5 occurring per second. Some of these will fail, due to
// rate limiting, but others will succeed.
func TestRaft_ChallengeSpam(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
	})
	defer cluster.Cleanup()

	// Execute 2 * MaxInFlightRequests, over a period that should allow some to proceed as the token bucket
	// refills.
	var someLaterFailed bool
	var someLaterSucceeded bool
	for n := 0; n < 2*vault.RaftInitialChallengeLimit; n++ {
		_, err := cluster.Cores[0].Client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
			"server_id": fmt.Sprintf("core-%d", n),
		})
		// First MaxInFlightRequests should succeed for sure
		if n < vault.RaftInitialChallengeLimit {
			require.NoError(t, err)
		} else {
			// slow down to twice the configured rps
			time.Sleep((1000 * time.Millisecond) / (2 * time.Duration(vault.RaftChallengesPerSecond)))
			if err != nil {
				require.Equal(t, 429, err.(*api.ResponseError).StatusCode)
				someLaterFailed = true
			} else {
				someLaterSucceeded = true
			}
		}
	}
	require.True(t, someLaterFailed)
	require.True(t, someLaterSucceeded)
}

func TestRaft_Join(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
	})
	defer cluster.Cleanup()

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()

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
				InmemCluster:                   true,
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

	seal, wrappers := vaultseal.NewTestSeal(nil)
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
		snap, readErr = io.ReadAll(r)
		wg.Done()
	}()

	wrappers[0].SetError(errors.New("seal failure"))
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

			cluster, _ := raftCluster(t, &RaftClusterOpts{
				DisablePerfStandby: tCaseLocal.DisablePerfStandby,
			})
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

			snap, err := io.ReadAll(resp.Body)
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

			cluster, _ := raftCluster(t, &RaftClusterOpts{
				DisablePerfStandby: tCaseLocal.DisablePerfStandby,
			})
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

			snap, err := io.ReadAll(resp.Body)
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

			snap2, err := io.ReadAll(resp.Body)
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

	snap, err := io.ReadAll(resp.Body)
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

	cluster, _ := raftCluster(t, &RaftClusterOpts{
		InmemCluster:         true,
		DisableFollowerJoins: true,
	})
	defer cluster.Cleanup()

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()

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

// TestRaftCluster_Removed creates a 3 node raft cluster and then removes one of
// the nodes. The test verifies that a write on the removed node errors, and that
// the removed node is sealed.
func TestRaftCluster_Removed(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, nil)
	defer cluster.Cleanup()

	follower := cluster.Cores[2]
	followerClient := follower.Client
	_, err := followerClient.Logical().Write("secret/foo", map[string]interface{}{
		"test": "data",
	})
	require.NoError(t, err)

	leaderClient := cluster.Cores[0].Client
	_, err = leaderClient.Logical().Write("/sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": follower.NodeID,
	})
	require.NoError(t, err)
	followerClient.SetCheckRedirect(func(request *http.Request, requests []*http.Request) error {
		require.Fail(t, "request caused a redirect", request.URL.Path)
		return fmt.Errorf("no redirects allowed")
	})
	configChanged := func() bool {
		config, err := leaderClient.Logical().Read("sys/storage/raft/configuration")
		require.NoError(t, err)
		cfg := config.Data["config"].(map[string]interface{})
		servers := cfg["servers"].([]interface{})
		return len(servers) == 2
	}
	// raft config changes happen async, so block until the config change is
	// applied
	require.Eventually(t, configChanged, 3*time.Second, 50*time.Millisecond)

	_, err = followerClient.Logical().Write("secret/foo", map[string]interface{}{
		"test": "other_data",
	})
	require.Error(t, err)
	require.Eventually(t, follower.Sealed, 3*time.Second, 250*time.Millisecond)
}

// TestRaftCluster_Removed_RaftConfig creates a 3 node raft cluster with an extremely long
// heartbeat interval, and then removes one of the nodes. The test verifies that
// removed node discovers that it has been removed (via not being present in the
// raft config) and seals.
func TestRaftCluster_Removed_RaftConfig(t *testing.T) {
	t.Parallel()
	conf, opts := raftClusterBuilder(t, nil)
	conf.ClusterHeartbeatInterval = 5 * time.Minute
	cluster := vault.NewTestCluster(t, conf, &opts)
	vault.TestWaitActive(t, cluster.Cores[0].Core)

	follower := cluster.Cores[2]
	followerClient := follower.Client
	_, err := followerClient.Logical().Write("secret/foo", map[string]interface{}{
		"test": "data",
	})
	require.NoError(t, err)

	_, err = cluster.Cores[0].Client.Logical().Write("/sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": follower.NodeID,
	})
	require.Eventually(t, follower.Sealed, 10*time.Second, 500*time.Millisecond)
}

// TestSysHealth_Raft creates a raft cluster and verifies that the health status
// is OK for a healthy follower. The test partitions one of the nodes so that it
// can't send request forwarding RPCs. The test verifies that the status
// endpoint  shows that HA isn't healthy. Finally, the test removes the
// partitioned follower and unpartitions it. The follower will learn that it has
// been removed, and should return the removed status.
func TestSysHealth_Raft(t *testing.T) {
	parseHealthBody := func(t *testing.T, resp *api.Response) *vaulthttp.HealthResponse {
		t.Helper()
		health := vaulthttp.HealthResponse{}
		defer resp.Body.Close()
		require.NoError(t, jsonutil.DecodeJSONFromReader(resp.Body, &health))
		return &health
	}

	opts := &vault.TestClusterOptions{
		HandlerFunc:        vaulthttp.Handler,
		NumCores:           3,
		InmemClusterLayers: true,
	}
	heartbeat := 500 * time.Millisecond
	teststorage.RaftBackendSetup(nil, opts)
	conf := &vault.CoreConfig{
		ClusterHeartbeatInterval: heartbeat,
	}
	vaultCluster := vault.NewTestCluster(t, conf, opts)
	defer vaultCluster.Cleanup()
	testhelpers.WaitForActiveNodeAndStandbys(t, vaultCluster)
	followerClient := vaultCluster.Cores[1].Client

	t.Run("healthy", func(t *testing.T) {
		resp, err := followerClient.Logical().ReadRawWithData("sys/health", map[string][]string{
			"perfstandbyok": {"true"},
			"standbyok":     {"true"},
		})
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, 200)
		r := parseHealthBody(t, resp)
		require.False(t, *r.RemovedFromCluster)
		require.True(t, *r.HAConnectionHealthy)
		require.Less(t, r.LastRequestForwardingHeartbeatMillis, 2*heartbeat.Milliseconds())
	})
	nl := vaultCluster.Cores[1].NetworkLayer()
	inmem, ok := nl.(*cluster.InmemLayer)
	require.True(t, ok)
	unpartition := inmem.Partition()

	t.Run("partition", func(t *testing.T) {
		time.Sleep(2 * heartbeat)
		var erroredResponse *api.Response
		// the node isn't able to send/receive heartbeats, so it will have
		// haunhealthy status.
		testhelpers.RetryUntil(t, 3*time.Second, func() error {
			resp, err := followerClient.Logical().ReadRawWithData("sys/health", map[string][]string{
				"perfstandbyok": {"true"},
				"standbyok":     {"true"},
			})
			if err == nil {
				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
				return errors.New("expected error")
			}
			if resp.StatusCode != 474 {
				resp.Body.Close()
				return fmt.Errorf("status code %d", resp.StatusCode)
			}
			erroredResponse = resp
			return nil
		})
		r := parseHealthBody(t, erroredResponse)
		require.False(t, *r.RemovedFromCluster)
		require.False(t, *r.HAConnectionHealthy)
		require.Greater(t, r.LastRequestForwardingHeartbeatMillis, 2*heartbeat.Milliseconds())

		// ensure haunhealthycode is respected
		resp, err := followerClient.Logical().ReadRawWithData("sys/health", map[string][]string{
			"perfstandbyok":   {"true"},
			"standbyok":       {"true"},
			"haunhealthycode": {"299"},
		})
		require.NoError(t, err)
		require.Equal(t, 299, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("remove and unpartition", func(t *testing.T) {
		leaderClient := vaultCluster.Cores[0].Client
		_, err := leaderClient.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
			"server_id": vaultCluster.Cores[1].NodeID,
		})
		require.NoError(t, err)
		unpartition()

		var erroredResponse *api.Response

		// now that the node can connect again, it will start getting the removed
		// error when trying to connect. The code should be removed
		testhelpers.RetryUntil(t, 10*time.Second, func() error {
			resp, err := followerClient.Logical().ReadRawWithData("sys/health", map[string][]string{
				"perfstandbyok": {"true"},
				"standbyok":     {"true"},
			})
			if err == nil {
				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
				return fmt.Errorf("expected error")
			}
			if resp.StatusCode != 530 {
				resp.Body.Close()
				return fmt.Errorf("status code %d", resp.StatusCode)
			}
			erroredResponse = resp
			return nil
		})
		r := parseHealthBody(t, erroredResponse)
		require.True(t, true, *r.RemovedFromCluster)
		// The HA connection health should either be nil or false. It's possible
		// for it to be false if we got the response in between the node marking
		// itself removed and sealing
		if r.HAConnectionHealthy != nil {
			require.False(t, *r.HAConnectionHealthy)
		}
	})
}

// TestRaftCluster_Removed_ReAdd creates a three node raft cluster and then
// removes one of the nodes. The removed follower tries to re-join, and the test
// verifies that it errors and cannot join.
func TestRaftCluster_Removed_ReAdd(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, nil)
	defer cluster.Cleanup()

	leader := cluster.Cores[0]
	follower := cluster.Cores[2]

	_, err := leader.Client.Logical().Write("/sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": follower.NodeID,
	})
	require.NoError(t, err)
	require.Eventually(t, follower.Sealed, 10*time.Second, 250*time.Millisecond)

	joinReq := &api.RaftJoinRequest{LeaderAPIAddr: leader.Address.String()}
	_, err = follower.Client.Sys().RaftJoin(joinReq)
	require.Error(t, err)
}
