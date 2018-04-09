package openpgp

import (
	"context"
	"github.com/hashicorp/vault/logical"
	"testing"
)

func TestPGP_SignVerify(t *testing.T) {
	var b *backend
	storage := &logical.InmemStorage{}

	b = Backend()

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test",
		Data: map[string]interface{}{
			"real_name": "Vault PGP test",
			"email":     "vault@example.com",
		},
	}
	req2 := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test2",
		Data: map[string]interface{}{
			"real_name": "Vault PGP test2",
			"email":     "vault@example.com",
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.HandleRequest(context.Background(), req2)
	if err != nil {
		t.Fatal(err)
	}

	signRequest := func(req *logical.Request, keyName string, errExpected bool, postpath string) string {
		req.Path = "sign/" + keyName + postpath
		response, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatal(err)
		}
		if response == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !response.IsError() {
				t.Fatalf("expected error response: %#v", *response)
			}
			return ""
		}
		if response.IsError() {
			t.Fatalf("not expected error response: %#v", *response)
		}
		value, ok := response.Data["signature"]
		if !ok {
			t.Fatalf("no signature found in response data: %#v", response.Data)
		}
		return value.(string)
	}

	verifyRequest := func(req *logical.Request, keyName string, errExpected, validSignature bool, signature string) {
		req.Path = "verify/" + keyName
		req.Data["signature"] = signature
		response, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatalf("error: %v, signature was %v", err, signature)
		}
		if errExpected {
			if response != nil && !response.IsError() {
				t.Fatalf("expected error response: %#v", *response)
			}
			return
		}
		if response == nil {
			t.Fatal("expected non-nil response")
		}
		if response.IsError() {
			t.Fatalf("not expected error response: %#v", *response)
		}
		value, ok := response.Data["valid"]
		if !ok {
			t.Fatalf("no valid key found in response data %#v", response.Data)
		}
		if validSignature && !value.(bool) {
			t.Fatalf("not expected failing signature verification %#v %#v", *req, *response)
		}
		if !validSignature && value.(bool) {
			t.Fatalf("expected failing signature verification %#v %#v", *req, *response)
		}
	}

	req.Data = map[string]interface{}{
		"input": "dGhlIHF1aWNrIGJyb3duIGZveA==",
	}

	// Test defaults
	signature := signRequest(req, "test", false, "")
	verifyRequest(req, "test", false, true, signature)
	verifyRequest(req, "test2", false, false, signature)

	// Test algorithm selection in path
	signature = signRequest(req, "test", false, "/sha2-224")
	verifyRequest(req, "test", false, true, signature)

	// Test algorithm selection in the data
	req.Data["algorithm"] = "sha2-224"
	signature = signRequest(req, "test", false, "")
	verifyRequest(req, "test", false, true, signature)

	req.Data["algorithm"] = "sha2-384"
	signature = signRequest(req, "test", false, "")
	verifyRequest(req, "test", false, true, signature)

	req.Data["algorithm"] = "sha2-512"
	signature = signRequest(req, "test", false, "")
	verifyRequest(req, "test", false, true, signature)

	req.Data["algorithm"] = "notexisting"
	signature = signRequest(req, "test", true, "")
	delete(req.Data, "algorithm")

	// Test format selection
	req.Data["format"] = "ascii-armor"
	signature = signRequest(req, "test", false, "")
	verifyRequest(req, "test", false, true, signature)

	// Test validation format mismatch
	req.Data["format"] = "ascii-armor"
	signature = signRequest(req, "test", false, "")
	req.Data["format"] = "base64"
	verifyRequest(req, "test", false, false, signature)

	// Test bad format
	req.Data["format"] = "notexisting"
	signRequest(req, "test", true, "")
	verifyRequest(req, "test", true, true, signature)
	delete(req.Data, "format")

	// Test non existent key
	signRequest(req, "notfound", true, "")
	verifyRequest(req, "notfound", true, false, signature)

	// Test bad input
	req.Data["input"] = "foobar"
	signRequest(req, "test", true, "")
	verifyRequest(req, "test", true, false, signature)
}
