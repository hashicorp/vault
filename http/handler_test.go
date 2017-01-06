package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestHandler_CacheControlNoStore(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/mounts", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(AuthHeaderName, token)
	req.Header.Set(WrapTTLHeaderName, "60s")

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp == nil {
		t.Fatalf("nil response")
	}

	actual := resp.Header.Get("Cache-Control")

	if actual == "" {
		t.Fatalf("missing 'Cache-Control' header entry in response writer")
	}

	if actual != "no-store" {
		t.Fatalf("bad: Cache-Control. Expected: 'no-store', Actual: %q", actual)
	}
}

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

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"secret/": map[string]interface{}{
				"description": "generic secret storage",
				"type":        "generic",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
				},
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
				},
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
				},
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["request_id"] = actual["request_id"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual: %#v\n", expected, actual)
	}
}

// We use this test to verify header auth wrapping
func TestSysMounts_headerAuth_Wrapped(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/mounts", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(AuthHeaderName, token)
	req.Header.Set(WrapTTLHeaderName, "60s")

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"request_id":     "",
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data":           nil,
		"wrap_info": map[string]interface{}{
			"ttl": json.Number("60"),
		},
		"warnings": nil,
		"auth":     nil,
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	actualToken, ok := actual["wrap_info"].(map[string]interface{})["token"]
	if !ok || actualToken == "" {
		t.Fatal("token missing in wrap info")
	}
	expected["wrap_info"].(map[string]interface{})["token"] = actualToken

	actualCreationTime, ok := actual["wrap_info"].(map[string]interface{})["creation_time"]
	if !ok || actualCreationTime == "" {
		t.Fatal("creation_time missing in wrap info")
	}
	expected["wrap_info"].(map[string]interface{})["creation_time"] = actualCreationTime

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual: %#v\n%T %T", expected, actual, actual["warnings"], actual["data"])
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

func TestHandler_error(t *testing.T) {
	w := httptest.NewRecorder()

	respondError(w, 500, errors.New("Test Error"))

	if w.Code != 500 {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	// The code inside of the error should override
	// the argument to respondError
	w2 := httptest.NewRecorder()
	e := logical.CodedError(403, "error text")

	respondError(w2, 500, e)

	if w2.Code != 403 {
		t.Fatalf("expected 403, got %d", w2.Code)
	}

	// vault.ErrSealed is a special case
	w3 := httptest.NewRecorder()

	respondError(w3, 400, vault.ErrSealed)

	if w3.Code != 503 {
		t.Fatalf("expected 503, got %d", w3.Code)
	}

}
