// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testSSHSecretsCreate tests SSH secrets engine creation
func testSSHSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Enable SSH secrets engine
	v.MustEnableSecretsEngine("ssh-create", &api.MountInput{Type: "ssh"})

	// Generate CA key pair
	caResp := v.MustWrite("ssh-create/config/ca", map[string]any{
		"generate_signing_key": true,
	})

	if caResp.Data == nil {
		t.Fatal("Expected CA generation response")
	}

	// Verify CA was created
	assertions := v.AssertSecret(caResp)
	assertions.Data().HasKeyExists("public_key")

	// Create an SSH role for signing certificates
	v.MustWrite("ssh-create/roles/test-role", map[string]any{
		"key_type":                "ca",
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"default_user":            "ubuntu",
		"ttl":                     "30m",
		"max_ttl":                 "24h",
	})

	// Verify role was created by reading it
	roleResp := v.MustRead("ssh-create/roles/test-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read SSH role configuration")
	}

	// Verify role properties
	roleAssertions := v.AssertSecret(roleResp)
	roleAssertions.Data().
		HasKey("key_type", "ca").
		HasKey("allow_user_certificates", true).
		HasKey("default_user", "ubuntu")

	t.Log("Successfully created SSH secrets engine with CA and role")
}

// testSSHSecretsRead tests SSH secrets engine read operations
func testSSHSecretsRead(t *testing.T, v *blackbox.Session) {
	// Enable SSH secrets engine
	v.MustEnableSecretsEngine("ssh-read", &api.MountInput{Type: "ssh"})

	// Generate CA
	v.MustWrite("ssh-read/config/ca", map[string]any{
		"generate_signing_key": true,
	})

	// Create a role
	v.MustWrite("ssh-read/roles/read-role", map[string]any{
		"key_type":                "ca",
		"allow_user_certificates": true,
		"allowed_users":           "testuser",
		"default_user":            "testuser",
		"ttl":                     "1h",
	})

	// Read the role configuration
	roleResp := v.MustRead("ssh-read/roles/read-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read SSH role configuration")
	}

	// Verify role properties
	assertions := v.AssertSecret(roleResp)
	assertions.Data().
		HasKey("key_type", "ca").
		HasKey("allow_user_certificates", true).
		HasKey("allowed_users", "testuser").
		HasKey("default_user", "testuser")

	// Read CA public key
	publicKeyResp := v.MustRead("ssh-read/config/ca")
	if publicKeyResp.Data == nil {
		t.Fatal("Expected to read CA public key")
	}

	assertions = v.AssertSecret(publicKeyResp)
	assertions.Data().HasKeyExists("public_key")

	t.Log("Successfully read SSH secrets engine configuration")
}

// testSSHSecretsDelete tests SSH secrets engine delete operations
func testSSHSecretsDelete(t *testing.T, v *blackbox.Session) {
	// Enable SSH secrets engine
	v.MustEnableSecretsEngine("ssh-delete", &api.MountInput{Type: "ssh"})

	// Generate CA
	v.MustWrite("ssh-delete/config/ca", map[string]any{
		"generate_signing_key": true,
	})

	// Create a role
	v.MustWrite("ssh-delete/roles/delete-role", map[string]any{
		"key_type":                "ca",
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"default_user":            "ubuntu",
	})

	// Verify role exists
	roleResp := v.MustRead("ssh-delete/roles/delete-role")
	if roleResp.Data == nil {
		t.Fatal("Expected role to exist before deletion")
	}

	// Delete the role
	_, err := v.Client.Logical().Delete("ssh-delete/roles/delete-role")
	if err != nil {
		t.Fatalf("Failed to delete SSH role: %v", err)
	}

	// Verify role is deleted
	deletedResp, err := v.Client.Logical().Read("ssh-delete/roles/delete-role")
	if err == nil && deletedResp != nil {
		t.Fatal("Expected role to be deleted, but it still exists")
	}

	t.Log("Successfully deleted SSH role")
}
