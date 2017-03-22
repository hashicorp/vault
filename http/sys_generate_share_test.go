package http

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/vault"
)

func TestSysGenerateShareAttempt_Status(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp, err := http.Get(addr + "/v1/sys/generate-share/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":         false,
		"progress":        json.Number("0"),
		"required":        json.Number("3"),
		"complete":        false,
		"key":             "",
		"pgp_fingerprint": "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysGenerateShareAttempt_Setup_PGP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-share/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-share/attempt")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":         true,
		"progress":        json.Number("0"),
		"required":        json.Number("3"),
		"complete":        false,
		"key":             "",
		"pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysGenerateShareAttempt_Cancel(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-share/attempt", map[string]interface{}{})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":         true,
		"progress":        json.Number("0"),
		"required":        json.Number("3"),
		"complete":        false,
		"key":             "",
		"pgp_fingerprint": "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	resp = testHttpDelete(t, token, addr+"/v1/sys/generate-share/attempt")
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/generate-share/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"started":         false,
		"progress":        json.Number("0"),
		"required":        json.Number("3"),
		"complete":        false,
		"key":             "",
		"pgp_fingerprint": "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysGenerateShare_badKey(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-share/update", map[string]interface{}{
		"key": "0123",
	})
	testResponseStatus(t, resp, 400)
}

func TestSysGenerateShare_ReAttemptUpdate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-share/attempt", map[string]interface{}{})
	testResponseStatus(t, resp, 200)

	resp = testHttpDelete(t, token, addr+"/v1/sys/generate-share/attempt")
	testResponseStatus(t, resp, 204)

	resp = testHttpPut(t, token, addr+"/v1/sys/generate-share/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})

	testResponseStatus(t, resp, 200)
}

func TestSysGenerateShare_Update_PGP(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-share/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})
	testResponseStatus(t, resp, 200)

	resp, err := http.Get(addr + "/v1/sys/generate-share/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	var rootGenerationStatus map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &rootGenerationStatus)

	var actual map[string]interface{}
	var expected map[string]interface{}
	for i, key := range keys {
		resp = testHttpPut(t, token, addr+"/v1/sys/generate-share/update", map[string]interface{}{
			"key": hex.EncodeToString(key),
		})

		actual = map[string]interface{}{}
		expected = map[string]interface{}{
			"complete":        false,
			"progress":        json.Number(fmt.Sprintf("%d", i+1)),
			"required":        json.Number(fmt.Sprintf("%d", len(keys))),
			"started":         true,
			"pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
		}
		if i+1 == len(keys) {
			expected["complete"] = true
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
	}

	if actual["key"] == nil {
		t.Fatalf("no master key share found in response")
	}
	expected["key"] = actual["key"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	decodedKeyBuf, err := pgpkeys.DecryptBytes(actual["key"].(string), pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if decodedKeyBuf == nil {
		t.Fatal("decoded key buffer is nil")
	}

	newShare := decodedKeyBuf.String()
	newShareBytes, err := base64.StdEncoding.DecodeString(newShare)

	keys[0] = newShareBytes

	for i, key := range keys {
		resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
			"key": hex.EncodeToString(key),
		})

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"sealed":   true,
			"t":        json.Number("3"),
			"n":        json.Number("3"),
			"progress": json.Number(fmt.Sprintf("%d", i+1)),
			"nonce":    "",
		}
		if i == len(keys)-1 {
			expected["sealed"] = false
			expected["progress"] = json.Number("0")
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
	}
}
