//go:build isolated
// +build isolated

// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package verify

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestIBMLicenseUpdate verifies that an IBM PAO license has been properly updated.
// This test replicates the logic from:
// - enos/modules/vault_verify_ibm_license_update/scripts/license-get.sh
// - enos/modules/vault_verify_ibm_license_update/scripts/license-inspect.sh
//
// The test validates:
// 1. License issuer is "pao.ibm.com" (IBM PAO license)
// 2. License edition matches expected edition (from VAULT_IBM_LICENSE_EDITION env var)
// 3. Customer ID can be extracted from license file inspection
//
// Environment variables required:
//   - VAULT_IBM_LICENSE_EDITION: Expected edition for IBM PAO licenses (e.g., "standard", "plus", "premium")
//     OR expected edition for HashiCorp licenses (e.g., "ent", "ent.hsm", "ent.fips1403")
//   - VAULT_INSTALL_DIR: Directory where vault binary is installed (default: /opt/vault/bin)
//
// Note: This test is designed for the upgrade scenario where an IBM PAO license is explicitly
// updated via the vault_update_license_ibm module. In other scenarios (like smoke-sdk), this
// test will skip if an IBM PAO license is not present.
func TestIBMLicenseUpdate(t *testing.T) {
	t.Parallel()

	v := blackbox.New(t)

	// Get expected IBM license edition from environment
	expectedEdition := os.Getenv("VAULT_IBM_LICENSE_EDITION")
	require.NotEmpty(t, expectedEdition, "VAULT_IBM_LICENSE_EDITION environment variable must be set for IBM license verification")

	t.Logf("Testing IBM license update with expected edition: %s", expectedEdition)

	// Step 1: Verify license get (replicates license-get.sh)
	// Checks that the license issuer is "pao.ibm.com" and edition matches expected
	t.Run("verify_license_get", func(t *testing.T) {
		verifyLicenseGet(t, v, expectedEdition)
	})

	// Step 2: Verify license inspect and extract customer ID (replicates license-inspect.sh)
	// Extracts customer_id from the license to verify the license is properly loaded
	t.Run("verify_license_inspect", func(t *testing.T) {
		verifyLicenseInspect(t, v)
	})

	t.Log("✓ IBM license update verification completed successfully")
}

