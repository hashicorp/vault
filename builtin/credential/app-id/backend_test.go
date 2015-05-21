package appId

import (
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_basic(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepMapAppId(t),
			testAccStepMapUserId(t),
			testAccLogin(t, ""),
			testAccLoginInvalid(t),
			testAccStepDeleteUserId(t),
			testAccLoginDeleted(t),
		},
	})
}

func TestBackend_cidr(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
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
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepMapAppIdDisplayName(t),
			testAccStepMapUserId(t),
			testAccLogin(t, "tubbin"),
			testAccLoginInvalid(t),
			testAccStepDeleteUserId(t),
			testAccLoginDeleted(t),
		},
	})
}

func testAccStepMapAppId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "map/app-id/foo",
		Data: map[string]interface{}{
			"value": "foo,bar",
		},
	}
}

func testAccStepMapAppIdDisplayName(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "map/app-id/foo",
		Data: map[string]interface{}{
			"display_name": "tubbin",
			"value":        "foo,bar",
		},
	}
}

func testAccStepMapUserId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
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
		Operation: logical.WriteOperation,
		Path:      "map/user-id/42",
		Data: map[string]interface{}{
			"value":      "foo",
			"cidr_block": cidr,
		},
	}
}

func testAccLogin(t *testing.T, display string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"app_id":  "foo",
			"user_id": "42",
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckMulti(
			logicaltest.TestCheckAuth([]string{"bar", "foo"}),
			logicaltest.TestCheckAuthDisplayName(display),
		),
	}
}

func testAccLoginCidr(t *testing.T, ip string, err bool) logicaltest.TestStep {
	check := logicaltest.TestCheckError()
	if !err {
		check = logicaltest.TestCheckAuth([]string{"bar", "foo"})
	}

	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
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
		Operation: logical.WriteOperation,
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
		Operation: logical.WriteOperation,
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
