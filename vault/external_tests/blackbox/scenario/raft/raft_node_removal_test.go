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

const (
	clusterHealthTimeout = 15 * time.Second
	raftHealthTimeout    = 5 * time.Second
)

// TestRaft_ClusterHealthVerification verifies raft cluster health.
// Node removal is handled by enos scenarios, not this test.
func TestRaft_ClusterHealthVerification(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Wait for a healthy cluster
	v.EventuallyClusterHealthyUnsealed(clusterHealthTimeout)

	// Check if using raft storage
	storage := v.MustGetConfigStorageType()
	if storage != "raft" {
		t.Log("skipping as cluster is not using integrated storage")
		return
	}

	// Verify raft cluster is healthy
	v.EventuallyRaftClusterHealthy(raftHealthTimeout)

	t.Log("Successfully verified raft cluster is healthy")
}

// TestRaft_RemovedNodeStatus verifies removed node status.
// Run on a removed node via enos scenario. Skips if node not removed.
func TestRaft_RemovedNodeStatus(t *testing.T) {
	v := setupRemovedNodeTest(t)
	testRemovedNodeStatus(t, v)
}

func testRemovedNodeStatus(t *testing.T, v *blackbox.Session) {
	t.Helper()
	verifyRemovedStatus(t, v)
	verifyRemovedHealth(t, v)
	verifyManualRejoinFails(t, v)
}

// TestRaft_RemovedNodeAfterRestart verifies removed status persists after restart.
// Run on a removed node after restart via enos scenario. Skips if node not removed.
func TestRaft_RemovedNodeAfterRestart(t *testing.T) {
	v := setupRemovedNodeTest(t)
	testRemovedNodeAfterRestart(t, v)
}

func testRemovedNodeAfterRestart(t *testing.T, v *blackbox.Session) {
	t.Helper()
	verifyRemovedStatus(t, v)
	verifyRemovedHealth(t, v)
}

// setupRemovedNodeTest creates a session and skips if node not removed.
func setupRemovedNodeTest(t *testing.T) *blackbox.Session {
	t.Helper()
	t.Parallel()
	v := blackbox.New(t)
	if !isNodeRemoved(t, v) {
		t.Skip("Skipping - node not removed")
	}
	return v
}

// isNodeRemoved checks if node has been removed from the cluster.
func isNodeRemoved(t *testing.T, v *blackbox.Session) bool {
	t.Helper()
	resp, err := v.Client.Sys().SealStatus()
	if err != nil {
		t.Fatalf("Failed to get seal status: %v", err)
	}
	return resp.RemovedFromCluster != nil && *resp.RemovedFromCluster
}

// verifyRemovedStatus checks vault status shows removed_from_cluster=true.
func verifyRemovedStatus(t *testing.T, v *blackbox.Session) {
	t.Helper()
	v.EventuallyWithTimeout(func() error {
		resp, err := v.Client.Sys().SealStatus()
		if err != nil {
			return fmt.Errorf("failed to get vault status: %w", err)
		}
		// Verify removed_from_cluster field is set and true
		if resp.RemovedFromCluster == nil || !*resp.RemovedFromCluster {
			return fmt.Errorf("status shows removed_from_cluster=%v, expected true", resp.RemovedFromCluster)
		}
		t.Log("✓ Vault status shows removed_from_cluster=true")
		return nil
	}, 15*time.Second)
}

// verifyRemovedHealth checks sys/health shows removed_from_cluster=true.
func verifyRemovedHealth(t *testing.T, v *blackbox.Session) {
	t.Helper()
	v.EventuallyWithTimeout(func() error {
		// Use custom status codes to avoid errors from sealed/uninitialized states
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
		// Verify removed_from_cluster field is set and true
		if healthResp.RemovedFromCluster == nil || !*healthResp.RemovedFromCluster {
			return fmt.Errorf("health shows removed_from_cluster=%v, expected true", healthResp.RemovedFromCluster)
		}
		t.Log("✓ Health endpoint shows removed_from_cluster=true")
		return nil
	}, 15*time.Second)
}

