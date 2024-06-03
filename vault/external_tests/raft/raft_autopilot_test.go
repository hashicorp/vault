// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rafttests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	autopilot "github.com/hashicorp/raft-autopilot"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestRaft_Autopilot_Disable(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		// Not setting EnableAutopilot here.
	})
	defer cluster.Cleanup()

	cli := cluster.Cores[0].Client
	state, err := cli.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Nil(t, nil, state)
}

// TestRaft_Autopilot_Stabilization_And_State verifies that nodes get promoted
// to be voters after the stabilization time has elapsed.  Also checks that
// the autopilot state is Healthy once all nodes are available.
func TestRaft_Autopilot_Stabilization_And_State(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		PhysicalFactoryConfig: map[string]interface{}{
			"performance_multiplier": "5",
		},
	})
	defer cluster.Cleanup()

	// Check that autopilot execution state is running
	client := cluster.Cores[0].Client
	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, true, state.Healthy)
	require.Len(t, state.Servers, 1)
	require.Equal(t, "core-0", state.Servers["core-0"].ID)
	require.Equal(t, "alive", state.Servers["core-0"].NodeStatus)
	require.Equal(t, "leader", state.Servers["core-0"].Status)

	writeConfig := func(config map[string]interface{}, expectError bool) {
		resp, err := client.Logical().Write("sys/storage/raft/autopilot/configuration", config)
		if expectError {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)
		require.Nil(t, resp)
	}

	writableConfig := map[string]interface{}{
		"last_contact_threshold":    "5s",
		"max_trailing_logs":         100,
		"server_stabilization_time": "10s",
	}
	writeConfig(writableConfig, false)

	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	stabilizationKickOffWaitDuration := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(stabilizationKickOffWaitDuration)

	joinAndStabilize(t, cluster.Cores[1], client, cluster, config, "core-1", 2)
	waitUntilVoter(t, 2*autopilot.DefaultReconcileInterval, client, "core-1")
	joinAndStabilize(t, cluster.Cores[2], client, cluster, config, "core-2", 3)
	waitUntilVoter(t, 2*autopilot.DefaultReconcileInterval, client, "core-2")
	state, err = client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, []string{"core-0", "core-1", "core-2"}, state.Voters)

	// Now make sure that after we seal and unseal a node, the current leader
	// remains leader, and that the cluster becomes healthy again.
	leader := state.Leader
	testhelpers.EnsureCoreSealed(t, cluster.Cores[1])
	time.Sleep(10 * time.Second)
	testhelpers.EnsureCoreUnsealed(t, cluster, cluster.Cores[1])

	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		state, err = client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		if state.Healthy && state.Leader == leader {
			break
		}
		time.Sleep(time.Second)
	}
	require.Equal(t, true, state.Healthy)
	require.Equal(t, leader, state.Leader)
}

