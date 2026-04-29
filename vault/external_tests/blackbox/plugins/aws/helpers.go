// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// =============================================================================
// AWS Helper Functions
// =============================================================================

// getPolicyArnByName finds and returns the ARN for an IAM policy by name.
func getPolicyArnByName(ctx context.Context, iamClient *iam.Client, policyName string) (string, error) {
	paginator := iam.NewListPoliciesPaginator(iamClient, &iam.ListPoliciesInput{Scope: "All"})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", err
		}
		for _, p := range page.Policies {
			if aws.ToString(p.PolicyName) == policyName {
				return aws.ToString(p.Arn), nil
			}
		}
	}
	return "", fmt.Errorf("policy %s not found", policyName)
}

// getRoleArnByName finds and returns the ARN for an IAM role by name.
func getRoleArnByName(ctx context.Context, iamClient *iam.Client, roleName string) (string, error) {
	paginator := iam.NewListRolesPaginator(iamClient, &iam.ListRolesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", err
		}
		for _, r := range page.Roles {
			if aws.ToString(r.RoleName) == roleName {
				return aws.ToString(r.Arn), nil
			}
		}
	}
	return "", fmt.Errorf("role %s not found", roleName)
}

// createTestIAMUser creates a test IAM user with DemoUser policy and returns credentials.
func createTestIAMUser(t *testing.T) (
	userName string,
	accessKeyID string,
	secretAccessKey string,
	demoUserPolicyArn string,
	assumedRoleArn string,
	awsRegion string,
) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}
	awsRegion = cfg.Region
	if awsRegion == "" {
		t.Fatalf("AWS region is empty in config")
	}
	iamClient := iam.NewFromConfig(cfg)
	stsClient := sts.NewFromConfig(cfg)

	// Get current AWS account identity (for unique name)
	caller, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		t.Fatalf("failed to get caller identity: %v", err)
	}
	accountID := aws.ToString(caller.Account)

	// Generate a random hex suffix for uniqueness
	const randomSuffixByteLength = 4
	suffix := make([]byte, randomSuffixByteLength)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatalf("failed to generate random suffix: %v", err)
	}
	hexSuffix := hex.EncodeToString(suffix)
	userName = fmt.Sprintf("demo-GitHubActions-%s-%s", accountID, hexSuffix)

	// Lookup DemoUser policy ARN
	demoUserPolicyArn, err = getPolicyArnByName(ctx, iamClient, "DemoUser")
	if err != nil {
		t.Fatalf("DemoUser policy not found: %v", err)
	}

	// Lookup vault-assumed-role-credentials-demo role ARN
	assumedRoleArn, err = getRoleArnByName(ctx, iamClient, "vault-assumed-role-credentials-demo")
	if err != nil {
		t.Fatalf("vault-assumed-role-credentials-demo role not found: %v", err)
	}

	// Create IAM user
	_, err = iamClient.CreateUser(ctx, &iam.CreateUserInput{
		UserName:            aws.String(userName),
		PermissionsBoundary: aws.String(demoUserPolicyArn),
	})
	if err != nil {
		t.Fatalf("failed to create IAM user: %v", err)
	}

	// Attach policy to user
	_, err = iamClient.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
		UserName:  aws.String(userName),
		PolicyArn: aws.String(demoUserPolicyArn),
	})
	if err != nil {
		t.Fatalf("failed to attach policy: %v", err)
	}

	// Create access key
	keyOut, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		t.Fatalf("failed to create access key: %v", err)
	}
	accessKeyID = aws.ToString(keyOut.AccessKey.AccessKeyId)
	secretAccessKey = aws.ToString(keyOut.AccessKey.SecretAccessKey)

	// IAM is eventually consistent; wait briefly before verifying the user is readable.
	t.Logf("Verifying IAM user %s exists...", userName)
	waitTime := 10 * time.Second
	verifyDeadline := time.Now().Add(waitTime * 2)
	var lastErr error
	for time.Now().Before(verifyDeadline) {
		time.Sleep(waitTime)
		_, lastErr = iamClient.GetUser(ctx, &iam.GetUserInput{UserName: aws.String(userName)})
		if lastErr == nil {
			break
		}
		t.Logf("IAM user %q not readable yet; retrying: %v", userName, lastErr)
	}
	if lastErr != nil {
		t.Fatalf("failed to verify IAM user %q: %v", userName, lastErr)
	}

	return userName, accessKeyID, secretAccessKey, demoUserPolicyArn, assumedRoleArn, awsRegion
}

