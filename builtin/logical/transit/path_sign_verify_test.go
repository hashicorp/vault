package transit

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

// The outcome of processing a request includes
// the possibility that the request is incomplete or incorrect,
// or that the request is well-formed but the signature (for verification)
// is invalid, or that the signature is valid, but the key is not.
type signOutcome struct {
	requestOk bool
	valid     bool
	keyValid  bool
}

func TestTransit_SignVerify_P256(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	// First create a key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo",
		Data: map[string]interface{}{
			"type": "ecdsa-p256",
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Now, change the key value to something we control
	p, _, err := b.lm.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Useful code to output a key for openssl verification
	/*
		{
			key := p.Keys[p.LatestVersion]
			keyBytes, _ := x509.MarshalECPrivateKey(&ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     key.X,
					Y:     key.Y,
				},
				D: key.D,
			})
			pemBlock := &pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: keyBytes,
			}
			pemBytes := pem.EncodeToMemory(pemBlock)
			t.Fatalf("X: %s, Y: %s, D: %s, marshaled: %s", key.X.Text(16), key.Y.Text(16), key.D.Text(16), string(pemBytes))
		}
	*/

	keyEntry := p.Keys[strconv.Itoa(p.LatestVersion)]
	_, ok := keyEntry.EC_X.SetString("7336010a6da5935113d26d9ea4bb61b3b8d102c9a8083ed432f9b58fd7e80686", 16)
	if !ok {
		t.Fatal("could not set X")
	}
	_, ok = keyEntry.EC_Y.SetString("4040aa31864691a8a9e7e3ec9250e85425b797ad7be34ba8df62bfbad45ebb0e", 16)
	if !ok {
		t.Fatal("could not set Y")
	}
	_, ok = keyEntry.EC_D.SetString("99e5569be8683a2691dfc560ca9dfa71e887867a3af60635a08a3e3655aba3ef", 16)
	if !ok {
		t.Fatal("could not set D")
	}
	p.Keys[strconv.Itoa(p.LatestVersion)] = keyEntry
	if err = p.Persist(context.Background(), storage); err != nil {
		t.Fatal(err)
	}
	req.Data = map[string]interface{}{
		"input": "dGhlIHF1aWNrIGJyb3duIGZveA==",
	}

	signRequest := func(req *logical.Request, errExpected bool, postpath string) string {
		t.Helper()
		req.Path = "sign/foo" + postpath
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !resp.IsError() {
				t.Fatalf("bad: should have gotten error response: %#v", *resp)
			}
			return ""
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		value, ok := resp.Data["signature"]
		if !ok {
			t.Fatalf("no signature key found in returned data, got resp data %#v", resp.Data)
		}
		return value.(string)
	}

	verifyRequest := func(req *logical.Request, errExpected bool, postpath, sig string) {
		t.Helper()
		req.Path = "verify/foo" + postpath
		req.Data["signature"] = sig
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatalf("got error: %v, sig was %v", err, sig)
		}
		if errExpected {
			if resp != nil && !resp.IsError() {
				t.Fatalf("bad: should have gotten error response: %#v", *resp)
			}
			return
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		value, ok := resp.Data["valid"]
		if !ok {
			t.Fatalf("no valid key found in returned data, got resp data %#v", resp.Data)
		}
		if !value.(bool) && !errExpected {
			t.Fatalf("verification failed; req was %#v, resp is %#v", *req, *resp)
		}
	}

	// Comparisons are against values generated via openssl

	// Test defaults -- sha2-256
	sig := signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	// Test a bad signature
	verifyRequest(req, true, "", sig[0:len(sig)-2])

	// Test a signature generated with the same key by openssl
	sig = `vault:v1:MEUCIAgnEl9V8P305EBAlz68Nq4jZng5fE8k6MactcnlUw9dAiEAvJVePg3dazW6MaW7lRAVtEz82QJDVmR98tXCl8Pc7DA=`
	verifyRequest(req, false, "", sig)

	// Test algorithm selection in the path
	sig = signRequest(req, false, "/sha2-224")
	verifyRequest(req, false, "/sha2-224", sig)

	// Reset and test algorithm selection in the data
	req.Data["hash_algorithm"] = "sha2-224"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["hash_algorithm"] = "sha2-384"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["prehashed"] = true
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)
	delete(req.Data, "prehashed")

	// Test marshaling selection
	// Bad value
	req.Data["marshaling_algorithm"] = "asn2"
	sig = signRequest(req, true, "")
	// Use the default, verify we can't validate with jws
	req.Data["marshaling_algorithm"] = "asn1"
	sig = signRequest(req, false, "")
	req.Data["marshaling_algorithm"] = "jws"
	verifyRequest(req, true, "", sig)
	// Sign with jws, verify we can validate
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)
	// If we change marshaling back to asn1 we shouldn't be able to verify
	delete(req.Data, "marshaling_algorithm")
	verifyRequest(req, true, "", sig)

	// Test 512 and save sig for later to ensure we can't validate once min
	// decryption version is set
	req.Data["hash_algorithm"] = "sha2-512"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	v1sig := sig

	// Test bad algorithm
	req.Data["hash_algorithm"] = "foobar"
	signRequest(req, true, "")

	// Test bad input
	req.Data["hash_algorithm"] = "sha2-256"
	req.Data["input"] = "foobar"
	signRequest(req, true, "")

	// Rotate and set min decryption version
	err = p.Rotate(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Rotate(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}

	p.MinDecryptionVersion = 2
	if err = p.Persist(context.Background(), storage); err != nil {
		t.Fatal(err)
	}

	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.Data["hash_algorithm"] = "sha2-256"
	// Make sure signing still works fine
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)
	// Now try the v1
	verifyRequest(req, true, "", v1sig)
}

