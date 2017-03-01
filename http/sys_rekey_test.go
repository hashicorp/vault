package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/vault"
)

// Test to check if the API errors out when wrong number of PGP keys are
// supplied for rekey
func TestSysRekeyInit_pgpKeysEntriesForRekey(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
		"pgp_keys":         []string{"pgpkey1"},
	})
	testResponseStatus(t, resp, 400)
}

// Test to check if the API errors out when pgp and wrap shares are supplied
func TestSysRekeyInit_pgpAndWrap(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    1,
		"secret_threshold": 1,
		"pgp_keys":         []string{pgpkeys.TestPubKey1},
		"wrap_shares":      true,
	})
	testResponseStatus(t, resp, 400)
}

func TestSysRekeyInit_Status(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp, err := http.Get(addr + "/v1/sys/rekey/init")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":          false,
		"t":                json.Number("0"),
		"n":                json.Number("0"),
		"progress":         json.Number("0"),
		"required":         json.Number("3"),
		"pgp_fingerprints": interface{}(nil),
		"backup":           false,
		"nonce":            "",
		"wrap_shares":      false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRekeyInit_Setup(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})
	testResponseStatus(t, resp, 200)

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":          true,
		"t":                json.Number("3"),
		"n":                json.Number("5"),
		"progress":         json.Number("0"),
		"required":         json.Number("3"),
		"pgp_fingerprints": interface{}(nil),
		"backup":           false,
		"wrap_shares":      false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["nonce"].(string) == "" {
		t.Fatalf("nonce was empty")
	}
	expected["nonce"] = actual["nonce"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	resp = testHttpGet(t, token, addr+"/v1/sys/rekey/init")

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"started":          true,
		"t":                json.Number("3"),
		"n":                json.Number("5"),
		"progress":         json.Number("0"),
		"required":         json.Number("3"),
		"pgp_fingerprints": interface{}(nil),
		"backup":           false,
		"wrap_shares":      false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["nonce"].(string) == "" {
		t.Fatalf("nonce was empty")
	}
	expected["nonce"] = actual["nonce"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRekeyInit_Cancel(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpDelete(t, token, addr+"/v1/sys/rekey/init")
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/rekey/init")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":          false,
		"t":                json.Number("0"),
		"n":                json.Number("0"),
		"progress":         json.Number("0"),
		"required":         json.Number("3"),
		"pgp_fingerprints": interface{}(nil),
		"backup":           false,
		"nonce":            "",
		"wrap_shares":      false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRekey_badKey(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/update", map[string]interface{}{
		"key": "0123",
	})
	testResponseStatus(t, resp, 400)
}

func TestSysRekey_Update(t *testing.T) {
	testSysRekey_Update(t, false)
	testSysRekey_Update(t, true)
}

func testSysRekey_Update(t *testing.T, wrapShares bool) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
		"wrap_shares":      wrapShares,
	})
	var rekeyStatus map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &rekeyStatus)

	var actual map[string]interface{}
	var expected map[string]interface{}

	for i, key := range keys {
		resp = testHttpPut(t, token, addr+"/v1/sys/rekey/update", map[string]interface{}{
			"nonce": rekeyStatus["nonce"].(string),
			"key":   hex.EncodeToString(key),
		})

		actual = map[string]interface{}{}
		expected = map[string]interface{}{
			"started":          true,
			"nonce":            rekeyStatus["nonce"].(string),
			"backup":           false,
			"pgp_fingerprints": interface{}(nil),
			"required":         json.Number("3"),
			"t":                json.Number("3"),
			"n":                json.Number("5"),
			"progress":         json.Number(fmt.Sprintf("%d", i+1)),
			"wrap_shares":      wrapShares,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)

		if i+1 == len(keys) {
			delete(expected, "started")
			delete(expected, "required")
			delete(expected, "t")
			delete(expected, "n")
			delete(expected, "progress")
			expected["complete"] = true
			expected["keys"] = actual["keys"]
			expected["keys_base64"] = actual["keys_base64"]
		}

		if i+1 < len(keys) && (actual["nonce"] == nil || actual["nonce"].(string) == "") {
			t.Fatalf("expected a nonce, i is %d, actual is %#v", i, actual)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected: \n%#v\nactual: \n%#v", expected, actual)
		}
	}

	retKeys := actual["keys"].([]interface{})
	if len(retKeys) != 5 {
		t.Fatalf("bad: %#v", retKeys)
	}
	keysB64 := actual["keys_base64"].([]interface{})
	if len(keysB64) != 5 {
		t.Fatalf("bad: %#v", keysB64)
	}

	if wrapShares {
		for _, key := range retKeys {
			if strings.Count(key.(string), ".") != 2 {
				t.Fatalf("expected JWT token, got: %s", key.(string))
			}
		}
	}

}

func TestSysRekey_ReInitUpdate(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpDelete(t, token, addr+"/v1/sys/rekey/init")
	testResponseStatus(t, resp, 204)

	resp = testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpPut(t, token, addr+"/v1/sys/rekey/update", map[string]interface{}{
		"key": hex.EncodeToString(keys[0]),
	})

	testResponseStatus(t, resp, 400)
}