func TestRaft_Autopilot_Configuration(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	configCheckFunc := func(config *api.AutopilotConfig) {
		conf, err := client.Sys().RaftAutopilotConfiguration()
		require.NoError(t, err)
		require.Equal(t, config, conf)
	}

	writeConfigFunc := func(config map[string]interface{}, expectError bool) {
		resp, err := client.Logical().Write("sys/storage/raft/autopilot/configuration", config)
		if expectError {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)
		require.Nil(t, resp)
	}

	// Ensure autopilot's default config has taken effect
	config := &api.AutopilotConfig{
		CleanupDeadServers:             false,
		DeadServerLastContactThreshold: 24 * time.Hour,
		LastContactThreshold:           10 * time.Second,
		MaxTrailingLogs:                1000,
		ServerStabilizationTime:        10 * time.Second,
	}
	configCheckFunc(config)

	// Update config
	writableConfig := map[string]interface{}{
		"cleanup_dead_servers":               true,
		"dead_server_last_contact_threshold": "100h",
		"last_contact_threshold":             "100s",
		"max_trailing_logs":                  100,
		"min_quorum":                         100,
		"server_stabilization_time":          "100s",
	}
	writeConfigFunc(writableConfig, false)

	// Ensure update has taken effect
	config.CleanupDeadServers = true
	config.DeadServerLastContactThreshold = 100 * time.Hour
	config.LastContactThreshold = 100 * time.Second
	config.MaxTrailingLogs = 100
	config.MinQuorum = 100
	config.ServerStabilizationTime = 100 * time.Second
	configCheckFunc(config)

	// Update some fields and leave the rest as it is.
	writableConfig = map[string]interface{}{
		"dead_server_last_contact_threshold": "50h",
		"max_trailing_logs":                  50,
		"server_stabilization_time":          "50s",
	}
	writeConfigFunc(writableConfig, false)

	// Check update
	config.DeadServerLastContactThreshold = 50 * time.Hour
	config.MaxTrailingLogs = 50
	config.ServerStabilizationTime = 50 * time.Second
	configCheckFunc(config)

	// Check error case
	writableConfig = map[string]interface{}{
		"min_quorum":                         2,
		"dead_server_last_contact_threshold": "48h",
	}
	writeConfigFunc(writableConfig, true)
	configCheckFunc(config)

	// Check dead server last contact threshold minimum
	writableConfig = map[string]interface{}{
		"cleanup_dead_servers":               true,
		"dead_server_last_contact_threshold": "5s",
	}
	writeConfigFunc(writableConfig, true)
	configCheckFunc(config)

	// Ensure that the configuration stays across reboots
	leaderCore := cluster.Cores[0]
	testhelpers.EnsureCoreSealed(t, cluster.Cores[0])
	cluster.UnsealCore(t, leaderCore)
	vault.TestWaitActive(t, leaderCore.Core)
	configCheckFunc(config)
}

// TestRaft_Autopilot_Stabilization_Delay verifies that if a node takes a long
// time to become ready, it doesn't get promoted to voter until then.
func TestRaft_Autopilot_Stabilization_Delay(t *testing.T) {
	t.Parallel()
	core2SnapshotDelay := 5 * time.Second
	conf, opts := raftClusterBuilder(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		PhysicalFactoryConfig: map[string]interface{}{
			"trailing_logs": "10",
		},
		PerNodePhysicalFactoryConfig: map[int]map[string]interface{}{
			2: {
				"snapshot_delay": core2SnapshotDelay.String(),
			},
		},
	})

	cluster := vault.NewTestCluster(t, conf, &opts)
	defer cluster.Cleanup()
	testhelpers.WaitForActiveNode(t, cluster)

	// Check that autopilot execution state is running
	client := cluster.Cores[0].Client
	state, err := client.Sys().RaftAutopilotState()
	require.NotNil(t, state)
	require.NoError(t, err)
	require.Equal(t, true, state.Healthy)
	require.Len(t, state.Servers, 1)
	require.Equal(t, "core-0", state.Servers["core-0"].ID)
	require.Equal(t, "alive", state.Servers["core-0"].NodeStatus)
	require.Equal(t, "leader", state.Servers["core-0"].Status)

	_, err = client.Logical().Write("sys/storage/raft/autopilot/configuration", map[string]interface{}{
		"server_stabilization_time": "5s",
	})
	require.NoError(t, err)

	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	stabilizationPadded := time.Duration(math.Ceil(1.25 * float64(config.ServerStabilizationTime)))
	time.Sleep(stabilizationPadded)

	cli := cluster.Cores[0].Client
	// Write more keys than snapshot_threshold
	for i := 0; i < 50; i++ {
		_, err := cli.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Take a snpashot, which should compact the raft log db, which should prevent
	// followers from getting logs and require that they instead apply a snapshot,
	// which should allow our snapshot_delay to come into play, which should result
	// in core2 coming online slower.
	err = client.Sys().RaftSnapshot(io.Discard)
	require.NoError(t, err)

	joinAndUnseal(t, cluster.Cores[1], cluster, false, false)
	joinAndUnseal(t, cluster.Cores[2], cluster, false, false)

	// Add an extra fudge factor, since once the snapshot delay completes it can
	// take time for the snapshot to actually be applied.
	core2shouldBeHealthyAt := time.Now().Add(core2SnapshotDelay).Add(stabilizationPadded).Add(5 * time.Second)

	// Wait for enough time for stabilization to complete if things were good
	// - but they're not good, due to our snapshot_delay.  So we fail if both
	// nodes are healthy.
	testhelpers.RetryUntil(t, stabilizationPadded, func() error {
		state, err := client.Sys().RaftAutopilotState()
		if err != nil {
			return err
		}
		core1healthy := state.Servers["core-1"] != nil && state.Servers["core-1"].Healthy
		core2healthy := state.Servers["core-2"] != nil && state.Servers["core-2"].Healthy

		if !core1healthy || core2healthy {
			return fmt.Errorf("expected health: core1=true and core2=false, got: core1=%v, core2=%v", core1healthy, core2healthy)
		}

		if diff := cmp.Diff(state.Voters, []string{"core-0", "core-1"}); len(diff) > 0 {
			return fmt.Errorf("expected core-0 and core-1 as voters, diff: %v", diff)
		}

		return nil
	})

	// Now we expect that after the snapshot_delay has elapsed, and enough
	// stabilization time subsequent to that has occurred, that autopilot will
	// deem core-2 healthy and promote it as a voter.
	testhelpers.RetryUntil(t, core2shouldBeHealthyAt.Sub(time.Now()), func() error {
		state, err = client.Sys().RaftAutopilotState()
		if err != nil {
			return err
		}
		if diff := cmp.Diff(state.Voters, []string{"core-0", "core-1", "core-2"}); len(diff) > 0 {
			return fmt.Errorf("expected all nodes as voters, diff: %v", diff)
		}

		return nil
	})
}