// verifyManualRejoinFails checks that raft join attempts fail on removed nodes.
func verifyManualRejoinFails(t *testing.T, v *blackbox.Session) {
	t.Helper()
	leaderAddr := v.Client.Address()
	if leaderAddr == "" {
		t.Skip("Leader address not available, skipping manual rejoin test")
	}

	// Attempt to rejoin the cluster - should fail for removed nodes
	req := v.Client.NewRequest("POST", "/v1/sys/storage/raft/join")
	req.BodyBytes = []byte(fmt.Sprintf(`{"leader_api_addr": "%s"}`, leaderAddr))

	resp, err := v.Client.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
		// Success response (2xx) means rejoin worked, which shouldn't happen
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			t.Fatal("Expected raft join to fail, but it succeeded")
		}
		t.Logf("✓ Raft join failed as expected with status code: %d", resp.StatusCode)
		return
	}

	// Verify error mentions removed node
	if !strings.Contains(err.Error(), "removed") {
		t.Fatalf("Expected error about removed node, got: %v", err)
	}
	t.Logf("✓ Raft join failed as expected with error: %v", err)
}

// TestRaft_UnsealFailsOnRemovedNode verifies unseal fails on removed nodes (Shamir seal only).
// Run on a removed node via enos scenario. Skips if node not removed or not Shamir seal.
func TestRaft_UnsealFailsOnRemovedNode(t *testing.T) {
	v := setupRemovedNodeTest(t)

	status, err := v.Client.Sys().SealStatus()
	if err != nil {
		t.Fatalf("Failed to get seal status: %v", err)
	}
	if status.Type != "shamir" {
		t.Skip("Skipping unseal test - only applicable for Shamir seal")
	}

	verifyUnsealFails(t, v)
}

// TestRaft_RemovedNodeShimVerification runs all removed node verifications in sequence.
// Equivalent to vault_verify_removed_node_shim enos module. Ensures all checks pass together
// and provides single test invocation for enos scenarios.
func TestRaft_RemovedNodeShimVerification(t *testing.T) {
	v := setupRemovedNodeTest(t)

	t.Log("Running comprehensive removed node verification")

	t.Log("Step 1: Verifying removed status...")
	verifyRemovedStatus(t, v)
	verifyRemovedHealth(t, v)

	t.Log("Step 2: Verifying manual rejoin fails...")
	verifyManualRejoinFails(t, v)

	status, err := v.Client.Sys().SealStatus()
	if err != nil {
		t.Fatalf("Failed to get seal status: %v", err)
	}

	if status.Type == "shamir" {
		t.Log("Step 3: Verifying unseal fails (Shamir seal)...")
		verifyUnsealFails(t, v)
	} else {
		t.Logf("Step 3: Skipping unseal verification (seal type: %s)", status.Type)
	}

	t.Log("✓ All removed node verifications passed")
}

// verifyUnsealFails checks that unseal operations fail on removed nodes.
func verifyUnsealFails(t *testing.T, v *blackbox.Session) {
	t.Helper()

	// Attempt unseal with dummy key - should fail on removed node
	req := v.Client.NewRequest("POST", "/v1/sys/unseal")
	req.BodyBytes = []byte(`{"key": "dummy-key-for-testing"}`)

	resp, err := v.Client.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
		// Success response (2xx) means unseal worked, which shouldn't happen
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			t.Fatal("Expected unseal to fail on removed node, but it succeeded")
		}
		t.Logf("✓ Unseal failed as expected with status code: %d", resp.StatusCode)
		return
	}

	// Verify error mentions removed node - check ResponseError type for better validation
	if respErr, ok := err.(*api.ResponseError); ok {
		errMsg := respErr.Error()
		if !strings.Contains(strings.ToLower(errMsg), "removed") &&
			!strings.Contains(strings.ToLower(errMsg), "not part of raft") {
			t.Fatalf("Expected error about removed node, got: %v (status code: %d)", errMsg, respErr.StatusCode)
		}
		t.Logf("✓ Unseal failed as expected with ResponseError (status %d): %v", respErr.StatusCode, errMsg)
		return
	}

	// Fallback for non-ResponseError types
	errMsg := err.Error()
	if !strings.Contains(strings.ToLower(errMsg), "removed") &&
		!strings.Contains(strings.ToLower(errMsg), "not part of raft") {
		t.Fatalf("Expected error about removed node, got: %v", err)
	}
	t.Logf("✓ Unseal failed as expected with error: %v", err)
}
