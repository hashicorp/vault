// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/testhelpers"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

var initSetup sync.Once

// This looks a bit curious. The policy document and the role document act
// as a logical intersection of policies. The role allows ec2:Describe*
// (among other permissions). This policy allows everything BUT
// ec2:DescribeAvailabilityZones. Thus, the logical intersection of the two
// is all ec2:Describe* EXCEPT ec2:DescribeAvailabilityZones, and so the
// describeAZs call should fail
const allowAllButDescribeAzs = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"NotAction": "ec2:DescribeAvailabilityZones",
			"Resource": "*"
		}
	]
}`

type mockIAMClient struct {
	iamAPI
}

func (m *mockIAMClient) CreateUser(_ context.Context, input *iam.CreateUserInput, _ ...func(*iam.Options)) (*iam.CreateUserOutput, error) {
	return nil, &smithy.GenericAPIError{Code: "Throttling", Message: ""}
}

func getBackend(t *testing.T) logical.Backend {
	cfg := logical.TestBackendConfig()
	cfg.System = &testSystemView{StaticSystemView: *logical.TestSystemView()}
	be, _ := Factory(t.Context(), cfg)
	return be
}

func TestAcceptanceBackend_basic(t *testing.T) {
	t.Parallel()
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWritePolicy(t, "test", testDynamoPolicy),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{listDynamoTablesTest}),
		},
	})
}

func TestAcceptanceBackend_IamUserWithPermissionsBoundary(t *testing.T) {
	t.Parallel()
	roleData := map[string]interface{}{
		"credential_type":          iamUserCred,
		"policy_arns":              adminAccessPolicyArn,
		"permissions_boundary_arn": iamPolicyArn,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{listIamUsersTest, describeAzsTestUnauthorized}),
		},
	})
}

func TestAcceptanceBackend_basicSTS(t *testing.T) {
	t.Parallel()
	awsAccountID, err := getAccountID(t.Context())
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}
	roleName := generateUniqueRoleName(t.Name())
	userName := generateUniqueUserName(t.Name())
	accessKey := &awsAccessKey{}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createUser(t, userName, accessKey)
			createRole(t, roleName, awsAccountID, []string{ec2PolicyArn}, nil)
			// Sleep sometime because AWS is eventually consistent
			// Both the createUser and createRole depend on this
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfigWithCreds(t, accessKey),
			testAccStepRotateRoot(accessKey),
			testAccStepWritePolicy(t, "test", testDynamoPolicy),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{listDynamoTablesTest}),
			testAccStepWriteArnPolicyRef(t, "test", ec2PolicyArn),
			testAccStepReadSTSWithArnPolicy(t, "test"),
			testAccStepWriteArnRoleRef(t, "test2", roleName, awsAccountID),
			testAccStepRead(t, "sts", "test2", []credentialTestFunc{describeInstancesTest}),
		},
		Teardown: func() error {
			if err := deleteTestRole(roleName); err != nil {
				return err
			}
			return deleteTestUser(accessKey, userName)
		},
	})
}

// TestBackend_policyCRUD tests the CRUD operations for a policy.
func TestBackend_policyCRUD(t *testing.T) {
	t.Parallel()
	compacted, err := compactJSON(testDynamoPolicy)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWritePolicy(t, "test", testDynamoPolicy),
			testAccStepReadPolicy(t, "test", compacted),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", ""),
		},
	})
}

func TestBackend_throttled(t *testing.T) {
	t.Parallel()
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(t.Context(), config); err != nil {
		t.Fatal(err)
	}

	connData := map[string]interface{}{
		"credential_type": "iam_user",
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/something",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(t.Context(), confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	// Mock the IAM API call to return a throttled response to the CreateUser API
	// call
	b.iamClient = &mockIAMClient{}

	credReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "creds/something",
		Storage:   config.StorageView,
	}

	credResp, err := b.HandleRequest(t.Context(), credReq)
	if err == nil {
		t.Fatalf("failed to trigger expected throttling error condition: resp:%#v", credResp)
	}
	rErr := credResp.Error()
	expected := "Error creating IAM user: api error Throttling: "
	if rErr.Error() != expected {
		t.Fatalf("error message did not match, expected (%s), got (%s)", expected, rErr.Error())
	}

	// verify the error we got back is returned with a http.StatusBadGateway
	code, err := logical.RespondErrorCommon(credReq, credResp, err)
	if err == nil {
		t.Fatal("expected error after running req/resp/err through RespondErrorCommon, got nil")
	}
	if code != http.StatusBadGateway {
		t.Fatalf("expected HTTP status 'bad gateway', got: (%d)", code)
	}
}

func testAccPreCheck(t *testing.T) {
	if !hasAWSCredentials() {
		t.Skip("Skipping because AWS credentials could not be resolved. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials for information on how to set up AWS credentials.")
	}

	initSetup.Do(func() {
		if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
			log.Println("[INFO] Test: Using us-west-2 as test region")
			os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
		}
	})
}

func hasAWSCredentials() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return false
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return false
	}

	return creds.HasKeys()
}

func getAccountID(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		return "", err
	}
	svc := sts.NewFromConfig(cfg)

	params := &sts.GetCallerIdentityInput{}
	res, err := svc.GetCallerIdentity(ctx, params)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("got nil response from GetCallerIdentity")
	}

	return *res.Account, nil
}

func createRole(t *testing.T, roleName, awsAccountID string, policyARNs, extraTrustPolicies []string) {
	t.Helper()

	ctx := t.Context()

	trustPolicyStmts := append([]string{
		fmt.Sprintf(`
		{
		  "Effect":"Allow",
		  "Principal": {
			  "AWS": "arn:aws:iam::%s:root"
		  },
		  "Action": [
			  "sts:AssumeRole",
			  "sts:SetSourceIdentity"
		  ]
		}`, awsAccountID),
	},
		extraTrustPolicies...)

	testRoleAssumePolicy := fmt.Sprintf(`{
      "Version": "2012-10-17",
      "Statement": [
%s
      ]
}
`, strings.Join(trustPolicyStmts, ","))

	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		t.Fatal(err)
	}
	svc := iam.NewFromConfig(awsConfig)

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(testRoleAssumePolicy),
		RoleName:                 aws.String(roleName),
		Path:                     aws.String("/"),
	}

	log.Printf("[INFO] AWS CreateRole: %s", roleName)
	output, err := svc.CreateRole(ctx, params)
	if err != nil {
		t.Fatalf("AWS CreateRole failed: %v", err)
	}

	for _, policyARN := range policyARNs {
		attachment := &iam.AttachRolePolicyInput{
			PolicyArn: aws.String(policyARN),
			RoleName:  output.Role.RoleName,
		}
		_, err = svc.AttachRolePolicy(ctx, attachment)
		if err != nil {
			t.Fatalf("AWS AttachRolePolicy failed: %v", err)
		}
	}
}

func createUser(t *testing.T, userName string, accessKey *awsAccessKey) {
	// The sequence of user creation actions is carefully chosen to minimize
	// impact of stolen IAM user credentials
	// 1. Create user, without any permissions or credentials. At this point,
	//	  nobody cares if creds compromised because this user can do nothing.
	// 2. Attach the timebomb policy. This grants no access but puts a time limit
	//	  on validity of compromised credentials. If this fails, nobody cares
	//	  because the user has no permissions to do anything anyway
	// 3. Attach the AdminAccess policy. The IAM user still has no credentials to
	//	  do anything
	// 4. Generate API creds to get an actual access key and secret key
	timebombPolicyTemplate := `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Deny",
			"Action": "*",
			"Resource": "*",
			"Condition": {
				"DateGreaterThan": {
					"aws:CurrentTime": "%s"
				}
			}
		}
	]
}
`
	validity := time.Duration(2 * time.Hour)
	expiry := time.Now().Add(validity)
	timebombPolicy := fmt.Sprintf(timebombPolicyTemplate, expiry.Format(time.RFC3339))
	ctx := t.Context()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		t.Fatal(err)
	}
	svc := iam.NewFromConfig(awsConfig)
	createUserInput := &iam.CreateUserInput{
		UserName: aws.String(userName),
	}
	log.Printf("[INFO] AWS CreateUser: %s", userName)
	if _, err := svc.CreateUser(ctx, createUserInput); err != nil {
		t.Fatalf("AWS CreateUser failed: %v", err)
	}

	putPolicyInput := &iam.PutUserPolicyInput{
		PolicyDocument: aws.String(timebombPolicy),
		PolicyName:     aws.String("SelfDestructionTimebomb"),
		UserName:       aws.String(userName),
	}
	_, err = svc.PutUserPolicy(ctx, putPolicyInput)
	if err != nil {
		t.Fatalf("AWS PutUserPolicy failed: %v", err)
	}

	attachUserPolicyInput := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		UserName:  aws.String(userName),
	}
	_, err = svc.AttachUserPolicy(ctx, attachUserPolicyInput)
	if err != nil {
		t.Fatalf("AWS AttachUserPolicy failed, %v", err)
	}

	createAccessKeyInput := &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	}
	createAccessKeyOutput, err := svc.CreateAccessKey(ctx, createAccessKeyInput)
	if err != nil {
		t.Fatalf("AWS CreateAccessKey failed: %v", err)
	}
	if createAccessKeyOutput == nil {
		t.Fatalf("AWS CreateAccessKey returned nil")
	}
	genAccessKey := createAccessKeyOutput.AccessKey

	accessKey.AccessKeyID = *genAccessKey.AccessKeyId
	accessKey.SecretAccessKey = *genAccessKey.SecretAccessKey
}

// Create an IAM Group and add an inline policy and managed policies if specified
func createGroup(t *testing.T, groupName string, inlinePolicy string, managedPolicies []string) {
	ctx := t.Context()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		t.Fatal(err)
	}
	svc := iam.NewFromConfig(awsConfig)
	createGroupInput := &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	}
	log.Printf("[INFO] AWS CreateGroup: %s", groupName)
	if _, err := svc.CreateGroup(ctx, createGroupInput); err != nil {
		t.Fatalf("AWS CreateGroup failed: %v", err)
	}

	if len(inlinePolicy) > 0 {
		putPolicyInput := &iam.PutGroupPolicyInput{
			PolicyDocument: aws.String(inlinePolicy),
			PolicyName:     aws.String("InlinePolicy"),
			GroupName:      aws.String(groupName),
		}
		_, err = svc.PutGroupPolicy(ctx, putPolicyInput)
		if err != nil {
			t.Fatalf("AWS PutGroupPolicy failed: %v", err)
		}
	}

	for _, mp := range managedPolicies {
		attachGroupPolicyInput := &iam.AttachGroupPolicyInput{
			PolicyArn: aws.String(mp),
			GroupName: aws.String(groupName),
		}
		_, err = svc.AttachGroupPolicy(ctx, attachGroupPolicyInput)
		if err != nil {
			t.Fatalf("AWS AttachGroupPolicy failed, %v", err)
		}
	}
}

func deleteTestRole(roleName string) error {
	awsConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		return err
	}
	svc := iam.NewFromConfig(awsConfig)
	listAttachmentsInput := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	paginator := iam.NewListAttachedRolePoliciesPaginator(svc, listAttachmentsInput)
	for paginator.HasMorePages() {
		result, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Printf("[WARN] AWS ListAttachedRolePolicies failed: %v", err)
			break
		}
		for _, policy := range result.AttachedPolicies {
			detachInput := &iam.DetachRolePolicyInput{
				PolicyArn: policy.PolicyArn,
				RoleName:  aws.String(roleName), // Required
			}
			_, err := svc.DetachRolePolicy(context.Background(), detachInput)
			if err != nil {
				log.Printf("[WARN] AWS DetachRolePolicy failed for policy %s: %v", *policy.PolicyArn, err)
			}
		}
	}

	params := &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	}

	log.Printf("[INFO] AWS DeleteRole: %s", roleName)
	_, err = svc.DeleteRole(context.Background(), params)
	if err != nil {
		log.Printf("[WARN] AWS DeleteRole failed: %v", err)
		return err
	}
	return nil
}

func deleteTestUser(accessKey *awsAccessKey, userName string) error {
	awsConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		return err
	}
	svc := iam.NewFromConfig(awsConfig)
	userDetachment := &iam.DetachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		UserName:  aws.String(userName),
	}
	if _, err := svc.DetachUserPolicy(context.Background(), userDetachment); err != nil {
		log.Printf("[WARN] AWS DetachUserPolicy failed: %v", err)
		return err
	}

	deleteAccessKeyInput := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accessKey.AccessKeyID),
		UserName:    aws.String(userName),
	}
	_, err = svc.DeleteAccessKey(context.Background(), deleteAccessKeyInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteAccessKey failed: %v", err)
		return err
	}

	deleteTestUserPolicyInput := &iam.DeleteUserPolicyInput{
		PolicyName: aws.String("SelfDestructionTimebomb"),
		UserName:   aws.String(userName),
	}
	_, err = svc.DeleteUserPolicy(context.Background(), deleteTestUserPolicyInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteUserPolicy failed: %v", err)
		return err
	}
	deleteTestUserInput := &iam.DeleteUserInput{
		UserName: aws.String(userName),
	}
	log.Printf("[INFO] AWS DeleteUser: %s", userName)
	_, err = svc.DeleteUser(context.Background(), deleteTestUserInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteUser failed: %v", err)
		return err
	}

	return nil
}

func deleteTestGroup(groupName string) error {
	awsConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
	)
	if err != nil {
		return err
	}
	svc := iam.NewFromConfig(awsConfig)

	// Detach any managed group policies
	getGroupsInput := &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(groupName),
	}
	getGroupsOutput, err := svc.ListAttachedGroupPolicies(context.Background(), getGroupsInput)
	if err != nil {
		log.Printf("[WARN] AWS ListAttachedGroupPolicies failed: %v", err)
		return err
	}
	for _, g := range getGroupsOutput.AttachedPolicies {
		detachGroupInput := &iam.DetachGroupPolicyInput{
			GroupName: aws.String(groupName),
			PolicyArn: g.PolicyArn,
		}
		if _, err := svc.DetachGroupPolicy(context.Background(), detachGroupInput); err != nil {
			log.Printf("[WARN] AWS DetachGroupPolicy failed: %v", err)
			return err
		}
	}

	// Remove any inline policies
	listGroupPoliciesInput := &iam.ListGroupPoliciesInput{
		GroupName: aws.String(groupName),
	}
	listGroupPoliciesOutput, err := svc.ListGroupPolicies(context.Background(), listGroupPoliciesInput)
	if err != nil {
		log.Printf("[WARN] AWS ListGroupPolicies failed: %v", err)
		return err
	}
	for _, g := range listGroupPoliciesOutput.PolicyNames {
		deleteGroupPolicyInput := &iam.DeleteGroupPolicyInput{
			GroupName:  aws.String(groupName),
			PolicyName: aws.String(g),
		}
		if _, err := svc.DeleteGroupPolicy(context.Background(), deleteGroupPolicyInput); err != nil {
			log.Printf("[WARN] AWS DeleteGroupPolicy failed: %v", err)
			return err
		}
	}

	// Delete the group
	deleteTestGroupInput := &iam.DeleteGroupInput{
		GroupName: aws.String(groupName),
	}
	log.Printf("[INFO] AWS DeleteGroup: %s", groupName)
	_, err = svc.DeleteGroup(context.Background(), deleteTestGroupInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteGroup failed: %v", err)
		return err
	}

	return nil
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/root",
		Data: map[string]interface{}{
			"region": os.Getenv("AWS_DEFAULT_REGION"),
		},
	}
}

func testAccStepConfigWithCreds(t *testing.T, accessKey *awsAccessKey) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/root",
		Data: map[string]interface{}{
			"region": os.Getenv("AWS_DEFAULT_REGION"),
		},
		PreFlight: func(req *logical.Request) error {
			// Values in Data above get eagerly evaluated due to the testing framework.
			// In particular, they get evaluated before accessKey gets set by CreateUser
			// and thus would fail. By moving to a closure in a PreFlight, we ensure that
			// the creds get evaluated lazily after they've been properly set
			req.Data["access_key"] = accessKey.AccessKeyID
			req.Data["secret_key"] = accessKey.SecretAccessKey
			return nil
		},
	}
}

func testAccStepRotateRoot(oldAccessKey *awsAccessKey) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("received nil response from config/rotate-root")
			}
			newAccessKeyID := resp.Data["access_key"].(string)
			if newAccessKeyID == oldAccessKey.AccessKeyID {
				return fmt.Errorf("rotate-root didn't rotate access key")
			}
			ctx := context.Background()
			// Save the old key ID before updating the struct so we can verify
			// the old key was deleted after rotation.
			oldKeyID := oldAccessKey.AccessKeyID
			oldAccessKey.AccessKeyID = newAccessKeyID
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
			awsConfig, err := config.LoadDefaultConfig(ctx,
				config.WithRegion("us-east-1"),
				config.WithHTTPClient(cleanhttp.DefaultClient()),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(oldKeyID, oldAccessKey.SecretAccessKey, "")),
			)
			if err != nil {
				return err
			}
			svc := sts.NewFromConfig(awsConfig)
			params := &sts.GetCallerIdentityInput{}
			if _, err := svc.GetCallerIdentity(ctx, params); err == nil {
				return fmt.Errorf("bad: old credentials succeeded after rotate")
			} else {
				var apiErr smithy.APIError
				if errors.As(err, &apiErr) {
					if apiErr.ErrorCode() != "InvalidClientTokenId" {
						return fmt.Errorf("Unknown error returned from AWS: %#v", apiErr)
					}
					return nil
				}
				return err
			}
		},
	}
}

func testAccStepRead(_ *testing.T, path, name string, credentialTests []credentialTestFunc) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path + "/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				AccessKey string `mapstructure:"access_key"`
				SecretKey string `mapstructure:"secret_key"`
				STSToken  string `mapstructure:"session_token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated credentials: %v", d)
			for _, test := range credentialTests {
				err := test(d.AccessKey, d.SecretKey, d.STSToken)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func testAccStepReadWithMFA(t *testing.T, path, name, mfaCode string, credentialTests []credentialTestFunc) logicaltest.TestStep {
	step := testAccStepRead(t, path, name, credentialTests)
	step.Data = map[string]interface{}{
		"mfa_code": mfaCode,
	}

	return step
}

func testAccStepReadSTSResponse(name string, maximumTTL time.Duration) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			if resp.Secret == nil {
				return fmt.Errorf("bad: nil Secret returned")
			}
			ttl := resp.Secret.TTL
			if ttl > maximumTTL {
				return fmt.Errorf("bad: ttl of %d greater than maximum of %d", ttl/time.Second, maximumTTL/time.Second)
			}
			return nil
		},
	}
}

