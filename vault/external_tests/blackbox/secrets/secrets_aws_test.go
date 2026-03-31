// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testAWSSecretsCreate tests AWS secrets engine creation
func testAWSSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Check if AWS credentials are available
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey == "" || secretKey == "" {
		t.Skip("AWS credentials not available - skipping AWS secrets engine test")
	}

	// Enable AWS secrets engine
	v.MustEnableSecretsEngine("aws-create", &api.MountInput{Type: "aws"})

	// Configure AWS secrets engine with root credentials
	v.MustWrite("aws-create/config/root", map[string]any{
		"access_key": accessKey,
		"secret_key": secretKey,
		"region":     "us-east-1",
	})

	// Create a role for generating credentials
	v.MustWrite("aws-create/roles/test-role", map[string]any{
		"credential_type": "iam_user",
		"policy_document": `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": "ec2:Describe*",
					"Resource": "*"
				}
			]
		}`,
	})

	// Verify role was created by reading it
	roleResp := v.MustRead("aws-create/roles/test-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read AWS role configuration")
	}

	t.Log("Successfully created AWS secrets engine with role")
}

// testAWSSecretsRead tests AWS secrets engine read operations
func testAWSSecretsRead(t *testing.T, v *blackbox.Session) {
	// Check if AWS credentials are available
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey == "" || secretKey == "" {
		t.Skip("AWS credentials not available - skipping AWS secrets engine test")
	}

	// Enable AWS secrets engine
	v.MustEnableSecretsEngine("aws-read", &api.MountInput{Type: "aws"})

	// Configure AWS secrets engine
	v.MustWrite("aws-read/config/root", map[string]any{
		"access_key": accessKey,
		"secret_key": secretKey,
		"region":     "us-west-2",
	})

	// Create a role
	v.MustWrite("aws-read/roles/read-role", map[string]any{
		"credential_type": "iam_user",
		"policy_document": `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": "s3:ListBucket",
					"Resource": "*"
				}
			]
		}`,
	})

	// Read the role configuration
	roleResp := v.MustRead("aws-read/roles/read-role")
	if roleResp.Data == nil {
		t.Fatal("Expected to read AWS role configuration")
	}

	// Verify role properties
	assertions := v.AssertSecret(roleResp)
	assertions.Data().
		HasKey("credential_type", "iam_user").
		HasKeyExists("policy_document")

	// Read root configuration (should not expose credentials)
	configResp := v.MustRead("aws-read/config/root")
	if configResp.Data == nil {
		t.Fatal("Expected to read AWS root configuration")
	}

	// Verify config properties (credentials should not be returned)
	configAssertions := v.AssertSecret(configResp)
	configAssertions.Data().HasKey("region", "us-west-2")

	t.Log("Successfully read AWS secrets engine configuration")
}
