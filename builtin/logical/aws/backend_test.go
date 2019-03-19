package aws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

var initSetup sync.Once

type mockIAMClient struct {
	iamiface.IAMAPI
}

func (m *mockIAMClient) CreateUser(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
	return nil, awserr.New("Throttling", "", nil)
}

func getBackend(t *testing.T) logical.Backend {
	be, _ := Factory(context.Background(), logical.TestBackendConfig())
	return be
}

func TestBackend_basic(t *testing.T) {
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

func TestBackend_basicSTS(t *testing.T) {
	t.Parallel()
	awsAccountID, err := getAccountID()
	if err != nil {
		t.Logf("Unable to retrive user via sts:GetCallerIdentity: %#v", err)
		t.Skip("Could not determine AWS account ID from sts:GetCallerIdentity for acceptance tests, skipping")
	}
	roleName := generateUniqueName(t.Name())
	userName := generateUniqueName(t.Name())
	accessKey := &awsAccessKey{}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createUser(t, userName, accessKey)
			createRole(t, roleName, awsAccountID)
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
			return teardown(accessKey, roleName, userName)
		},
	})
}

func TestBackend_policyCrud(t *testing.T) {
	t.Parallel()
	compacted, err := compactJSON(testDynamoPolicy)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
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

	b := Backend()
	if err := b.Setup(context.Background(), config); err != nil {
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

	resp, err := b.HandleRequest(context.Background(), confReq)
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

	credResp, err := b.HandleRequest(context.Background(), credReq)
	if err == nil {
		t.Fatalf("failed to trigger expected throttling error condition: resp:%#v", credResp)
	}
	rErr := credResp.Error()
	expected := "Error creating IAM user: Throttling: "
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
	initSetup.Do(func() {
		if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
			log.Println("[INFO] Test: Using us-west-2 as test region")
			os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
		}
	})
}

func getAccountID() (string, error) {
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := sts.New(session.New(awsConfig))

	params := &sts.GetCallerIdentityInput{}
	res, err := svc.GetCallerIdentity(params)

	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("got nil response from GetCallerIdentity")
	}

	return *res.Account, nil
}

func createRole(t *testing.T, roleName, awsAccountID string) {
	const testRoleAssumePolicy = `{
      "Version": "2012-10-17",
      "Statement": [
          {
              "Effect":"Allow",
              "Principal": {
                  "AWS": "arn:aws:iam::%s:root"
              },
              "Action": "sts:AssumeRole"
           }
      ]
}
`
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := iam.New(session.New(awsConfig))
	trustPolicy := fmt.Sprintf(testRoleAssumePolicy, awsAccountID)

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		RoleName:                 aws.String(roleName),
		Path:                     aws.String("/"),
	}

	log.Printf("[INFO] AWS CreateRole: %s", roleName)
	_, err := svc.CreateRole(params)

	if err != nil {
		t.Fatalf("AWS CreateRole failed: %v", err)
	}

	attachment := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(ec2PolicyArn),
		RoleName:  aws.String(roleName), // Required
	}
	_, err = svc.AttachRolePolicy(attachment)

	if err != nil {
		t.Fatalf("AWS CreateRole failed: %v", err)
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
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := iam.New(session.New(awsConfig))

	createUserInput := &iam.CreateUserInput{
		UserName: aws.String(userName),
	}
	log.Printf("[INFO] AWS CreateUser: %s", userName)
	_, err := svc.CreateUser(createUserInput)
	if err != nil {
		t.Fatalf("AWS CreateUser failed: %v", err)
	}

	putPolicyInput := &iam.PutUserPolicyInput{
		PolicyDocument: aws.String(timebombPolicy),
		PolicyName:     aws.String("SelfDestructionTimebomb"),
		UserName:       aws.String(userName),
	}
	_, err = svc.PutUserPolicy(putPolicyInput)
	if err != nil {
		t.Fatalf("AWS PutUserPolicy failed: %v", err)
	}

	attachUserPolicyInput := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		UserName:  aws.String(userName),
	}
	_, err = svc.AttachUserPolicy(attachUserPolicyInput)
	if err != nil {
		t.Fatalf("AWS AttachUserPolicy failed, %v", err)
	}

	createAccessKeyInput := &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	}
	createAccessKeyOutput, err := svc.CreateAccessKey(createAccessKeyInput)
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