func describeInstancesTest(accessKey, secretKey, token string) error {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		return err
	}
	client := ec2.NewFromConfig(awsConfig)
	log.Printf("[WARN] Verifying that the generated credentials work with ec2:DescribeInstances...")
	return retryUntilSuccess(func() error {
		_, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
		return err
	})
}

func describeAzsTestUnauthorized(accessKey, secretKey, token string) error {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		return err
	}
	client := ec2.NewFromConfig(awsConfig)
	log.Printf("[WARN] Verifying that the generated credentials don't work with ec2:DescribeAvailabilityZones...")
	return retryUntilSuccess(func() error {
		_, err := client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{})
		// Need to make sure AWS authenticates the generated credentials but does not authorize the operation
		if err == nil {
			return fmt.Errorf("operation succeeded when expected failure")
		}
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "UnauthorizedOperation" {
				return nil
			}
		}
		return err
	})
}

func assertCreatedIAMUser(accessKey, secretKey, token string) error {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		return err
	}
	client := iam.NewFromConfig(awsConfig)
	log.Printf("[WARN] Checking if IAM User is created properly...")
	userOutput, err := client.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		return err
	}

	if *userOutput.User.Path != "/path/" {
		return fmt.Errorf("bad: got: %#v\nexpected: %#v", userOutput.User.Path, "/path/")
	}

	return nil
}

