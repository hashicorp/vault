package rafttests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/go-hclog"

	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
	vaultcluster "github.com/hashicorp/vault/vault/cluster"
)

func raftClusterWithAutopilot(t testing.TB, joinNodes bool) *vault.TestCluster {
	conf := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}

	inmemCluster, err := vaultcluster.NewInmemLayerCluster("inmem-cluster", 3, hclog.New(&hclog.LoggerOptions{
		Mutex: &sync.Mutex{},
		Level: hclog.Trace,
		Name:  "inmem-cluster",
	}))
	if err != nil {
		t.Fatal(err)
	}

	var opts = vault.TestClusterOptions{
		HandlerFunc:   vaulthttp.Handler,
		ClusterLayers: inmemCluster,
	}

	teststorage.RaftBackendWithAutopilotSetup(conf, &opts)

	if !joinNodes {
		opts.SetupFunc = nil
	}

	cluster := vault.NewTestCluster(t, conf, &opts)
	cluster.Start()
	vault.TestWaitActive(t, cluster.Cores[0].Core)

	return cluster
}

func TestRaft_Autopilot(t *testing.T) {
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

func TestRaft_Autopilot_DefaultConfiguration(t *testing.T) {
	// Create a single node raft cluster with autopilot enabled
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
		require.Equal(t, conf, config)
	}

	config := &api.AutopilotConfig{
		CleanupDeadServers:          false,
		LastContactFailureThreshold: 24 * time.Hour,
		LastContactThreshold:        10 * time.Second,
		MaxTrailingLogs:             1000,
		MinQuorum:                   3,
		ServerStabilizationTime:     10 * time.Second,
	}
	configCheckFunc(config)
}
