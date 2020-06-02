package aws

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/vault/sdk/testing/stepwise"
	dockerDriver "github.com/hashicorp/vault/sdk/testing/stepwise/drivers/docker"
)

// TEST_AWS_SECRET_KEY
// TEST_AWS_ACCESS_KEY

func TestAccBackend_Stepwise_basic(t *testing.T) {
	t.Parallel()
	driverOptions := &stepwise.DriverOptions{
		Name:       "aws-sec",
		PluginType: stepwise.PluginTypeSecrets,
		PluginName: "aws",
		MountPath:  "aws-sec",
	}
	roleName := "vault-stepwise-role"
	stepwise.Run(t, stepwise.Case{
		PreCheck: func() { testAccStepwisePreCheck(t) },
		Driver:   dockerDriver.NewDockerDriver("aws", driverOptions),
		Steps: []stepwise.Step{
			testAccStepwiseConfig(t),
			testAccStepwiseWritePolicy(t, roleName, testDynamoPolicy),
			testAccStepwiseRead(t, "creds", roleName, []credentialTestFunc{listDynamoTablesTest}),
		},
	})
}

func testAccStepwiseConfig(t *testing.T) stepwise.Step {
	return stepwise.Step{
		Operation: logical.UpdateOperation,
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
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"policy_document": policy,
			"credential_type": "iam_user",
		},
	}
}

func testAccStepwiseRead(t *testing.T, path, name string, credentialTests []credentialTestFunc) stepwise.Step {
	return stepwise.Step{
		Operation: logical.ReadOperation,
		Path:      path + "/" + name,
		Check: func(resp *api.Secret, err error) error {
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

func testAccStepwisePreCheck(t *testing.T) {
	initSetup.Do(func() {
		if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
			log.Println("[INFO] Test: Using us-west-2 as test region")
			os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
		}

		// TEST_AWS_SECRET_KEY
		// TEST_AWS_ACCESS_KEY
		// Ensure test variables are set
		if v := os.Getenv("TEST_AWS_SECRET_KEY"); v == "" {
			t.Fatal("TEST_AWS_SECRET_KEY not set")
		}
	})
}