func listIamUsersTest(accessKey, secretKey, token string) error {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		return err
	}
	client := iam.NewFromConfig(awsConfig)
	log.Printf("[WARN] Verifying that the generated credentials work with iam:ListUsers...")
	return retryUntilSuccess(func() error {
		_, err := client.ListUsers(ctx, &iam.ListUsersInput{})
		return err
	})
}

func listDynamoTablesTest(accessKey, secretKey, token string) error {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(awsConfig)
	log.Printf("[WARN] Verifying that the generated credentials work with dynamodb:ListTables...")
	return retryUntilSuccess(func() error {
		_, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
		return err
	})
}

func listS3BucketsTest(accessKey, secretKey, token string) error {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(cleanhttp.DefaultClient()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(awsConfig)
	log.Printf("[WARN] Verifying that the generated credentials work with s3:ListBuckets...")
	return retryUntilSuccess(func() error {
		_, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		return err
	})
}

func retryUntilSuccess(op func() error) error {
	retryCount := 0
	success := false
	var err error
	for !success && retryCount < 10 {
		err = op()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
		retryCount++
	}
	return err
}

func testAccStepReadSTSWithArnPolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "sts/" + name,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp.Data["error"] !=
				"attempted to retrieve iam_user credentials through the sts path; this is not allowed for legacy roles" {
				t.Fatalf("bad: %v", resp)
			}
			return nil
		},
	}
}

