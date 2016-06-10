package http

import (
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
		"policies": []interface{}{"default", "response-wrapping", "root"},
		"keys":     []interface{}{"default", "response-wrapping", "root"},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
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
		"name":  "root",
		"rules": "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
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
		"rules": ``,
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/policy")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"policies": []interface{}{"default", "foo", "response-wrapping", "root"},
		"keys":     []interface{}{"default", "foo", "response-wrapping", "root"},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
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
		"rules": ``,
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/policy/foo")
	testResponseStatus(t, resp, 204)

	// Also attempt to delete these since they should not be allowed (ignore
	// responses, if they exist later that's sufficient)
	resp = testHttpDelete(t, token, addr+"/v1/sys/policy/default")
	resp = testHttpDelete(t, token, addr+"/v1/sys/policy/response-wrapping")

	resp = testHttpGet(t, token, addr+"/v1/sys/policy")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"policies": []interface{}{"default", "response-wrapping", "root"},
		"keys":     []interface{}{"default", "response-wrapping", "root"},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}