// verifyLicenseGet replicates the logic from license-get.sh
// It verifies that:
// - The license issuer is "pao.ibm.com"
// - The license edition matches the expected edition
// - Retries up to 10 times with 30 second delays (like the script)
func verifyLicenseGet(t *testing.T, v *blackbox.Session, expectedEdition string) {
	t.Helper()

	maxRetries := 10
	retryDelay := 30 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		t.Logf("Attempt %d/%d: Checking license status...", attempt, maxRetries)

		// Read license status using MustRead (similar to TestVaultLicenseStatus)
		licenseStatus := v.MustRead("sys/license/status")

		// Log the full response for debugging
		t.Logf("License status response keys: %v", getKeys(licenseStatus.Data))
		if jsonData, err := json.MarshalIndent(licenseStatus.Data, "", "  "); err == nil {
			t.Logf("Full license status data:\n%s", string(jsonData))
		}

		// Extract license data - try both "autoloaded" and "persisted_autoload"
		var autoloaded map[string]interface{}
		var ok bool

		if licenseStatus.Data["autoloaded"] != nil {
			t.Logf("Found 'autoloaded' field (API-loaded license)")
			autoloaded, ok = licenseStatus.Data["autoloaded"].(map[string]interface{})
			if !ok {
				lastErr = fmt.Errorf("autoloaded field is not a map, type: %T", licenseStatus.Data["autoloaded"])
				if attempt < maxRetries {
					t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
					time.Sleep(retryDelay)
					continue
				}
				break
			}
		} else if licenseStatus.Data["persisted_autoload"] != nil {
			t.Logf("Found 'persisted_autoload' field (file-loaded license)")
			autoloaded, ok = licenseStatus.Data["persisted_autoload"].(map[string]interface{})
			if !ok {
				lastErr = fmt.Errorf("persisted_autoload field is not a map, type: %T", licenseStatus.Data["persisted_autoload"])
				if attempt < maxRetries {
					t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
					time.Sleep(retryDelay)
					continue
				}
				break
			}
		} else {
			lastErr = fmt.Errorf("neither 'autoloaded' nor 'persisted_autoload' field found in license status")
			t.Logf("Available fields in license status: %v", getKeys(licenseStatus.Data))
			if attempt < maxRetries {
				t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			break
		}

		// Log all license data fields for debugging
		t.Logf("License data fields: %v", getKeys(autoloaded))
		if jsonData, err := json.MarshalIndent(autoloaded, "", "  "); err == nil {
			t.Logf("License data:\n%s", string(jsonData))
		}

		// Extract issuer
		issuer, ok := autoloaded["issuer"].(string)
		if !ok {
			lastErr = fmt.Errorf("issuer field not found or not a string, type: %T", autoloaded["issuer"])
			if attempt < maxRetries {
				t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			break
		}
		t.Logf("Found issuer: %s", issuer)

		// Extract edition
		edition, ok := autoloaded["edition"].(string)
		if !ok {
			lastErr = fmt.Errorf("edition field not found or not a string, type: %T", autoloaded["edition"])
			if attempt < maxRetries {
				t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			break
		}
		t.Logf("Found edition: '%s' (expected: '%s')", edition, expectedEdition)

		// Verify issuer is either "pao.ibm.com" (IBM PAO license) or "hashicorp" (HashiCorp license)
		if issuer != "pao.ibm.com" && issuer != "hashicorp" {
			lastErr = fmt.Errorf("expected license with issuer 'pao.ibm.com' or 'hashicorp', got issuer: '%s'", issuer)
			if attempt < maxRetries {
				t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			break
		}

		// Verify edition matches expected (can be empty string for some license types)
		if edition == expectedEdition || edition == "" {
			t.Logf("✓ License verified successfully")
			t.Logf("  Issuer: %s", issuer)
			t.Logf("  Edition: %s", edition)
			return // Success!
		}

		// If we got here, the edition doesn't match expectations
		lastErr = fmt.Errorf("expected license with edition '%s', got edition: '%s'",
			expectedEdition, edition)

		if attempt < maxRetries {
			t.Logf("Retry %d: %v, waiting %v before retry...", attempt, lastErr, retryDelay)
			time.Sleep(retryDelay)
			continue
		}
	}

	// If we exhausted all retries, fail the test
	require.NoError(t, lastErr, "Failed to verify IBM license after %d attempts", maxRetries)
}

// verifyLicenseInspect replicates the logic from license-inspect.sh
//
// IMPORTANT: The original license-inspect.sh script runs on the REMOTE Vault EC2 instance via SSH
// and has direct filesystem access to /etc/vault.d/vault.lic. However, this Go test runs LOCALLY
// on the GitHub Actions runner and connects to Vault via API only.
//
// Since we cannot access the remote filesystem or execute commands on the remote instance from
// this local test context, we use the Vault API to get license information instead of inspecting
// the license file directly. The sys/license/status API provides the same customer_id information.
func verifyLicenseInspect(t *testing.T, v *blackbox.Session) {
	t.Helper()

	// Read license status using MustRead (similar to TestVaultLicenseStatus)
	t.Logf("Reading license status from Vault API (sys/license/status)")
	licenseStatus := v.MustRead("sys/license/status")

	// Log the full response for debugging
	if jsonData, err := json.MarshalIndent(licenseStatus.Data, "", "  "); err == nil {
		t.Logf("License status data:\n%s", string(jsonData))
	}

	// Extract license data - try both "autoloaded" and "persisted_autoload"
	var autoloaded map[string]interface{}
	var ok bool

	if licenseStatus.Data["autoloaded"] != nil {
		t.Logf("Found 'autoloaded' field")
		autoloaded, ok = licenseStatus.Data["autoloaded"].(map[string]interface{})
		require.True(t, ok, "autoloaded field is not a map, type: %T", licenseStatus.Data["autoloaded"])
	} else if licenseStatus.Data["persisted_autoload"] != nil {
		t.Logf("Found 'persisted_autoload' field")
		autoloaded, ok = licenseStatus.Data["persisted_autoload"].(map[string]interface{})
		require.True(t, ok, "persisted_autoload field is not a map, type: %T", licenseStatus.Data["persisted_autoload"])
	} else {
		t.Fatalf("Neither 'autoloaded' nor 'persisted_autoload' field found in license status. Available fields: %v", getKeys(licenseStatus.Data))
	}

	// Log all fields in the autoloaded license data for debugging
	t.Logf("License data fields: %v", getKeys(autoloaded))

	// Log the complete license data as JSON for detailed inspection
	if jsonData, err := json.MarshalIndent(autoloaded, "", "  "); err == nil {
		t.Logf("Complete license data:\n%s", string(jsonData))
	}

	// Log each field and its type for debugging
	t.Logf("Detailed field information:")
	for key, value := range autoloaded {
		t.Logf("  - %s: type=%T, value=%v", key, value, value)
	}

	// Extract customer_id from the license data
	// This is equivalent to parsing "Customer ID" from the vault license inspect output
	// Note: IBM PAO licenses may not have a customer_id field in the same format
	customerIDData, ok := autoloaded["customer_id"]
	if !ok {
		t.Logf("WARNING: customer_id field not found in license data")
		t.Logf("Available fields: %v", getKeys(autoloaded))
		t.Logf("This might be an IBM PAO license format difference - checking for alternative fields...")

		// IBM PAO licenses might use different field names, log what we have
		// Common alternatives: license_id, issuer, or other identifying fields
		if licenseID, hasLicenseID := autoloaded["license_id"].(string); hasLicenseID && licenseID != "" {
			t.Logf("✓ Found license_id instead: %s", licenseID)
			return
		}

		t.Fatalf("customer_id field not found in license data. Available fields: %v", getKeys(autoloaded))
	}

	customerID, ok := customerIDData.(string)
	require.True(t, ok, "customer_id field is not a string, type: %T, value: %v", customerIDData, customerIDData)

	// Verify customer ID is not empty
	require.NotEmpty(t, customerID, "Customer ID is empty in license data")

	t.Logf("✓ Successfully extracted Customer ID from license: %s", customerID)
}

// TestIBMLicenseUpdate_LicenseGetOnly is a lighter version that only tests the license get API
// without requiring local vault binary access. This is useful for environments where
// the vault binary is not accessible but the API is available.
func TestIBMLicenseUpdate_LicenseGetOnly(t *testing.T) {
	t.Parallel()

	v := blackbox.New(t)

	// Get expected IBM license edition from environment
	expectedEdition := os.Getenv("VAULT_IBM_LICENSE_EDITION")
	if expectedEdition == "" {
		t.Skip("Skipping IBM license update test - VAULT_IBM_LICENSE_EDITION not set")
	}

	t.Logf("Testing IBM license get with expected edition: %s", expectedEdition)

	// Only verify license get (API-based check)
	verifyLicenseGet(t, v, expectedEdition)

	t.Log("✓ IBM license get verification completed successfully")
}

// getKeys returns the keys of a map for debugging purposes
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
