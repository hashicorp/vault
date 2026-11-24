//go:build isolated
// +build isolated

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package verify

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestDRReplicationStatus verifies DR replication status between primary and secondary clusters.
// This test validates:
// - Replication state is not idle
// - Connection status is not disconnected
// - Cluster state matches expected values (running for primary, stream-wals for secondary)
// - Known primary cluster addresses contain the expected primary leader address
func TestDRReplicationStatus(t *testing.T) {
	v := blackbox.New(t, blackbox.WithoutNamespace())

	// Get expected primary leader address from environment
	primaryLeaderAddr := os.Getenv("PRIMARY_LEADER_ADDR")
	if primaryLeaderAddr == "" {
		t.Fatal("PRIMARY_LEADER_ADDR environment variable not set")
	}

	// Read DR replication status
	drStatus, err := v.Client.Logical().Read("sys/replication/dr/status")
	if err != nil {
		t.Fatalf("Failed to read DR replication status: %v", err)
	}
	if drStatus == nil || drStatus.Data == nil {
		t.Fatal("Failed to read DR replication status")
	}

	// Check cluster state is not idle
	state, ok := drStatus.Data["state"].(string)
	if !ok {
		t.Fatal("DR replication state field not found or not a string")
	}
	if state == "idle" {
		t.Fatalf("DR replication cluster state is idle")
	}
	t.Logf("DR replication state: %s", state)

	// Get connection mode (primary or secondary)
	mode, ok := drStatus.Data["mode"].(string)
	if !ok {
		t.Fatal("DR replication mode field not found or not a string")
	}
	t.Logf("DR replication mode: %s", mode)

	if mode == "primary" {
		// Verify primary cluster
		secondaries, ok := drStatus.Data["secondaries"].([]interface{})
		if !ok || len(secondaries) == 0 {
			t.Fatal("No secondaries found in DR replication status")
		}

		secondary := secondaries[0].(map[string]interface{})
		connectionStatus, ok := secondary["connection_status"].(string)
		if !ok {
			t.Fatal("Secondary connection_status field not found or not a string")
		}
		if connectionStatus == "disconnected" {
			t.Fatalf("Secondary connection status is disconnected")
		}
		t.Logf("Secondary connection status: %s", connectionStatus)

		// Verify primary is in running state
		if state != "running" {
			t.Fatalf("Primary cluster state is not running, got: %s", state)
		}
		t.Log("Primary cluster is in running state")

	} else if mode == "secondary" {
		// Verify secondary cluster
		primaries, ok := drStatus.Data["primaries"].([]interface{})
		if !ok || len(primaries) == 0 {
			t.Fatal("No primaries found in DR replication status")
		}

		primary := primaries[0].(map[string]interface{})
		connectionStatus, ok := primary["connection_status"].(string)
		if !ok {
			t.Fatal("Primary connection_status field not found or not a string")
		}
		if connectionStatus == "disconnected" {
			t.Fatalf("Primary connection status is disconnected")
		}
		t.Logf("Primary connection status: %s", connectionStatus)

		// Verify secondary is in stream-wals state
		if state != "stream-wals" {
			t.Fatalf("Secondary cluster state is not stream-wals, got: %s", state)
		}
		t.Log("Secondary cluster is in stream-wals state")

		// Verify known primary cluster addresses contain the primary leader
		knownPrimaryAddrs, ok := drStatus.Data["known_primary_cluster_addrs"].([]interface{})
		if !ok {
			t.Fatal("known_primary_cluster_addrs field not found or not an array")
		}

		found := false
		for _, addr := range knownPrimaryAddrs {
			addrStr, ok := addr.(string)
			if !ok {
				continue
			}
			if strings.Contains(addrStr, primaryLeaderAddr) {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("Primary leader address %s not found in known_primary_cluster_addrs: %v",
				primaryLeaderAddr, knownPrimaryAddrs)
		}
		t.Logf("Verified primary leader address %s is in known_primary_cluster_addrs", primaryLeaderAddr)
	} else {
		t.Fatalf("Unexpected DR replication mode: %s", mode)
	}

	t.Log("Successfully verified DR replication status")
}

// TestDRReplicationStatusOutput outputs the full DR replication status as JSON for consumption by Terraform.
// This test is designed to be used by the vault_verify_dr_replication module to capture the replication
// status and make it available as Terraform outputs.
func TestDRReplicationStatusOutput(t *testing.T) {
	v := blackbox.New(t, blackbox.WithoutNamespace())

	// Read DR replication status
	drStatus, err := v.Client.Logical().Read("sys/replication/dr/status")
	if err != nil {
		t.Fatalf("Failed to read DR replication status: %v", err)
	}
	if drStatus == nil || drStatus.Data == nil {
		t.Fatal("Failed to read DR replication status")
	}

	// Marshal the status data to JSON
	jsonBytes, err := json.Marshal(drStatus.Data)
	if err != nil {
		t.Fatalf("Failed to marshal DR replication status to JSON: %v", err)
	}

	// Output the JSON to stdout so it can be captured by the test runner
	fmt.Println(string(jsonBytes))
}
