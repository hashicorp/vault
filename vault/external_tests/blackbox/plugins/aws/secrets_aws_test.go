// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestAWS_GenerateNewUser tests AWS secrets engine credential generation.
func TestAWS_GenerateNewUser(t *testing.T) {
	t.Parallel()
	skipIfNoAWSCredentials(t)
	v := blackbox.New(t)

	// Create test IAM user for Vault configuration
	userName, tempAccessKeyId, tempSecretAccessKey, demoUserPolicyArn, _, _ := createTestIAMUser(t)
	t.Logf("Created test IAM user: %s", userName)

	// Track generated credentials for cleanup
	var newAccessKey string
	t.Cleanup(func() {
		if newAccessKey != "" {
			t.Logf("Cleanup: deleting IAM user created by Vault with access key: %s", newAccessKey)
			deleteIAMUserByAccessKey(t, newAccessKey)
		}
		t.Logf("Cleanup: deleting IAM user by initial access key: %s", tempAccessKeyId)
		deleteIAMUserByAccessKey(t, tempAccessKeyId)
	})

	// Enable and configure AWS secrets engine
	path := setupAWSSecretsEngine(t, v, tempAccessKeyId, tempSecretAccessKey, getAwsUsernameTemplate(userName))

	// Create Vault role for credential generation
	roleName := "aws-enos-role"
	createVaultAWSRole(t, v, path, roleName, demoUserPolicyArn)
	verifyRoleExists(t, v, path, roleName)

	// Verify username template was configured
	t.Logf("Reading root config to verify username template is set correctly")
	rootUser := v.MustRead(fmt.Sprintf("%s/config/root", path))
	if rootUser == nil || rootUser.Data == nil {
		t.Fatalf("Expected to read root config, got nil: %#v", rootUser)
	}
	if val, ok := rootUser.Data["username_template"]; !ok || val == nil {
		t.Fatalf("username_template missing in root config: %#v", rootUser)
	}

	// Generate new IAM user credentials via Vault
	t.Logf("Generating new credentials for IAM user using role: %s", roleName)
	newUser := v.MustRead(fmt.Sprintf("%s/creds/%s", path, roleName))
	if newUser == nil || newUser.Data == nil {
		t.Fatalf("Failed to generate new credentials for IAM user: %s", roleName)
	}
	if val, ok := newUser.Data["access_key"]; !ok || val == nil || val == tempAccessKeyId {
		t.Fatalf("The new access key is empty or is matching the old one: %v", val)
	}

	// Extract and save access key for cleanup
	var ok bool
	newAccessKey, ok = newUser.Data["access_key"].(string)
	if !ok || newAccessKey == "" {
		t.Fatalf("Could not extract access_key from new credentials: %v", newUser.Data["access_key"])
	}
	t.Logf("Captured Vault-created access key for cleanup: %s", newAccessKey)
}

// TestAWS_CreateDeleteVaultAwsRole tests Vault AWS role lifecycle.
func TestAWS_CreateDeleteVaultAwsRole(t *testing.T) {
	t.Parallel()
	skipIfNoAWSCredentials(t)
	v := blackbox.New(t)

	// Create test IAM user for Vault configuration
	userName, tempAccessKeyId, tempSecretAccessKey, demoUserPolicyArn, _, _ := createTestIAMUser(t)
	t.Logf("Created test IAM user: %s", userName)
	t.Cleanup(func() {
		t.Logf("Cleanup: deleting IAM user by initial access key: %s", tempAccessKeyId)
		deleteIAMUserByAccessKey(t, tempAccessKeyId)
	})

	// Enable and configure AWS secrets engine
	path := setupAWSSecretsEngine(t, v, tempAccessKeyId, tempSecretAccessKey, "")

	// Create and verify Vault role exists
	roleName := "aws-enos-role"
	createVaultAWSRole(t, v, path, roleName, demoUserPolicyArn)
	verifyRoleExists(t, v, path, roleName)

	// Delete role and verify it's gone
	t.Logf("Deleting Vault AWS role: %s", roleName)
	v.MustDelete(fmt.Sprintf("%s/roles/%s", path, roleName))
	t.Logf("Role deleted at path: %s", roleName)

	verifyRoleDeleted(t, v, path, roleName)
}