func testAccStepWritePolicy(t *testing.T, name string, policy string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"policy": policy,
		},
	}
}

func testAccStepDeletePolicy(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadPolicy(t *testing.T, name string, value string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if value == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			expected := map[string]interface{}{
				"policy_arns":              []string(nil),
				"role_arns":                []string(nil),
				"policy_document":          value,
				"credential_type":          strings.Join([]string{iamUserCred, federationTokenCred}, ","),
				"default_sts_ttl":          int64(0),
				"max_sts_ttl":              int64(0),
				"user_path":                "",
				"permissions_boundary_arn": "",
				"iam_groups":               []string(nil),
				"iam_tags":                 map[string]string(nil),
				"mfa_serial_number":        "",
				"session_tags":             map[string]string(nil),
				"external_id":              "",
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("bad: got: %#v\nexpected: %#v", resp.Data, expected)
			}
			return nil
		},
	}
}

const testDynamoPolicy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1426528957000",
            "Effect": "Allow",
            "Action": [
                "dynamodb:List*"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
`

const testS3Policy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:Get*",
                "s3:List*"
            ],
            "Resource": "*"
        }
    ]
}`

const (
	adminAccessPolicyArn = "arn:aws:iam::aws:policy/AdministratorAccess"
	ec2PolicyArn         = "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"
	iamPolicyArn         = "arn:aws:iam::aws:policy/IAMReadOnlyAccess"
	dynamoPolicyArn      = "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"
)

func testAccStepWriteRole(t *testing.T, name string, data map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data:      data,
	}
}

func testAccStepReadRole(t *testing.T, name string, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if expected == nil {
					return nil
				}
				return fmt.Errorf("bad: nil response")
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("bad: got %#v\nexpected: %#v", resp.Data, expected)
			}
			return nil
		},
	}
}

func testAccStepWriteArnPolicyRef(t *testing.T, name string, arn string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"arn": ec2PolicyArn,
		},
	}
}

func TestAcceptanceBackend_basicPolicyArnRef(t *testing.T) {
	t.Parallel()
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteArnPolicyRef(t, "test", ec2PolicyArn),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest}),
		},
	})
}

func TestAcceptanceBackend_iamUserManagedInlinePoliciesGroups(t *testing.T) {
	t.Parallel()
	compacted, err := compactJSON(testDynamoPolicy)
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	groupName := generateUniqueGroupName(t.Name())
	roleData := map[string]interface{}{
		"policy_document": testDynamoPolicy,
		"policy_arns":     []string{ec2PolicyArn, iamPolicyArn},
		"iam_groups":      []string{groupName},
		"credential_type": iamUserCred,
		"user_path":       "/path/",
	}
	expectedRoleData := map[string]interface{}{
		"policy_document":          compacted,
		"policy_arns":              []string{ec2PolicyArn, iamPolicyArn},
		"credential_type":          iamUserCred,
		"role_arns":                []string(nil),
		"default_sts_ttl":          int64(0),
		"max_sts_ttl":              int64(0),
		"user_path":                "/path/",
		"permissions_boundary_arn": "",
		"iam_groups":               []string{groupName},
		"iam_tags":                 map[string]string(nil),
		"mfa_serial_number":        "",
		"session_tags":             map[string]string(nil),
		"external_id":              "",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createGroup(t, groupName, testS3Policy, []string{})
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepReadRole(t, "test", expectedRoleData),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest, listIamUsersTest, listDynamoTablesTest, assertCreatedIAMUser, listS3BucketsTest}),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{describeInstancesTest, listIamUsersTest, listDynamoTablesTest, listS3BucketsTest}),
		},
		Teardown: func() error {
			return deleteTestGroup(groupName)
		},
	})
}

