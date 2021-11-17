package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestHelp(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, "", addr+"/v1/sys/mounts?help=1")
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatal("expected permission denied with no token")
	}

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts?help=1")

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if _, ok := actual["help"]; !ok {
		t.Fatalf("bad: %#v", actual)
	}
}
