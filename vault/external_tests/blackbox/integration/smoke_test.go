// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package integration

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestStepdownAndLeaderElection tests raft leadership changes by forcing the
// current leader to step down and verifying that a new leader is elected
// successfully, while ensuring the cluster remains healthy throughout the
// process.
//
// NOTE: We do not run this test in parallel to avoid making other tests flaky
// during leader elections.
func TestStepdownAndLeaderElection(t *testing.T) {
	v := blackbox.New(t)

	// Wait for a healthy cluster
	v.EventuallyClusterHealthyUnsealed(15 * time.Second)

	// Check cluster size to determine expected behavior
	nodeCount := v.MustGetClusterNodeCount()
	t.Logf("Cluster has %d nodes", nodeCount)

	// Get current leader before step down
	initialLeader := v.MustGetCurrentLeader()
	t.Logf("Initial leader: %s", initialLeader)

	// Force leader to step down
	v.MustStepDownLeader()

	// Wait for a new leader to be active and the cluster to be healthy
	v.EventuallyClusterHealthyUnsealed(1 * time.Minute)

	// Get current leader before step down
	newLeader := v.MustGetCurrentLeader()
	t.Logf("New leader: %s", initialLeader)

	// For multi-node clusters, verify new leader is different from initial leader
	// For single-node clusters, just verify it's healthy again
	if nodeCount > 1 {
		require.NotEqual(t, initialLeader, newLeader, "Expected new leader to be different from initial leader")
		t.Logf("Successfully elected new leader: %s (was: %s)", newLeader, initialLeader)
	} else {
		t.Logf("Single-node cluster successfully recovered with leader: %s", newLeader)
	}
}
