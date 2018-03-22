package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func getBackend(t *testing.T) logical.Backend {
	be, _ := Factory(context.Background(), logical.TestBackendConfig())
	return be
}

func TestBackend_basic(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWritePolicy(t, "test", testPolicy),
			testAccStepReadUser(t, "test"),
		},
	})
}

func TestBackend_basicSTS(t *testing.T) {
	accessKey := &awsAccessKey{}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createUser(t, accessKey)
			createRole(t)
			// Sleep sometime because AWS is eventually consistent
			// Both the createUser and createRole depend on this
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)
		},
		Backend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfigWithCreds(t, accessKey),
			testAccStepWritePolicy(t, "test", testPolicy),
			testAccStepReadSTS(t, "test"),
			testAccStepWriteArnPolicyRef(t, "test", testPolicyArn),
			testAccStepReadSTSWithArnPolicy(t, "test"),
			testAccStepWriteArnRoleRef(t, testRoleName),
			testAccStepReadSTS(t, testRoleName),
		},
		Teardown: func() error {
			return teardown(accessKey)
		},
	})
}

func TestBackend_policyCrud(t *testing.T) {
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, []byte(testPolicy)); err != nil {
		t.Fatalf("bad: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWritePolicy(t, "test", testPolicy),
			testAccStepReadPolicy(t, "test", compacted.String()),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", ""),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
		log.Println("[INFO] Test: Using us-west-2 as test region")
		os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
	}

	if v := os.Getenv("AWS_ACCOUNT_ID"); v == "" {
		accountId, err := getAccountId()
		if err != nil {
			t.Fatalf("AWS_ACCOUNT_ID could not be read from iam:GetUser for acceptance tests: %#v", err)
		}
		log.Printf("[INFO] Test: Used %s as AWS_ACCOUNT_ID", accountId)
		os.Setenv("AWS_ACCOUNT_ID", accountId)
	}
}

func getAccountId() (string, error) {
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

const testRoleName = "Vault-Acceptance-Test-AWS-Assume-Role"

func createRole(t *testing.T) {
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
	trustPolicy := fmt.Sprintf(testRoleAssumePolicy, os.Getenv("AWS_ACCOUNT_ID"))

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		RoleName:                 aws.String(testRoleName),
		Path:                     aws.String("/"),
	}

	log.Printf("[INFO] AWS CreateRole: %s", testRoleName)
	_, err := svc.CreateRole(params)

	if err != nil {
		t.Fatalf("AWS CreateRole failed: %v", err)
	}

	attachment := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(testPolicyArn),
		RoleName:  aws.String(testRoleName), // Required
	}
	_, err = svc.AttachRolePolicy(attachment)

	if err != nil {
		t.Fatalf("AWS CreateRole failed: %v", err)
	}
}

const testUserName = "Vault-Acceptance-Test-AWS-FederationToken"

func createUser(t *testing.T, accessKey *awsAccessKey) {
	// The sequence of user creation actions is carefully chosen to minimize
	// impact of stolen IAM user credentials
	// 1. Create user, without any permissions or credentials. At this point,
	//	  nobody cares if creds compromised because this user can do nothing.
	// 2. Attach the timebomb policy. This grants no access but puts a time limit
	//	  on validitity of compromised credentials. If this fails, nobody cares
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
		UserName: aws.String(testUserName),
	}
	log.Printf("[INFO] AWS CreateUser: %s", testUserName)
	_, err := svc.CreateUser(createUserInput)
	if err != nil {
		t.Fatalf("AWS CreateUser failed: %v", err)
	}

	putPolicyInput := &iam.PutUserPolicyInput{
		PolicyDocument: aws.String(timebombPolicy),
		PolicyName:     aws.String("SelfDestructionTimebomb"),
		UserName:       aws.String(testUserName),
	}
	_, err = svc.PutUserPolicy(putPolicyInput)
	if err != nil {
		t.Fatalf("AWS PutUserPolicy failed: %v", err)
	}

	attachUserPolicyInput := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		UserName:  aws.String(testUserName),
	}
	_, err = svc.AttachUserPolicy(attachUserPolicyInput)
	if err != nil {
		t.Fatalf("AWS AttachUserPolicy failed, %v", err)
	}

	createAccessKeyInput := &iam.CreateAccessKeyInput{
		UserName: aws.String(testUserName),
	}
	createAccessKeyOutput, err := svc.CreateAccessKey(createAccessKeyInput)
	if err != nil {
		t.Fatalf("AWS CreateAccessKey failed: %v", err)
	}
	if createAccessKeyOutput == nil {
		t.Fatalf("AWS CreateAccessKey returned nil")
	}
	genAccessKey := createAccessKeyOutput.AccessKey

	accessKey.AccessKeyId = *genAccessKey.AccessKeyId
	accessKey.SecretAccessKey = *genAccessKey.SecretAccessKey
}

