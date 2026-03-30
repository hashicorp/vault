// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package replication

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestReplicationStatus verifies replication status for both DR and performance replication
func TestReplicationStatus(t *testing.T) {
	v := blackbox.New(t)

	// Read replication status with proper nil checks
	drStatus := v.MustRead("sys/replication/dr/status")
	if drStatus == nil || drStatus.Data == nil {
		t.Log("DR replication not available or not configured - skipping DR replication check")
	} else {
		if drMode, ok := drStatus.Data["mode"]; ok {
			t.Logf("DR replication mode: %v", drMode)
		} else {
			t.Log("DR replication mode not available")
		}
	}

	prStatus := v.MustRead("sys/replication/performance/status")
	if prStatus == nil || prStatus.Data == nil {
		t.Log("Performance replication not available or not configured - skipping performance replication check")
	} else {
		if prMode, ok := prStatus.Data["mode"]; ok {
			t.Logf("Performance replication mode: %v", prMode)
		} else {
			t.Log("Performance replication mode not available")
		}
	}

	t.Log("Successfully verified replication status endpoints are accessible")
}
