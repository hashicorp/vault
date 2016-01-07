package userpass

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_TTLDurations(t *testing.T) {
	sysTTL := time.Hour * 10
	sysMaxTTL := time.Hour * 20
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
			DefaultLeaseTTLVal: sysTTL,
			MaxLeaseTTLVal:     sysMaxTTL,
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

func testAccStepLogin(t *testing.T, user string, pass string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
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
