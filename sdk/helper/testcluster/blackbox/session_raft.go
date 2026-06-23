// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
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

func (s *Session) raftConfig() (*api.Secret, error) {
	// TODO
	// query the autopilot state endpoint and verify that all nodes are healthy according to autopilot
	var state *api.Secret
	return state, s.Req(
		func(c *api.Client) error {
			var err error
			state, err = c.Logical().Read("sys/storage/raft/autopilot/state")
			return err
		},
		WithClientRootNamespace(),
		WithClientTimeout(2*time.Second),
	)
}

// EventuallyRaftClusterHealthy verifies that the raft cluster eventually becomes
// healthy regardless of node count.
func (s *Session) EventuallyRaftClusterHealthy(timeout time.Duration) {
	s.t.Helper()

	s.EventuallyWithTimeout(
		func() error {
			healthy, err := s.autopilotStateHealthy()
			if err != nil {
				return err
			}

			if !healthy {
				return fmt.Errorf("expected raft state to be healthy: got %t", healthy)
			}

			return nil
		}, timeout,
	)

	// TODO: Perhaps move the server config and voter checks into the retry above.

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

// AssertRaftClusterHealthy verifies that the raft cluster is healthy regardless of node count
// This is a more flexible alternative to AssertRaftStable for cases where you don't know
// or don't care about the exact cluster size, just that it's working properly.
func (s *Session) AssertRaftClusterHealthy() {
	s.t.Helper()

	// First verify autopilot reports the cluster as healthy
	s.AssertAutopilotHealthy()

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

	_, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Write("sys/storage/raft/remove-peer", map[string]any{
			"server_id": nodeID,
		})
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

// MustGetClusterNodeCount returns the number of nodes in the cluster
func (s *Session) MustGetClusterNodeCount() int {
	s.t.Helper()

	count, err := s.getClusterNodeCount()
	require.NoError(s.t, err)

	return count
}

// MustGetNonLeaderNode returns a non-leader node ID from the raft cluster
func (s *Session) MustGetNonLeaderNode() string {
	s.t.Helper()

	// Get current leader
	leader := s.MustGetCurrentLeader()

	// Get raft configuration
	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/storage/raft/configuration")
	})
	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	configData, ok := secret.Data["config"].(map[string]any)
	require.True(s.t, ok, "Could not parse raft config data")

	serversData, ok := configData["servers"].([]any)
	require.True(s.t, ok, "Could not parse raft servers data")

	// Find a non-leader node
	for _, server := range serversData {
		serverMap, ok := server.(map[string]any)
		if !ok {
			continue
		}

		nodeID, ok := serverMap["node_id"].(string)
		if !ok {
			continue
		}

		// Check if this is not the leader
		address, ok := serverMap["address"].(string)
		if ok && address != leader {
			return nodeID
		}
	}

	s.t.Fatal("Could not find a non-leader node in the cluster")
	return ""
}
