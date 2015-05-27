package http

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

func TestLogical(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// WRITE
	resp := testHttpPut(t, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 204)

	// READ
	resp, err := http.Get(addr + "/v1/secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"renewable":      false,
		"lease_duration": float64((30 * 24 * time.Hour) / time.Second),
		"data": map[string]interface{}{
			"data": "bar",
		},
		"auth": nil,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual, "lease_id")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v %#v", actual, expected)
	}

	// DELETE
	resp = testHttpDelete(t, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 204)

	resp, err = http.Get(addr + "/v1/secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 404)
}

func TestLogical_noExist(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp, err := http.Get(addr + "/v1/secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 404)
}

func TestLogical_StandbyRedirect(t *testing.T) {
	ln1, addr1 := TestListener(t)
	defer ln1.Close()
	ln2, addr2 := TestListener(t)
	defer ln2.Close()

	// Create an HA Vault
	inm := physical.NewInmemHA()
	conf := &vault.CoreConfig{
		Physical:      inm,
		AdvertiseAddr: addr1,
		DisableMlock:  true,
	}
	core1, err := vault.NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	key, root := vault.TestCoreInit(t, core1)
	if _, err := core1.Unseal(vault.TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
	}

	// Create a second HA Vault
	conf2 := &vault.CoreConfig{
		Physical:      inm,
		AdvertiseAddr: addr2,
		DisableMlock:  true,
	}
	core2, err := vault.NewCore(conf2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, err := core2.Unseal(vault.TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
	}

	TestServerWithListener(t, ln1, addr1, core1)
	TestServerWithListener(t, ln2, addr2, core2)
	TestServerAuth(t, addr1, root)

	// WRITE to STANDBY
	resp := testHttpPut(t, addr2+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 307)

	//// READ to standby
	resp, err = http.Get(addr2 + "/v1/auth/token/lookup-self")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"renewable":      false,
		"lease_duration": float64(0),
		"data": map[string]interface{}{
			"meta":         nil,
			"num_uses":     float64(0),
			"path":         "auth/token/root",
			"policies":     []interface{}{"root"},
			"display_name": "root",
			"id":           root,
		},
		"auth": nil,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual, "lease_id")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v %#v", actual, expected)
	}

	//// DELETE to standby
	resp = testHttpDelete(t, addr2+"/v1/secret/foo")
	testResponseStatus(t, resp, 307)
}

func TestLogical_CreateToken(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// WRITE
	resp := testHttpPut(t, addr+"/v1/auth/token/create", map[string]interface{}{
		"data": "bar",
	})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": float64(0),
		"data":           nil,
		"auth": map[string]interface{}{
			"policies":       []interface{}{"root"},
			"metadata":       nil,
			"lease_duration": float64(0),
			"renewable":      false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual["auth"].(map[string]interface{}), "client_token")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v %#v", actual, expected)
	}

	// Should not get auth cookie
	if cookies := resp.Cookies(); len(cookies) != 0 {
		t.Fatalf("should not get cookies: %#v", cookies)
	}
}

func TestLogical_RawHTTP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type": "http",
	})
	testResponseStatus(t, resp, 204)

	// Get the raw response
	resp, err := http.Get(addr + "/v1/foo/raw")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 200)

	// Test the headers
	if resp.Header.Get("Content-Type") != "plain/text" {
		t.Fatalf("Bad: %#v", resp.Header)
	}

	// Get the body
	body := new(bytes.Buffer)
	io.Copy(body, resp.Body)
	if string(body.Bytes()) != "hello world" {
		t.Fatalf("Bad: %s", body.Bytes())
	}
}