func teardown(accessKey *awsAccessKey) error {
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := iam.New(session.New(awsConfig))

	attachment := &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(testPolicyArn),
		RoleName:  aws.String(testRoleName), // Required
	}
	_, err := svc.DetachRolePolicy(attachment)
	if err != nil {
		log.Printf("[WARN] AWS DetachRolePolicy failed: %v", err)
		return err
	}

	params := &iam.DeleteRoleInput{
		RoleName: aws.String(testRoleName),
	}

	log.Printf("[INFO] AWS DeleteRole: %s", testRoleName)
	_, err = svc.DeleteRole(params)

	if err != nil {
		log.Printf("[WARN] AWS DeleteRole failed: %v", err)
		return err
	}

	userDetachment := &iam.DetachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		UserName:  aws.String(testUserName),
	}
	_, err = svc.DetachUserPolicy(userDetachment)
	if err != nil {
		log.Printf("[WARN] AWS DetachUserPolicy failed: %v", err)
		return err
	}

	deleteAccessKeyInput := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accessKey.AccessKeyId),
		UserName:    aws.String(testUserName),
	}
	_, err = svc.DeleteAccessKey(deleteAccessKeyInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteAccessKey failed: %v", err)
		return err
	}

	deleteUserPolicyInput := &iam.DeleteUserPolicyInput{
		PolicyName: aws.String("SelfDestructionTimebomb"),
		UserName:   aws.String(testUserName),
	}
	_, err = svc.DeleteUserPolicy(deleteUserPolicyInput)
	if err != nil {
		log.Printf("[WARN] AWS DeleteUserPolicy failed: %v", err)
		return err
	}
	deleteUserInput := &iam.DeleteUserInput{
		UserName: aws.String(testUserName),
	}
	log.Printf("[INFO] AWS DeleteUser: %s", testUserName)
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
			req.Data["access_key"] = accessKey.AccessKeyId
			req.Data["secret_key"] = accessKey.SecretAccessKey
			return nil
		},
	}
}

func testAccStepReadUser(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				AccessKey string `mapstructure:"access_key"`
				SecretKey string `mapstructure:"secret_key"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated credentials: %v", d)

			// Build a client and verify that the credentials work
			creds := credentials.NewStaticCredentials(d.AccessKey, d.SecretKey, "")
			awsConfig := &aws.Config{
				Credentials: creds,
				Region:      aws.String("us-east-1"),
				HTTPClient:  cleanhttp.DefaultClient(),
			}
			client := ec2.New(session.New(awsConfig))

			log.Printf("[WARN] Verifying that the generated credentials work...")
			retryCount := 0
			success := false
			var err error
			for !success && retryCount < 10 {
				_, err = client.DescribeInstances(&ec2.DescribeInstancesInput{})
				if err == nil {
					return nil
				}
				time.Sleep(time.Second)
				retryCount++
			}

			return err
		},
	}
}

func testAccStepReadSTS(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "sts/" + name,
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

			// Build a client and verify that the credentials work
			creds := credentials.NewStaticCredentials(d.AccessKey, d.SecretKey, d.STSToken)
			awsConfig := &aws.Config{
				Credentials: creds,
				Region:      aws.String("us-east-1"),
				HTTPClient:  cleanhttp.DefaultClient(),
			}
			client := ec2.New(session.New(awsConfig))

			log.Printf("[WARN] Verifying that the generated credentials work...")
			_, err := client.DescribeInstances(&ec2.DescribeInstancesInput{})
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func testAccStepReadSTSWithArnPolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "sts/" + name,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp.Data["error"] !=
				"Can't generate STS credentials for a managed policy; use a role to assume or an inline policy instead" {
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
			"policy": testPolicy,
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

			var d struct {
				Policy string `mapstructure:"policy"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Policy != value {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

const testPolicy = `
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1426528957000",
            "Effect": "Allow",
            "Action": [
                "ec2:*"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
`

const testPolicyArn = "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"

func testAccStepWriteArnPolicyRef(t *testing.T, name string, arn string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"arn": testPolicyArn,
		},
	}
}

func TestBackend_basicPolicyArnRef(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteArnPolicyRef(t, "test", testPolicyArn),
			testAccStepReadUser(t, "test"),
		},
	})
}

func TestBackend_policyArnCrud(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteArnPolicyRef(t, "test", testPolicyArn),
			testAccStepReadArnPolicy(t, "test", testPolicyArn),
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

			var d struct {
				Policy string `mapstructure:"arn"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Policy != value {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

func testAccStepWriteArnRoleRef(t *testing.T, roleName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + roleName,
		Data: map[string]interface{}{
			"arn": fmt.Sprintf("arn:aws:iam::%s:role/%s", os.Getenv("AWS_ACCOUNT_ID"), roleName),
		},
	}
}

type awsAccessKey struct {
	AccessKeyId     string
	SecretAccessKey string
}
