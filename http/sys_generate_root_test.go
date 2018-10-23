package http

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/vault"
)

func TestSysGenerateRootAttempt_Status(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp, err := http.Get(addr + "/v1/sys/generate-root/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            false,
		"progress":           json.Number("0"),
		"required":           json.Number("3"),
		"complete":           false,
		"encoded_token":      "",
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
		"nonce":              "",
		"otp_length":         json.Number("26"),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["otp"] = actual["otp"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestSysGenerateRootAttempt_Setup_OTP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", nil)
	testResponseStatus(t, resp, 200)

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            true,
		"progress":           json.Number("0"),
		"required":           json.Number("3"),
		"complete":           false,
		"encoded_token":      "",
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
		"otp_length":         json.Number("26"),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["nonce"].(string) == "" {
		t.Fatalf("nonce was empty")
	}
	expected["nonce"] = actual["nonce"]
	expected["otp"] = actual["otp"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-root/attempt")

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"started":            true,
		"progress":           json.Number("0"),
		"required":           json.Number("3"),
		"complete":           false,
		"encoded_token":      "",
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
		"otp":                "",
		"otp_length":         json.Number("26"),
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

func TestSysGenerateRootAttempt_Setup_PGP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-root/attempt")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            true,
		"progress":           json.Number("0"),
		"required":           json.Number("3"),
		"complete":           false,
		"encoded_token":      "",
		"encoded_root_token": "",
		"pgp_fingerprint":    "816938b8a29146fbe245dd29e7cbaf8e011db793",
		"otp":                "",
		"otp_length":         json.Number("26"),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["nonce"].(string) == "" {
		t.Fatalf("nonce was empty")
	}
	expected["nonce"] = actual["nonce"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestSysGenerateRootAttempt_Cancel(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", nil)

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            true,
		"progress":           json.Number("0"),
		"required":           json.Number("3"),
		"complete":           false,
		"encoded_token":      "",
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
		"otp_length":         json.Number("26"),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["nonce"].(string) == "" {
		t.Fatalf("nonce was empty")
	}
	expected["nonce"] = actual["nonce"]
	expected["otp"] = actual["otp"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}

	resp = testHttpDelete(t, token, addr+"/v1/sys/generate-root/attempt")
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/generate-root/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"started":            false,
		"progress":           json.Number("0"),
		"required":           json.Number("3"),
		"complete":           false,
		"encoded_token":      "",
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
		"nonce":              "",
		"otp":                "",
		"otp_length":         json.Number("26"),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysGenerateRoot_badKey(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/update", map[string]interface{}{
		"key": "0123",
	})
	testResponseStatus(t, resp, 400)
}

func TestSysGenerateRoot_ReAttemptUpdate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", nil)
	testResponseStatus(t, resp, 200)

	resp = testHttpDelete(t, token, addr+"/v1/sys/generate-root/attempt")
	testResponseStatus(t, resp, 204)

	resp = testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})

	testResponseStatus(t, resp, 200)
}

func TestSysGenerateRoot_Update_OTP(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", map[string]interface{}{})
	var rootGenerationStatus map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &rootGenerationStatus)
	otp := rootGenerationStatus["otp"].(string)

	var actual map[string]interface{}
	var expected map[string]interface{}
	for i, key := range keys {
		resp = testHttpPut(t, token, addr+"/v1/sys/generate-root/update", map[string]interface{}{
			"nonce": rootGenerationStatus["nonce"].(string),
			"key":   hex.EncodeToString(key),
		})

		actual = map[string]interface{}{}
		expected = map[string]interface{}{
			"complete":        false,
			"nonce":           rootGenerationStatus["nonce"].(string),
			"progress":        json.Number(fmt.Sprintf("%d", i+1)),
			"required":        json.Number(fmt.Sprintf("%d", len(keys))),
			"started":         true,
			"pgp_fingerprint": "",
			"otp":             "",
			"otp_length":      json.Number("0"),
		}
		if i+1 == len(keys) {
			expected["complete"] = true
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
	}

	if actual["encoded_token"] == nil || actual["encoded_token"] == "" {
		t.Fatalf("no encoded token found in response")
	}
	if actual["encoded_root_token"] == nil || actual["encoded_root-token"] == "" {
		t.Fatalf("no encoded root token found in response")
	}
	expected["encoded_token"] = actual["encoded_token"]
	expected["encoded_root_token"] = actual["encoded_root_token"]
	expected["encoded_token"] = actual["encoded_token"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(expected["encoded_token"].(string))
	if err != nil {
		t.Fatal(err)
	}

	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		t.Fatal(err)
	}
	newRootToken := string(tokenBytes)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"id":               newRootToken,
		"display_name":     "root",
		"meta":             interface{}(nil),
		"num_uses":         json.Number("0"),
		"policies":         []interface{}{"root"},
		"orphan":           true,
		"creation_ttl":     json.Number("0"),
		"ttl":              json.Number("0"),
		"path":             "auth/token/root",
		"explicit_max_ttl": json.Number("0"),
		"expire_time":      nil,
		"entity_id":        "",
		"type":             "service",
	}

	resp = testHttpGet(t, newRootToken, addr+"/v1/auth/token/lookup-self")
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["creation_time"] = actual["data"].(map[string]interface{})["creation_time"]
	expected["accessor"] = actual["data"].(map[string]interface{})["accessor"]

	if !reflect.DeepEqual(actual["data"], expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual["data"])
	}
}

