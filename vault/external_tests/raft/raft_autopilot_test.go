package rafttests

import (
	"context"
	"fmt"
	"math"
	"reflect"
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
	"github.com/kr/pretty"
	testingintf "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

func TestRaft_Autopilot_Disable(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
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
	cluster := raftCluster(t, &RaftClusterOpts{
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
	cluster := raftCluster(t, &RaftClusterOpts{
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

	join(t, cluster.Cores[1], client, cluster)
	join(t, cluster.Cores[2], client, cluster)

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
	cluster := raftCluster(t, &RaftClusterOpts{
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

	join(t, cluster.Cores[1], client, cluster)
	join(t, cluster.Cores[2], client, cluster)

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
	join(t, core, client, cluster)
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

func join(t *testing.T, core *vault.TestClusterCore, client *api.Client, cluster *vault.TestCluster) {
	t.Helper()
	_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
		{
			LeaderAPIAddr: client.Address(),
			TLSConfig:     cluster.Cores[0].TLSConfig,
			Retry:         true,
		},
	}, false)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	cluster.UnsealCore(t, core)
}
