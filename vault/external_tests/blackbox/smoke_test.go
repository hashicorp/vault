// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestEnosSmoke performs comprehensive smoke testing for Enos scenarios,
// verifying cluster health, replication status, raft stability, and basic
// KV operations with authentication. This test validates core functionality.
func TestEnosSmoke(t *testing.T) {
	v := blackbox.New(t)

	v.AssertUnsealedAny()
	v.AssertDRReplicationStatus("primary")
	v.AssertPerformanceReplicationStatus("disabled")
	v.AssertRaftStable(3, false)
	v.AssertRaftHealthy()

	// Setup using common utilities
	bob := SetupStandardKVUserpass(v, "secret", "bob", "lol")

	// Write and verify standard test data
	v.MustWriteKV2("secret", "app-config", StandardKVData)

	secret := bob.MustReadKV2("secret", "app-config")
	AssertKVData(t, bob, secret, StandardKVData)
}

// TestStepdownAndLeaderElection tests raft leadership changes by forcing the current
// leader to step down and verifying that a new leader is elected successfully,
// while ensuring the cluster remains healthy throughout the process.
func TestStepdownAndLeaderElection(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy raft cluster first
	v.AssertRaftClusterHealthy()

	// Check cluster size to determine expected behavior
	nodeCount := v.GetClusterNodeCount()
	t.Logf("Cluster has %d nodes", nodeCount)

	// Get current leader before step down
	initialLeader := v.MustGetCurrentLeader()
	t.Logf("Initial leader: %s", initialLeader)

	// Force leader to step down
	v.MustStepDownLeader()

	// Wait for new leader election (with timeout)
	v.WaitForNewLeader(initialLeader, 120)

	// Verify cluster is still healthy after leader change/recovery
	v.AssertRaftClusterHealthy()

	// For multi-node clusters, verify new leader is different from initial leader
	// For single-node clusters, just verify it's healthy again
	newLeader := v.MustGetCurrentLeader()
	if nodeCount > 1 {
		if newLeader == initialLeader {
			t.Fatalf("Expected new leader to be different from initial leader %s, got %s", initialLeader, newLeader)
		}
		t.Logf("Successfully elected new leader: %s (was: %s)", newLeader, initialLeader)
	} else {
		t.Logf("Single-node cluster successfully recovered with leader: %s", newLeader)
	}
}
