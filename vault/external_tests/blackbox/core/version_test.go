// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestVaultVersion verifies Vault version endpoint accessibility and response
func TestVaultVersion(t *testing.T) {
	v := blackbox.New(t)

	// Read the sys/seal-status endpoint which should contain version info
	sealStatus := v.MustRead("sys/seal-status")
	if sealStatus.Data["version"] == nil {
		t.Fatal("Could not retrieve version from sys/seal-status")
	}

	t.Logf("Vault version: %v", sealStatus.Data["version"])
}
