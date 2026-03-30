// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testLDAPRootCredentialRollbackSuccess tests successful rollback when rotation fails
// Converts: secrets-rollback-invalid-config.sh
// Scenario: Rotation fails with invalid LDAP endpoint, old password is preserved
func testLDAPRootCredentialRollbackSuccess(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("Failed to create LDAP domain in CI: %v", err)
		}
		t.Skipf("LDAP domain creation not available: %v", err)
	}
	defer cleanup()

	// Create admin user in isolated domain
	adminUser := "admin-rollback-success"
	adminPassword := "initial-password-123"
	if err := CreateLDAPUser(t, ldapConfig, adminUser, adminPassword); err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	// Configure LDAP secrets engine with valid config
	mount := "ldap-rollback-success"
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "ldap"})

	adminDN := fmt.Sprintf("uid=%s,%s", adminUser, ldapConfig.UserDN)
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})

	// Verify baseline LDAP authentication works
	if !verifyLDAPAuth(t, ldapConfig, adminDN, adminPassword) {
		t.Fatal("Baseline LDAP authentication failed")
	}
	t.Log("✓ Baseline LDAP authentication successful")

	// Poison Vault config with unreachable LDAP endpoint (port 9999)
	badURL := strings.Replace(ldapConfig.URL, ldapConfig.URL[strings.LastIndex(ldapConfig.URL, ":")+1:], "9999", 1)
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      badURL,
		"userdn":   ldapConfig.UserDN,
	})
	t.Logf("Poisoned config with unreachable endpoint: %s", badURL)

	// Attempt rotation (should fail)
	_, err = v.Client.Logical().Write(mount+"/rotate-root", nil)
	if err == nil {
		t.Fatal("Expected rotation to fail with invalid endpoint, but it succeeded")
	}
	t.Logf("✓ Rotation failed as expected: %v", err)

	// Restore valid config
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})

	// Verify old password still works (rollback success) with retry
	v.Eventually(func() error {
		if !verifyLDAPAuth(t, ldapConfig, adminDN, adminPassword) {
			return fmt.Errorf("LDAP authentication failed")
		}
		return nil
	})
	t.Log("✓ ROLLBACK SUCCESS: Old password preserved after failed rotation")

	// Verify Vault can read LDAP config
	configResp := v.MustRead(mount + "/config")
	if configResp.Data == nil {
		t.Fatal("Failed to read LDAP config after recovery")
	}
	t.Log("✓ Vault reconnected successfully")
}

// testLDAPRootCredentialRollbackFailure tests critical failure scenarios
// Converts: secrets-rollback-creds-mismatch.sh
// Scenario: Rotation with wrong binddn AND wrong password (credential mismatch)
func testLDAPRootCredentialRollbackFailure(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("Failed to create LDAP domain in CI: %v", err)
		}
		t.Skipf("LDAP domain creation not available: %v", err)
	}
	defer cleanup()

	// Create admin user in isolated domain
	adminUser := "admin-rollback-failure"
	adminPassword := "initial-password-456"
	if err := CreateLDAPUser(t, ldapConfig, adminUser, adminPassword); err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	// Configure LDAP secrets engine with valid config
	mount := "ldap-rollback-failure"
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "ldap"})

	adminDN := fmt.Sprintf("uid=%s,%s", adminUser, ldapConfig.UserDN)
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})

	// Verify baseline LDAP authentication works
	if !verifyLDAPAuth(t, ldapConfig, adminDN, adminPassword) {
		t.Fatal("Baseline LDAP authentication failed")
	}
	t.Log("✓ Baseline LDAP authentication successful")

	// Poison Vault config with wrong binddn AND wrong password (credential mismatch)
	badBindDN := fmt.Sprintf("cn=nonexistent-admin,%s", ldapConfig.BaseDN)
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   badBindDN,
		"bindpass": "intentionally-wrong-password",
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})
	t.Logf("Poisoned config with bad credentials: binddn=%s", badBindDN)

	// Attempt rotation (should fail with credential mismatch)
	_, rotationErr := v.Client.Logical().Write(mount+"/rotate-root", nil)
	if rotationErr == nil {
		t.Fatal("Expected rotation to fail with credential mismatch, but it succeeded")
	}
	t.Logf("✓ Rotation failed as expected with credential mismatch: %v", rotationErr)

	// Restore correct config for recovery
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})

	// Post-recovery validation
	configResp := v.MustRead(mount + "/config")
	if configResp.Data == nil {
		t.Fatal("CRITICAL: Vault failed to reconnect after recovery")
	}
	t.Log("✓ Vault reconnected after recovery")

	// Verify LDAP authentication after recovery
	if !verifyLDAPAuth(t, ldapConfig, adminDN, adminPassword) {
		t.Fatal("CRITICAL: LDAP authentication failed after recovery")
	}
	t.Log("✓ LDAP authentication successful after recovery")
}

