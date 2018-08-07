package http

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysInternal_UIMounts(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Get original tune values, ensure that listing_visibility is not set
	resp := testHttpGet(t, "", addr+"/v1/sys/internal/ui/mounts")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data": map[string]interface{}{
			"auth":   map[string]interface{}{},
			"secret": map[string]interface{}{},
		},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Mount-tune the listing_visibility
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"listing_visibility": "unauth",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]interface{}{
		"listing_visibility": "unauth",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, "", addr+"/v1/sys/internal/ui/mounts")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data": map[string]interface{}{
			"secret": map[string]interface{}{
				"secret/": map[string]interface{}{
					"type":        "kv",
					"description": "key/value secret storage",
					"options":     map[string]interface{}{"version": "1"},
				},
			},
			"auth": map[string]interface{}{
				"token/": map[string]interface{}{
					"type":        "token",
					"description": "token based credentials",
					"options":     interface{}(nil),
				},
			},
		},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}