// Similar to TestBackend_iamUserManagedInlinePoliciesGroups() but managing
// policies only with groups
func TestAcceptanceBackend_iamUserGroups(t *testing.T) {
	t.Parallel()
	group1Name := generateUniqueGroupName(t.Name())
	group2Name := generateUniqueGroupName(t.Name())
	roleData := map[string]interface{}{
		"iam_groups":      []string{group1Name, group2Name},
		"credential_type": iamUserCred,
		"user_path":       "/path/",
	}
	expectedRoleData := map[string]interface{}{
		"policy_document":          "",
		"policy_arns":              []string(nil),
		"credential_type":          iamUserCred,
		"role_arns":                []string(nil),
		"default_sts_ttl":          int64(0),
		"max_sts_ttl":              int64(0),
		"user_path":                "/path/",
		"permissions_boundary_arn": "",
		"iam_groups":               []string{group1Name, group2Name},
		"iam_tags":                 map[string]string(nil),
		"mfa_serial_number":        "",
		"session_tags":             map[string]string(nil),
		"external_id":              "",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createGroup(t, group1Name, testS3Policy, []string{ec2PolicyArn, iamPolicyArn})
			createGroup(t, group2Name, testDynamoPolicy, []string{})
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepReadRole(t, "test", expectedRoleData),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest, listIamUsersTest, listDynamoTablesTest, assertCreatedIAMUser, listS3BucketsTest}),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{describeInstancesTest, listIamUsersTest, listDynamoTablesTest, listS3BucketsTest}),
		},
		Teardown: func() error {
			if err := deleteTestGroup(group1Name); err != nil {
				return err
			}
			return deleteTestGroup(group2Name)
		},
	})
}

func TestAcceptanceBackend_AssumedRoleWithPolicyDoc(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueRoleName(t.Name())

	awsAccountID, err := getAccountID(t.Context())
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}
	roleData := map[string]interface{}{
		"policy_document": allowAllButDescribeAzs,
		"role_arns":       []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, roleName)},
		"credential_type": assumedRoleCred,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createRole(t, roleName, awsAccountID, []string{ec2PolicyArn}, nil)
			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{describeInstancesTest, describeAzsTestUnauthorized}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest, describeAzsTestUnauthorized}),
		},
		Teardown: func() error {
			return deleteTestRole(roleName)
		},
	})
}

func TestAcceptanceBackend_AssumedRoleWithPolicyARN(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueRoleName(t.Name())

	awsAccountID, err := getAccountID(t.Context())
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}
	roleData := map[string]interface{}{
		"policy_arns":     iamPolicyArn,
		"role_arns":       []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, roleName)},
		"credential_type": assumedRoleCred,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createRole(t, roleName, awsAccountID, []string{ec2PolicyArn, iamPolicyArn}, nil)
			log.Printf("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{listIamUsersTest, describeAzsTestUnauthorized}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{listIamUsersTest, describeAzsTestUnauthorized}),
		},
		Teardown: func() error {
			return deleteTestRole(roleName)
		},
	})
}

func TestAcceptanceBackend_AssumedRoleWithGroups(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueRoleName(t.Name())
	groupName := generateUniqueGroupName(t.Name())

	awsAccountID, err := getAccountID(t.Context())
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}

	roleData := map[string]interface{}{
		"iam_groups":      []string{groupName},
		"role_arns":       []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, roleName)},
		"credential_type": assumedRoleCred,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createRole(t, roleName, awsAccountID, []string{ec2PolicyArn}, nil)
			createGroup(t, groupName, allowAllButDescribeAzs, []string{})
			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{describeInstancesTest, describeAzsTestUnauthorized}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest, describeAzsTestUnauthorized}),
		},
		Teardown: func() error {
			if err := deleteTestGroup(groupName); err != nil {
				return err
			}
			return deleteTestRole(roleName)
		},
	})
}

// TestAcceptanceBackend_AssumedRoleWithSessionTags tests that session tags are
// passed to the assumed role.
func TestAcceptanceBackend_AssumedRoleWithSessionTags(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueRoleName(t.Name())
	awsAccountID, err := getAccountID(t.Context())
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}

	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, roleName)
	roleData := map[string]interface{}{
		"policy_document": allowAllButDescribeAzs,
		"role_arns":       []string{roleARN},
		"credential_type": assumedRoleCred,
		"session_tags": map[string]string{
			"foo": "bar",
			"baz": "qux",
		},
	}

	// allowSessionTagsPolicy allows the role to tag the session, it needs to be
	// included in the trust policy.
	allowSessionTagsPolicy := fmt.Sprintf(`
		{
			"Sid": "AllowPassSessionTagsAndTransitive",
			"Effect": "Allow",
			"Action": "sts:TagSession",
			"Principal": {
				  "AWS": "arn:aws:iam::%s:root"
			}
		}
`, awsAccountID)

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createRole(t, roleName, awsAccountID, []string{ec2PolicyArn}, []string{allowSessionTagsPolicy})
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{describeInstancesTest, describeAzsTestUnauthorized}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest, describeAzsTestUnauthorized}),
		},
		Teardown: func() error {
			return deleteTestRole(roleName)
		},
	})
}

func TestAcceptanceBackend_FederationTokenWithPolicyARN(t *testing.T) {
	t.Parallel()
	userName := generateUniqueUserName(t.Name())
	accessKey := &awsAccessKey{}

	roleData := map[string]interface{}{
		"policy_arns":     dynamoPolicyArn,
		"credential_type": federationTokenCred,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createUser(t, userName, accessKey)
			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfigWithCreds(t, accessKey),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{listDynamoTablesTest, describeAzsTestUnauthorized}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{listDynamoTablesTest, describeAzsTestUnauthorized}),
		},
		Teardown: func() error {
			return deleteTestUser(accessKey, userName)
		},
	})
}

