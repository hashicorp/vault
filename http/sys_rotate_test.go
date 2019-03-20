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

	resp := testHttpPost(t, token, addr+"/v1/sys/rotate", map[string]interface{}{})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/key-status")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"term": json.Number("2"),
		},
		"term": json.Number("2"),
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	actualInstallTime, ok := actual["data"].(map[string]interface{})["install_time"]
	if !ok || actualInstallTime == "" {
		t.Fatal("install_time missing in data")
	}
	expected["data"].(map[string]interface{})["install_time"] = actualInstallTime
	expected["install_time"] = actualInstallTime

	expected["request_id"] = actual["request_id"]

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}
