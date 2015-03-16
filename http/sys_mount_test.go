package http

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysMounts(t *testing.T) {
	core, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/mounts")
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
