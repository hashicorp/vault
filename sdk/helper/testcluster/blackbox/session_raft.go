// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) AssertRaftStable(numNodes int, allowNonVoters bool) {
	s.t.Helper()

	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/storage/raft/configuration")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	assertions := s.AssertSecret(secret).
		Data().
		GetMap("config").
		GetSlice("servers").
		Length(numNodes)

	if !allowNonVoters {
		assertions.AllHaveKey("voter", true)
	}
}

func (s *Session) AssertRaftHealthy() {
	s.t.Helper()

	// query the autopilot state endpoint and verify that all nodes are healthy according to autopilot
	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	_ = s.AssertSecret(secret).
		Data().
		HasKey("healthy", true)
}

// AssertRaftClusterHealthy verifies that the raft cluster is healthy regardless of node count
// This is a more flexible alternative to AssertRaftStable for cases where you don't know
// or don't care about the exact cluster size, just that it's working properly.
func (s *Session) AssertRaftClusterHealthy() {
	s.t.Helper()

	// First verify autopilot reports the cluster as healthy
	s.AssertRaftHealthy()

	// Get raft configuration to ensure we have at least one node
	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/storage/raft/configuration")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	// Verify we have at least one server configured
	servers := s.AssertSecret(secret).
		Data().
		GetMap("config").
		GetSlice("servers")

	// Ensure we have at least 1 server
	if len(servers.data) < 1 {
		s.t.Fatal("Expected at least 1 raft server, got 0")
	}

	// Verify that we have at least one voter in the cluster
	hasVoter := false
	for _, server := range servers.data {
		if serverMap, ok := server.(map[string]any); ok {
			if voter, exists := serverMap["voter"]; exists {
				if voterBool, ok := voter.(bool); ok && voterBool {
					hasVoter = true
					break
				}
			}
		}
	}

	if !hasVoter {
		s.t.Fatal("Expected at least one voter in the raft cluster")
	}
}

func (s *Session) MustRaftRemovePeer(nodeID string) {
	s.t.Helper()

	_, err := s.Client.Logical().Write("sys/storage/raft/remove-peer", map[string]any{
		"server_id": nodeID,
	})
	require.NoError(s.t, err)
}

func (s *Session) AssertRaftPeerRemoved(nodeID string) {
	s.t.Helper()

	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/storage/raft/configuration")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	_ = s.AssertSecret(secret).
		Data().
		GetMap("config").
		GetSlice("servers").
		NoneHaveKeyVal("node_id", nodeID)
}

// MustGetCurrentLeader returns the current leader's node ID
func (s *Session) MustGetCurrentLeader() string {
	s.t.Helper()

	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/leader")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	leaderAddress, ok := secret.Data["leader_address"].(string)
	require.True(s.t, ok, "leader_address not found or not a string")
	require.NotEmpty(s.t, leaderAddress, "leader_address is empty")

	return leaderAddress
}

// MustStepDownLeader forces the current leader to step down
func (s *Session) MustStepDownLeader() {
	s.t.Helper()

	_, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Write("sys/step-down", nil)
	})

	require.NoError(s.t, err)
}

// GetClusterNodeCount returns the number of nodes in the raft cluster
func (s *Session) GetClusterNodeCount() int {
	s.t.Helper()

	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/storage/raft/configuration")
	})
	if err != nil {
		s.t.Logf("Failed to read raft configuration: %v", err)
		return 0
	}

	if secret == nil {
		s.t.Log("Raft configuration response was nil")
		return 0
	}

	configData, ok := secret.Data["config"].(map[string]any)
	if !ok {
		s.t.Log("Could not parse raft config data")
		return 0
	}

	serversData, ok := configData["servers"].([]any)
	if !ok {
		s.t.Log("Could not parse raft servers data")
		return 0
	}

	return len(serversData)
}

