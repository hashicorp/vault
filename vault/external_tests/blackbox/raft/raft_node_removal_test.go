// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestRaft_ClusterHealthVerification verifies raft cluster health
// Note: This test does NOT perform node removal. Node removal is handled by the
// enos scenario infrastructure. This test simply verifies the cluster is healthy
// and using raft storage, similar to other raft tests in this package.
func TestRaft_ClusterHealthVerification(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Wait for a healthy cluster
	v.EventuallyClusterHealthyUnsealed(15 * time.Second)

	// Check if using raft storage
	storage := v.MustGetConfigStorageType()
	if storage != "raft" {
		t.Log("skipping as cluster is not using integrated storage")
		return
	}

	// Verify raft cluster is healthy
	v.EventuallyRaftClusterHealthy(5 * time.Second)

	t.Log("Successfully verified raft cluster is healthy")
}

// TestRaft_RemovedNodeStatus tests the status of a removed node
// This test should ONLY be run ON a removed node itself (via enos scenario)
// It verifies that the removed node correctly reports its status
// If the node is not removed, the test is skipped
func TestRaft_RemovedNodeStatus(t *testing.T) {
	v := setupRemovedNodeTest(t)

	testRemovedNodeStatus(t, v)
}

// testRemovedNodeStatus runs all removed node status checks
func testRemovedNodeStatus(t *testing.T, v *blackbox.Session) {
	t.Helper()
	verifyRemovedStatus(t, v)
	verifyRemovedHealth(t, v)
	verifyManualRejoinFails(t, v)
}

// TestRaft_RemovedNodeAfterRestart verifies removed status persists after restart
// This test should ONLY be run ON a removed node after it has been restarted
// If the node is not removed, the test is skipped
func TestRaft_RemovedNodeAfterRestart(t *testing.T) {
	v := setupRemovedNodeTest(t)

	testRemovedNodeAfterRestart(t, v)
}

// testRemovedNodeAfterRestart verifies removed status persists after restart
func testRemovedNodeAfterRestart(t *testing.T, v *blackbox.Session) {
	t.Helper()
	verifyRemovedStatus(t, v)
	verifyRemovedHealth(t, v)
}

// setupRemovedNodeTest sets up a test for a removed node
// It checks if the node is removed and skips if not
func setupRemovedNodeTest(t *testing.T) *blackbox.Session {
	t.Helper()
	t.Parallel()
	v := blackbox.New(t)
	if !isNodeRemoved(t, v) {
		t.Skip("Skipping - node not removed")
	}
	return v
}

// isNodeRemoved checks if the current node has been removed from the cluster
// Returns true if the node is removed, fails the test on error
func isNodeRemoved(t *testing.T, v *blackbox.Session) bool {
	t.Helper()

	// Get vault status
	resp, err := v.Client.Sys().SealStatus()
	if err != nil {
		t.Fatalf("Failed to get seal status: %v", err)
	}

	// Check if removed_from_cluster is set and true
	return resp.RemovedFromCluster != nil && *resp.RemovedFromCluster
}

// verifyRemovedStatus verifies that vault status shows removed_from_cluster=true
func verifyRemovedStatus(t *testing.T, v *blackbox.Session) {
	t.Helper()

	v.EventuallyWithTimeout(func() error {
		// Get vault status
		resp, err := v.Client.Sys().SealStatus()
		if err != nil {
			return fmt.Errorf("failed to get vault status: %w", err)
		}

		// Check if removed_from_cluster is true
		if resp.RemovedFromCluster == nil || !*resp.RemovedFromCluster {
			return fmt.Errorf("status shows removed_from_cluster=%v, expected true", resp.RemovedFromCluster)
		}

		t.Log("✓ Vault status shows removed_from_cluster=true")
		return nil
	}, 15*time.Second)
}

// verifyRemovedHealth verifies that sys/health shows removed_from_cluster=true
func verifyRemovedHealth(t *testing.T, v *blackbox.Session) {
	t.Helper()

	v.EventuallyWithTimeout(func() error {
		// Query sys/health with custom status codes to avoid errors
		req := v.Client.NewRequest("GET", "/v1/sys/health")
		req.Params.Set("sealedcode", "299")
		req.Params.Set("uninitcode", "299")
		req.Params.Set("removedcode", "299")

		resp, err := v.Client.RawRequest(req)
		if err != nil {
			return fmt.Errorf("failed to get health status: %w", err)
		}
		defer resp.Body.Close()

		var healthResp api.HealthResponse
		if err := resp.DecodeJSON(&healthResp); err != nil {
			return fmt.Errorf("failed to decode health response: %w", err)
		}

		// Check if removed_from_cluster is true
		if healthResp.RemovedFromCluster == nil || !*healthResp.RemovedFromCluster {
			return fmt.Errorf("health shows removed_from_cluster=%v, expected true", healthResp.RemovedFromCluster)
		}

		t.Log("✓ Health endpoint shows removed_from_cluster=true")
		return nil
	}, 15*time.Second)
}

// verifyManualRejoinFails verifies that manual raft join attempts fail on removed nodes
func verifyManualRejoinFails(t *testing.T, v *blackbox.Session) {
	t.Helper()

	// Get the leader address - in enos this would be passed via environment
	leaderAddr := v.Client.Address()
	if leaderAddr == "" {
		t.Skip("Leader address not available, skipping manual rejoin test")
	}

	// Attempt to join the raft cluster - this should fail
	req := v.Client.NewRequest("POST", "/v1/sys/storage/raft/join")
	req.BodyBytes = []byte(fmt.Sprintf(`{"leader_api_addr": "%s"}`, leaderAddr))

	resp, err := v.Client.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
		// If no error, check the status code - should not be 2xx
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			t.Fatal("Expected raft join to fail, but it succeeded")
		}
		t.Logf("✓ Raft join failed as expected with status code: %d", resp.StatusCode)
		return
	}

	// Error is expected - verify it mentions removed node
	if !strings.Contains(err.Error(), "removed") {
		t.Fatalf("Expected error about removed node, got: %v", err)
	}
	t.Logf("✓ Raft join failed as expected with error: %v", err)
}

// TestRaft_UnsealFailsOnRemovedNode verifies that unseal fails on removed nodes (Shamir seal only)
// This test should ONLY be run ON a removed node (via enos scenario)
// If the node is not removed, the test is skipped
func TestRaft_UnsealFailsOnRemovedNode(t *testing.T) {
	v := setupRemovedNodeTest(t)

	// Check if this is a Shamir seal
	status, err := v.Client.Sys().SealStatus()
	if err != nil {
		t.Fatalf("Failed to get seal status: %v", err)
	}

	if status.Type != "shamir" {
		t.Skip("Skipping unseal test - only applicable for Shamir seal")
	}

	// Attempt to unseal - this should fail on a removed node
	// Note: We use a dummy key since we don't have the actual unseal keys
	req := v.Client.NewRequest("POST", "/v1/sys/unseal")
	req.BodyBytes = []byte(`{"key": "dummy-key-for-testing"}`)

	resp, err := v.Client.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
		// Check if the response indicates failure
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			t.Fatal("Expected unseal to fail on removed node, but it succeeded")
		}
		t.Logf("✓ Unseal failed as expected with status code: %d", resp.StatusCode)
		return
	}

	// Error is expected - verify it mentions removed node
	if !strings.Contains(err.Error(), "removed") {
		t.Fatalf("Expected error about removed node, got: %v", err)
	}
	t.Logf("✓ Unseal failed as expected with error: %v", err)
}
