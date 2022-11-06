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

	actual := map[string]any{}
	expected := map[string]any{
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data": map[string]any{
			"auth":   map[string]any{},
			"secret": map[string]any{},
		},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Mount-tune the listing_visibility
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]any{
		"listing_visibility": "unauth",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]any{
		"listing_visibility": "unauth",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, "", addr+"/v1/sys/internal/ui/mounts")
	testResponseStatus(t, resp, 200)

	actual = map[string]any{}
	expected = map[string]any{
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data": map[string]any{
			"secret": map[string]any{
				"secret/": map[string]any{
					"type":        "kv",
					"description": "key/value secret storage",
					"options":     map[string]any{"version": "1"},
				},
			},
			"auth": map[string]any{
				"token/": map[string]any{
					"type":        "token",
					"description": "token based credentials",
					"options":     any(nil),
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
