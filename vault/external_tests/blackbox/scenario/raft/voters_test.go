//go:build scenario
// +build scenario

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestRaftVoters verifies that all nodes in the raft cluster are voters
func TestRaftVoters(t *testing.T) {
	v := blackbox.New(t)

	// Wait for a healthy cluster
	v.EventuallyClusterHealthyUnsealed(15 * time.Second)

	storage := v.MustGetConfigStorageType()
	if storage != "raft" {
		t.Log("skipping as cluster is not using integrated storage")
		return
	}

	v.EventuallyRaftClusterHealthy(5 * time.Second)

	t.Log("Successfully verified raft cluster is healthy with at least one voter")
}
