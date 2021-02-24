package rafttests

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kr/pretty"

	"github.com/hashicorp/vault/api"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
)

func TestRaft_Autopilot_Disable(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.EqualValues(t, "not-running", state.ExecutionStatus)
}

func TestRaft_Autopilot_ServerStabilization(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	waitTime := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(waitTime)

	joinFunc := func(core *vault.TestClusterCore) {
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
			{
				LeaderAPIAddr: client.Address(),
				TLSConfig:     cluster.Cores[0].TLSConfig,
				Retry:         true,
			},
		}, false)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		cluster.UnsealCore(t, core)
	}

	joinFunc(cluster.Cores[1])
	joinFunc(cluster.Cores[2])

	testhelpers.VerifyRaftPeers(t, client, map[string]bool{
		"core-0": true,
		"core-1": true,
		"core-2": true,
	})

	deadline := time.Now().Add(20 * time.Second)
	success := false
	healthy := false

	var state *api.AutopilotState
	for time.Now().Before(deadline) {
		state, err := client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		if state.Healthy {
			healthy = true
		}

		if healthy && len(state.Voters) == 3 {
			success = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !healthy {
		t.Fatalf("servers failed to become healthy ")
	}

	if !success {
		t.Fatalf("servers failed to promote followers; state: %#v\n", state)
	}
}

func TestRaft_Autopilot_ServerStabilization_UnstableServer(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	waitTime := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(waitTime)

	joinFunc := func(core *vault.TestClusterCore) {
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
			{
				LeaderAPIAddr: client.Address(),
				TLSConfig:     cluster.Cores[0].TLSConfig,
				Retry:         true,
			},
		}, false)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		cluster.UnsealCore(t, core)
	}

	joinFunc(cluster.Cores[1])
	joinFunc(cluster.Cores[2])

	// Stabilization time + twice the reconcile time is what we will be waiting for.
	// Keep the server unstable till then.
	waitTime = time.Duration(float64(config.ServerStabilizationTime)) + 20*time.Second

	deadline := time.Now().Add(waitTime)
	success := false
	healthy := false

	var state *api.AutopilotState
	for time.Now().Before(deadline) {
		testhelpers.EnsureCoreSealed(t, cluster.Cores[1])
		cluster.UnsealCore(t, cluster.Cores[1])

		state, err := client.Sys().RaftAutopilotState()
		require.NoError(t, err)

		fmt.Printf("=====vishal: state: %# v\n", pretty.Formatter(state))

		if state.Healthy {
			healthy = true
		}

		if healthy && len(state.Voters) == 3 {
			success = true
			break
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Printf("=====vishal: final state: %# v\n", pretty.Formatter(state))
	fmt.Printf("=====vishal: success: %#v\n", success)
	if !healthy {
		t.Fatalf("servers failed to become healthy ")
	}

	if !success {
		t.Fatalf("servers failed to promote followers; state: %#v\n", state)
	}
}

func TestRaft_Autopilot_HCLConfiguration(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		PhysicalFactoryConfig: map[string]interface{}{
			"autopilot": `[{"cleanup_dead_servers":true,"last_contact_threshold":"500s","left_server_last_contact_threshold":"500h","max_trailing_logs":500,"min_quorum":500,"server_stabilization_time":"500s"}]`,
		},
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	configCheckFunc := func(config *api.AutopilotConfig) {
		conf, err := client.Sys().RaftAutopilotConfiguration()
		require.NoError(t, err)
		require.Equal(t, config, conf)
	}

	config := &api.AutopilotConfig{
		CleanupDeadServers:             true,
		LeftServerLastContactThreshold: 500 * time.Hour,
		LastContactThreshold:           500 * time.Second,
		MaxTrailingLogs:                500,
		MinQuorum:                      500,
		ServerStabilizationTime:        500 * time.Second,
	}
	configCheckFunc(config)

	// Ensure that the configuration stays across reboots
	leaderCore := cluster.Cores[0]
	testhelpers.EnsureCoreSealed(t, cluster.Cores[0])
	cluster.UnsealCore(t, leaderCore)
	vault.TestWaitActive(t, leaderCore.Core)
	configCheckFunc(config)

	// Only set some values, expect others to assume default values
	cluster = raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		PhysicalFactoryConfig: map[string]interface{}{
			"autopilot": `[{"cleanup_dead_servers":true,"server_stabilization_time":"500s"}]`,
		},
	})
	defer cluster.Cleanup()

	config.LeftServerLastContactThreshold = 24 * time.Hour
	config.LastContactThreshold = 10 * time.Second
	config.MaxTrailingLogs = 1000
	config.MinQuorum = 3

	client = cluster.Cores[0].Client
	configCheckFunc(config)

	leaderCore = cluster.Cores[0]
	testhelpers.EnsureCoreSealed(t, leaderCore)
	cluster.UnsealCore(t, leaderCore)
	vault.TestWaitActive(t, leaderCore.Core)
	configCheckFunc(config)
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
		LeftServerLastContactThreshold: 24 * time.Hour,
		LastContactThreshold:           10 * time.Second,
		MaxTrailingLogs:                1000,
		MinQuorum:                      3,
		ServerStabilizationTime:        10 * time.Second,
	}
	configCheckFunc(config)

	// Update config
	writableConfig := map[string]interface{}{
		"cleanup_dead_servers":               true,
		"left_server_last_contact_threshold": "100h",
		"last_contact_threshold":             "100s",
		"max_trailing_logs":                  100,
		"min_quorum":                         100,
		"server_stabilization_time":          "100s",
	}
	writeConfigFunc(writableConfig, false)

	// Ensure update has taken effect
	config.CleanupDeadServers = true
	config.LeftServerLastContactThreshold = 100 * time.Hour
	config.LastContactThreshold = 100 * time.Second
	config.MaxTrailingLogs = 100
	config.MinQuorum = 100
	config.ServerStabilizationTime = 100 * time.Second
	configCheckFunc(config)

	// Update some fields and leave the rest as it is.
	writableConfig = map[string]interface{}{
		"left_server_last_contact_threshold": "50h",
		"max_trailing_logs":                  50,
		"server_stabilization_time":          "50s",
	}
	writeConfigFunc(writableConfig, false)

	// Check update
	config.LeftServerLastContactThreshold = 50 * time.Hour
	config.MaxTrailingLogs = 50
	config.ServerStabilizationTime = 50 * time.Second
	configCheckFunc(config)

	// Check error case
	writableConfig = map[string]interface{}{
		"min_quorum":                         2,
		"left_server_last_contact_threshold": "48h",
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

func TestRaft_Autopilot_State(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
	})
	defer cluster.Cleanup()

	// Check that autopilot execution state is running
	client := cluster.Cores[0].Client
	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, state.ExecutionStatus, api.AutopilotRunning)
	require.Equal(t, state.Healthy, true)
	require.Len(t, state.Servers, 1)
	require.Equal(t, state.Servers["core-0"].ID, "core-0")
	require.Equal(t, state.Servers["core-0"].NodeStatus, "alive")
	require.Equal(t, state.Servers["core-0"].Status, "leader")

	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	waitTime := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(waitTime)

	joinAndStabilizeFunc := func(core *vault.TestClusterCore, nodeID string, numServers int) {
		joinFunc := func(core *vault.TestClusterCore) {
			_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
				{
					LeaderAPIAddr: client.Address(),
					TLSConfig:     cluster.Cores[0].TLSConfig,
					Retry:         true,
				},
			}, false)
			require.NoError(t, err)
			time.Sleep(2 * time.Second)
			cluster.UnsealCore(t, core)
		}
		joinFunc(core)

		state, err = client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		require.Equal(t, state.Healthy, false)
		require.Len(t, state.Servers, numServers)
		require.Equal(t, state.Servers[nodeID].Healthy, false)
		require.Equal(t, state.Servers[nodeID].NodeStatus, "alive")
		require.Equal(t, state.Servers[nodeID].Status, "non-voter")

		// Wait till the stabilization period is over
		waitTime = time.Duration(float64(config.ServerStabilizationTime))
		time.Sleep(waitTime)
		state, err = client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		require.Equal(t, state.Servers[nodeID].Healthy, true)

		// Now that the server is stable, wait for autopilot to reconcile and
		// promotion to happen. Reconcile interval is 10 seconds. Bound it by
		// doubling.
		deadline := time.Now().Add(20 * time.Second)
		failed := true
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
			t.Fatalf("failed to promote node: id: %#v: state:%# v\n", nodeID, pretty.Formatter(state))
		}
	}
	joinAndStabilizeFunc(cluster.Cores[1], "core-1", 2)
	joinAndStabilizeFunc(cluster.Cores[2], "core-2", 3)
	state, err = client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, state.Voters, []string{"core-0", "core-1", "core-2"})
}
