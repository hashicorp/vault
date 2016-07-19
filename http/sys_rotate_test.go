package http

import (
	"encoding/json"
	"reflect"
	"testing"

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
		"term": json.Number("2"),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual, "install_time")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nexpected: %#v\nactual: %#v", expected, actual)
	}
}
