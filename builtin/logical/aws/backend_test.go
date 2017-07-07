package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func getBackend(t *testing.T) logical.Backend {
	be, _ := Factory(logical.TestBackendConfig())
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
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
			createRole(t)
		},
		Backend: getBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWritePolicy(t, "test", testPolicy),
			testAccStepReadSTS(t, "test"),
			testAccStepWriteArnPolicyRef(t, "test", testPolicyArn),
			testAccStepReadSTSWithArnPolicy(t, "test"),
			testAccStepWriteArnRoleRef(t, testRoleName),
			testAccStepReadSTS(t, testRoleName),
		},
		Teardown: teardown,
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
	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}

	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
		log.Println("[INFO] Test: Using us-west-2 as test region")
		os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
	}

	if v := os.Getenv("AWS_ACCOUNT_ID"); v == "" {
		accountId, err := getAccountId()
		if err != nil {
			t.Fatal("AWS_ACCOUNT_ID could not be read from iam:GetUser for acceptance tests")
		}
		log.Printf("[INFO] Test: Used %s as AWS_ACCOUNT_ID", accountId)
		os.Setenv("AWS_ACCOUNT_ID", accountId)
	}
}

func getAccountId() (string, error) {
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"")

	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
	}
	svc := iam.New(session.New(awsConfig))

	params := &iam.GetUserInput{}
	res, err := svc.GetUser(params)

	if err != nil {
		return "", err
	}

	// split "arn:aws:iam::012345678912:user/username"
	accountId := strings.Split(*res.User.Arn, ":")[4]
	return accountId, nil
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
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")

	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
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

	// Sleep sometime because AWS is eventually consistent
	log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
	time.Sleep(10 * time.Second)
}

func teardown() error {
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")

	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
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

	return nil
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/root",
		Data: map[string]interface{}{
			"access_key": os.Getenv("AWS_ACCESS_KEY_ID"),
			"secret_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"region":     os.Getenv("AWS_DEFAULT_REGION"),
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

			// Sleep sometime because AWS is eventually consistent
			log.Println("[WARN] Sleeping for 10 seconds waiting for AWS...")
			time.Sleep(10 * time.Second)

			// Build a client and verify that the credentials work
			creds := credentials.NewStaticCredentials(d.AccessKey, d.SecretKey, "")
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
