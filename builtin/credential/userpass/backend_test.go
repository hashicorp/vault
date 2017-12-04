package userpass

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

const (
	testSysTTL    = time.Hour * 10
	testSysMaxTTL = time.Hour * 20
)

func TestBackend_TTLDurations(t *testing.T) {
	data1 := map[string]interface{}{
		"password": "password",
		"policies": "root",
		"ttl":      "21h",
		"max_ttl":  "11h",
	}
	data2 := map[string]interface{}{
		"password": "password",
		"policies": "root",
		"ttl":      "10h",
		"max_ttl":  "21h",
	}
	data3 := map[string]interface{}{
		"password": "password",
		"policies": "root",
		"ttl":      "10h",
		"max_ttl":  "10h",
	}
	data4 := map[string]interface{}{
		"password": "password",
		"policies": "root",
		"ttl":      "11h",
		"max_ttl":  "5h",
	}
	data5 := map[string]interface{}{
		"password": "password",
	}
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testUsersWrite(t, "test", data1, true),
			testUsersWrite(t, "test", data2, true),
			testUsersWrite(t, "test", data3, false),
			testUsersWrite(t, "test", data4, false),
			testLoginWrite(t, "test", data5, false),
			testLoginWrite(t, "wrong", data5, true),
		},
	})
}

func TestBackend_basic(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepUser(t, "web2", "password", "foo"),
			testAccStepUser(t, "web3", "password", "foo"),
			testAccStepList(t, []string{"web", "web2", "web3"}),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
		},
	})
}

func TestBackend_userCrud(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

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

func TestBackend_userCreateOperation(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testUserCreateOperation(t, "web", "password", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
		},
	})
}

func TestBackend_passwordUpdate(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
			testUpdatePassword(t, "web", "newpassword"),
			testAccStepLogin(t, "web", "newpassword", []string{"default", "foo"}),
		},
	})

}

func TestBackend_policiesUpdate(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
			testUpdatePolicies(t, "web", "foo,bar"),
			testAccStepReadUser(t, "web", "bar,foo"),
			testAccStepLogin(t, "web", "password", []string{"bar", "default", "foo"}),
		},
	})

}

func testUpdatePassword(t *testing.T, user, password string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user + "/password",
		Data: map[string]interface{}{
			"password": password,
		},
	}
}

func testUpdatePolicies(t *testing.T, user, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user + "/policies",
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testUsersWrite(t *testing.T, user string, data map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user,
		Data:      data,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp == nil && expectError {
				return fmt.Errorf("Expected error but received nil")
			}
			return nil
		},
	}
}

func testLoginWrite(t *testing.T, user string, data map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data:      data,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp == nil && expectError {
				return fmt.Errorf("Expected error but received nil")
			}
			return nil
		},
	}
}

func testAccStepList(t *testing.T, users []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "users",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("Got error response: %#v", *resp)
			}

			exp := []string{"web", "web2", "web3"}
			if !reflect.DeepEqual(exp, resp.Data["keys"].([]string)) {
				return fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", exp, resp.Data["keys"])
			}
			return nil
		},
	}
}

func testAccStepLogin(t *testing.T, user string, pass string, policies []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": pass,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(policies),
	}
}

func testUserCreateOperation(
	t *testing.T, name string, password string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"password": password,
			"policies": policies,
		},
	}
}

func testAccStepUser(
	t *testing.T, name string, password string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
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
