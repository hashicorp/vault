package http

import (
	"encoding/hex"
	"net/http"
	"reflect"
	"testing"
)

func TestSysSealStatus(t *testing.T) {
	core := testCore(t)
	testCoreInit(t, core)
	ln, addr := testServer(t, core)
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

func TestSysSeal(t *testing.T) {
	core := testCore(t)
	testCoreInit(t, core)
	ln, addr := testServer(t, core)
	defer ln.Close()

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
	core := testCore(t)
	ln, addr := testServer(t, core)
	defer ln.Close()

	keys := testCoreInit(t, core)
	if _, err := core.Unseal(keys[0]); err != nil {
		t.Fatalf("err: %s", err)
	}

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
	core := testCore(t)
	keys := testCoreInit(t, core)
	ln, addr := testServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, addr+"/v1/sys/unseal", map[string]interface{}{
		"key": hex.EncodeToString(keys[0]),
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
	core := testCore(t)
	testCoreInit(t, core)
	ln, addr := testServer(t, core)
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
