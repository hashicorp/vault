// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft_binary

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	autopilot "github.com/hashicorp/raft-autopilot"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	rafttest "github.com/hashicorp/vault/vault/external_tests/raft"
	"github.com/stretchr/testify/require"
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

// removeRaftNode removes a node from the raft configuration using the leader client
// and removes the docker node.
func removeRaftNode(t *testing.T, node *docker.DockerClusterNode, client *api.Client, serverID string) {
	t.Helper()
	_, err := client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id": serverID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// clean up the cluster nodes. Note that the node is not removed from the ClusterNodes
	node.Cleanup()
}

// stabilizeAndPromote makes sure the given node ID is among the voters using
// autoPilot state
func stabilizeAndPromote(t *testing.T, client *api.Client, nodeID string) {
	t.Helper()
	deadline := time.Now().Add(2 * autopilot.DefaultReconcileInterval)
	failed := true
	var state *api.AutopilotState
	var err error
	for time.Now().Before(deadline) {
		state, err = client.Sys().RaftAutopilotState()
		// If the state endpoint gets called during a leader election, we'll get an error about
		// there not being an active cluster node. Rather than erroring out of this loop, just
		// ignore the error and keep trying. It should resolve in a few seconds. There's a
		// deadline after all, so it's not like this loop will continue indefinitely.
		if err != nil {
			if strings.Contains(err.Error(), "active cluster node not found") {
				continue
			}

			t.Fatal(err)
		}

		if state != nil && state.Servers != nil && state.Servers[nodeID].Status == "voter" {
			failed = false
			break
		}
		time.Sleep(1 * time.Second)
	}

	if failed {
		t.Fatalf("autopilot failed to promote node: id: %#v: state:%# v\n", nodeID, state)
	}
}