// testLDAPRootCredentialAutomaticRollbackOnFailure tests automatic rollback mechanism
// Converts: secrets-rollback-transactional.sh
// Scenario: System consistency after mid-rotation failure with automatic rollback
func testLDAPRootCredentialAutomaticRollbackOnFailure(t *testing.T, v *blackbox.Session) {
	// Create isolated LDAP domain for this test
	cleanup, ldapConfig, err := PrepareTestLDAPDomain(t, v, isCI())
	if err != nil {
		if isCI() {
			t.Fatalf("Failed to create LDAP domain in CI: %v", err)
		}
		t.Skipf("LDAP domain creation not available: %v", err)
	}
	defer cleanup()

	// Create admin user in isolated domain
	adminUser := "admin-auto-rollback"
	adminPassword := "initial-password-789"
	if err := CreateLDAPUser(t, ldapConfig, adminUser, adminPassword); err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	// Configure LDAP secrets engine with valid config
	mount := "ldap-auto-rollback"
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "ldap"})

	adminDN := fmt.Sprintf("uid=%s,%s", adminUser, ldapConfig.UserDN)
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})

	// Verify baseline LDAP health
	if !verifyLDAPAuth(t, ldapConfig, adminDN, adminPassword) {
		t.Fatal("Baseline LDAP authentication failed")
	}
	t.Log("✓ Baseline LDAP authentication successful")

	// Verify LDAP config is readable before rotation
	configResp := v.MustRead(mount + "/config")
	if configResp.Data == nil {
		t.Fatal("LDAP config not readable before rotation")
	}
	t.Log("✓ LDAP config readable before rotation")

	// Pre-poison config BEFORE rotation attempt (deterministic test)
	badBindDN := fmt.Sprintf("cn=nonexistent,%s", ldapConfig.BaseDN)
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   badBindDN,
		"bindpass": "wrong-password",
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})
	t.Log("Pre-poisoned config with invalid credentials before rotation")

	// Attempt rotation with poisoned config (should fail immediately)
	_, rotationErr := v.Client.Logical().Write(mount+"/rotate-root", nil)
	if rotationErr == nil {
		t.Fatal("Expected rotation to fail with poisoned config, but it succeeded")
	}
	t.Logf("✓ Rotation failed as expected with poisoned config: %v", rotationErr)

	// Restore valid config for recovery
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   adminDN,
		"bindpass": adminPassword,
		"url":      ldapConfig.URL,
		"userdn":   ldapConfig.UserDN,
	})
	t.Log("Restored valid config after failed rotation")

	// Verify old password still works (rollback success) with retry
	v.Eventually(func() error {
		if !verifyLDAPAuth(t, ldapConfig, adminDN, adminPassword) {
			return fmt.Errorf("LDAP authentication failed")
		}
		return nil
	})
	t.Log("✓ AUTOMATIC ROLLBACK SUCCESS: Old password still works after failed rotation")

	// Verify Vault is operational after recovery
	v.Eventually(func() error {
		resp, err := v.Client.Logical().Read(mount + "/config")
		if err != nil {
			return fmt.Errorf("cannot read LDAP config: %w", err)
		}
		if resp == nil || resp.Data == nil {
			return fmt.Errorf("LDAP config data is nil")
		}
		return nil
	})
	t.Log("✓ Vault operational: Can read LDAP config after recovery")
}

// verifyLDAPAuth verifies LDAP authentication using ldapwhoami
// Returns true if authentication succeeds, false otherwise
// Note: Callers should use v.Eventually() for retry logic
func verifyLDAPAuth(t *testing.T, config *LDAPDomainConfig, bindDN, password string) bool {
	t.Helper()

	// Use SetupURL (public IP) for authentication checks from GitHub runner
	// config.URL contains private IP which is only accessible from Vault cluster

	cmd := exec.Command("ldapwhoami",
		"-x",
		"-H", config.SetupURL,
		"-D", bindDN,
		"-w", password,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("LDAP authentication failed for %s: %v, output: %s", bindDN, err, string(output))
		return false
	}

	return true
}

// TestLDAPRootCredentialRollbackWorkflows runs all rollback workflow tests
// This is the main test function that gets triggered by enos-scenario-plugin.hcl
func TestLDAPRootCredentialRollbackWorkflows(t *testing.T) {
	t.Run("RollbackSuccess", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testLDAPRootCredentialRollbackSuccess(t, v)
	})

	t.Run("RollbackFailure", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testLDAPRootCredentialRollbackFailure(t, v)
	})

	t.Run("AutomaticRollback", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testLDAPRootCredentialAutomaticRollbackOnFailure(t, v)
	})
}
