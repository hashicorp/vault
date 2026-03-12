// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testLDAPSecretsCreate tests LDAP secrets engine creation
func testLDAPSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Check if LDAP server configuration is available from integration host
	ldapServer := os.Getenv("LDAP_SERVER")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	if ldapServer == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Skip("LDAP server configuration not available - skipping LDAP secrets engine test")
	}

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

// testLDAPSecretsRead tests LDAP secrets engine read operations
func testLDAPSecretsRead(t *testing.T, v *blackbox.Session) {
	// Check if LDAP server configuration is available from integration host
	ldapServer := os.Getenv("LDAP_SERVER")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	if ldapServer == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Skip("LDAP server configuration not available - skipping LDAP secrets engine test")
	}

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

	t.Log("Successfully read LDAP secrets engine configuration")
}

// testLDAPSecretsDelete tests LDAP secrets engine delete operations
func testLDAPSecretsDelete(t *testing.T, v *blackbox.Session) {
	// Check if LDAP server configuration is available from integration host
	ldapServer := os.Getenv("LDAP_SERVER")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	if ldapServer == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Skip("LDAP server configuration not available - skipping LDAP secrets engine test")
	}

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
		"service_account_names": []string{"svc-delete"},
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