// stabilize makes sure the cluster is in a healthy state using autopilot state
func stabilize(t *testing.T, client *api.Client) {
	t.Helper()
	deadline := time.Now().Add(2 * autopilot.DefaultReconcileInterval)
	healthy := false
	for time.Now().Before(deadline) {
		state, err := client.Sys().RaftAutopilotState()
		require.NoError(t, err)
		if state.Healthy {
			healthy = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	if !healthy {
		t.Fatalf("cluster failed to stabilize")
	}
}

// TestDocker_LogStore_Boltdb_To_Raftwal_And_Back runs 3 node cluster leveraging boltDB
// as the logStore, then migrates the cluster to raft-wal logStore and back.
// This shows raft-wal does not lose data.
// There is no migration procedure for individual nodes.
// The correct procedure is destroying existing raft-boltdb nodes and starting brand-new
// nodes that use raft-wal (and vice-versa)
// Having a cluster of mixed nodes, some using raft-boltdb and some using raft-wal, is not a problem.
func TestDocker_LogStore_Boltdb_To_Raftwal_And_Back(t *testing.T) {
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
			},
		},
	}
	cluster := docker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()

	rafttest.Raft_Configuration_Test(t, cluster)

	leaderNode := cluster.GetActiveClusterNode()
	leaderClient := leaderNode.APIClient()

	err := leaderClient.Sys().MountWithContext(context.TODO(), "kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	val := 0
	writeKV := func(client *api.Client, num int) {
		t.Helper()
		start := val
		for i := start; i < start+num; i++ {
			if _, err = leaderClient.Logical().WriteWithContext(context.TODO(), fmt.Sprintf("kv/foo-%d", i), map[string]interface{}{
				"bar": val,
			}); err != nil {
				t.Fatal(err)
			}
			val++
		}
	}

	readKV := func(client *api.Client) {
		t.Helper()
		for i := 0; i < val; i++ {
			secret, err := client.Logical().Read(fmt.Sprintf("kv/foo-%d", i))
			if err != nil {
				t.Fatal(err)
			}
			if secret == nil || secret.Data == nil {
				t.Fatal("failed to read the value")
			}
		}
	}
	// writing then reading some data
	writeKV(leaderClient, 10)
	readKV(leaderClient)

	if opts.ClusterOptions.VaultNodeConfig.StorageOptions == nil {
		opts.ClusterOptions.VaultNodeConfig.StorageOptions = make(map[string]string, 0)
	}
	// adding three new nodes with raft-wal as their log store
	opts.ClusterOptions.VaultNodeConfig.StorageOptions["raft_wal"] = "true"
	for i := 0; i < 3; i++ {
		if err := cluster.AddNode(context.TODO(), opts); err != nil {
			t.Fatal(err)
		}
	}
	// check raft config contain 6 nodes
	rafttest.Raft_Configuration_Test(t, cluster)

	// write data before removing two ndoes
	writeKV(leaderClient, 10)

	// remove two nodes using boltDB
	removeRaftNode(t, cluster.ClusterNodes[1], leaderClient, "core-1")
	removeRaftNode(t, cluster.ClusterNodes[2], leaderClient, "core-2")

	// all data should still be readable after removing two nodes
	readKV(leaderClient)

	err = testhelpers.VerifyRaftPeers(t, leaderClient, map[string]bool{
		"core-0": true,
		"core-3": true,
		"core-4": true,
		"core-5": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// stabilize the cluster, wait for the autopilot to promote new nodes to Voter
	stabilizeAndPromote(t, leaderClient, "core-3")
	stabilizeAndPromote(t, leaderClient, "core-4")
	stabilizeAndPromote(t, leaderClient, "core-5")

	// step down leader and remove the node afterwards
	// this will remove the container without explicitly calling stepdown
	cluster.ClusterNodes[0].Stop()

	// get new leader node
	leaderNode = cluster.GetActiveClusterNode()
	leaderClient = leaderNode.APIClient()

	// remove the old leader, which was using boltDB
	// this will remove it from raft configuration though the container was already removed
	removeRaftNode(t, cluster.ClusterNodes[0], leaderClient, "core-0")

	// check if the cluster is stable
	stabilize(t, leaderClient)

	// write some more data and read all data again.
	// Here, only raft-wal is in use in the cluster.
	writeKV(leaderClient, 10)
	readKV(leaderClient)

	// going back to boltdb. Adding three nodes to the cluster having them use boltDB
	opts.ClusterOptions.VaultNodeConfig.StorageOptions["raft_wal"] = "false"
	for i := 0; i < 3; i++ {
		if err := cluster.AddNode(context.TODO(), opts); err != nil {
			t.Fatal(err)
		}
	}
	// make sure the new nodes are promoted to voter
	stabilizeAndPromote(t, leaderClient, "core-6")
	stabilizeAndPromote(t, leaderClient, "core-7")
	stabilizeAndPromote(t, leaderClient, "core-8")

	err = testhelpers.VerifyRaftPeers(t, leaderClient, map[string]bool{
		"core-3": true,
		"core-4": true,
		"core-5": true,
		"core-6": true,
		"core-7": true,
		"core-8": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// get the leader node again just in case
	leaderNode = cluster.GetActiveClusterNode()
	leaderClient = leaderNode.APIClient()

	// remove all the nodes that use raft-wal except the one that is the leader.
	removingNodes := []string{"core-3", "core-4", "core-5"}
	var raftWalLeader string
	// if the leaderNode.NodeID exists in removingNodes, keep track of that
	if strutil.StrListContains(removingNodes, leaderNode.NodeID) {
		raftWalLeader = leaderNode.NodeID
	}

	// remove all nodes except the leader, if the leader is a node with raft-wal as its log store
	for _, node := range cluster.ClusterNodes {
		if node.NodeID == leaderNode.NodeID || !strutil.StrListContains(removingNodes, node.NodeID) {
			continue
		}
		removeRaftNode(t, node, leaderClient, node.NodeID)
	}

	// write and read data again after removing two or three nodes
	writeKV(leaderClient, 10)
	readKV(leaderClient)

	// remove the old leader that uses raft-wal as its logStore if it has not been removed
	if raftWalLeader != "" {
		oldLeader := leaderNode

		// remove the node
		leaderNode.Stop()

		// get new leader node
		leaderNode = cluster.GetActiveClusterNode()
		leaderClient = leaderNode.APIClient()

		// remove the old leader from the raft configuration
		removeRaftNode(t, oldLeader, leaderClient, raftWalLeader)
	}

	// make sure the cluster is healthy
	stabilize(t, leaderClient)

	// write some data again and read all data from the beginning
	writeKV(leaderClient, 10)
	readKV(leaderClient)
}

// TestRaft_LogStore_Migration_Snapshot checks migration from a boltDB to raftwal
// by performing a snapshot restore from one cluster to another, and checking no data loss
func TestRaft_LogStore_Migration_Snapshot(t *testing.T) {
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
			NumCores: 1,
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
			},
		},
	}
	cluster := docker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()
	rafttest.Raft_Configuration_Test(t, cluster)

	leaderNode := cluster.GetActiveClusterNode()
	leaderClient := leaderNode.APIClient()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := leaderClient.Sys().MountWithContext(ctx, "kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}
	val := 1
	for i := 0; i < 10; i++ {
		_, err = leaderClient.Logical().WriteWithContext(ctx, fmt.Sprintf("kv/foo-%d", i), map[string]interface{}{
			"bar": val,
		})
		val++
	}

	readKV := func(client *api.Client) {
		for i := 0; i < 10; i++ {
			secret, err := client.Logical().Read(fmt.Sprintf("kv/foo-%d", i))
			if err != nil {
				t.Fatal(err)
			}
			if secret == nil || secret.Data == nil {
				t.Fatal("failed to read the value")
			}
		}
	}
	readKV(leaderClient)

	// Take a snapshot
	buf := new(bytes.Buffer)
	err = leaderClient.Sys().RaftSnapshot(buf)
	if err != nil {
		t.Fatal(err)
	}
	snap, err := io.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(snap) == 0 {
		t.Fatal("no snapshot returned")
	}

	// start a new cluster with raft-wal as its logStore
	if opts.ClusterOptions.VaultNodeConfig.StorageOptions == nil {
		opts.ClusterOptions.VaultNodeConfig.StorageOptions = make(map[string]string, 0)
	}
	opts.ClusterOptions.VaultNodeConfig.StorageOptions["raft_wal"] = "true"

	// caching the old cluster's barrier keys
	oldBarrierKeys := cluster.GetBarrierKeys()
	// clean up the old cluster as there is no further use to it
	cluster.Cleanup()

	// Start a new cluster, set the old cluster's barrier keys as its own, and restore
	// the snapshot from the old cluster
	newCluster := docker.NewTestDockerCluster(t, opts)
	defer newCluster.Cleanup()

	// get the leader client
	newLeaderNode := newCluster.GetActiveClusterNode()
	newLeaderClient := newLeaderNode.APIClient()

	// set the barrier keys to the old cluster so that we could restore the snapshot
	newCluster.SetBarrierKeys(oldBarrierKeys)

	// Restore snapshot
	err = newLeaderClient.Sys().RaftSnapshotRestore(bytes.NewReader(snap), true)
	if err != nil {
		t.Fatal(err)
	}

	if err = testcluster.UnsealNode(ctx, newCluster, 0); err != nil {
		t.Fatal(err)
	}
	testcluster.WaitForActiveNode(ctx, newCluster)

	// generate a root token as the unseal keys have changed
	rootToken, err := testcluster.GenerateRoot(newCluster, testcluster.GenerateRootRegular)
	if err != nil {
		t.Fatal(err)
	}
	newLeaderClient.SetToken(rootToken)

	// stabilize the cluster
	stabilize(t, newLeaderClient)
	// check all data exists
	readKV(newLeaderClient)
}