func TestAcceptanceBackend_FederationTokenWithGroups(t *testing.T) {
	t.Parallel()
	userName := generateUniqueUserName(t.Name())
	groupName := generateUniqueGroupName(t.Name())
	accessKey := &awsAccessKey{}

	// IAM policy where Statement is a single element, not a list
	iamSingleStatementPolicy := `{
		"Version": "2012-10-17",
		"Statement": {
			"Effect": "Allow",
			"Action": [
				"s3:Get*",
				"s3:List*"
			],
			"Resource": "*"
		}
	}`

	roleData := map[string]interface{}{
		"iam_groups":      []string{groupName},
		"policy_document": iamSingleStatementPolicy,
		"credential_type": federationTokenCred,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createUser(t, userName, accessKey)
			createGroup(t, groupName, "", []string{dynamoPolicyArn})
			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfigWithCreds(t, accessKey),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{listDynamoTablesTest, describeAzsTestUnauthorized, listS3BucketsTest}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{listDynamoTablesTest, describeAzsTestUnauthorized, listS3BucketsTest}),
		},
		Teardown: func() error {
			if err := deleteTestGroup(groupName); err != nil {
				return err
			}
			return deleteTestUser(accessKey, userName)
		},
	})
}

// TestAcceptanceBackend_SessionToken
func TestAcceptanceBackend_SessionToken(t *testing.T) {
	t.Parallel()
	userName := generateUniqueUserName(t.Name())
	accessKey := &awsAccessKey{}

	roleData := map[string]interface{}{
		"credential_type": sessionTokenCred,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createUser(t, userName, accessKey)
			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfigWithCreds(t, accessKey),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{listDynamoTablesTest}),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{listDynamoTablesTest}),
		},
		Teardown: func() error {
			return deleteTestUser(accessKey, userName)
		},
	})
}

// Running this test requires a pre-made IAM user that has the necessary access permissions set
// and a set MFA device. This device serial number along with the other associated values must
// be set to the environment variables in the function below.
// For this reason, the test is currently a manually run-only acceptance test.
func TestAcceptanceBackend_SessionTokenWithMFA(t *testing.T) {
	t.Parallel()

	serial, found := os.LookupEnv("AWS_TEST_MFA_SERIAL_NUMBER")
	if !found {
		t.Skipf("AWS_TEST_MFA_SERIAL_NUMBER not set, skipping")
	}
	code, found := os.LookupEnv("AWS_TEST_MFA_CODE")
	if !found {
		t.Skipf("AWS_TEST_MFA_CODE not set, skipping")
	}
	accessKeyID, found := os.LookupEnv("AWS_TEST_MFA_USER_ACCESS_KEY")
	if !found {
		t.Skipf("AWS_TEST_MFA_USER_ACCESS_KEY not set, skipping")
	}
	secretKey, found := os.LookupEnv("AWS_TEST_MFA_USER_SECRET_KEY")
	if !found {
		t.Skipf("AWS_TEST_MFA_USER_SECRET_KEY not set, skipping")
	}

	accessKey := &awsAccessKey{}
	accessKey.AccessKeyID = accessKeyID
	accessKey.SecretAccessKey = secretKey

	roleData := map[string]interface{}{
		"credential_type":   sessionTokenCred,
		"mfa_serial_number": serial,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfigWithCreds(t, accessKey),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepReadWithMFA(t, "sts", "test", code, []credentialTestFunc{listDynamoTablesTest}),
			testAccStepReadWithMFA(t, "creds", "test", code, []credentialTestFunc{listDynamoTablesTest}),
		},
	})
}

func TestAcceptanceBackend_RoleDefaultSTSTTL(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueRoleName(t.Name())
	minAwsAssumeRoleDuration := 900
	awsAccountID, err := getAccountID(t.Context())
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}
	roleData := map[string]interface{}{
		"role_arns":       []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, roleName)},
		"credential_type": assumedRoleCred,
		"default_sts_ttl": minAwsAssumeRoleDuration,
		"max_sts_ttl":     minAwsAssumeRoleDuration,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createRole(t, roleName, awsAccountID, []string{ec2PolicyArn}, nil)
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepReadSTSResponse("test", time.Duration(minAwsAssumeRoleDuration)*time.Second), // allow a little slack
		},
		Teardown: func() error {
			return deleteTestRole(roleName)
		},
	})
}

// TestBackend_policyArnCRUD test the CRUD operations for policy ARNs.
func TestBackend_policyArnCRUD(t *testing.T) {
	t.Parallel()
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteArnPolicyRef(t, "test", ec2PolicyArn),
			testAccStepReadArnPolicy(t, "test", ec2PolicyArn),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadArnPolicy(t, "test", ""),
		},
	})
}

func testAccStepReadArnPolicy(t *testing.T, name string, value string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if value == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			expected := map[string]interface{}{
				"policy_arns":              []string{value},
				"role_arns":                []string(nil),
				"policy_document":          "",
				"credential_type":          iamUserCred,
				"default_sts_ttl":          int64(0),
				"max_sts_ttl":              int64(0),
				"user_path":                "",
				"permissions_boundary_arn": "",
				"iam_groups":               []string(nil),
				"iam_tags":                 map[string]string(nil),
				"mfa_serial_number":        "",
				"session_tags":             map[string]string(nil),
				"external_id":              "",
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("bad: got: %#v\nexpected: %#v", resp.Data, expected)
			}

			return nil
		},
	}
}

func testAccStepWriteArnRoleRef(t *testing.T, vaultRoleName, awsRoleName, awsAccountID string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + vaultRoleName,
		Data: map[string]interface{}{
			"arn": fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, awsRoleName),
		},
	}
}

// TestBackend_iamGroupsCRUD tests CRUD operations for IAM groups.
func TestBackend_iamGroupsCRUD(t *testing.T) {
	t.Parallel()
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteIamGroups(t, "test", []string{"group1", "group2"}),
			testAccStepReadIamGroups(t, "test", []string{"group1", "group2"}),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadIamGroups(t, "test", []string{}),
		},
	})
}

func testAccStepWriteIamGroups(t *testing.T, name string, groups []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"credential_type": iamUserCred,
			"iam_groups":      groups,
		},
	}
}