func deleteTestRole(roleName string) error {
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := iam.New(session.New(awsConfig))

	attachment := &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(ec2PolicyArn),
		RoleName:  aws.String(roleName), // Required
	}
	_, err := svc.DetachRolePolicy(attachment)
	if err != nil {
		log.Printf("[WARN] AWS DetachRolePolicy failed: %v", err)
		return err
	}

	params := &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	}

	log.Printf("[INFO] AWS DeleteRole: %s", roleName)
	_, err = svc.DeleteRole(params)

	if err != nil {
		log.Printf("[WARN] AWS DeleteRole failed: %v", err)
		return err
	}
	return nil
}

func teardown(accessKey *awsAccessKey, roleName, userName string) error {

	if err := deleteTestRole(roleName); err != nil {
		return err
	}
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := iam.New(session.New(awsConfig))

	userDetachment := &iam.DetachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		UserName:  aws.String(userName),
	}
	_, err := svc.DetachUserPolicy(userDetachment)
	if err != nil {
		log.Printf("[WARN] AWS DetachUserPolicy failed: %v", err)
		return err
	}

	deleteAccessKeyInput := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accessKey.AccessKeyID),
		UserName:    aws.String(userName),
	}
	_, err = svc.DeleteAccessKey(deleteAccessKeyInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteAccessKey failed: %v", err)
		return err
	}

	deleteUserPolicyInput := &iam.DeleteUserPolicyInput{
		PolicyName: aws.String("SelfDestructionTimebomb"),
		UserName:   aws.String(userName),
	}
	_, err = svc.DeleteUserPolicy(deleteUserPolicyInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteUserPolicy failed: %v", err)
		return err
	}
	deleteUserInput := &iam.DeleteUserInput{
		UserName: aws.String(userName),
	}
	log.Printf("[INFO] AWS DeleteUser: %s", userName)
	_, err = svc.DeleteUser(deleteUserInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteUser failed: %v", err)
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
			awsConfig := &aws.Config{
				Region:      aws.String("us-east-1"),
				HTTPClient:  cleanhttp.DefaultClient(),
				Credentials: credentials.NewStaticCredentials(oldAccessKey.AccessKeyID, oldAccessKey.SecretAccessKey, ""),
			}
			// sigh....
			oldAccessKey.AccessKeyID = newAccessKeyID
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
			svc := sts.New(session.New(awsConfig))
			params := &sts.GetCallerIdentityInput{}
			_, err := svc.GetCallerIdentity(params)
			if err == nil {
				return fmt.Errorf("bad: old credentials succeeded after rotate")
			}
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() != "InvalidClientTokenId" {
					return fmt.Errorf("Unknown error returned from AWS: %#v", aerr)
				}
				return nil
			}
			return err
		},
	}
}

func testAccStepRead(t *testing.T, path, name string, credentialTests []credentialTestFunc) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path + "/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				AccessKey string `mapstructure:"access_key"`
				SecretKey string `mapstructure:"secret_key"`
				STSToken  string `mapstructure:"security_token"`
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

func testAccStepReadTTL(name string, maximumTTL time.Duration) logicaltest.TestStep {
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
	creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
	}
	client := ec2.New(session.New(awsConfig))
	log.Printf("[WARN] Verifying that the generated credentials work with ec2:DescribeInstances...")
	return retryUntilSuccess(func() error {
		_, err := client.DescribeInstances(&ec2.DescribeInstancesInput{})
		return err
	})
}

func describeAzsTestUnauthorized(accessKey, secretKey, token string) error {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
	}
	client := ec2.New(session.New(awsConfig))
	log.Printf("[WARN] Verifying that the generated credentials don't work with ec2:DescribeAvailabilityZones...")
	return retryUntilSuccess(func() error {
		_, err := client.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
		// Need to make sure AWS authenticates the generated credentials but does not authorize the operation
		if err == nil {
			return fmt.Errorf("operation succeeded when expected failure")
		}
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "UnauthorizedOperation" {
				return nil
			}
		}
		return err
	})
}

