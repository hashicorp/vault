package http

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/vault"
)

func TestSysRootGenerationInit_Status(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp, err := http.Get(addr + "/v1/sys/root-generation/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            false,
		"progress":           float64(0),
		"required":           float64(1),
		"complete":           false,
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["nonce"] = actual["nonce"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRootGenerationInit_Setup_OTP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	otpBytes, err := xor.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"otp": otp,
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/root-generation/attempt")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            true,
		"progress":           float64(0),
		"required":           float64(1),
		"complete":           false,
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["nonce"] = actual["nonce"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRootGenerationInit_Setup_PGP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/root-generation/attempt")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            true,
		"progress":           float64(0),
		"required":           float64(1),
		"complete":           false,
		"encoded_root_token": "",
		"pgp_fingerprint":    "816938b8a29146fbe245dd29e7cbaf8e011db793",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["nonce"] = actual["nonce"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRootGenerationInit_Cancel(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	otpBytes, err := xor.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"otp": otp,
	})

	resp = testHttpDelete(t, token, addr+"/v1/sys/root-generation/attempt")
	testResponseStatus(t, resp, 204)

	resp, err = http.Get(addr + "/v1/sys/root-generation/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"started":            false,
		"progress":           float64(0),
		"required":           float64(1),
		"complete":           false,
		"encoded_root_token": "",
		"pgp_fingerprint":    "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["nonce"] = actual["nonce"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}
}

func TestSysRootGeneration_badKey(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	otpBytes, err := xor.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/update", map[string]interface{}{
		"key": "0123",
		"otp": otp,
	})
	testResponseStatus(t, resp, 400)
}

func TestSysRootGeneration_ReInitUpdate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	otpBytes, err := xor.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)
	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"otp": otp,
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/root-generation/attempt")
	testResponseStatus(t, resp, 204)

	resp = testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})

	testResponseStatus(t, resp, 204)
}

func TestSysRootGeneration_Update_OTP(t *testing.T) {
	core, master, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	otpBytes, err := xor.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"otp": otp,
	})
	testResponseStatus(t, resp, 204)

	// We need to get the nonce first before we update
	resp, err = http.Get(addr + "/v1/sys/root-generation/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	var rootGenerationStatus map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &rootGenerationStatus)

	resp = testHttpPut(t, token, addr+"/v1/sys/root-generation/update", map[string]interface{}{
		"nonce": rootGenerationStatus["nonce"].(string),
		"key":   hex.EncodeToString(master),
	})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"complete":        true,
		"nonce":           rootGenerationStatus["nonce"].(string),
		"progress":        float64(1),
		"required":        float64(1),
		"started":         true,
		"pgp_fingerprint": "",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	if actual["encoded_root_token"] == nil {
		t.Fatalf("no encoded root token found in response")
	}
	expected["encoded_root_token"] = actual["encoded_root_token"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	decodedToken, err := xor.XORBase64(otp, actual["encoded_root_token"].(string))
	if err != nil {
		t.Fatal(err)
	}
	newRootToken, err := uuid.FormatUUID(decodedToken)
	if err != nil {
		t.Fatal(err)
	}

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"id":           newRootToken,
		"display_name": "root",
		"meta":         interface{}(nil),
		"num_uses":     float64(0),
		"policies":     []interface{}{"root"},
		"orphan":       true,
		"ttl":          float64(0),
		"path":         "auth/token/root",
	}

	resp = testHttpGet(t, newRootToken, addr+"/v1/auth/token/lookup-self")
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["creation_time"] = actual["data"].(map[string]interface{})["creation_time"]

	if !reflect.DeepEqual(actual["data"], expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual["data"])
	}
}

func TestSysRootGeneration_Update_PGP(t *testing.T) {
	core, master, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/root-generation/attempt", map[string]interface{}{
		"pgp_key": pgpkeys.TestPubKey1,
	})
	testResponseStatus(t, resp, 204)

	// We need to get the nonce first before we update
	resp, err := http.Get(addr + "/v1/sys/root-generation/attempt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	var rootGenerationStatus map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &rootGenerationStatus)

	resp = testHttpPut(t, token, addr+"/v1/sys/root-generation/update", map[string]interface{}{
		"nonce": rootGenerationStatus["nonce"].(string),
		"key":   hex.EncodeToString(master),
	})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"complete":        true,
		"nonce":           rootGenerationStatus["nonce"].(string),
		"progress":        float64(1),
		"required":        float64(1),
		"started":         true,
		"pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	if actual["encoded_root_token"] == nil {
		t.Fatalf("no encoded root token found in response")
	}
	expected["encoded_root_token"] = actual["encoded_root_token"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
	}

	decodedTokenBuf, err := pgpkeys.DecryptBytes(actual["encoded_root_token"].(string), pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if decodedTokenBuf == nil {
		t.Fatal("decoded root token buffer is nil")
	}

	newRootToken := decodedTokenBuf.String()

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"id":           newRootToken,
		"display_name": "root",
		"meta":         interface{}(nil),
		"num_uses":     float64(0),
		"policies":     []interface{}{"root"},
		"orphan":       true,
		"ttl":          float64(0),
		"path":         "auth/token/root",
	}

	resp = testHttpGet(t, newRootToken, addr+"/v1/auth/token/lookup-self")
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["creation_time"] = actual["data"].(map[string]interface{})["creation_time"]

	if !reflect.DeepEqual(actual["data"], expected) {
		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual["data"])
	}
}