func testAccStepReadIamGroups(t *testing.T, name string, groups []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if len(groups) == 0 {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			expected := map[string]interface{}{
				"policy_arns":              []string(nil),
				"role_arns":                []string(nil),
				"policy_document":          "",
				"credential_type":          iamUserCred,
				"default_sts_ttl":          int64(0),
				"max_sts_ttl":              int64(0),
				"user_path":                "",
				"permissions_boundary_arn": "",
				"iam_groups":               groups,
				"iam_tags":                 map[string]string(nil),
				"mfa_serial_number":        "",
				"session_tags":             map[string]string(nil),
				"external_id":              "",
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("bad: got: %#v\nexpected: %#v", resp.Data, expected)
			}

			return nil
		},
	}
}

// TestBackend_iamTagsCRUD tests the CRUD operations for IAM tags.
func TestBackend_iamTagsCRUD(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteIamTags(t, "test", map[string]string{"key1": "value1", "key2": "value2"}),
			testAccStepReadIamTags(t, "test", map[string]string{"key1": "value1", "key2": "value2"}),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadIamTags(t, "test", map[string]string{}),
		},
	})
}

func testAccStepWriteIamTags(t *testing.T, name string, tags map[string]string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"credential_type": iamUserCred,
			"iam_tags":        tags,
		},
	}
}

func testAccStepReadIamTags(t *testing.T, name string, tags map[string]string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if len(tags) == 0 {
					return nil
				}

				return fmt.Errorf("vault response not received")
			}

			expected := map[string]interface{}{
				"policy_arns":              []string(nil),
				"role_arns":                []string(nil),
				"policy_document":          "",
				"credential_type":          iamUserCred,
				"default_sts_ttl":          int64(0),
				"max_sts_ttl":              int64(0),
				"user_path":                "",
				"permissions_boundary_arn": "",
				"iam_groups":               []string(nil),
				"iam_tags":                 tags,
				"mfa_serial_number":        "",
				"session_tags":             map[string]string(nil),
				"external_id":              "",
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("bad: got: %#v\nexpected: %#v", resp.Data, expected)
			}

			return nil
		},
	}
}

// TestBackend_stsSessionTagsCRUD tests the CRUD operations for STS session tags.
func TestBackend_stsSessionTagsCRUD(t *testing.T) {
	t.Parallel()

	tagParams0 := map[string]string{"tag1": "value1", "tag2": "value2"}
	tagParams1 := map[string]string{"tag1": "value1", "tag2": "value4", "tag3": "value3"}

	// list of tags in the form of "key=value"
	tagParamsList0 := []string{"key1=value1", "key2=value2"}
	tagParamsList0Expect := map[string]string{"key1": "value1", "key2": "value2"}
	tagParamsList1 := []string{"key1=value2", "key3=value4"}
	tagParamsList1Expect := map[string]string{"key1": "value2", "key3": "value4"}

	type testCase struct {
		name        string
		expectTags  []map[string]string
		tagsParams  []any
		externalIDs []string
	}

	for _, tt := range []testCase{
		{
			name: "mapped-only",
			tagsParams: []any{
				tagParams0,
				map[string]string{},
				tagParams1,
			},
			expectTags: []map[string]string{
				tagParams0,
				{},
				tagParams1,
			},
			externalIDs: []string{"foo", "", "bar"},
		},
		{
			name: "string-list-only",
			tagsParams: []any{
				tagParamsList0,
				tagParamsList1,
			},
			expectTags: []map[string]string{
				tagParamsList0Expect,
				tagParamsList1Expect,
			},
			externalIDs: []string{"foo"},
		},
		{
			name: "mixed-param-types",
			tagsParams: []any{
				tagParams0,
				tagParamsList0,
				tagParams1,
				tagParamsList1,
			},
			expectTags: []map[string]string{
				tagParams0,
				tagParamsList0Expect,
				tagParams1,
				tagParamsList1Expect,
			},
			externalIDs: []string{"foo", "bar"},
		},
		{
			name: "unset-tags",
			tagsParams: []any{
				tagParams0,
				map[string]string{},
			},
			expectTags: []map[string]string{
				tagParams0,
				{},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			steps := []logicaltest.TestStep{
				testAccStepConfig(t),
			}

			if len(tt.tagsParams) != len(tt.expectTags) {
				t.Fatalf("invalid test case: test case params and expect must have the same length")
			}

			// lastNonEmptyExternalID is used to store the last non-empty external ID for the
			// test case. The value will is expected to be set on the role. Setting the value
			// to an empty string has no effect on update operations.
			var lastNonEmptyExternalID string
			for idx, params := range tt.tagsParams {
				var externalID string
				if len(tt.externalIDs) > idx {
					externalID = tt.externalIDs[idx]
				}
				if externalID != "" {
					lastNonEmptyExternalID = externalID
				}
				steps = append(steps, testAccStepWriteSTSSessionTags(t, tt.name, params, externalID))
				steps = append(steps, testAccStepReadSTSSessionTags(t, tt.name, tt.expectTags[idx], lastNonEmptyExternalID, false))
			}
			steps = append(
				steps,
				testAccStepDeletePolicy(t, tt.name),
				testAccStepReadSTSSessionTags(t, tt.name, nil, "", true),
			)
			logicaltest.Test(t, logicaltest.TestCase{
				AcceptanceTest: false,
				LogicalBackend: getBackend(t),
				Steps:          steps,
			})
		})
	}
}

func testAccStepWriteSTSSessionTags(t *testing.T, name string, tags any, externalID string) logicaltest.TestStep {
	t.Helper()

	data := map[string]interface{}{
		"credential_type": assumedRoleCred,
		"session_tags":    tags,
	}
	if externalID != "" {
		data["external_id"] = externalID
	}
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data:      data,
	}
}

func testAccStepReadSTSSessionTags(t *testing.T, name string, tags any, externalID string, expectNilResp bool) logicaltest.TestStep {
	t.Helper()

	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if expectNilResp {
					return nil
				}
				return fmt.Errorf("vault response not received")
			}

			expected := map[string]interface{}{
				"policy_arns":              []string(nil),
				"role_arns":                []string(nil),
				"policy_document":          "",
				"credential_type":          assumedRoleCred,
				"default_sts_ttl":          int64(0),
				"max_sts_ttl":              int64(0),
				"user_path":                "",
				"permissions_boundary_arn": "",
				"iam_groups":               []string(nil),
				"iam_tags":                 map[string]string(nil),
				"mfa_serial_number":        "",
				"session_tags":             tags,
				"external_id":              externalID,
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("bad: got: %#v\nexpected: %#v", resp.Data, expected)
			}

			return nil
		},
	}
}

