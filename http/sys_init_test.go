package http

import (
	"encoding/hex"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysInit_get(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	{
		// Pre-init
		resp, err := http.Get(addr + "/v1/sys/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"initialized": false,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: %#v", actual)
		}
	}

	vault.TestCoreInit(t, core)

	{
		// Post-init
		resp, err := http.Get(addr + "/v1/sys/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"initialized": true,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: %#v", actual)
		}
	}
}

func TestSysInit_put(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	keysRaw, ok := actual["keys"]
	if !ok {
		t.Fatalf("no keys: %#v", actual)
	}

	if _, ok := actual["root_token"]; !ok {
		t.Fatal("no root token")
	}

	for _, key := range keysRaw.([]interface{}) {
		keySlice, err := hex.DecodeString(key.(string))
		if err != nil {
			t.Fatalf("bad: %s", err)
		}

		if _, err := core.Unseal(keySlice); err != nil {
			t.Fatalf("bad: %s", err)
		}
	}

	seal, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if seal {
		t.Fatal("should not be sealed")
	}
}
