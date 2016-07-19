package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func TestLogical(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// WRITE
	resp := testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 204)

	// READ
	// Bad token should return a 403
	resp = testHttpGet(t, token+"bad", addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 403)

	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	var actual map[string]interface{}
	var nilWarnings interface{}
	expected := map[string]interface{}{
		"renewable":      false,
		"lease_duration": json.Number(strconv.Itoa(int((30 * 24 * time.Hour) / time.Second))),
		"data": map[string]interface{}{
			"data": "bar",
		},
		"auth":      nil,
		"wrap_info": nil,
		"warnings":  nilWarnings,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual, "lease_id")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nactual:\n%#v\nexpected:\n%#v", actual, expected)
	}

	// DELETE
	resp = testHttpDelete(t, token, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 404)
}

func TestLogical_noExist(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 404)
}

func TestLogical_StandbyRedirect(t *testing.T) {
	ln1, addr1 := TestListener(t)
	defer ln1.Close()
	ln2, addr2 := TestListener(t)
	defer ln2.Close()

	// Create an HA Vault
	inmha := physical.NewInmemHA(logger)
	conf := &vault.CoreConfig{
		Physical:      inmha,
		HAPhysical:    inmha,
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

	// Attempt to fix raciness in this test by giving the first core a chance
	// to grab the lock
	time.Sleep(time.Second)

	// Create a second HA Vault
	conf2 := &vault.CoreConfig{
		Physical:      inmha,
		HAPhysical:    inmha,
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
	resp := testHttpPut(t, root, addr2+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 307)

	//// READ to standby
	resp = testHttpGet(t, root, addr2+"/v1/auth/token/lookup-self")
	var actual map[string]interface{}
	var nilWarnings interface{}
	expected := map[string]interface{}{
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data": map[string]interface{}{
			"meta":             nil,
			"num_uses":         json.Number("0"),
			"path":             "auth/token/root",
			"policies":         []interface{}{"root"},
			"display_name":     "root",
			"orphan":           true,
			"id":               root,
			"ttl":              json.Number("0"),
			"creation_ttl":     json.Number("0"),
			"role":             "",
			"explicit_max_ttl": json.Number("0"),
		},
		"warnings":  nilWarnings,
		"wrap_info": nil,
		"auth":      nil,
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	actualDataMap := actual["data"].(map[string]interface{})
	delete(actualDataMap, "creation_time")
	delete(actualDataMap, "accessor")
	actual["data"] = actualDataMap
	delete(actual, "lease_id")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got %#v; expected %#v", actual, expected)
	}

	//// DELETE to standby
	resp = testHttpDelete(t, root, addr2+"/v1/secret/foo")
	testResponseStatus(t, resp, 307)
}

func TestLogical_CreateToken(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// WRITE
	resp := testHttpPut(t, token, addr+"/v1/auth/token/create", map[string]interface{}{
		"data": "bar",
	})

	var actual map[string]interface{}
	var nilWarnings interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data":           nil,
		"wrap_info":      nil,
		"auth": map[string]interface{}{
			"policies":       []interface{}{"root"},
			"metadata":       nil,
			"lease_duration": json.Number("0"),
			"renewable":      true,
		},
		"warnings": nilWarnings,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual["auth"].(map[string]interface{}), "client_token")
	delete(actual["auth"].(map[string]interface{}), "accessor")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nexpected:\n%#v\nactual:\n%#v", expected, actual)
	}
}

func TestLogical_RawHTTP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type": "http",
	})
	testResponseStatus(t, resp, 204)

	// Get the raw response
	resp = testHttpGet(t, token, addr+"/v1/foo/raw")
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
