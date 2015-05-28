package http

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysRotate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, addr+"/v1/sys/rotate", map[string]interface{}{})
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/key-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"term": float64(2),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual, "install_time")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
