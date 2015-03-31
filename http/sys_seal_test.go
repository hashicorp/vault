package http

import (
	"encoding/hex"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysSealStatus(t *testing.T) {
	core := vault.TestCore(t)
	vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"sealed":   true,
		"t":        float64(1),
		"n":        float64(1),
		"progress": float64(0),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysSealStatus_uninit(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 400)
}

func TestSysSeal(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, addr+"/v1/sys/seal", nil)
	testResponseStatus(t, resp, 204)

	check, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !check {
		t.Fatal("should be sealed")
	}
}

func TestSysSeal_unsealed(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, addr+"/v1/sys/seal", nil)
	testResponseStatus(t, resp, 204)

	check, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !check {
		t.Fatal("should be sealed")
	}
}

func TestSysUnseal(t *testing.T) {
	core := vault.TestCore(t)
	key, _ := vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, addr+"/v1/sys/unseal", map[string]interface{}{
		"key": hex.EncodeToString(key),
	})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"sealed":   false,
		"t":        float64(1),
		"n":        float64(1),
		"progress": float64(0),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysUnseal_badKey(t *testing.T) {
	core := vault.TestCore(t)
	vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, addr+"/v1/sys/unseal", map[string]interface{}{
		"key": "0123",
	})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"sealed":   true,
		"t":        float64(1),
		"n":        float64(1),
		"progress": float64(0),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