func TestSysGenerateRoot_Update_PGP(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/generate-root/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})
	testResponseStatus(t, resp, 200)

	// We need to get the nonce first before we update
	resp, err := http.Get(addr + "/v1/sys/generate-root/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	var rootGenerationStatus map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &rootGenerationStatus)

	var actual map[string]interface{}
	var expected map[string]interface{}
	for i, key := range keys {
		resp = testHttpPut(t, token, addr+"/v1/sys/generate-root/update", map[string]interface{}{
			"nonce": rootGenerationStatus["nonce"].(string),
			"key":   hex.EncodeToString(key),
		})

		actual = map[string]interface{}{}
		expected = map[string]interface{}{
			"complete":        false,
			"nonce":           rootGenerationStatus["nonce"].(string),
			"progress":        json.Number(fmt.Sprintf("%d", i+1)),
			"required":        json.Number(fmt.Sprintf("%d", len(keys))),
			"started":         true,
			"pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
			"otp":             "",
			"otp_length":      json.Number("0"),
		}
		if i+1 == len(keys) {
			expected["complete"] = true
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
	}

	if actual["encoded_token"] == nil || actual["encoded_token"] == "" {
		t.Fatalf("no encoded token found in response")
	}
	if actual["encoded_root_token"] == nil || actual["encoded_root-token"] == "" {
		t.Fatalf("no encoded root token found in response")
	}
	expected["encoded_token"] = actual["encoded_token"]
	expected["encoded_root_token"] = actual["encoded_root_token"]
	expected["encoded_token"] = actual["encoded_token"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	decodedTokenBuf, err := pgpkeys.DecryptBytes(actual["encoded_token"].(string), pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if decodedTokenBuf == nil {
		t.Fatal("decoded root token buffer is nil")
	}

	newRootToken := decodedTokenBuf.String()

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"id":               newRootToken,
		"display_name":     "root",
		"meta":             interface{}(nil),
		"num_uses":         json.Number("0"),
		"policies":         []interface{}{"root"},
		"orphan":           true,
		"creation_ttl":     json.Number("0"),
		"ttl":              json.Number("0"),
		"path":             "auth/token/root",
		"explicit_max_ttl": json.Number("0"),
		"expire_time":      nil,
		"entity_id":        "",
		"type":             "service",
	}

	resp = testHttpGet(t, newRootToken, addr+"/v1/auth/token/lookup-self")
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["creation_time"] = actual["data"].(map[string]interface{})["creation_time"]
	expected["accessor"] = actual["data"].(map[string]interface{})["accessor"]

	if diff := deep.Equal(actual["data"], expected); diff != nil {
		t.Fatal(diff)
	}
}
