package http

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/vault"
)

func TestSysRotate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/rotate", map[string]any{})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/key-status")

	var actual map[string]any
	expected := map[string]any{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]any{
			"term": json.Number("2"),
		},
		"term": json.Number("2"),
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	for _, field := range []string{"install_time", "encryptions"} {
		actualVal, ok := actual["data"].(map[string]any)[field]
		if !ok || actualVal == "" {
			t.Fatal(field, " missing in data")
		}
		expected["data"].(map[string]any)[field] = actualVal
		expected[field] = actualVal
	}

	expected["request_id"] = actual["request_id"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}
