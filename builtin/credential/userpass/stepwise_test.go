package userpass

import (
	"fmt"
	"reflect"
	"testing"

	stepwise "github.com/hashicorp/vault-testing-stepwise"
	dockerEnvironment "github.com/hashicorp/vault-testing-stepwise/environments/docker"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/mitchellh/mapstructure"
)

func TestAccBackend_stepwise_UserCrud(t *testing.T) {
	customPluginName := "my-userpass"
	envOptions := &stepwise.MountOptions{
		RegistryName:    customPluginName,
		PluginType:      api.PluginTypeCredential,
		PluginName:      "userpass",
		MountPathPrefix: customPluginName,
	}
	stepwise.Run(t, stepwise.Case{
		Environment: dockerEnvironment.NewEnvironment(customPluginName, envOptions),
		Steps: []stepwise.Step{
			testAccStepwiseUser(t, "web", "password", "foo"),
			testAccStepwiseReadUser(t, "web", "foo"),
			testAccStepwiseDeleteUser(t, "web"),
			testAccStepwiseReadUser(t, "web", ""),
		},
	})
}

func testAccStepwiseUser(
	t *testing.T, name string, password string, policies string,
) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"password": password,
			"policies": policies,
		},
	}
}

func testAccStepwiseDeleteUser(t *testing.T, name string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.DeleteOperation,
		Path:      "users/" + name,
	}
}

func testAccStepwiseReadUser(t *testing.T, name string, policies string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ReadOperation,
		Path:      "users/" + name,
		Assert: func(resp *api.Secret, err error) error {
			if resp == nil {
				if policies == "" {
					return nil
				}

				return fmt.Errorf("unexpected nil response")
			}

			var d struct {
				Policies []string `mapstructure:"policies"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			expectedPolicies := policyutil.ParsePolicies(policies)
			if !reflect.DeepEqual(d.Policies, expectedPolicies) {
				return fmt.Errorf("Actual policies: %#v\nExpected policies: %#v", d.Policies, expectedPolicies)
			}

			return nil
		},
	}
}
