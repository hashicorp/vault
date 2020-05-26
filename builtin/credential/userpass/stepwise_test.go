package userpass

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/mitchellh/mapstructure"
	"github.com/y0ssar1an/q"

	"github.com/hashicorp/vault/sdk/testing/stepwise"
	dockerDriver "github.com/hashicorp/vault/sdk/testing/stepwise/drivers/docker"
)

func TestBackend_stepwise_UserCrud(t *testing.T) {
	driverOptions := &stepwise.DriverOptions{
		Name:       "userpass23",
		PluginType: stepwise.PluginTypeCredential,
		PluginName: "userpass",
		MountPath:  "userpass23",
	}
	q.Q("testing:", stepwise.PluginTypeCredential.String())
	q.Q("do testing:", driverOptions.PluginType.String())
	stepwise.Run(t, stepwise.Case{
		Driver: dockerDriver.NewDockerDriver("userpass23", driverOptions),
		Steps: []stepwise.Step{
			testAccStepwiseUser(t, "web", "password", "foo"),
			testAccStepwiseReadUser(t, "web", "foo"),
			testAccStepwiseDeleteUser(t, "web"),
			testAccStepwiseReadUser(t, "web", ""),
		},
	})
}

func testAccStepwiseUser(
	t *testing.T, name string, password string, policies string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"password": password,
			"policies": policies,
		},
	}
}

func testAccStepwiseDeleteUser(t *testing.T, n string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.DeleteOperation,
		Path:      "users/" + n,
	}
}

func testAccStepwiseReadUser(t *testing.T, name string, policies string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ReadOperation,
		Path:      "users/" + name,
		Check: func(resp *api.Secret, err error) error {
			if resp == nil {
				if policies == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Policies []string `mapstructure:"policies"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if !reflect.DeepEqual(d.Policies, policyutil.ParsePolicies(policies)) {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}