func TestRaft_AutoPilot_Peersets_Equivalent(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		InmemCluster:         true,
		EnableAutopilot:      true,
		DisableFollowerJoins: true,
	})
	defer cluster.Cleanup()
	testhelpers.WaitForActiveNode(t, cluster)

	// Create a very large stabilization time so we can test the state between
	// joining and promotions
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/storage/raft/autopilot/configuration", map[string]interface{}{
		"server_stabilization_time": "1h",
	})
	require.NoError(t, err)

	joinAsVoterAndUnseal(t, cluster.Cores[1], cluster)
	joinAsVoterAndUnseal(t, cluster.Cores[2], cluster)

	deadline := time.Now().Add(10 * time.Second)
	var core0Peers, core1Peers, core2Peers []raft.Peer
	for time.Now().Before(deadline) {
		// Make sure all nodes have an equivalent configuration
		core0Peers, err = cluster.Cores[0].UnderlyingRawStorage.(*raft.RaftBackend).Peers(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		core1Peers, err = cluster.Cores[1].UnderlyingRawStorage.(*raft.RaftBackend).Peers(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		core2Peers, err = cluster.Cores[2].UnderlyingRawStorage.(*raft.RaftBackend).Peers(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if len(core0Peers) == 3 && reflect.DeepEqual(core0Peers, core1Peers) && reflect.DeepEqual(core1Peers, core2Peers) {
			break
		}
		time.Sleep(time.Second)
	}
	require.Equal(t, core0Peers, core1Peers)
	require.Equal(t, core1Peers, core2Peers)
}

// TestRaft_VotersStayVoters ensures that autopilot doesn't demote a node just
// because it hasn't been heard from in some time.
func TestRaft_VotersStayVoters(t *testing.T) {
	t.Parallel()
	reconcileInterval := 300 * time.Millisecond
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		PhysicalFactoryConfig: map[string]interface{}{
			"performance_multiplier":       "5",
			"autopilot_reconcile_interval": reconcileInterval.String(),
			"autopilot_update_interval":    "100ms",
		},
		VersionMap: map[int]string{
			0: version.Version,
			1: version.Version,
			2: version.Version,
		},
	})
	defer cluster.Cleanup()
	testhelpers.WaitForActiveNode(t, cluster)

	client := cluster.Cores[0].Client

	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)
	joinAndStabilize(t, cluster.Cores[1], client, cluster, config, "core-1", 2)
	waitUntilVoter(t, 2*reconcileInterval, client, "core-1")
	joinAndStabilize(t, cluster.Cores[2], client, cluster, config, "core-2", 3)
	waitUntilVoter(t, 2*reconcileInterval, client, "core-1")

	errIfNonVotersExist := func() error {
		t.Helper()
		resp, err := client.Sys().RaftAutopilotState()
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range resp.Servers {
			if v.Status == "non-voter" {
				return fmt.Errorf("node %q is a non-voter", k)
			}
		}
		return nil
	}
	testhelpers.RetryUntil(t, 10*time.Second, errIfNonVotersExist)

	// Core0 is the leader, sealing it will both cause an election - and the
	// new leader won't have seen any heartbeats initially - and create a "down"
	// node that won't be sending heartbeats.
	testhelpers.EnsureCoreSealed(t, cluster.Cores[0])
	time.Sleep(config.ServerStabilizationTime + 2*time.Second)
	client = cluster.Cores[1].Client
	err = errIfNonVotersExist()
	require.NoError(t, err)
}

