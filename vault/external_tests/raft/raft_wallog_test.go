package rafttests

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/vault"
	testingintf "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math"
	"path/filepath"
	"testing"
	"time"
)

func TestRaft_WALLog(t *testing.T) {
	archivePath := t.TempDir()
	conf, opts := teststorage.ClusterSetup(nil, nil, teststorage.RaftBackendSetup)
	conf.DisableAutopilot = false
	opts.InmemClusterLayers = true
	opts.KeepStandbysSealed = true
	opts.SetupFunc = nil
	opts.PhysicalFactory = func(t testingintf.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
		config := map[string]interface{}{
			"snapshot_threshold":           "50",
			"trailing_logs":                "100",
			"autopilot_reconcile_interval": "1s",
			"archive_path":                 archivePath,
			"performance_multiplier":       "1",
			"logstore":                     "raft-wal",
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
		"server_stabilization_time": "3s",
	})
	require.NoError(t, err)

	config, err := client.Sys().RaftAutopilotConfiguration()
	require.NoError(t, err)

	// Wait for 110% of the stabilization time to add nodes
	stabilizationKickOffWaitDuration := time.Duration(math.Ceil(1.1 * float64(config.ServerStabilizationTime)))
	time.Sleep(stabilizationKickOffWaitDuration)

	cli := cluster.Cores[0].Client
	for i := 0; i < 25; i++ {
		_, err := cli.Logical().Write(fmt.Sprintf("secret/%d", i), map[string]interface{}{
			"test": "data",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	joinFunc := func(core *vault.TestClusterCore) {
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), []*raft.LeaderJoinInfo{
			{
				LeaderAPIAddr: client.Address(),
				TLSConfig:     cluster.Cores[0].TLSConfig,
			},
		}, false)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		cluster.UnsealCore(t, core)
	}

	joinFunc(cluster.Cores[1])
	joinFunc(cluster.Cores[2])

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		state, err = client.Sys().RaftAutopilotState()
		if err != nil {
			t.Fatal(err)
		}
		if strutil.EquivalentSlices(state.Voters, []string{"core-0", "core-1", "core-2"}) {
			break
		}
	}
	require.Equal(t, state.Voters, []string{"core-0", "core-1", "core-2"})

	ents, err := ioutil.ReadDir(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	var logs []string
	for _, ent := range ents {
		logs = append(logs, ent.Name())
	}
	t.Logf("log files: %v", logs)

	lastIndex, _, _ := cluster.Cores[0].Core.GetRaftIndexes()
	cluster.Cleanup()

	newRaftDir := t.TempDir()
	_, err = raft.CopyDir(filepath.Join(newRaftDir, "raft", "logs"), archivePath)
	if err != nil {
		t.Fatal(err)
	}

	conf, opts = teststorage.ClusterSetup(nil, nil, teststorage.RaftBackendSetup)
	conf.DisableAutopilot = false
	opts.InmemClusterLayers = true
	opts.KeepStandbysSealed = true
	opts.SetupFunc = nil
	opts.SkipInit = true
	opts.NumCores = 1
	conf.DisableAutopilot = true
	opts.PhysicalFactory = func(t testingintf.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
		config := map[string]interface{}{
			"performance_multiplier": "1",
			"logstore":               "raft-wal",
			"lastindex":              fmt.Sprintf("%d", lastIndex),
			"path":                   newRaftDir,
		}
		return teststorage.MakeRaftBackend(t, coreIdx, logger, config)
	}

	// Prepare peers.json
	type RecoveryPeer struct {
		ID       string `json:"id"`
		Address  string `json:"address"`
		NonVoter bool   `json:"non_voter"`
	}

	// Leave out node 1 during recovery
	peersList := make([]*RecoveryPeer, 0, 3)
	peersList = append(peersList, &RecoveryPeer{
		ID:       "newnode",
		Address:  "inmem-clusters_node_0",
		NonVoter: false,
	})

	peersJSONBytes, err := jsonutil.EncodeJSON(peersList)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(filepath.Join(newRaftDir, "raft"), "peers.json"), peersJSONBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	oldcluster := cluster
	cluster = vault.NewTestCluster(t, conf, opts)
	cluster.BarrierKeys = oldcluster.BarrierKeys
	cluster.Start()
	defer cluster.Cleanup()
	if err := cluster.UnsealCoresWithError(false); err != nil {
		t.Fatal(err)
	}
	testhelpers.WaitForActiveNode(t, cluster)

	token := cli.Token()
	cli = cluster.Cores[0].Client
	cli.SetToken(token)
	secret, err := cli.Logical().List("secret")
	if len(secret.Data["keys"].([]interface{})) != 25 {
		t.Fatal("didnt' find 25 keys under secret/")
	}
}