// deleteIAMUserByAccessKey finds and deletes the IAM user owning the specified access key.
func deleteIAMUserByAccessKey(t *testing.T, targetAccessKeyID string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}
	iamClient := iam.NewFromConfig(cfg)

	paginator := iam.NewListUsersPaginator(iamClient, &iam.ListUsersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			t.Fatalf("failed to list IAM users: %v", err)
		}
		for _, user := range page.Users {
			userName := aws.ToString(user.UserName)
			if !strings.Contains(userName, "demo-GitHubActions") {
				continue
			}
			// List all access keys for this user
			keyPaginator := iam.NewListAccessKeysPaginator(iamClient, &iam.ListAccessKeysInput{
				UserName: &userName,
			})
			for keyPaginator.HasMorePages() {
				keyPage, err := keyPaginator.NextPage(ctx)
				if err != nil {
					t.Logf("warning: failed to list access keys for user %q: %v", userName, err)
					continue
				}
				for _, key := range keyPage.AccessKeyMetadata {
					accessKeyId := aws.ToString(key.AccessKeyId)
					if accessKeyId == targetAccessKeyID {
						// Found the user with the target access key. Detach managed policies first,
						// then delete all access keys, then delete the user.
						policyPaginator := iam.NewListAttachedUserPoliciesPaginator(iamClient, &iam.ListAttachedUserPoliciesInput{
							UserName: &userName,
						})
						for policyPaginator.HasMorePages() {
							policyPage, err := policyPaginator.NextPage(ctx)
							if err != nil {
								t.Logf("warning: failed to list attached policies for user %q: %v", userName, err)
								continue
							}
							for _, policy := range policyPage.AttachedPolicies {
								policyArn := aws.ToString(policy.PolicyArn)
								if _, err := iamClient.DetachUserPolicy(ctx, &iam.DetachUserPolicyInput{
									UserName:  &userName,
									PolicyArn: &policyArn,
								}); err != nil {
									t.Logf("warning: failed to detach policy %q from user %q: %v", policyArn, userName, err)
								}
							}
						}
						keyPaginator2 := iam.NewListAccessKeysPaginator(iamClient, &iam.ListAccessKeysInput{
							UserName: &userName,
						})
						for keyPaginator2.HasMorePages() {
							keyPage2, err := keyPaginator2.NextPage(ctx)
							if err != nil {
								t.Logf("warning: failed to list access keys for cleanup on user %q: %v", userName, err)
								continue
							}
							for _, key2 := range keyPage2.AccessKeyMetadata {
								accessKeyId2 := aws.ToString(key2.AccessKeyId)
								if _, err := iamClient.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
									UserName:    &userName,
									AccessKeyId: &accessKeyId2,
								}); err != nil {
									t.Logf("warning: failed to delete access key %q for user %q: %v", accessKeyId2, userName, err)
								}
							}
						}
						// Delete the user
						if _, err := iamClient.DeleteUser(ctx, &iam.DeleteUserInput{
							UserName: &userName,
						}); err != nil {
							t.Logf("warning: failed to delete user %q: %v", userName, err)
						}
						return
					}
				}
			}
		}
	}
}

