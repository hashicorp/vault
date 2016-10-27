package transit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_HMAC(t *testing.T) {
	var b *backend
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}

	b = Backend(&logical.BackendConfig{
		StorageView: storage,
		System:      sysView,
	})

	// First create a key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo",
	}
	_, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	// Now, change the key value to something we control
	p, lock, err := b.lm.GetPolicyShared(storage, "foo")
	if err != nil {
		t.Fatal(err)
	}
	// We don't care as we're the only one using this
	lock.RUnlock()
	keyEntry := p.Keys[p.LatestVersion]
	keyEntry.HMACKey = []byte("01234567890123456789012345678901")
	p.Keys[p.LatestVersion] = keyEntry
	if err = p.Persist(storage); err != nil {
		t.Fatal(err)
	}

	req.Path = "hmac/foo"
	req.Data = map[string]interface{}{
		"input": "dGhlIHF1aWNrIGJyb3duIGZveA==",
	}

	doRequest := func(req *logical.Request, errExpected bool, expected string) {
		path := req.Path
		defer func() { req.Path = path }()

		resp, err := b.HandleRequest(req)
		if err != nil && !errExpected {
			panic(fmt.Sprintf("%v", err))
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
			}
			return
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		value, ok := resp.Data["hmac"]
		if !ok {
			t.Fatalf("no hmac key found in returned data, got resp data %#v", resp.Data)
		}
		if value.(string) != expected {
			panic(fmt.Sprintf("mismatched hashes; expected %s, got resp data %#v", expected, resp.Data))
		}

		// Now verify
		req.Path = strings.Replace(req.Path, "hmac", "verify", -1)
		req.Data["hmac"] = value.(string)
		resp, err = b.HandleRequest(req)
		if err != nil {
			t.Fatalf("%v: %v", err, resp)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Data["valid"].(bool) == false {
			panic(fmt.Sprintf("error validating hmac;\nreq:\n%#v\nresp:\n%#v", *req, *resp))
		}
	}

	// Comparisons are against values generated via openssl

	// Test defaults -- sha2-256
	doRequest(req, false, "vault:v1:UcBvm5VskkukzZHlPgm3p5P/Yr/PV6xpuOGZISya3A4=")

	// Test algorithm selection in the path
	req.Path = "hmac/foo/sha2-224"
	doRequest(req, false, "vault:v1:3p+ZWVquYDvu2dSTCa65Y3fgoMfIAc6fNaBbtg==")

	// Reset and test algorithm selection in the data
	req.Path = "hmac/foo"
	req.Data["algorithm"] = "sha2-224"
	doRequest(req, false, "vault:v1:3p+ZWVquYDvu2dSTCa65Y3fgoMfIAc6fNaBbtg==")

	req.Data["algorithm"] = "sha2-384"
	doRequest(req, false, "vault:v1:jDB9YXdPjpmr29b1JCIEJO93IydlKVfD9mA2EO9OmJtJQg3QAV5tcRRRb7IQGW9p")

	req.Data["algorithm"] = "sha2-512"
	doRequest(req, false, "vault:v1:PSXLXvkvKF4CpU65e2bK1tGBZQpcpCEM32fq2iUoiTyQQCfBcGJJItQ+60tMwWXAPQrC290AzTrNJucGrr4GFA==")

	// Test returning as base64
	req.Data["format"] = "base64"
	doRequest(req, false, "vault:v1:PSXLXvkvKF4CpU65e2bK1tGBZQpcpCEM32fq2iUoiTyQQCfBcGJJItQ+60tMwWXAPQrC290AzTrNJucGrr4GFA==")

	req.Data["algorithm"] = "foobar"
	doRequest(req, true, "")

	req.Data["algorithm"] = "sha2-256"
	req.Data["input"] = "foobar"
	doRequest(req, true, "")
	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="

	// Rotate
	err = p.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}
	keyEntry = p.Keys[2]
	// Set to another value we control
	keyEntry.HMACKey = []byte("12345678901234567890123456789012")
	p.Keys[2] = keyEntry
	if err = p.Persist(storage); err != nil {
		t.Fatal(err)
	}

	doRequest(req, false, "vault:v2:Dt+mO/B93kuWUbGMMobwUNX5Wodr6dL3JH4DMfpQ0kw=")

	// Verify a previous version
	req.Path = "verify/foo"

	req.Data["hmac"] = "vault:v1:UcBvm5VskkukzZHlPgm3p5P/Yr/PV6xpuOGZISya3A4="
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("%v: %v", err, resp)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Data["valid"].(bool) == false {
		t.Fatalf("error validating hmac\nreq\n%#v\nresp\n%#v", *req, *resp)
	}

	// Try a bad value
	req.Data["hmac"] = "vault:v1:UcBvm4VskkukzZHlPgm3p5P/Yr/PV6xpuOGZISya3A4="
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("%v: %v", err, resp)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Data["valid"].(bool) {
		t.Fatalf("expected error validating hmac")
	}

	// Set min decryption version, attempt to verify
	p.MinDecryptionVersion = 2
	if err = p.Persist(storage); err != nil {
		t.Fatal(err)
	}

	req.Data["hmac"] = "vault:v1:UcBvm5VskkukzZHlPgm3p5P/Yr/PV6xpuOGZISya3A4="
	resp, err = b.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected an error, got response %#v", resp)
	}
	if err != logical.ErrInvalidRequest {
		t.Fatalf("expected invalid request error, got %v", err)
	}
}
