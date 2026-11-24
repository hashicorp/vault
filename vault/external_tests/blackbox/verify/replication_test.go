// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package verify

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestReplicationAvailability verifies replication status based on Vault edition.
// For CE: replication mode should be "disabled"
// For ENT: DR and performance replication should be available
func TestReplicationAvailability(t *testing.T) {
	v := blackbox.New(t)

	edition := os.Getenv("VAULT_EDITION")
	if edition == "" {
		t.Skip("VAULT_EDITION not set, skipping edition-specific replication checks")
	}

	// Read overall replication status (must use root namespace)
	status, err := v.WithRootNamespace(func() (*api.Secret, error) {
		return v.Client.Logical().Read("sys/replication/status")
	})
	if err != nil {
		t.Fatalf("Failed to read replication status: %v", err)
	}
	if status == nil || status.Data == nil {
		t.Fatal("Failed to read replication status: response was nil")
	}

	// Log the full status for debugging
	t.Logf("Replication status for edition %s: %+v", edition, status.Data)

	if edition == "ce" {
		// For CE, replication mode should be disabled
		mode, ok := status.Data["mode"].(string)
		if !ok {
			t.Fatal("replication mode field not found or not a string")
		}
		if mode != "disabled" {
			t.Fatalf("replication data mode is not disabled for CE release! Got: %s", mode)
		}
		t.Log("Successfully verified replication is disabled for CE edition")
	} else {
		// For ENT, DR and performance replication should be available
		drData, drOk := status.Data["dr"]
		if !drOk {
			t.Fatalf("DR replication field not found for ENT release %s! Full status: %+v", edition, status.Data)
		}
		if drData == nil {
			t.Fatalf("DR replication data is nil for ENT release %s! Full status: %+v", edition, status.Data)
		}

		perfData, perfOk := status.Data["performance"]
		if !perfOk {
			t.Fatalf("Performance replication field not found for ENT release %s! Full status: %+v", edition, status.Data)
		}
		if perfData == nil {
			t.Fatalf("Performance replication data is nil for ENT release %s! Full status: %+v", edition, status.Data)
		}

		t.Logf("Successfully verified DR and performance replication are available for %s edition", edition)
	}
}
