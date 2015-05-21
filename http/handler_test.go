package http

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

// We use this test to verify header auth
func TestSysMounts_headerAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/mounts", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(AuthHeaderName, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestHandler_sealed(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	core.Seal(token)

	resp, err := http.Get(addr + "/v1/secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 503)
}
