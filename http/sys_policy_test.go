package http

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysPolicies(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/policy")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"policies": []interface{}{"default", "root"},
			"keys":     []interface{}{"default", "root"},
		},
		"policies": []interface{}{"default", "root"},
		"keys":     []interface{}{"default", "root"},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}

func TestSysReadPolicy(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/policy/root")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"name":  "root",
			"rules": "",
		},
		"name":  "root",
		"rules": "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}

func TestSysWritePolicy(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/policy/foo", map[string]interface{}{
		"rules": `path "*" { capabilities = ["read"] }`,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpGet(t, token, addr+"/v1/sys/policy")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"policies": []interface{}{"default", "foo", "root"},
			"keys":     []interface{}{"default", "foo", "root"},
		},
		"policies": []interface{}{"default", "foo", "root"},
		"keys":     []interface{}{"default", "foo", "root"},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	resp = testHttpPost(t, token, addr+"/v1/sys/policy/response-wrapping", map[string]interface{}{
		"rules": ``,
	})
	testResponseStatus(t, resp, 400)
}

func TestSysDeletePolicy(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/policy/foo", map[string]interface{}{
		"rules": `path "*" { capabilities = ["read"] }`,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpDelete(t, token, addr+"/v1/sys/policy/foo")
	testResponseStatus(t, resp, 204)

	// Also attempt to delete these since they should not be allowed (ignore
	// responses, if they exist later that's sufficient)
	resp = testHttpDelete(t, token, addr+"/v1/sys/policy/default")
	resp = testHttpDelete(t, token, addr+"/v1/sys/policy/response-wrapping")

	resp = testHttpGet(t, token, addr+"/v1/sys/policy")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"policies": []interface{}{"default", "root"},
			"keys":     []interface{}{"default", "root"},
		},
		"policies": []interface{}{"default", "root"},
		"keys":     []interface{}{"default", "root"},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}