func validatePublicKey(t *testing.T, in string, sig string, pubKeyRaw []byte, expectValid bool, postpath string, b *backend) {
	t.Helper()
	input, _ := base64.StdEncoding.DecodeString(in)
	splitSig := strings.Split(sig, ":")
	signature, _ := base64.StdEncoding.DecodeString(splitSig[2])
	valid := ed25519.Verify(ed25519.PublicKey(pubKeyRaw), input, signature)
	if valid != expectValid {
		t.Fatalf("status of signature: expected %v. Got %v", valid, expectValid)
	}
	if !valid {
		return
	}

	keyReadReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "keys/" + postpath,
	}
	keyReadResp, err := b.HandleRequest(context.Background(), keyReadReq)
	if err != nil {
		t.Fatal(err)
	}
	val := keyReadResp.Data["keys"].(map[string]map[string]interface{})[strings.TrimPrefix(splitSig[1], "v")]
	var ak asymKey
	if err := mapstructure.Decode(val, &ak); err != nil {
		t.Fatal(err)
	}
	if ak.PublicKey != "" {
		t.Fatal("got non-empty public key")
	}
	keyReadReq.Data = map[string]interface{}{
		"context": "abcd",
	}
	keyReadResp, err = b.HandleRequest(context.Background(), keyReadReq)
	if err != nil {
		t.Fatal(err)
	}
	val = keyReadResp.Data["keys"].(map[string]map[string]interface{})[strings.TrimPrefix(splitSig[1], "v")]
	if err := mapstructure.Decode(val, &ak); err != nil {
		t.Fatal(err)
	}
	if ak.PublicKey != base64.StdEncoding.EncodeToString(pubKeyRaw) {
		t.Fatalf("got incorrect public key; got %q, expected %q\nasymKey struct is\n%#v", ak.PublicKey, pubKeyRaw, ak)
	}
}

