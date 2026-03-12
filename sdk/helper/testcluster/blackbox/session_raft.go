// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/backoff"
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
// cluster to become healthy again after stepdown. Uses reasonable timeouts to detect race conditions early.
func (s *Session) WaitForNewLeader(initialLeader string, timeoutSeconds int) {
	s.t.Helper()

	// Use reasonable timeout - if it takes more than a few seconds, there's likely a race condition
	if timeoutSeconds > 10 {
		s.t.Logf("Warning: timeout of %d seconds is quite high, consider investigating potential race conditions", timeoutSeconds)
	}

	// Check cluster size to handle single-node case
	nodeCount := s.GetClusterNodeCount()
	if nodeCount <= 1 {
		s.t.Logf("Single-node cluster detected, waiting for cluster to recover after stepdown...")

		// Use backoff helper for single-node recovery
		b := backoff.NewBackoff(20, 100*time.Millisecond, 1*time.Second) // Max ~10 seconds with backoff

		err := b.Retry(func() error {
			secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
				return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
			})
			if err != nil {
				return err
			}
			if secret == nil {
				return errors.New("no autopilot state returned")
			}

			healthy, ok := secret.Data["healthy"].(bool)
			if !ok {
				return errors.New("autopilot healthy status not found")
			}
			if !healthy {
				return errors.New("cluster not yet healthy")
			}

			return nil
		})
		if err != nil {
			s.t.Fatalf("Single-node cluster failed to recover: %v", err)
		}

		s.t.Log("Single-node cluster has recovered and is healthy")
		return
	}

	// Multi-node cluster logic - wait for actual leader change
	s.t.Logf("Multi-node cluster detected, waiting for new leader election...")

	// Phase 1: Wait for new leader (should be fast)
	leaderBackoff := backoff.NewBackoff(20, 100*time.Millisecond, 500*time.Millisecond) // Max ~5 seconds
	var currentLeader string

	err := leaderBackoff.Retry(func() error {
		secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
			return s.Client.Logical().Read("sys/leader")
		})
		if err != nil {
			return err
		}
		if secret == nil {
			return errors.New("no leader data returned")
		}

		leaderAddress, ok := secret.Data["leader_address"].(string)
		if !ok || leaderAddress == "" {
			return errors.New("no leader address found")
		}

		if leaderAddress == initialLeader {
			return errors.New("still waiting for new leader")
		}

		currentLeader = leaderAddress
		return nil
	})
	if err != nil {
		s.t.Fatalf("Failed to elect new leader: %v", err)
	}

	s.t.Logf("New leader elected: %s (was: %s)", currentLeader, initialLeader)

	// Phase 2: Wait for cluster health (should also be fast)
	healthBackoff := backoff.NewBackoff(20, 100*time.Millisecond, 500*time.Millisecond) // Max ~5 seconds

	err = healthBackoff.Retry(func() error {
		secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
			return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
		})
		if err != nil {
			return err
		}
		if secret == nil {
			return errors.New("no autopilot state returned")
		}

		healthy, ok := secret.Data["healthy"].(bool)
		if !ok {
			return errors.New("autopilot healthy status not found")
		}
		if !healthy {
			return errors.New("cluster not yet healthy")
		}

		return nil
	})
	if err != nil {
		s.t.Fatalf("Cluster failed to become healthy with new leader: %v", err)
	}

	s.t.Logf("Cluster is now healthy with new leader: %s", currentLeader)
}

// AssertClusterHealthy verifies that the cluster is healthy, with fallback for managed environments
// like HCP where raft APIs may not be accessible. This is the recommended method for general
// cluster health checks in blackbox tests. Uses backoff helper for reasonable retry logic.
func (s *Session) AssertClusterHealthy() {
	s.t.Helper()

	// Use backoff helper for cluster readiness checks
	b := backoff.NewBackoff(15, 200*time.Millisecond, 2*time.Second) // Max ~15 seconds with backoff

	err := b.Retry(func() error {
		// Try raft-based health check first (works for self-managed clusters)
		secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
			return s.Client.Logical().Read("sys/storage/raft/autopilot/state")
		})

		if err == nil && secret != nil {
			// Check if autopilot reports healthy
			if healthy, ok := secret.Data["healthy"].(bool); ok && healthy {
				// Raft API is available and healthy, use full raft health check
				s.AssertRaftClusterHealthy()
				return nil
			} else if ok && !healthy {
				return errors.New("cluster not yet healthy according to autopilot")
			}
		}

		// Raft API not accessible or no healthy status - check basic connectivity
		sealStatus, err := s.WithRootNamespace(func() (*api.Secret, error) {
			return s.Client.Logical().Read("sys/seal-status")
		})
		if err != nil {
			return err
		}

		if sealStatus == nil {
			return errors.New("seal status response was nil")
		}

		// Verify cluster is unsealed
		sealed, ok := sealStatus.Data["sealed"].(bool)
		if !ok {
			return errors.New("could not determine seal status")
		}

		if sealed {
			return errors.New("cluster is sealed")
		}

		// If we get here, cluster is unsealed and responsive
		if secret != nil {
			s.t.Log("Cluster health verified (self-managed environment)")
		} else {
			s.t.Log("Cluster health verified (managed environment - raft APIs not accessible)")
		}
		return nil
	})
	if err != nil {
		s.t.Fatalf("Cluster health check failed: %v", err)
	}
}
