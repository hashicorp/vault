package appId

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_basic(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
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

// Verify that we are able to update from non-salted (<0.2) to
// using a Salt for the paths
func TestBackend_upgradeToSalted(t *testing.T) {
	inm := new(logical.InmemStorage)

	// Create some fake keys
	se, _ := logical.StorageEntryJSON("struct/map/app-id/foo",
		map[string]string{"value": "test"})
	inm.Put(se)
	se, _ = logical.StorageEntryJSON("struct/map/user-id/bar",
		map[string]string{"value": "foo"})
	inm.Put(se)

	// Initialize the backend, this should do the automatic upgrade
	conf := &logical.BackendConfig{
		StorageView: inm,
	}
	backend, err := Factory(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the keys have been upgraded
	out, err := inm.Get("struct/map/app-id/foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("unexpected key")
	}
	out, err = inm.Get("struct/map/user-id/bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("unexpected key")
	}

	// Backend should still be able to resolve
	req := logical.TestRequest(t, logical.ReadOperation, "map/app-id/foo")
	req.Storage = inm
	resp, err := backend.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["value"] != "test" {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "map/user-id/bar")
	req.Storage = inm
	resp, err = backend.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["value"] != "foo" {
		t.Fatalf("bad: %#v", resp)
	}
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
		if resp.Auth.LeaseOptions.TTL.String() != "720h0m0s" {
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
		if resp.Auth.LeaseOptions.TTL.String() != "720h0m0s" {
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
