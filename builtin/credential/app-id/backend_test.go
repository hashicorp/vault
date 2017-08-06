package appId

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_basic(t *testing.T) {
	var b *backend
	var err error
	var storage logical.Storage
	factory := func(conf *logical.BackendConfig) (logical.Backend, error) {
		b, err = Backend(conf)
		if err != nil {
			t.Fatal(err)
		}
		storage = conf.StorageView
		if err := b.Setup(conf); err != nil {
			return nil, err
		}
		return b, nil
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: factory,
		Steps: []logicaltest.TestStep{
			testAccStepMapAppId(t),
			testAccStepMapUserId(t),
			testAccLogin(t, ""),
			testAccLoginAppIDInPath(t, ""),
			testAccLoginInvalid(t),
			testAccStepDeleteUserId(t),
			testAccLoginDeleted(t),
		},
	})

	req := &logical.Request{
		Path:      "map/app-id",
		Operation: logical.ListOperation,
		Storage:   storage,
	}
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	salt, err := b.Salt()
	if err != nil {
		t.Fatal(err)
	}
	if keys[0] != salt.SaltID("foo") {
		t.Fatal("value was improperly salted")
	}
}

func TestBackend_cidr(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepMapAppIdDisplayName(t),
			testAccStepMapUserIdCidr(t, "192.168.1.0/16"),
			testAccLoginCidr(t, "192.168.1.5", false),
			testAccLoginCidr(t, "10.0.1.5", true),
			testAccLoginCidr(t, "", true),
		},
	})
}

func TestBackend_displayName(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepMapAppIdDisplayName(t),
			testAccStepMapUserId(t),
			testAccLogin(t, "tubbin"),
			testAccLoginAppIDInPath(t, "tubbin"),
			testAccLoginInvalid(t),
			testAccStepDeleteUserId(t),
			testAccLoginDeleted(t),
		},
	})
}

func testAccStepMapAppId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/app-id/foo",
		Data: map[string]interface{}{
			"value": "foo,bar",
		},
	}
}

func testAccStepMapAppIdDisplayName(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/app-id/foo",
		Data: map[string]interface{}{
			"display_name": "tubbin",
			"value":        "foo,bar",
		},
	}
}

func testAccStepMapUserId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/user-id/42",
		Data: map[string]interface{}{
			"value": "foo",
		},
	}
}

func testAccStepDeleteUserId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "map/user-id/42",
	}
}

func testAccStepMapUserIdCidr(t *testing.T, cidr string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/user-id/42",
		Data: map[string]interface{}{
			"value":      "foo",
			"cidr_block": cidr,
		},
	}
}

func testAccLogin(t *testing.T, display string) logicaltest.TestStep {
	checkTTL := func(resp *logical.Response) error {
		if resp.Auth.LeaseOptions.TTL.String() != "768h0m0s" {
			return fmt.Errorf("invalid TTL")
		}
		return nil
	}
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"app_id":  "foo",
			"user_id": "42",
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckMulti(
			logicaltest.TestCheckAuth([]string{"bar", "default", "foo"}),
			logicaltest.TestCheckAuthDisplayName(display),
			checkTTL,
		),
	}
}

func testAccLoginAppIDInPath(t *testing.T, display string) logicaltest.TestStep {
	checkTTL := func(resp *logical.Response) error {
		if resp.Auth.LeaseOptions.TTL.String() != "768h0m0s" {
			return fmt.Errorf("invalid TTL")
		}
		return nil
	}
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/foo",
		Data: map[string]interface{}{
			"user_id": "42",
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckMulti(
			logicaltest.TestCheckAuth([]string{"bar", "default", "foo"}),
			logicaltest.TestCheckAuthDisplayName(display),
			checkTTL,
		),
	}
}

func testAccLoginCidr(t *testing.T, ip string, err bool) logicaltest.TestStep {
	check := logicaltest.TestCheckError()
	if !err {
		check = logicaltest.TestCheckAuth([]string{"bar", "default", "foo"})
	}

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"app_id":  "foo",
			"user_id": "42",
		},
		ErrorOk:         err,
		Unauthenticated: true,
		RemoteAddr:      ip,

		Check: check,
	}
}

func testAccLoginInvalid(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"app_id":  "foo",
			"user_id": "48",
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}

func testAccLoginDeleted(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"app_id":  "foo",
			"user_id": "42",
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}
