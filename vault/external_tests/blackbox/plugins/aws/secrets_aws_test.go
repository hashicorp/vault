// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestAWS_GenerateNewUser verifies AWS secrets engine can generate IAM user credentials.
func TestAWS_GenerateNewUser(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		t.Log("AWS credentials not available - skipping AWS secrets engine test")
		t.Skip("AWS credentials not available - skipping AWS secrets engine test")
	}

	t.Logf("Creating test IAM user via helpers.go...")
	userName, tempAccessKeyId, tempSecretAccessKey, demoUserPolicyArn, _, _ := createTestIAMUser(t)
	t.Logf("Created test IAM user: %s", userName)
	var newAccessKey string
	t.Cleanup(func() {
		if newAccessKey != "" {
			t.Logf("Cleanup: deleting IAM user created by Vault with access key: %s", newAccessKey)
			deleteIAMUserByAccessKey(t, newAccessKey)
		}

		t.Logf("Cleanup: deleting IAM user by initial access key: %s", tempAccessKeyId)
		deleteIAMUserByAccessKey(t, tempAccessKeyId)
	})

	path := fmt.Sprintf("aws-test-%d", time.Now().UnixNano())
	t.Logf("Enabling AWS secrets engine at path: %s", path)
	v.MustEnableSecretsEngine(path, &api.MountInput{Type: "aws"})

	t.Logf("Configuring AWS secrets engine with root credentials and username template for user: %s", userName)
	v.MustWrite(fmt.Sprintf("%s/config/root", path), map[string]any{
		"access_key":        tempAccessKeyId,
		"secret_key":        tempSecretAccessKey,
		"region":            "us-east-1",
		"username_template": getAwsUsernameTemplate(userName),
	})

	roleName := "aws-enos-role"
	t.Logf("Creating Vault AWS role: %s", roleName)
	v.MustWrite(fmt.Sprintf("%s/roles/%s", path, roleName), map[string]any{
		"credential_type":          "iam_user",
		"permissions_boundary_arn": demoUserPolicyArn,
		"policy_document":          getAllowDescribeRegionsPolicy(),
	})

	t.Logf("Reading and verifying AWS role configuration for role: %s", roleName)
	roleResp := v.MustRead(fmt.Sprintf("%s/roles/%s", path, roleName))
	if roleResp.Data == nil {
		t.Fatal("Expected to read AWS role configuration")
	}

	t.Logf("Listing AWS roles at path: %s/roles", path)
	rolesList := v.MustList(fmt.Sprintf("%s/roles", path))
	if rolesList == nil || rolesList.Data == nil {
		t.Fatal("No AWS roles created! (rolesList is nil or Data is nil)")
	}
	roleKeys, ok := rolesList.Data["keys"].([]interface{})
	if !ok || len(roleKeys) == 0 {
		t.Fatal("No AWS roles created! (rolesList.Data['keys'] is empty or not a slice)")
	}
	t.Logf("Found AWS roles: %v", roleKeys)

	t.Logf("Reading root config to verify username template is set correctly")
	rootUser := v.MustRead(fmt.Sprintf("%s/config/root", path))
	if rootUser == nil || rootUser.Data == nil {
		t.Fatalf("Expected to read root config, got nil: %#v", rootUser)
	}
	if val, ok := rootUser.Data["username_template"]; !ok || val == nil {
		t.Fatalf("username_template missing in root config: %#v", rootUser)
	}

	t.Logf("Generating new credentials for IAM user using role: %s", roleName)
	newUser := v.MustRead(fmt.Sprintf("%s/creds/%s", path, roleName))
	if newUser == nil || newUser.Data == nil {
		t.Fatalf("Failed to generate new credentials for IAM user: %s", roleName)
	}
	if val, ok := newUser.Data["access_key"]; !ok || val == nil || val == tempAccessKeyId {
		t.Fatalf("The new access key is empty or is matching the old one: %v", val)
	}

	newAccessKey, ok = newUser.Data["access_key"].(string)
	if !ok || newAccessKey == "" {
		t.Fatalf("Could not extract access_key from new credentials: %v", newUser.Data["access_key"])
	}
	t.Logf("Captured Vault-created access key for cleanup: %s", newAccessKey)
}