// WaitForNewLeader waits for a new leader to be elected that is different from initialLeader
// and for the cluster to become healthy. For single-node clusters, it just waits for the
// cluster to become healthy again after stepdown.
func (s *Session) WaitForNewLeader(initialLeader string, timeoutSeconds int) {
	s.t.Helper()

	// Check cluster size to handle single-node case
	nodeCount := s.GetClusterNodeCount()
	if nodeCount <= 1 {
		s.t.Logf("Single-node cluster detected, waiting for cluster to recover after stepdown...")

		// For single-node clusters, just wait for the same leader to come back and be healthy
		timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				s.t.Fatalf("Timeout waiting for single-node cluster to recover after %d seconds", timeoutSeconds)
			case <-ticker.C:
				// Check if cluster is healthy again
				secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
					return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
				})
				if err != nil {
					s.t.Logf("Error reading autopilot state: %v, retrying...", err)
					continue
				}

				if secret == nil {
					s.t.Logf("No autopilot state returned, retrying...")
					continue
				}

				healthy, ok := secret.Data["healthy"].(bool)
				if !ok {
					s.t.Logf("Autopilot healthy status not found, retrying...")
					continue
				}

				if healthy {
					s.t.Log("Single-node cluster has recovered and is healthy")
					return
				} else {
					s.t.Logf("Single-node cluster not yet healthy, waiting...")
				}
			}
		}
	}

	// Multi-node cluster logic - wait for actual leader change
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	newLeaderFound := false
	var currentLeader string

	for {
		select {
		case <-timeout:
			if newLeaderFound {
				s.t.Fatalf("Timeout waiting for cluster to become healthy after %d seconds (new leader: %s)", timeoutSeconds, currentLeader)
			} else {
				s.t.Fatalf("Timeout waiting for new leader election after %d seconds", timeoutSeconds)
			}
		case <-ticker.C:
			// First, check if a new leader has been elected
			if !newLeaderFound {
				secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
					return s.Client.Logical().Read("sys/leader")
				})
				if err != nil {
					s.t.Logf("Error reading leader status: %v, retrying...", err)
					continue
				}

				if secret == nil {
					s.t.Logf("No leader data returned, retrying...")
					continue
				}

				leaderAddress, ok := secret.Data["leader_address"].(string)
				if !ok || leaderAddress == "" {
					s.t.Logf("No leader address found, retrying...")
					continue
				}

				if leaderAddress != initialLeader {
					s.t.Logf("New leader elected: %s (was: %s)", leaderAddress, initialLeader)
					currentLeader = leaderAddress
					newLeaderFound = true
				} else {
					s.t.Logf("Still waiting for new leader, current: %s", leaderAddress)
					continue
				}
			}

			// Once we have a new leader, wait for cluster to be healthy
			if newLeaderFound {
				secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
					return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
				})
				if err != nil {
					s.t.Logf("Error reading autopilot state: %v, retrying...", err)
					continue
				}

				if secret == nil {
					s.t.Logf("No autopilot state returned, retrying...")
					continue
				}

				healthy, ok := secret.Data["healthy"].(bool)
				if !ok {
					s.t.Logf("Autopilot healthy status not found, retrying...")
					continue
				}

				if healthy {
					s.t.Logf("Cluster is now healthy with new leader: %s", currentLeader)
					return
				} else {
					s.t.Logf("Cluster not yet healthy, waiting...")
				}
			}
		}
	}
}

// AssertClusterHealthy verifies that the cluster is healthy, with fallback for managed environments
// like HCP where raft APIs may not be accessible. This is the recommended method for general
// cluster health checks in blackbox tests. It includes retry logic for Docker environments
// where the cluster may not be immediately ready.
func (s *Session) AssertClusterHealthy() {
	s.t.Helper()

	// For Docker environments, wait for the cluster to be ready with retry logic
	maxRetries := 30
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Try raft-based health check first (works for self-managed clusters)
		secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
			return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
		})

		if err == nil && secret != nil {
			// Check if autopilot reports healthy
			if healthy, ok := secret.Data["healthy"].(bool); ok && healthy {
				// Raft API is available and healthy, use full raft health check
				s.AssertRaftClusterHealthy()
				return
			} else if ok && !healthy {
				// Raft API available but not healthy yet, retry if we have attempts left
				if attempt < maxRetries {
					s.t.Logf("Cluster not yet healthy (attempt %d/%d), waiting %v...", attempt, maxRetries, retryDelay)
					time.Sleep(retryDelay)
					continue
				} else {
					s.t.Fatalf("Cluster failed to become healthy after %d attempts", maxRetries)
				}
			}
		}

		// Raft API not accessible or no healthy status - check basic connectivity
		sealStatus, err := s.WithRootNamespace(func() (*api.Secret, error) {
			return s.Client.Logical().Read("sys/seal-status")
		})
		if err != nil {
			if attempt < maxRetries {
				s.t.Logf("Failed to read seal status (attempt %d/%d): %v, retrying in %v...", attempt, maxRetries, err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			require.NoError(s.t, err, "Failed to read seal status - cluster may be unreachable")
		}

		if sealStatus == nil {
			if attempt < maxRetries {
				s.t.Logf("Seal status response was nil (attempt %d/%d), retrying in %v...", attempt, maxRetries, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			require.NotNil(s.t, sealStatus, "Seal status response was nil")
		}

		// Verify cluster is unsealed
		sealed, ok := sealStatus.Data["sealed"].(bool)
		if !ok {
			if attempt < maxRetries {
				s.t.Logf("Could not determine seal status (attempt %d/%d), retrying in %v...", attempt, maxRetries, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			require.True(s.t, ok, "Could not determine seal status")
		}

		if sealed {
			if attempt < maxRetries {
				s.t.Logf("Cluster is sealed (attempt %d/%d), retrying in %v...", attempt, maxRetries, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			require.False(s.t, sealed, "Cluster is sealed")
		}

		// If we get here, cluster is unsealed and responsive
		if secret != nil {
			s.t.Log("Cluster health verified (self-managed environment)")
		} else {
			s.t.Log("Cluster health verified (managed environment - raft APIs not accessible)")
		}
		return
	}
}
