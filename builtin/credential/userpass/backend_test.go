package userpass

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepLogin(t, "web", "password"),
		},
	})
}

func TestBackend_userCrud(t *testing.T) {
	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepDeleteUser(t, "web"),
			testAccStepReadUser(t, "web", ""),
		},
	})
}

func testAccStepLogin(t *testing.T, user string, pass string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": pass,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth([]string{"foo"}),
	}
}

func testAccStepUser(
	t *testing.T, name string, password string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"password": password,
			"policies": policies,
		},
	}
}

func testAccStepDeleteUser(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "users/" + n,
	}
}

func testAccStepReadUser(t *testing.T, name string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "users/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if policies == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Policies string `mapstructure:"policies"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Policies != policies {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}
