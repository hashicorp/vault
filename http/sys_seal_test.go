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

	resp := testHttpPut(t, token, addr+"/v1/sys/seal", nil)
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

	resp := testHttpPut(t, token, addr+"/v1/sys/seal", nil)
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

	resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
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

	resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
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

func TestSysUnseal_Reset(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	thresh := 3
	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": thresh,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	keysRaw, ok := actual["keys"]
	if !ok {
		t.Fatalf("no keys: %#v", actual)
	}
	for i, key := range keysRaw.([]interface{}) {
		if i > thresh-2 {
			break
		}

		resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
			"key": key.(string),
		})

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"sealed":   true,
			"t":        float64(3),
			"n":        float64(5),
			"progress": float64(i + 1),
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected:\n%#v\nactual:\n%#v\n", expected, actual)
		}
	}

	resp = testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
		"reset": true,
	})

	actual = map[string]interface{}{}
	expected := map[string]interface{}{
		"sealed":   true,
		"t":        float64(3),
		"n":        float64(5),
		"progress": float64(0),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected:\n%#v\nactual:\n%#v\n", expected, actual)
	}
}
