// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package verify

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestVaultServerVersion verifies the Vault server version via sys/version-history API
// This test runs from CI/GitHub runners and connects to the Vault cluster via API
func TestVaultServerVersion(t *testing.T) {
	t.Parallel()

	version := os.Getenv("VAULT_VERSION")
	if version == "" {
		t.Fatal("VAULT_VERSION environment variable is required")
	}

	buildDate := os.Getenv("VAULT_BUILD_DATE")
	if buildDate == "" {
		t.Fatal("VAULT_BUILD_DATE environment variable is required")
	}

	v := blackbox.New(t)
	replicationMode, err := v.GetDRReplicationMode()
	require.NoError(t, err)
	t.Logf("Found DR Replication mode: %s", replicationMode)
	switch replicationMode {
	case "primary", "secondary":
		t.Skip("Skipping server version check on DR cluster as path is not available in that mode")
	default:
		v.AssertVersion(version)
		v.AssertBuildDate(version, buildDate)
	}
}
