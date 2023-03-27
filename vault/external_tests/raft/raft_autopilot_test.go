// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rafttests

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	autopilot "github.com/hashicorp/raft-autopilot"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
	"github.com/kr/pretty"
	testingintf "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

func TestRaft_Autopilot_Disable(t *testing.T) {
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

func TestRaft_Autopilot_Stabilization_And_State(t *testing.T) {
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

	joinAndStabilizeAndPromote(t, cluster.Cores[1], client, cluster, config, "core-1", 2)
	joinAndStabilizeAndPromote(t, cluster.Cores[2], client, cluster, config, "core-2", 3)
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
	conf, opts := teststorage.ClusterSetup(nil, nil, teststorage.RaftBackendSetup)
	conf.DisableAutopilot = false
	opts.InmemClusterLayers = true
	opts.KeepStandbysSealed = true
	opts.SetupFunc = nil
	timeToHealthyCore2 := 5 * time.Second
	opts.PhysicalFactory = func(t testingintf.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
		config := map[string]interface{}{
			"snapshot_threshold":           "50",
			"trailing_logs":                "100",
			"autopilot_reconcile_interval": "1s",
			"autopilot_update_interval":    "500ms",
			"snapshot_interval":            "1s",
		}
		if coreIdx == 2 {
			config["snapshot_delay"] = timeToHealthyCore2.String()
		}
		return teststorage.MakeRaftBackend(t, coreIdx, logger, config)
	}

	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
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
	stabilizationKickOffWaitDuration := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(stabilizationKickOffWaitDuration)

	cli := cluster.Cores[0].Client
	// Write more keys than snapshot_threshold
	for i := 0; i < 250; i++ {
		_, err := cli.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	joinAndUnseal(t, cluster.Cores[1], cluster, false, false)
	joinAndUnseal(t, cluster.Cores[2], cluster, false, false)

	core2shouldBeHealthyAt := time.Now().Add(timeToHealthyCore2)

	stabilizationWaitDuration := time.Duration(1.25 * float64(config.ServerStabilizationTime))
	deadline := time.Now().Add(stabilizationWaitDuration)
	var core1healthy, core2healthy bool
	for time.Now().Before(deadline) {
		state, err := client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		core1healthy = state.Servers["core-1"] != nil && state.Servers["core-1"].Healthy
		core2healthy = state.Servers["core-2"] != nil && state.Servers["core-2"].Healthy
		time.Sleep(1 * time.Second)
	}
	if !core1healthy || core2healthy {
		t.Fatalf("expected health: core1=true and core2=false, got: core1=%v, core2=%v", core1healthy, core2healthy)
	}

	time.Sleep(2 * time.Second) // wait for reconciliation
	state, err = client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, []string{"core-0", "core-1"}, state.Voters)

	for time.Now().Before(core2shouldBeHealthyAt) {
		state, err := client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		core2healthy = state.Servers["core-2"].Healthy
		time.Sleep(1 * time.Second)
		t.Log(core2healthy)
	}

	deadline = time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		state, err = client.Sys().RaftAutopilotState()
		if err != nil {
			t.Fatal(err)
		}
		if strutil.EquivalentSlices(state.Voters, []string{"core-0", "core-1", "core-2"}) {
			break
		}
	}
	require.Equal(t, []string{"core-0", "core-1", "core-2"}, state.Voters)
}

