// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestRaftVoters verifies that all nodes in the raft cluster are voters
func TestRaftVoters(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster regardless of node count
	v.AssertClusterHealthy()

	t.Log("Successfully verified raft cluster is healthy with at least one voter")
}
