package rafttests

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
)

func TestRaft_Autopilot_ServerStabilization(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
	})
	defer cluster.Cleanup()

	// Wait 11s before trying to add nodes: the autopilot ServerStabilization time
	// is 10s, and autopilot.State.ServerStabilizationTime basically ignores server
	// stabilization for promotion purposes until the autopilot node has been
	// running for 110% of the ServerStabilization config setting.
	time.Sleep(11 * time.Second)

	joinFunc := func(core *vault.TestClusterCore) {
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
			{
				LeaderAPIAddr: cluster.Cores[0].Client.Address(),
				TLSConfig:     cluster.Cores[0].TLSConfig,
				Retry:         true,
			},
		}, false)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		cluster.UnsealCore(t, core)
	}

	client := cluster.Cores[0].Client
	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, state.ExecutionStatus, api.AutopilotRunning)

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
		t.Fatalf("servers failed to promote followers; state: %#v", state)
	}
}

func TestRaft_Autopilot_HCLConfiguration(t *testing.T) {
	cluster := raftCluster(t, &RaftClusterOpts{
		DisableFollowerJoins: true,
		InmemCluster:         true,
		EnableAutopilot:      true,
		AutopilotHCLValue:    `[{"cleanup_dead_servers":true,"last_contact_threshold":"500s","last_contact_failure_threshold":"500h","max_trailing_logs":500,"min_quorum":500,"server_stabilization_time":"500s"}]`,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	configCheckFunc := func(config *api.AutopilotConfig) {
		conf, err := client.Sys().RaftAutopilotConfiguration()
		require.NoError(t, err)
		require.Equal(t, config, conf)
	}

	config := &api.AutopilotConfig{
		CleanupDeadServers:          true,
		LastContactFailureThreshold: 500 * time.Hour,
		LastContactThreshold:        500 * time.Second,
		MaxTrailingLogs:             500,
		MinQuorum:                   500,
		ServerStabilizationTime:     500 * time.Second,
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
		AutopilotHCLValue:    `[{"cleanup_dead_servers":true,"server_stabilization_time":"500s"}]`,
	})
	defer cluster.Cleanup()

	config.LastContactFailureThreshold = 24 * time.Hour
	config.LastContactThreshold = 10 * time.Second
	config.MaxTrailingLogs = 1000
	config.MinQuorum = 3

	client = cluster.Cores[0].Client
	configCheckFunc(config)

	leaderCore = cluster.Cores[0]
	testhelpers.EnsureCoreSealed(t, cluster.Cores[0])
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

	// Check that autopilot execution state is running
	client := cluster.Cores[0].Client
	state, err := client.Sys().RaftAutopilotState()
	require.NoError(t, err)
	require.Equal(t, state.ExecutionStatus, api.AutopilotRunning)
	require.Equal(t, state.Healthy, true)

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
		CleanupDeadServers:          false,
		LastContactFailureThreshold: 24 * time.Hour,
		LastContactThreshold:        10 * time.Second,
		MaxTrailingLogs:             1000,
		MinQuorum:                   3,
		ServerStabilizationTime:     10 * time.Second,
	}
	configCheckFunc(config)

	// Update config
	writableConfig := map[string]interface{}{
		"cleanup_dead_servers":           true,
		"last_contact_failure_threshold": "100h",
		"last_contact_threshold":         "100s",
		"max_trailing_logs":              100,
		"min_quorum":                     100,
		"server_stabilization_time":      "100s",
	}
	writeConfigFunc(writableConfig, false)

	// Ensure update has taken effect
	config.CleanupDeadServers = true
	config.LastContactFailureThreshold = 100 * time.Hour
	config.LastContactThreshold = 100 * time.Second
	config.MaxTrailingLogs = 100
	config.MinQuorum = 100
	config.ServerStabilizationTime = 100 * time.Second
	configCheckFunc(config)

	// Update some fields and leave the rest as it is.
	writableConfig = map[string]interface{}{
		"last_contact_failure_threshold": "50h",
		"max_trailing_logs":              50,
		"server_stabilization_time":      "50s",
	}
	writeConfigFunc(writableConfig, false)

	// Check update
	config.LastContactFailureThreshold = 50 * time.Hour
	config.MaxTrailingLogs = 50
	config.ServerStabilizationTime = 50 * time.Second
	configCheckFunc(config)

	// Check error case
	writableConfig = map[string]interface{}{
		"min_quorum":                     2,
		"last_contact_failure_threshold": "48h",
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