// TestRaft_Autopilot_DeadServerCleanup tests that dead servers are correctly
// removed by Vault and autopilot when a node stops and a replacement node joins.
// The expected behavior is that removing a node from a 3 node cluster wouldn't
// remove it from Raft until a replacement voter had joined and stabilized/been promoted.
func TestRaft_Autopilot_DeadServerCleanup(t *testing.T) {
	t.Parallel()
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		NumCores:             4,
	})
	defer cluster.Cleanup()
	testhelpers.WaitForActiveNode(t, cluster)

	// Join 2 extra nodes manually, store the 3rd for later
	leader := cluster.Cores[0]
	core1 := cluster.Cores[1]
	core2 := cluster.Cores[2]
	core3 := cluster.Cores[3]
	joinAsVoterAndUnseal(t, core1, cluster)
	joinAsVoterAndUnseal(t, core2, cluster)
	// Do not join node 3
	testhelpers.WaitForNodesExcludingSelectedStandbys(t, cluster, 3)

	config, err := leader.Client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)
	require.True(t, isHealthyAfterStabilization(t, leader, config.ServerStabilizationTime))

	// Ensure Autopilot has the aggressive settings
	config.CleanupDeadServers = true
	config.ServerStabilizationTime = 5 * time.Second
	config.DeadServerLastContactThreshold = 1 * time.Minute
	config.MaxTrailingLogs = 10
	config.LastContactThreshold = 10 * time.Second
	config.MinQuorum = 3
	config.DisableUpgradeMigration = true

	// We can't use Client.Sys().PutRaftAutopilotConfiguration(config) in OSS as disable_upgrade_migration isn't in OSS
	b, err := json.Marshal(&config)
	require.NoError(t, err)
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	require.NoError(t, err)
	if !constants.IsEnterprise {
		delete(m, "disable_upgrade_migration")
	}
	_, err = leader.Client.Logical().Write("sys/storage/raft/autopilot/configuration", m)
	require.NoError(t, err)

	// Observe for healthy state
	state, err := leader.Client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.True(t, state.Healthy)

	// Kill a node (core-2)
	cluster.StopCore(t, 2)
	// Wait for just over the dead server threshold to ensure the core is classed as 'dead'
	time.Sleep(config.DeadServerLastContactThreshold + 2*time.Second)

	// Observe for an unhealthy state (but we still have 3 voters according to Raft)
	state, err = leader.Client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.False(t, state.Healthy)
	require.Len(t, state.Voters, 3)

	// Join node 3 now
	joinAsVoterAndUnseal(t, core3, cluster)

	// Stabilization time
	require.True(t, isHealthyAfterStabilization(t, leader, config.ServerStabilizationTime))

	// Observe for healthy and contains 3 correct voters
	state, err = leader.Client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.True(t, state.Healthy)
	require.Len(t, state.Voters, 3)
	require.Contains(t, state.Voters, "core-0")
	require.Contains(t, state.Voters, "core-1")
	require.NotContains(t, state.Voters, "core-2")
	require.Contains(t, state.Voters, "core-3")
}

func joinAndStabilizeAndPromote(t *testing.T, core *vault.TestClusterCore, client *api.Client, cluster *vault.TestCluster, config *api.AutopilotConfig, nodeID string, numServers int) {
	t.Helper()
	joinAndStabilize(t, core, client, cluster, config, nodeID, numServers)

	// Now that the server is stable, wait for autopilot to reconcile and
	// promotion to happen. Reconcile interval is 10 seconds. Bound it by
	// doubling.

	waitUntilVoter(t, 2*autopilot.DefaultReconcileInterval, client, nodeID)
}

