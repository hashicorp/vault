package userpass

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"crypto/tls"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

const (
	testSysTTL    = time.Hour * 10
	testSysMaxTTL = time.Hour * 20
)

func TestBackend_TTL(t *testing.T) {
	var resp *logical.Response
	var err error

	storage := &logical.InmemStorage{}

	config := logical.TestBackendConfig()
	config.StorageView = storage

	ctx := context.Background()

	b, err := Factory(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatalf("failed to create backend")
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"password": "testpassword",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp.Data["ttl"].(float64) != 0 && resp.Data["max_ttl"].(float64) != 0 {
		t.Fatalf("bad: ttl and max_ttl are not set correctly")
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"ttl":     "5m",
			"max_ttl": "10m",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp.Data["ttl"].(float64) != 300 && resp.Data["max_ttl"].(float64) != 600 {
		t.Fatalf("bad: ttl and max_ttl are not set correctly")
	}
}

func TestBackend_basic(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		CredentialBackend: b,
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
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepDeleteUser(t, "web"),
			testAccStepReadUser(t, "web", ""),
		},
	})
}

func TestBackend_userCreateOperation(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testUserCreateOperation(t, "web", "password", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
		},
	})
}

func TestBackend_passwordUpdate(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		CredentialBackend: b,
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
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		CredentialBackend: b,
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

func testAccStepList(t *testing.T, users []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "users",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("got error response: %#v", *resp)
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

		Check:     logicaltest.TestCheckAuth(policies),
		ConnState: &tls.ConnectionState{},
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