func listIamUsersTest(accessKey, secretKey, token string) error {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
	}
	client := iam.New(session.New(awsConfig))
	log.Printf("[WARN] Verifying that the generated credentials work with iam:ListUsers...")
	return retryUntilSuccess(func() error {
		_, err := client.ListUsers(&iam.ListUsersInput{})
		return err
	})
}

func listDynamoTablesTest(accessKey, secretKey, token string) error {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
	}
	client := dynamodb.New(session.New(awsConfig))
	log.Printf("[WARN] Verifying that the generated credentials work with dynamodb:ListTables...")
	return retryUntilSuccess(func() error {
		_, err := client.ListTables(&dynamodb.ListTablesInput{})
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
				"policy_arns":      []string(nil),
				"role_arns":        []string(nil),
				"policy_document":  value,
				"credential_types": []string{iamUserCred, federationTokenCred},
				"default_sts_ttl":  int64(0),
				"max_sts_ttl":      int64(0),
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

const ec2PolicyArn = "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"
const iamPolicyArn = "arn:aws:iam::aws:policy/IAMReadOnlyAccess"

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

func TestBackend_basicPolicyArnRef(t *testing.T) {
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

func TestBackend_iamUserManagedInlinePolicies(t *testing.T) {
	t.Parallel()
	compacted, err := compactJSON(testDynamoPolicy)
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	roleData := map[string]interface{}{
		"policy_document": testDynamoPolicy,
		"policy_arns":     []string{ec2PolicyArn, iamPolicyArn},
		"credential_type": iamUserCred,
	}
	expectedRoleData := map[string]interface{}{
		"policy_document":  compacted,
		"policy_arns":      []string{ec2PolicyArn, iamPolicyArn},
		"credential_types": []string{iamUserCred},
		"role_arns":        []string(nil),
		"default_sts_ttl":  int64(0),
		"max_sts_ttl":      int64(0),
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepReadRole(t, "test", expectedRoleData),
			testAccStepRead(t, "creds", "test", []credentialTestFunc{describeInstancesTest, listIamUsersTest, listDynamoTablesTest}),
			testAccStepRead(t, "sts", "test", []credentialTestFunc{describeInstancesTest, listIamUsersTest, listDynamoTablesTest}),
		},
	})
}

func TestBackend_AssumedRoleWithPolicyDoc(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueName(t.Name())
	// This looks a bit curious. The policy document and the role document act
	// as a logical intersection of policies. The role allows ec2:Describe*
	// (among other permissions). This policy allows everything BUT
	// ec2:DescribeAvailabilityZones. Thus, the logical intersection of the two
	// is all ec2:Describe* EXCEPT ec2:DescribeAvailabilityZones, and so the
	// describeAZs call should fail
	allowAllButDescribeAzs := `
{
	"Version": "2012-10-17",
	"Statement": [{
			"Effect": "Allow",
			"NotAction": "ec2:DescribeAvailabilityZones",
			"Resource": "*"
	}]
}
`
	awsAccountID, err := getAccountID()
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
			createRole(t, roleName, awsAccountID)
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

func TestBackend_RoleDefaultSTSTTL(t *testing.T) {
	t.Parallel()
	roleName := generateUniqueName(t.Name())
	minAwsAssumeRoleDuration := 900
	awsAccountID, err := getAccountID()
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
			createRole(t, roleName, awsAccountID)
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		LogicalBackend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteRole(t, "test", roleData),
			testAccStepReadTTL("test", time.Duration(minAwsAssumeRoleDuration)*time.Second), // allow a little slack
		},
		Teardown: func() error {
			return deleteTestRole(roleName)
		},
	})
}

func TestBackend_policyArnCrud(t *testing.T) {
	t.Parallel()
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
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
				"policy_arns":      []string{value},
				"role_arns":        []string(nil),
				"policy_document":  "",
				"credential_types": []string{iamUserCred},
				"default_sts_ttl":  int64(0),
				"max_sts_ttl":      int64(0),
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

func generateUniqueName(prefix string) string {
	return testhelpers.RandomWithPrefix(prefix)
}

type awsAccessKey struct {
	AccessKeyID     string
	SecretAccessKey string
}

type credentialTestFunc func(string, string, string) error
