// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft_binary

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	rafttest "github.com/hashicorp/vault/vault/external_tests/raft"
)

// TestRaft_Configuration_Docker is a variant of TestRaft_Configuration that
// uses docker containers for the vault nodes.
func TestRaft_Configuration_Docker(t *testing.T) {
	t.Parallel()
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}
	opts := &docker.DockerClusterOptions{
		ImageRepo: "hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag:    "latest",
		VaultBinary: binary,
		ClusterOptions: testcluster.ClusterOptions{
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
				// If you want the test to run faster locally, you could
				// uncomment this performance_multiplier change.
				//StorageOptions: map[string]string{
				//	"performance_multiplier": "1",
				//},
			},
		},
	}
	cluster := docker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()
	rafttest.Raft_Configuration_Test(t, cluster)

	if err := cluster.AddNode(context.TODO(), opts); err != nil {
		t.Fatal(err)
	}
	rafttest.Raft_Configuration_Test(t, cluster)
}