func generateUniqueRoleName(prefix string) string {
	return generateUniqueName(prefix, 64)
}

func generateUniqueUserName(prefix string) string {
	return generateUniqueName(prefix, 64)
}

func generateUniqueGroupName(prefix string) string {
	return generateUniqueName(prefix, 128)
}

func generateUniqueName(prefix string, maxLength int) string {
	name := testhelpers.RandomWithPrefix(prefix)
	if len(name) > maxLength {
		return name[:maxLength]
	}
	return name
}

type awsAccessKey struct {
	AccessKeyID     string
	SecretAccessKey string
}

type credentialTestFunc func(string, string, string) error

// TestAcceptanceBackend_MaxRetriesZero verifies that max_retries=0 is applied to the
// AWS SDK config rather than silently ignored. Before the SDK v2 migration fix, the
// condition was `> 0`, meaning 0 was treated as "unset" and the SDK applied its default
// retry count. The fix uses `>= 0` so 0 is passed through as an explicit zero.
//
// Required env vars: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY (or any credentials
// resolvable by the default AWS credential chain), AWS_DEFAULT_REGION.
func TestAcceptanceBackend_MaxRetriesZero(t *testing.T) {
	t.Parallel()
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			{
				Operation: logical.UpdateOperation,
				Path:      "config/root",
				Data: map[string]interface{}{
					"region":      os.Getenv("AWS_DEFAULT_REGION"),
					"max_retries": 0,
				},
			},
			{
				Operation: logical.ReadOperation,
				Path:      "config/root",
				Check: func(resp *logical.Response) error {
					if resp == nil {
						return fmt.Errorf("expected response but got nil")
					}
					v, ok := resp.Data["max_retries"]
					if !ok {
						return fmt.Errorf("max_retries not present in response")
					}
					if v.(int) != 0 {
						return fmt.Errorf("expected max_retries=0, got %v", v)
					}
					return nil
				},
			},
		},
	})
}

// TestAcceptanceBackend_WIF_IAMCredentials verifies that the AWS secrets engine
// can assume a role via Workload Identity Federation and issue valid IAM user
// credentials. This exercises the live-refresh token path in getRootIAMConfig
// fixed during the SDK v2 migration.
//
// Required env vars:
//   - AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY — base credentials with sts:AssumeRoleWithWebIdentity
//   - AWS_DEFAULT_REGION
//   - AWS_WIF_ROLE_ARN — ARN of the IAM role to assume via WIF
//   - AWS_WIF_TOKEN_AUDIENCE — audience claim expected on the identity token
func TestAcceptanceBackend_WIF_IAMCredentials(t *testing.T) {
	t.Parallel()

	roleARN := os.Getenv("AWS_WIF_ROLE_ARN")
	audience := os.Getenv("AWS_WIF_TOKEN_AUDIENCE")

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			if roleARN == "" {
				t.Skip("AWS_WIF_ROLE_ARN not set, skipping WIF acceptance test")
			}
			if audience == "" {
				t.Skip("AWS_WIF_TOKEN_AUDIENCE not set, skipping WIF acceptance test")
			}
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			{
				Operation: logical.UpdateOperation,
				Path:      "config/root",
				Data: map[string]interface{}{
					"region":                  os.Getenv("AWS_DEFAULT_REGION"),
					"role_arn":                roleARN,
					"identity_token_audience": audience,
					"identity_token_ttl":      3600,
				},
			},
			testAccStepWritePolicy(t, "wif-test", testDynamoPolicy),
			testAccStepRead(t, "creds", "wif-test", []credentialTestFunc{listDynamoTablesTest}),
		},
	})
}

// TestAcceptanceBackend_WIF_STSCredentials verifies that the AWS secrets engine
// can assume a role via Workload Identity Federation and issue valid STS credentials.
// This exercises the live-refresh token path in getRootSTSConfigs fixed during the
// SDK v2 migration.
//
// Required env vars:
//   - AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY — base credentials with sts:AssumeRoleWithWebIdentity
//   - AWS_DEFAULT_REGION
//   - AWS_WIF_ROLE_ARN — ARN of the IAM role to assume via WIF
//   - AWS_WIF_TOKEN_AUDIENCE — audience claim expected on the identity token
//   - AWS_WIF_ASSUME_ROLE_ARN — ARN of the role to assume for the sts/ credential type
func TestAcceptanceBackend_WIF_STSCredentials(t *testing.T) {
	t.Parallel()

	roleARN := os.Getenv("AWS_WIF_ROLE_ARN")
	audience := os.Getenv("AWS_WIF_TOKEN_AUDIENCE")
	assumeRoleARN := os.Getenv("AWS_WIF_ASSUME_ROLE_ARN")

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			if roleARN == "" {
				t.Skip("AWS_WIF_ROLE_ARN not set, skipping WIF acceptance test")
			}
			if audience == "" {
				t.Skip("AWS_WIF_TOKEN_AUDIENCE not set, skipping WIF acceptance test")
			}
			if assumeRoleARN == "" {
				t.Skip("AWS_WIF_ASSUME_ROLE_ARN not set, skipping WIF STS acceptance test")
			}
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			{
				Operation: logical.UpdateOperation,
				Path:      "config/root",
				Data: map[string]interface{}{
					"region":                  os.Getenv("AWS_DEFAULT_REGION"),
					"role_arn":                roleARN,
					"identity_token_audience": audience,
					"identity_token_ttl":      3600,
				},
			},
			{
				Operation: logical.UpdateOperation,
				Path:      "roles/wif-sts-test",
				Data: map[string]interface{}{
					"credential_type": assumedRoleCred,
					"role_arns":       []string{assumeRoleARN},
				},
			},
			testAccStepRead(t, "sts", "wif-sts-test", []credentialTestFunc{describeInstancesTest}),
		},
	})
}
