package transit

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/testing/stepwise"
	dockerDriver "github.com/hashicorp/vault/sdk/testing/stepwise/drivers/docker"
	"github.com/mitchellh/mapstructure"
)

// TestBackend_basic_derived_docker is an example test using the Docker Driver
func TestBackend_basic_derived_docker(t *testing.T) {
	// decryptData := make(map[string]interface{})
	stepwise.Run(t, stepwise.Case{
		Driver: dockerDriver.NewDockerDriver("transit"),
		Steps: []stepwise.Step{
			testAccStepwiseListPolicy(t, "test", true),
			testAccStepwiseWritePolicy(t, "test", true),
			testAccStepwiseListPolicy(t, "test", false),
			// testAccStepReadPolicy(t, "test", false, true),
			// testAccStepEncryptContext(t, "test", testPlaintext, "my-cool-context", decryptData),
			// testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			// testAccStepEnableDeletion(t, "test"),
			// testAccStepDeletePolicy(t, "test"),
			// testAccStepReadPolicy(t, "test", true, true),
		},
	})
}

func testAccStepwiseWritePolicy(t *testing.T, name string, derived bool) stepwise.Step {
	ts := stepwise.Step{
		// Operation: logical.UpdateOperation,
		Operation: stepwise.WriteOperation,
		Path:      "keys/" + name,
		Data: map[string]interface{}{
			"derived": derived,
		},
	}
	if os.Getenv("TRANSIT_ACC_KEY_TYPE") == "CHACHA" {
		ts.Data["type"] = "chacha20-poly1305"
	}
	return ts
}

func testAccStepwiseListPolicy(t *testing.T, name string, expectNone bool) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ListOperation,
		Path:      "keys",
		Check: func(resp *api.Secret) error {
			if resp == nil {
				return fmt.Errorf("missing response")
			}
			if expectNone {
				keysRaw, ok := resp.Data["keys"]
				if ok || keysRaw != nil {
					return fmt.Errorf("response data when expecting none")
				}
				return nil
			}
			if len(resp.Data) == 0 {
				return fmt.Errorf("no data returned")
			}

			var d struct {
				Keys []string `mapstructure:"keys"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if len(d.Keys) > 0 && d.Keys[0] != name {
				return fmt.Errorf("bad name: %#v", d)
			}
			if len(d.Keys) != 1 {
				return fmt.Errorf("only 1 key expected, %d returned", len(d.Keys))
			}
			return nil
		},
	}
}
