package http

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysAudit(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/audit/noop", map[string]any{
		"type": "noop",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/audit")

	var actual map[string]any
	expected := map[string]any{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]any{
			"noop/": map[string]any{
				"path":        "noop/",
				"type":        "noop",
				"description": "",
				"options":     map[string]any{},
				"local":       false,
			},
		},
		"noop/": map[string]any{
			"path":        "noop/",
			"type":        "noop",
			"description": "",
			"options":     map[string]any{},
			"local":       false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["request_id"] = actual["request_id"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:\n%#v actual:\n%#v\n", expected, actual)
	}
}

func TestSysDisableAudit(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/audit/foo", map[string]any{
		"type": "noop",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/audit/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/audit")

	var actual map[string]any
	expected := map[string]any{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data":           map[string]any{},
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["request_id"] = actual["request_id"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nactual:   %#v\nexpected: %#v\n", actual, expected)
	}
}

func TestSysAuditHash(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/audit/noop", map[string]any{
		"type": "noop",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/audit-hash/noop", map[string]any{
		"input": "bar",
	})

	var actual map[string]any
	expected := map[string]any{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]any{
			"hash": "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
		},
		"hash": "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["request_id"] = actual["request_id"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:\n%#v\n, got:\n%#v\n", expected, actual)
	}
}