func waitUntilVoter(t *testing.T, timeout time.Duration, client *api.Client, nodeID string) {
	t.Helper()

	// Now that the server is stable, wait for autopilot to reconcile and
	// promotion to happen. Reconcile interval is 10 seconds. Bound it by
	// doubling.
	testhelpers.RetryUntil(t, timeout, func() error {
		state, err := client.Sys().RaftAutopilotState()
		if err != nil {
			return err
		}
		if state.Servers[nodeID].Status != "voter" {
			return fmt.Errorf("autopilot failed to promote node: id: %#v: state:%# v\n", nodeID, pretty.Formatter(state))
		}
		return nil
	})
}

func joinAndStabilize(t *testing.T, core *vault.TestClusterCore, client *api.Client, cluster *vault.TestCluster, config *api.AutopilotConfig, nodeID string, numServers int) {
	t.Helper()
	joinAndUnseal(t, core, cluster, false, false)

	testhelpers.RetryUntil(t, config.ServerStabilizationTime, func() error {
		state, err := client.Sys().RaftAutopilotState()
		if err != nil {
			return err
		}
		if len(state.Servers) != numServers {
			return fmt.Errorf("expected %d servers, got %d", numServers, len(state.Servers))
		}
		if !state.Healthy {
			return fmt.Errorf("autopilot unhealthy: %# v", pretty.Formatter(state))
		}
		ss, ok := state.Servers[nodeID]
		if !ok {
			return fmt.Errorf("node %q not present", nodeID)
		}

		if ss.NodeStatus != "alive" {
			return fmt.Errorf("expected node %s to be alive, but NodeStatus=%q", nodeID, ss.NodeStatus)
		}

		return nil
	})
}

// joinAsVoterAndUnseal joins the specified core to the specified cluster as a voter and unseals it.
// It will wait (up to a timeout) for the core to be fully unsealed before returning
func joinAsVoterAndUnseal(t *testing.T, core *vault.TestClusterCore, cluster *vault.TestCluster) {
	t.Helper()
	joinAndUnseal(t, core, cluster, false, true)
}

// joinAndUnseal joins the specified core to the specified cluster and unseals it.
// You can specify if the core should be joined as a voter/non-voter,
// and whether to wait (up to a timeout) for the core to be unsealed before returning.
func joinAndUnseal(t *testing.T, core *vault.TestClusterCore, cluster *vault.TestCluster, nonVoter bool, waitForUnseal bool) {
	t.Helper()
	leaderIdx, err := testcluster.LeaderNode(context.Background(), cluster)
	require.NoError(t, err)
	leader := cluster.Cores[leaderIdx]

	resp, err := core.Client.Sys().RaftJoin(&api.RaftJoinRequest{
		LeaderAPIAddr:    leader.Client.Address(),
		LeaderCACert:     string(cluster.CACertPEM),
		LeaderClientCert: string(cluster.CACertPEM),
		LeaderClientKey:  string(cluster.CAKeyPEM),
		Retry:            false,
		NonVoter:         nonVoter,
	})
	require.NoError(t, err)
	require.True(t, resp.Joined)

	cluster.UnsealCore(t, core)
	if waitForUnseal {
		waitForCoreUnseal(t, core)
	}
}

// waitForCoreUnseal waits until the specified core is unsealed.
// It fails the calling test if the deadline has elapsed and the core is still sealed.
func waitForCoreUnseal(t *testing.T, core *vault.TestClusterCore) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if !core.Sealed() {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("expected core %v to unseal before deadline but it has not", core.NodeID)
}

// isHealthyAfterStabilization will use the supplied leader core to query the
// health of Raft Autopilot just after the specified deadline.
func isHealthyAfterStabilization(t *testing.T, leaderCore *vault.TestClusterCore, stabilizationTime time.Duration) bool {
	t.Helper()
	timeoutGrace := 2 * time.Second
	time.Sleep(stabilizationTime + timeoutGrace)
	state, err := leaderCore.Client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.NotNil(t, state)
	return state.Healthy
}
