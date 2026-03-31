// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestVaultLicenseStatus verifies Vault license status response
func TestVaultLicenseStatus(t *testing.T) {
	v := blackbox.New(t)

	// Read the sys/license/status endpoint which should contain license info
	licenseStatus := v.MustRead("sys/license/status")
	if licenseStatus.Data["autoloaded"] == nil {
		t.Fatal("Could not get license details from sys/license/status")
	}
	autoloaded := licenseStatus.Data["autoloaded"].(map[string]interface{})
	if autoloaded["license_id"].(string) == "" {
		t.Fatal("Could not retrieve license_id from sys/license/status")
	}
}
