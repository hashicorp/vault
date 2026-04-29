// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testPKISecretsCreate tests PKI secrets engine creation
func testPKISecretsCreate(t *testing.T, v *blackbox.Session) {
	// Enable PKI secrets engine
	v.MustEnableSecretsEngine("pki-create", &api.MountInput{Type: "pki"})

	// Configure max TTL for the mount
	err := v.Client.Sys().TuneMount("pki-create", api.MountConfigInput{
		MaxLeaseTTL: "87600h",
	})
	if err != nil {
		t.Fatalf("Failed to tune PKI mount: %v", err)
	}

	// Generate root CA
	rootResp := v.MustWrite("pki-create/root/generate/internal", map[string]any{
		"common_name": "test-root-ca.example.com",
		"ttl":         "8760h",
		"key_type":    "rsa",
		"key_bits":    2048,
	})

	if rootResp.Data == nil {
		t.Fatal("Expected root CA generation response")
	}

	// Verify root CA was created
	assertions := v.AssertSecret(rootResp)
	assertions.Data().
		HasKeyExists("certificate").
		HasKeyExists("issuing_ca").
		HasKeyExists("serial_number")

	// Create a role for issuing certificates
	v.MustWrite("pki-create/roles/test-role", map[string]any{
		"allowed_domains":  []string{"example.com"},
		"allow_subdomains": true,
		"max_ttl":          "72h",
		"key_type":         "rsa",
		"key_bits":         2048,
	})

	// Verify role was created by reading it
	roleResp := v.MustRead("pki-create/roles/test-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read role configuration")
	}

	t.Log("Successfully created PKI secrets engine with root CA and role")
}

// testPKISecretsRead tests PKI secrets engine read operations
func testPKISecretsRead(t *testing.T, v *blackbox.Session) {
	// Setup PKI engine with root CA
	roleName := v.MustSetupPKIRoot("pki-read")

	// Read the role configuration
	roleResp := v.MustRead("pki-read/roles/" + roleName)
	if roleResp.Data == nil {
		t.Fatal("Expected to read role configuration")
	}

	// Verify role properties
	assertions := v.AssertSecret(roleResp)
	assertions.Data().
		HasKey("allow_subdomains", true).
		HasKeyExists("max_ttl").
		HasKeyExists("allowed_domains")

	// Read CA certificate
	caResp := v.MustRead("pki-read/cert/ca")
	if caResp.Data == nil {
		t.Fatal("Expected to read CA certificate")
	}

	assertions = v.AssertSecret(caResp)
	assertions.Data().HasKeyExists("certificate")

	t.Log("Successfully read PKI secrets engine configuration and certificates")
}

// testPKISecretsDelete tests PKI secrets engine delete operations
func testPKISecretsDelete(t *testing.T, v *blackbox.Session) {
	roleName := v.MustSetupPKIRoot("pki-delete")
	roleResp := v.MustRead("pki-delete/roles/" + roleName)
	if roleResp.Data == nil {
		t.Fatal("Expected role to exist before deletion")
	}
	_, err := v.Client.Logical().Delete("pki-delete/roles/" + roleName)
	if err != nil {
		t.Fatalf("Failed to delete PKI role: %v", err)
	}
	deletedResp, err := v.Client.Logical().Read("pki-delete/roles/" + roleName)
	if err == nil && deletedResp != nil {
		t.Fatal("Expected role to be deleted, but it still exists")
	}
	t.Logf("Successfully deleted PKI role: %s", roleName)
}
