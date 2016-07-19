package http

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"token/": map[string]interface{}{
			"description": "token based credentials",
			"type":        "token",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestSysEnableAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/auth/foo", map[string]interface{}{
		"type":        "noop",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "noop",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
		},
		"token/": map[string]interface{}{
			"description": "token based credentials",
			"type":        "token",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestSysDisableAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/auth/foo", map[string]interface{}{
		"type":        "noop",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/auth/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"token/": map[string]interface{}{
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
			"description": "token based credentials",
			"type":        "token",
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}
