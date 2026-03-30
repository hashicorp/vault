// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testLDAPSecretsCreate tests LDAP secrets engine creation with isolated domain
func testLDAPSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("Failed to create LDAP domain in CI: %v", err)
		}
		t.Skipf("LDAP domain creation not available: %v", err)
	}
	defer cleanup()

	// Create test user in isolated domain
	if err := CreateLDAPUser(t, ldapConfig, "enos", "password123"); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Configure LDAP secrets engine with isolated domain
	SetupLDAPSecretsEngineWithConfig(t, v, "ldap-create", ldapConfig)

	// Create a static role for password rotation using the user in our isolated domain
	v.MustWrite("ldap-create/static-role/test-role", map[string]any{
		"username":        "enos",
		"dn":              "uid=enos," + ldapConfig.UserDN,
		"rotation_period": "24h",
	})

	// Verify role was created by reading it
	roleResp := v.MustRead("ldap-create/static-role/test-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read LDAP static role configuration")
	}

	t.Log("Successfully created LDAP secrets engine with static role in isolated domain")
}

// testLDAPSecretsRead tests LDAP secrets engine read operations with isolated domain
func testLDAPSecretsRead(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("Failed to create LDAP domain in CI: %v", err)
		}
		t.Skipf("LDAP domain creation not available: %v", err)
	}
	defer cleanup()

	// Create service account users in isolated domain
	if err := CreateLDAPUser(t, ldapConfig, "svc-account-1", "password1"); err != nil {
		t.Fatalf("Failed to create service account 1: %v", err)
	}
	if err := CreateLDAPUser(t, ldapConfig, "svc-account-2", "password2"); err != nil {
		t.Fatalf("Failed to create service account 2: %v", err)
	}

	// Configure LDAP secrets engine with isolated domain
	SetupLDAPSecretsEngineWithConfig(t, v, "ldap-read", ldapConfig)

	// Create a library set for service account management
	// Use Eventually wrapper since Vault needs to verify service accounts exist in LDAP
	WriteLibrarySetWithRetry(t, v, "ldap-read/library/test-set", map[string]any{
		"service_account_names":        []string{"svc-account-1", "svc-account-2"},
		"ttl":                          "10h",
		"max_ttl":                      "20h",
		"disable_check_in_enforcement": false,
	})

	// Read the library set configuration
	libraryResp := v.MustRead("ldap-read/library/test-set")
	if libraryResp.Data == nil {
		t.Fatal("Expected to read LDAP library set configuration")
	}

	// Verify library set properties
	assertions := v.AssertSecret(libraryResp)
	assertions.Data().
		HasKeyExists("service_account_names").
		HasKeyExists("ttl").
		HasKeyExists("max_ttl")

	// Read configuration (should not expose bind password)
	configResp := v.MustRead("ldap-read/config")
	if configResp.Data == nil {
		t.Fatal("Expected to read LDAP configuration")
	}

	t.Log("Successfully read LDAP secrets engine configuration in isolated domain")
}

// testLDAPSecretsDelete tests LDAP secrets engine delete operations with isolated domain
func testLDAPSecretsDelete(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("Failed to create LDAP domain in CI: %v", err)
		}
		t.Skipf("LDAP domain creation not available: %v", err)
	}
	defer cleanup()

	// Create service account user in isolated domain
	if err := CreateLDAPUser(t, ldapConfig, "svc-delete", "password"); err != nil {
		t.Fatalf("Failed to create service account: %v", err)
	}

	// Configure LDAP secrets engine with isolated domain
	SetupLDAPSecretsEngineWithConfig(t, v, "ldap-delete", ldapConfig)

	// Create a library set
	// Use Eventually wrapper since Vault needs to verify service accounts exist in LDAP
	WriteLibrarySetWithRetry(t, v, "ldap-delete/library/delete-set", map[string]any{
		"service_account_names": []string{"svc-delete"},
		"ttl":                   "1h",
	})

	// Verify library set exists
	libraryResp := v.MustRead("ldap-delete/library/delete-set")
	if libraryResp.Data == nil {
		t.Fatal("Expected library set to exist before deletion")
	}

	// Delete the library set
	_, err = v.Client.Logical().Delete("ldap-delete/library/delete-set")
	if err != nil {
		t.Fatalf("Failed to delete LDAP library set: %v", err)
	}

	// Verify library set is deleted
	deletedResp, err := v.Client.Logical().Read("ldap-delete/library/delete-set")
	if err == nil && deletedResp != nil {
		t.Fatal("Expected library set to be deleted, but it still exists")
	}

	t.Log("Successfully deleted LDAP library set in isolated domain")
}

// TestLDAPSecretsEngineComprehensive runs comprehensive LDAP secrets engine tests
// Each subtest now uses isolated domains and can run in parallel
func TestLDAPSecretsEngineComprehensive(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testLDAPSecretsCreate(t, v)
	})
	t.Run("Read", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testLDAPSecretsRead(t, v)
	})
	t.Run("Delete", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testLDAPSecretsDelete(t, v)
	})
}
