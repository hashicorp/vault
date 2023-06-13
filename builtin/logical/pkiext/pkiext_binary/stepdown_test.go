// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pkiext_binary

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

func TestStepDown(t *testing.T) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	expectedNodes := 3

	opts := &docker.DockerClusterOptions{
		ImageRepo: "docker.mirror.hashicorp.services/hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag:    "latest",
		VaultBinary: binary,
		ClusterOptions: testcluster.ClusterOptions{
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
			},
			NumCores: expectedNodes,
		},
	}

	cluster := docker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()

	ctx := context.Background()

	_, err := testcluster.WaitForActiveNode(ctx, cluster)
	require.NoError(t, err, "failed waiting for active node")

	activeNode, _, err := testcluster.GetActiveAndStandbys(ctx, cluster)
	require.NoError(t, err, "failed getting active and standby nodes")

	// Simply wait for all nodes to appear within ha status report
	testhelpers.RetryUntil(t, 1*time.Minute, func() error {
		haStatus, _ := activeNode.APIClient().Sys().HAStatus()

		numNodes := len(haStatus.Nodes)
		if numNodes != expectedNodes {
			return fmt.Errorf("expected %d nodes within ha status got %d", expectedNodes, numNodes)
		}

		return nil
	})

	// Wait for everything to be healthy based on Raft Auto Pilot
	testhelpers.RetryUntil(t, 2*time.Minute, func() error {
		state, err := activeNode.APIClient().Sys().RaftAutopilotState()
		if err != nil {
			return err
		}

		t.Logf("Raft AutoPilotState top-level healthy: %v Leader: %s\n", state.Healthy, state.Leader)

		if !state.Healthy {
			return fmt.Errorf("raft auto pilot top-level state is not healthy")
		}

		for node, nodeState := range state.Servers {
			t.Logf("Node: %s State: %v", node, nodeState)
			if !nodeState.Healthy {
				return fmt.Errorf("raft auto pilot node state for %s is not healthy", node)
			}
		}

		return nil
	})

	// Uncomment to get a passing test
	// time.Sleep(30 * time.Second)

	t.Logf("Sealing active node")
	err = activeNode.APIClient().Sys().Seal()
	require.NoError(t, err, "failed stepping down node")

	activeNode, _, err = testcluster.GetActiveAndStandbys(ctx, cluster)
	require.NoError(t, err, "failed getting active and standby nodes after sealing")
}
