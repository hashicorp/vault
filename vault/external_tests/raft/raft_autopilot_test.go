package rafttests

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/kr/pretty"

	autopilot "github.com/hashicorp/raft-autopilot"

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
		// Not setting EnableAutopilot here.
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Nil(t, nil, state)
}

func TestRaft_Autopilot_Stabilization_And_State(t *testing.T) {
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
	require.Equal(t, api.AutopilotRunning, state.ExecutionStatus)
	require.Equal(t, true, state.Healthy)
	require.Len(t, state.Servers, 1)
	require.Equal(t, "core-0", state.Servers["core-0"].ID)
	require.Equal(t, "alive", state.Servers["core-0"].NodeStatus)
	require.Equal(t, "leader", state.Servers["core-0"].Status)

	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	stabilizationKickOffWaitDuration := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(stabilizationKickOffWaitDuration)

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
			time.Sleep(1 * time.Second)
			cluster.UnsealCore(t, core)
		}
		joinFunc(core)

		state, err = client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		require.Equal(t, false, state.Healthy)
		require.Len(t, state.Servers, numServers)
		require.Equal(t, false, state.Servers[nodeID].Healthy)
		require.Equal(t, "alive", state.Servers[nodeID].NodeStatus)
		require.Equal(t, "non-voter", state.Servers[nodeID].Status)

		// Wait till the stabilization period is over
		stabilizationWaitDuration := time.Duration(float64(config.ServerStabilizationTime))
		deadline := time.Now().Add(stabilizationWaitDuration)
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

		// Now that the server is stable, wait for autopilot to reconcile and
		// promotion to happen. Reconcile interval is 10 seconds. Bound it by
		// doubling.
		deadline = time.Now().Add(2 * autopilot.DefaultReconcileInterval)
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
			t.Fatalf("autopilot failed to promote node: id: %#v: state:%# v\n", nodeID, pretty.Formatter(state))
		}
	}
	joinAndStabilizeFunc(cluster.Cores[1], "core-1", 2)
	joinAndStabilizeFunc(cluster.Cores[2], "core-2", 3)
	state, err = client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, []string{"core-0", "core-1", "core-2"}, state.Voters)
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
		MinQuorum:                      3,
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
