// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testKMIPSecretsCreate tests KMIP secrets engine creation
func testKMIPSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Check if this is Vault Enterprise (KMIP is enterprise-only)
	edition := os.Getenv("VAULT_EDITION")
	if edition == "ce" || edition == "" {
		t.Skip("KMIP secrets engine is only available in Vault Enterprise")
	}

	// Enable KMIP secrets engine
	v.MustEnableSecretsEngine("kmip-create", &api.MountInput{Type: "kmip"})

	// Configure KMIP secrets engine
	v.MustWrite("kmip-create/config", map[string]any{
		"listen_addrs":           []string{"0.0.0.0:5696"},
		"server_hostnames":       []string{"localhost"},
		"default_tls_client_ttl": "24h",
	})

	// Create a KMIP scope
	v.MustWrite("kmip-create/scope/test-scope", map[string]any{})

	// Create a KMIP role within the scope
	v.MustWrite("kmip-create/scope/test-scope/role/test-role", map[string]any{
		"operation_all": true,
	})

	// Verify role was created by reading it
	roleResp := v.MustRead("kmip-create/scope/test-scope/role/test-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read KMIP role configuration")
	}

	t.Log("Successfully created KMIP secrets engine with scope and role")
}

// testKMIPSecretsRead tests KMIP secrets engine read operations
func testKMIPSecretsRead(t *testing.T, v *blackbox.Session) {
	// Check if this is Vault Enterprise (KMIP is enterprise-only)
	edition := os.Getenv("VAULT_EDITION")
	if edition == "ce" || edition == "" {
		t.Skip("KMIP secrets engine is only available in Vault Enterprise")
	}

	// Enable KMIP secrets engine
	v.MustEnableSecretsEngine("kmip-read", &api.MountInput{Type: "kmip"})

	// Configure KMIP secrets engine
	v.MustWrite("kmip-read/config", map[string]any{
		"listen_addrs":     []string{"0.0.0.0:5697"},
		"server_hostnames": []string{"localhost"},
	})

	// Create a scope
	v.MustWrite("kmip-read/scope/read-scope", map[string]any{})

	// Create a role
	v.MustWrite("kmip-read/scope/read-scope/role/read-role", map[string]any{
		"operation_activate": true,
		"operation_create":   true,
		"operation_get":      true,
	})

	// Read the role configuration
	roleResp := v.MustRead("kmip-read/scope/read-scope/role/read-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read KMIP role configuration")
	}

	// Verify role properties
	assertions := v.AssertSecret(roleResp)
	assertions.Data().
		HasKey("operation_activate", true).
		HasKey("operation_create", true).
		HasKey("operation_get", true)

	// Note: Reading individual scopes is not supported (returns 405)
	// The KMIP API only supports listing scopes, not reading them individually

	// Read KMIP configuration
	configResp := v.MustRead("kmip-read/config")
	if configResp.Data == nil {
		t.Fatal("Expected to read KMIP configuration")
	}

	t.Log("Successfully read KMIP secrets engine configuration")
}
