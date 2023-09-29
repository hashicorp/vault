// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"os"
	"sync"
	"testing"

	stepwise "github.com/hashicorp/vault-testing-stepwise"
	dockerEnvironment "github.com/hashicorp/vault-testing-stepwise/environments/docker"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

var stepwiseSetup sync.Once

func TestAccBackend_Stepwise_basic(t *testing.T) {
	t.Parallel()
	envOptions := &stepwise.MountOptions{
		RegistryName:    "aws-sec",
		PluginType:      api.PluginTypeSecrets,
		PluginName:      "aws",
		MountPathPrefix: "aws-sec",
	}
	roleName := "vault-stepwise-role"
	stepwise.Run(t, stepwise.Case{
		Precheck:    func() { testAccStepwisePreCheck(t) },
		Environment: dockerEnvironment.NewEnvironment("aws", envOptions),
		Steps: []stepwise.Step{
			testAccStepwiseConfig(t),
			testAccStepwiseWritePolicy(t, roleName, testDynamoPolicy),
			testAccStepwiseRead(t, "creds", roleName, []credentialTestFunc{listDynamoTablesTest}),
		},
	})
}

func testAccStepwiseConfig(t *testing.T) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "config/root",
		Data: map[string]interface{}{
			"region":     os.Getenv("AWS_DEFAULT_REGION"),
			"access_key": os.Getenv("TEST_AWS_ACCESS_KEY"),
			"secret_key": os.Getenv("TEST_AWS_SECRET_KEY"),
		},
	}
}

func testAccStepwiseWritePolicy(t *testing.T, name string, policy string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"policy_document": policy,
			"credential_type": "iam_user",
		},
	}
}

func testAccStepwiseRead(t *testing.T, path, name string, credentialTests []credentialTestFunc) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ReadOperation,
		Path:      path + "/" + name,
		Assert: func(resp *api.Secret, err error) error {
			if err != nil {
				return err
			}
			var d struct {
				AccessKey string `mapstructure:"access_key"`
				SecretKey string `mapstructure:"secret_key"`
				STSToken  string `mapstructure:"security_token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			t.Logf("[WARN] Generated credentials: %v", d)
			for _, testFunc := range credentialTests {
				err := testFunc(d.AccessKey, d.SecretKey, d.STSToken)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func testAccStepwisePreCheck(t *testing.T) {
	stepwiseSetup.Do(func() {
		if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
			t.Logf("[INFO] Test: Using us-west-2 as test region")
			os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
		}

		// Ensure test variables are set
		if v := os.Getenv("TEST_AWS_ACCESS_KEY"); v == "" {
			t.Skip("TEST_AWS_ACCESS_KEY not set")
		}
		if v := os.Getenv("TEST_AWS_SECRET_KEY"); v == "" {
			t.Skip("TEST_AWS_SECRET_KEY not set")
		}
	})
}