func TestRaft_AutoPilot_Peersets_Equivalent(t *testing.T) {
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
	cluster, _ := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		PhysicalFactoryConfig: map[string]interface{}{
			"performance_multiplier":       "5",
			"autopilot_reconcile_interval": "300ms",
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
	joinAndStabilizeAndPromote(t, cluster.Cores[1], client, cluster, config, "core-1", 2)
	joinAndStabilizeAndPromote(t, cluster.Cores[2], client, cluster, config, "core-2", 3)

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
	time.Sleep(30 * time.Second)
	client = cluster.Cores[1].Client
	errIfNonVotersExist()
}

// TestDelegate_KnownServers_NodeType tests that KnownServers returns servers will the correct NodeType
// based on the number of voters that should exist within the cluster. The test is fairly rudimentary in
// that it tests the servers after setup, but doesn't deal with new nodes being added or failed nodes.
// The reason for the test is that raft-autopilot now uses the NodeType in order to determine if a node
// is a potential voter (i.e. a non-voter that could become a voter, or an existing voter).
// Related Jira: https://hashicorp.atlassian.net/browse/VAULT-14048
func TestDelegate_KnownServers_NodeType(t *testing.T) {
	tests := map[string]struct {
		numNodes  int
		numVoters int
	}{
		"3-nodes-3-voters": {numNodes: 3, numVoters: 3},
		"4-nodes-3-voters": {numNodes: 4, numVoters: 3},
		"5-nodes-3-voters": {numNodes: 5, numVoters: 3},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Some involved setup, so we can control the nodes being added to the cluster and have access to the leader backend.
			conf, opts := teststorage.ClusterSetup(nil, nil, teststorage.RaftBackendSetup)
			conf.DisableAutopilot = false
			opts.NumCores = tc.numNodes
			opts.SetupFunc = nil
			cluster := vault.NewTestCluster(t, conf, opts)
			cluster.Start()
			defer cluster.Cleanup()
			leader, addressProvider := setupLeaderAndUnseal(t, cluster)

			// Add the other nodes
			voters := 1
			for i := 1; i < len(cluster.Cores); i++ {
				core := cluster.Cores[i]
				core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
				// Only add the number of nodes specified as 'voters' the rest won't be.
				if voters < tc.numVoters {
					joinAsVoterAndUnseal(t, core, cluster)
				} else {
					joinAsNonVoterAndUnseal(t, core, cluster)
				}
				voters++
			}
			testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

			// Actual code related to testing the return values of the KnownServers func
			backend := leader.GetRaftBackend()
			d := raft.NewDelegate(backend)
			servers := d.KnownServers()
			voters = 0
			for id, srv := range servers {
				fmt.Printf("known servers: id: %v, node type: %v\n", id, srv.NodeType)
				if srv.NodeType == "voter" {
					voters++
				}
			}
			require.Equal(t, tc.numVoters, voters, "known servers reported a number of voters that wasn't expected")
		})
	}
}

// TestRaft_Autopilot_DeadServerCleanup tests that dead servers are correctly removed by Vault and autopilot when a node stops and a replacement node joins.
// The expected behavior is that removing a node from a 3 node cluster wouldn't remove it from Raft until a replacement voter had joined and stabilized/been promoted.
func TestRaft_Autopilot_DeadServerCleanup(t *testing.T) {
	// Some involved setup, so we can control the nodes being added to the cluster and have access to the leader backend.
	conf, opts := teststorage.ClusterSetup(nil, nil, teststorage.RaftBackendSetup)
	conf.DisableAutopilot = false
	opts.NumCores = 4
	opts.SetupFunc = nil
	opts.PhysicalFactoryConfig = map[string]interface{}{
		"autopilot_reconcile_interval": "300ms",
		"autopilot_update_interval":    "100ms",
	}

	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()
	leader, addressProvider := setupLeaderAndUnseal(t, cluster)

	// Join 2 extra nodes manually, store the 3rd for later
	core1 := cluster.Cores[1]
	core2 := cluster.Cores[2]
	core3 := cluster.Cores[3]
	core1.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
	core2.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
	core3.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
	joinAsVoterAndUnseal(t, core1, cluster)
	joinAsVoterAndUnseal(t, core2, cluster)
	core3, cluster.Cores = cluster.Cores[len(cluster.Cores)-1], cluster.Cores[:len(cluster.Cores)-1]
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

	config, err := leader.Client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)
	require.True(t, isStableByDeadline(t, leader, config.ServerStabilizationTime))

	// Ensure Autopilot has the aggressive settings
	config.CleanupDeadServers = true
	config.ServerStabilizationTime = 5 * time.Second
	config.DeadServerLastContactThreshold = 10 * time.Second
	config.MaxTrailingLogs = 10
	config.LastContactThreshold = 10 * time.Second
	config.MinQuorum = 3
	err = leader.Client.Sys().PutRaftAutopilotConfiguration(config)
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

	// Push the spare core back and join the node
	cluster.Cores = append(cluster.Cores, core3)
	joinAsVoterAndUnseal(t, core3, cluster)

	// Stabilization time
	require.True(t, isStableByDeadline(t, leader, config.ServerStabilizationTime))

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
	joinAndStabilize(t, core, client, cluster, config, nodeID, numServers)

	// Now that the server is stable, wait for autopilot to reconcile and
	// promotion to happen. Reconcile interval is 10 seconds. Bound it by
	// doubling.
	deadline := time.Now().Add(2 * autopilot.DefaultReconcileInterval)
	failed := true
	var err error
	var state *api.AutopilotState
	for time.Now().Before(deadline) {
		state, err = client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		if state.Servers[nodeID].Status == "voter" {
			failed = false
			break
		}
		time.Sleep(1 * time.Second)
	}

	if failed {
		t.Fatalf("autopilot failed to promote node: id: %#v: state:%# v\n", nodeID, pretty.Formatter(state))
	}
}