// getAwsUsernameTemplate builds a Vault username template for AWS credential generation.
func getAwsUsernameTemplate(awsUserName string) string {
	const prefix = `{{ if (eq .Type "STS") }}{{ printf "`
	const stsSuffix = `-%s-%s" (random 20) (unix_time) | truncate 32 }}{{ else }}{{ printf "`
	const iamUserSuffix = `-%s-%s" (unix_time) (random 20) | truncate 60 }}{{ end }}`
	return prefix + awsUserName + stsSuffix + awsUserName + iamUserSuffix
}

// getAllowDescribeRegionsPolicy returns an IAM policy allowing ec2:DescribeRegions.
func getAllowDescribeRegionsPolicy() string {
	return `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["ec2:DescribeRegions"],
      "Resource": ["*"]
    }
  ]
}`
}

// =============================================================================
// Vault AWS Secrets Engine Helpers
// =============================================================================

// skipIfNoAWSCredentials skips the test if AWS credentials are missing.
func skipIfNoAWSCredentials(t *testing.T) {
	t.Helper()
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		t.Skip("AWS credentials not available - skipping AWS secrets engine test")
	}
}

// setupAWSSecretsEngine enables and configures AWS secrets engine, returns mount path.
func setupAWSSecretsEngine(t *testing.T, v *blackbox.Session, accessKeyID, secretAccessKey string, usernameTemplate string) string {
	t.Helper()

	path := fmt.Sprintf("aws-test-%d", time.Now().UnixNano())
	t.Logf("Enabling AWS secrets engine at path: %s", path)
	v.MustEnableSecretsEngine(path, &api.MountInput{Type: "aws"})

	config := map[string]any{
		"access_key": accessKeyID,
		"secret_key": secretAccessKey,
		"region":     "us-east-1",
	}

	if usernameTemplate != "" {
		config["username_template"] = usernameTemplate
		t.Logf("Configuring AWS secrets engine with username template")
	} else {
		t.Logf("Configuring AWS secrets engine with root credentials")
	}

	v.MustWrite(fmt.Sprintf("%s/config/root", path), config)
	return path
}

// createVaultAWSRole creates a Vault AWS role with IAM user credential type.
func createVaultAWSRole(t *testing.T, v *blackbox.Session, path, roleName, policyArn string) {
	t.Helper()

	t.Logf("Creating Vault AWS role: %s", roleName)
	v.MustWrite(fmt.Sprintf("%s/roles/%s", path, roleName), map[string]any{
		"credential_type":          "iam_user",
		"permissions_boundary_arn": policyArn,
		"policy_document":          getAllowDescribeRegionsPolicy(),
	})
}

// verifyRoleExists checks that a Vault AWS role exists in the role list.
func verifyRoleExists(t *testing.T, v *blackbox.Session, path, roleName string) {
	t.Helper()

	roleList := v.MustList(fmt.Sprintf("%s/roles", path))
	if roleList == nil || roleList.Data == nil {
		t.Fatalf("failed to list roles at path %s", path)
	}

	roleKeys, ok := roleList.Data["keys"].([]interface{})
	if !ok || len(roleKeys) == 0 {
		t.Fatalf("no roles found at path %s", path)
	}

	for _, key := range roleKeys {
		if keyStr, ok := key.(string); ok && keyStr == roleName {
			return // Role found
		}
	}

	t.Fatalf("role %q not found in list: %v", roleName, roleKeys)
}

// verifyRoleDeleted checks that a Vault AWS role no longer exists in the role list.
func verifyRoleDeleted(t *testing.T, v *blackbox.Session, path, roleName string) {
	t.Helper()

	t.Logf("Verifying role %q was deleted", roleName)
	rolesList := v.MustList(fmt.Sprintf("%s/roles", path))
	if rolesList != nil && rolesList.Data != nil {
		if keys, ok := rolesList.Data["keys"].([]interface{}); ok {
			for _, key := range keys {
				if keyStr, ok := key.(string); ok && keyStr == roleName {
					t.Fatalf("role %q still exists after deletion", roleName)
				}
			}
		}
	}
	t.Logf("Successfully verified role %q was deleted", roleName)
}