func TestTransit_SignVerify_ED25519(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	// First create a key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo",
		Data: map[string]interface{}{
			"type": "ed25519",
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Now create a derived key"
	req = &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/bar",
		Data: map[string]interface{}{
			"type":    "ed25519",
			"derived": true,
		},
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Get the keys for later
	fooP, _, err := b.lm.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	barP, _, err := b.lm.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	signRequest := func(req *logical.Request, errExpected bool, postpath string) []string {
		t.Helper()
		// Delete any key that exists in the request
		delete(req.Data, "public_key")
		req.Path = "sign/" + postpath
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil {
			if !errExpected {
				t.Fatal(err)
			}
			return nil
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if resp.IsError() {
				return nil
			}
			t.Fatalf("bad: expected error response, got: %#v", *resp)
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		// memoize any pubic key
		if key, ok := resp.Data["public_key"]; ok {
			req.Data["public_key"] = key
		}
		// batch_input supplied
		if _, ok := req.Data["batch_input"]; ok {
			batchRequestItems := req.Data["batch_input"].([]batchRequestSignItem)

			batchResults, ok := resp.Data["batch_results"]
			if !ok {
				t.Fatalf("no batch_results in returned data, got resp data %#v", resp.Data)
			}
			batchResponseItems := batchResults.([]batchResponseSignItem)
			if len(batchResponseItems) != len(batchRequestItems) {
				t.Fatalf("Expected %d items in response. Got %d: %#v", len(batchRequestItems), len(batchResponseItems), resp)
			}
			if len(batchRequestItems) == 0 {
				return nil
			}
			ret := make([]string, len(batchRequestItems))
			for i, v := range batchResponseItems {
				ret[i] = v.Signature
			}
			return ret
		}

		// input supplied
		value, ok := resp.Data["signature"]
		if !ok {
			t.Fatalf("no signature key found in returned data, got resp data %#v", resp.Data)
		}
		return []string{value.(string)}
	}

	verifyRequest := func(req *logical.Request, errExpected bool, outcome []signOutcome, postpath string, sig []string, attachSig bool) {
		t.Helper()
		req.Path = "verify/" + postpath
		if _, ok := req.Data["batch_input"]; ok && attachSig {
			batchRequestItems := req.Data["batch_input"].([]batchRequestSignItem)
			if len(batchRequestItems) != len(sig) {
				t.Fatalf("number of requests in batch(%d) != number of signatures(%d)", len(batchRequestItems), len(sig))
			}
			for i, v := range sig {
				batchRequestItems[i]["signature"] = v
			}
		} else if attachSig {
			req.Data["signature"] = sig[0]
		}
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatalf("got error: %v, sig was %v", err, sig)
		}
		if errExpected {
			if resp != nil && !resp.IsError() {
				t.Fatalf("bad: expected error response, got: %#v\n%#v", *resp, req)
			}
			return
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}

		// batch_input field supplied
		if _, ok := req.Data["batch_input"]; ok {
			batchRequestItems := req.Data["batch_input"].([]batchRequestSignItem)

			batchResults, ok := resp.Data["batch_results"]
			if !ok {
				t.Fatalf("no batch_results in returned data, got resp data %#v", resp.Data)
			}
			batchResponseItems := batchResults.([]batchResponseVerifyItem)
			if len(batchResponseItems) != len(batchRequestItems) {
				t.Fatalf("Expected %d items in response. Got %d: %#v", len(batchRequestItems), len(batchResponseItems), resp)
			}
			if len(batchRequestItems) == 0 {
				return
			}
			for i, v := range batchResponseItems {
				if v.Error != "" && outcome[i].requestOk {
					t.Fatalf("verification failed; req was %#v, resp is %#v", *req, *resp)
				}
				if v.Error != "" {
					continue
				}
				if v.Valid != outcome[i].valid {
					t.Fatalf("verification failed; req was %#v, resp is %#v", *req, *resp)
				}
				if !v.Valid {
					continue
				}
				if pubKeyRaw, ok := req.Data["public_key"]; ok {
					validatePublicKey(t, batchRequestItems[i]["input"], sig[i], pubKeyRaw.([]byte), outcome[i].keyValid, postpath, b)
				}
			}
			return
		}

		// input field supplied
		value, ok := resp.Data["valid"]
		if !ok {
			t.Fatalf("no valid key found in returned data, got resp data %#v", resp.Data)
		}
		valid := value.(bool)
		if valid != outcome[0].valid {
			t.Fatalf("verification failed; req was %#v, resp is %#v", *req, *resp)
		}
		if !valid {
			return
		}

		if pubKeyRaw, ok := req.Data["public_key"]; ok {
			validatePublicKey(t, req.Data["input"].(string), sig[0], pubKeyRaw.([]byte), outcome[0].keyValid, postpath, b)
		}
	}

	req.Data = map[string]interface{}{
		"input":   "dGhlIHF1aWNrIGJyb3duIGZveA==",
		"context": "abcd",
	}

	outcome := []signOutcome{{requestOk: true, valid: true, keyValid: true}}
	// Test defaults
	sig := signRequest(req, false, "foo")
	verifyRequest(req, false, outcome, "foo", sig, true)

	sig = signRequest(req, false, "bar")
	verifyRequest(req, false, outcome, "bar", sig, true)

	// Verify with incorrect key
	outcome[0].valid = false
	verifyRequest(req, false, outcome, "foo", sig, true)

	// Verify with missing signatures
	delete(req.Data, "signature")
	verifyRequest(req, true, outcome, "foo", sig, false)

	// Test a bad signature
	badsig := sig[0]
	badsig = badsig[:len(badsig)-2]
	verifyRequest(req, true, outcome, "bar", []string{badsig}, true)

	v1sig := sig

	// Test a missing context
	delete(req.Data, "context")
	sig = signRequest(req, true, "bar")

	// Rotate and set min decryption version
	err = fooP.Rotate(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}
	err = fooP.Rotate(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}
	fooP.MinDecryptionVersion = 2
	if err = fooP.Persist(context.Background(), storage); err != nil {
		t.Fatal(err)
	}
	err = barP.Rotate(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}
	err = barP.Rotate(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}
	barP.MinDecryptionVersion = 2
	if err = barP.Persist(context.Background(), storage); err != nil {
		t.Fatal(err)
	}

	req.Data = map[string]interface{}{
		"input":   "dGhlIHF1aWNrIGJyb3duIGZveA==",
		"context": "abcd",
	}

	// Make sure signing still works fine
	sig = signRequest(req, false, "foo")
	outcome[0].valid = true
	verifyRequest(req, false, outcome, "foo", sig, true)
	// Now try the v1
	verifyRequest(req, true, outcome, "foo", v1sig, true)

	// Repeat with the other key
	sig = signRequest(req, false, "bar")
	verifyRequest(req, false, outcome, "bar", sig, true)
	verifyRequest(req, true, outcome, "bar", v1sig, true)

	// Test Batch Signing
	batchInput := []batchRequestSignItem{
		{"context": "abcd", "input": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
		{"context": "efgh", "input": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
	}

	req.Data = map[string]interface{}{
		"batch_input": batchInput,
	}

	outcome = []signOutcome{{requestOk: true, valid: true, keyValid: true}, {requestOk: true, valid: true, keyValid: true}}

	sig = signRequest(req, false, "foo")
	verifyRequest(req, false, outcome, "foo", sig, true)

	goodsig := signRequest(req, false, "bar")
	verifyRequest(req, false, outcome, "bar", goodsig, true)

	// key doesn't match signatures
	outcome[0].valid = false
	outcome[1].valid = false
	verifyRequest(req, false, outcome, "foo", goodsig, true)

	// Test a bad signature
	badsig = sig[0]
	badsig = badsig[:len(badsig)-2]
	// matching key, but first signature is corrupted
	outcome[0].requestOk = false
	outcome[1].valid = true
	verifyRequest(req, false, outcome, "bar", []string{badsig, goodsig[1]}, true)

	// Verify with missing signatures
	outcome[0].valid = false
	outcome[1].valid = false
	delete(batchInput[0], "signature")
	delete(batchInput[1], "signature")
	verifyRequest(req, true, outcome, "foo", sig, false)

	// Test missing context
	batchInput = []batchRequestSignItem{
		{"context": "abcd", "input": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
		{"input": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
	}

	req.Data = map[string]interface{}{
		"batch_input": batchInput,
	}

	sig = signRequest(req, false, "bar")

	outcome[0].requestOk = true
	outcome[0].valid = true
	outcome[1].requestOk = false
	verifyRequest(req, false, outcome, "bar", goodsig, true)

	// Test incorrect context
	batchInput = []batchRequestSignItem{
		{"context": "abca", "input": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
		{"context": "efga", "input": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
	}
	req.Data = map[string]interface{}{
		"batch_input": batchInput,
	}

	outcome[0].requestOk = true
	outcome[0].valid = false
	outcome[1].requestOk = true
	outcome[1].valid = false
	verifyRequest(req, false, outcome, "bar", goodsig, true)
}
