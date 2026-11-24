//go:build isolated
// +build isolated

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requireLDAPAvailable verifies LDAP server connectivity using testify Eventually
func requireLDAPAvailable(t *testing.T, timeout, interval time.Duration) {
	t.Helper()

	// Use public IP for external connectivity testing
	ldapServerPublic := os.Getenv("LDAP_URL_PUBLIC")
	require.NotEmpty(t, ldapServerPublic, "LDAP_URL_PUBLIC environment variable not set")

	u, err := url.Parse(ldapServerPublic)
	require.NoError(t, err, "Failed to parse LDAP URL: %s", ldapServerPublic)

	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		d := &net.Dialer{}
		conn, err := d.DialContext(t.Context(), "tcp", u.Host)
		require.NoError(ct, err)
		require.NoError(ct, conn.Close())
	}, timeout, interval, "LDAP server not available at %s", u.Host)

	t.Logf("LDAP server connectivity verified at %s", u.Host)
}

// TestLDAP_StaticRoleCreate tests LDAP secrets engine creation with static roles.
func TestLDAP_StaticRoleCreate(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Check if LDAP server configuration is available from integration host
	ldapServer := os.Getenv("LDAP_URL_PRIVATE")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	if ldapServer == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Skip("LDAP server configuration not available - skipping LDAP secrets engine test")
	}

	// Verify LDAP server is ready before proceeding
	requireLDAPAvailable(t, 1*time.Minute, 2*time.Second)

	// Enable LDAP secrets engine
	v.MustEnableSecretsEngine("ldap-create", &api.MountInput{Type: "ldap"})

	// Configure LDAP secrets engine with integration server details
	v.MustWrite("ldap-create/config", map[string]any{
		"binddn":   ldapBindDN,
		"bindpass": ldapBindPass,
		"url":      ldapServer,
		"userdn":   "ou=users,dc=enos,dc=com",
		"userattr": "uid",
	})

	// Create a static role for password rotation using the user created by the integration setup
	v.MustWrite("ldap-create/static-role/test-role", map[string]any{
		"username":        "enos",
		"dn":              "uid=enos,ou=users,dc=enos,dc=com",
		"rotation_period": "24h",
	})

	// Verify role was created by reading it
	roleResp := v.MustRead("ldap-create/static-role/test-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read LDAP static role configuration")
	}

	t.Log("Successfully created LDAP secrets engine with static role")
}

// TestLDAP_LibrarySetRead tests LDAP secrets engine read operations with library sets.
func TestLDAP_LibrarySetRead(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Check if LDAP server configuration is available from integration host
	ldapServer := os.Getenv("LDAP_URL_PRIVATE")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	if ldapServer == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Skip("LDAP server configuration not available - skipping LDAP secrets engine test")
	}

	// Verify LDAP server is ready before proceeding
	requireLDAPAvailable(t, 1*time.Minute, 2*time.Second)
	serviceAccounts := []string{"svc-account-1", "svc-account-2"}

	// Enable LDAP secrets engine
	v.MustEnableSecretsEngine("ldap-read", &api.MountInput{Type: "ldap"})

	// Configure LDAP secrets engine with integration server details
	v.MustWrite("ldap-read/config", map[string]any{
		"binddn":   ldapBindDN,
		"bindpass": ldapBindPass,
		"url":      ldapServer,
		"userdn":   "ou=users,dc=enos,dc=com",
		"userattr": "uid",
	})

	// Create a library set for service account management
	v.MustWrite("ldap-read/library/test-set", map[string]any{
		"service_account_names":        serviceAccounts,
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

	t.Log("Successfully read LDAP secrets engine configuration")
}

// TestLDAP_LibrarySetDelete tests LDAP secrets engine delete operations with library sets.
func TestLDAP_LibrarySetDelete(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Check if LDAP server configuration is available from integration host
	ldapServer := os.Getenv("LDAP_URL_PRIVATE")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	if ldapServer == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Skip("LDAP server configuration not available - skipping LDAP secrets engine test")
	}

	// Verify LDAP server is ready before proceeding
	requireLDAPAvailable(t, 1*time.Minute, 2*time.Second)
	serviceAccounts := []string{"svc-delete"}

	// Enable LDAP secrets engine
	v.MustEnableSecretsEngine("ldap-delete", &api.MountInput{Type: "ldap"})

	// Configure LDAP secrets engine with integration server details
	v.MustWrite("ldap-delete/config", map[string]any{
		"binddn":   ldapBindDN,
		"bindpass": ldapBindPass,
		"url":      ldapServer,
		"userdn":   "ou=users,dc=enos,dc=com",
		"userattr": "uid",
	})

	// Create a library set
	v.MustWrite("ldap-delete/library/delete-set", map[string]any{
		"service_account_names": serviceAccounts,
		"ttl":                   "1h",
	})

	// Verify library set exists
	libraryResp := v.MustRead("ldap-delete/library/delete-set")
	if libraryResp.Data == nil {
		t.Fatal("Expected library set to exist before deletion")
	}

	// Delete the library set
	_, err := v.Client.Logical().Delete("ldap-delete/library/delete-set")
	if err != nil {
		t.Fatalf("Failed to delete LDAP library set: %v", err)
	}

	// Verify library set is deleted
	deletedResp, err := v.Client.Logical().Read("ldap-delete/library/delete-set")
	if err == nil && deletedResp != nil {
		t.Fatal("Expected library set to be deleted, but it still exists")
	}

	t.Log("Successfully deleted LDAP library set")
}

// testLDAPSecretsCreate tests LDAP secrets engine creation with isolated domain
func testLDAPSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("LDAP domain creation failed in CI: %v", err)
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
			t.Fatalf("LDAP domain creation failed in CI: %v", err)
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
			t.Fatalf("LDAP domain creation failed in CI: %v", err)
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