func joinAndStabilize(t *testing.T, core *vault.TestClusterCore, client *api.Client, cluster *vault.TestCluster, config *api.AutopilotConfig, nodeID string, numServers int) {
	t.Helper()
	joinAndUnseal(t, core, cluster, false, false)
	time.Sleep(2 * time.Second)

	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, false, state.Healthy)
	require.Len(t, state.Servers, numServers)
	require.Equal(t, false, state.Servers[nodeID].Healthy)
	require.Equal(t, "alive", state.Servers[nodeID].NodeStatus)
	require.Equal(t, "non-voter", state.Servers[nodeID].Status)

	// Wait till the stabilization period is over
	deadline := time.Now().Add(config.ServerStabilizationTime)
	healthy := false
	for time.Now().Before(deadline) {
		state, err := client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		if state.Healthy {
			healthy = true
		}
		time.Sleep(1 * time.Second)
	}
	if !healthy {
		t.Fatalf("cluster failed to stabilize")
	}
}

// joinAsVoterAndUnseal joins the specified core to the specified cluster as a voter and unseals it.
// It will wait (up to a timeout) for the core to be fully unsealed before returning
func joinAsVoterAndUnseal(t *testing.T, core *vault.TestClusterCore, cluster *vault.TestCluster) {
	joinAndUnseal(t, core, cluster, false, true)
}

// joinAsNonVoterAndUnseal joins the specified core to the specified cluster as a non-voter and unseals it.
// It will wait (up to a timeout) for the core to be fully unsealed before returning
func joinAsNonVoterAndUnseal(t *testing.T, core *vault.TestClusterCore, cluster *vault.TestCluster) {
	joinAndUnseal(t, core, cluster, true, true)
}

// joinAndUnseal joins the specified core to the specified cluster and unseals it.
// You can specify if the core should be joined as a voter/non-voter, and whether to wait (up to a timeout) for the core to be unsealed before returning.
func joinAndUnseal(t *testing.T, core *vault.TestClusterCore, cluster *vault.TestCluster, nonVoter bool, waitForUnseal bool) {
	leader, leaderAddr := clusterLeader(t, cluster)
	_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
		{
			LeaderAPIAddr: leaderAddr,
			TLSConfig:     leader.TLSConfig(),
			Retry:         true,
		},
	}, nonVoter)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	cluster.UnsealCore(t, core)
	if waitForUnseal {
		waitForCoreUnseal(t, core)
	}
}

// clusterLeader gets the leader node and its address from the specified cluster
func clusterLeader(t *testing.T, cluster *vault.TestCluster) (*vault.TestClusterCore, string) {
	for _, core := range cluster.Cores {
		isLeader, addr, _, err := core.Leader()
		require.NoError(t, err)
		if isLeader {
			return core, addr
		}
	}

	t.Fatal("unable to find leader")
	return nil, ""
}

// setupLeaderAndUnseal configures and unseals the leader node.
// It will wait until the node is active before returning the core and the address of the leader.
func setupLeaderAndUnseal(t *testing.T, cluster *vault.TestCluster) (*vault.TestClusterCore, *testhelpers.TestRaftServerAddressProvider) {
	leader, _ := clusterLeader(t, cluster)

	// Lots of tests seem to do this when they deal with a TestRaftServerAddressProvider, it makes the test work rather than error out.
	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}
	testhelpers.EnsureCoreSealed(t, leader)
	leader.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
	cluster.UnsealCore(t, leader)
	vault.TestWaitActive(t, leader.Core)

	return leader, addressProvider
}

// waitForCoreUnseal waits until the specified core is unsealed.
// It fails the calling test if the deadline has elapsed and the core is still sealed.
func waitForCoreUnseal(t *testing.T, core *vault.TestClusterCore) {
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if !core.Sealed() {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("expected core %v to unseal before deadline but it has not", core.NodeID)
}

// isStableByDeadline will use the supplied leader core to query the health of Raft Autopilot up until just after the supplied deadline.
// It will return when the Raft Autopilot state is reported as healthy, or just after the specified stabilization time.
func isStableByDeadline(t *testing.T, leaderCore *vault.TestClusterCore, stabilizationTime time.Duration) bool {
	timeoutGrace := 2 * time.Second
	deadline := time.Now().Add(stabilizationTime).Add(timeoutGrace)
	for time.Now().Before(deadline) {
		state, err := leaderCore.Client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		if state.Healthy {
			return true
		}
		time.Sleep(1 * time.Second)
	}

	return false
}
