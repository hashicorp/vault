package transit

import (
	"encoding/base64"
	"strings"
	"testing"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

func TestTransit_SignVerify_P256(t *testing.T) {
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
		Data: map[string]interface{}{
			"type": "ecdsa-p256",
		},
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

	keyEntry := p.Keys[p.LatestVersion]
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
	p.Keys[p.LatestVersion] = keyEntry
	if err = p.Persist(storage); err != nil {
		t.Fatal(err)
	}
	req.Data = map[string]interface{}{
		"input": "dGhlIHF1aWNrIGJyb3duIGZveA==",
	}

	signRequest := func(req *logical.Request, errExpected bool, postpath string) string {
		req.Path = "sign/foo" + postpath
		resp, err := b.HandleRequest(req)
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
		req.Path = "verify/foo" + postpath
		req.Data["signature"] = sig
		resp, err := b.HandleRequest(req)
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
	req.Data["algorithm"] = "sha2-224"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["algorithm"] = "sha2-384"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["prehashed"] = true
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)
	delete(req.Data, "prehashed")

	// Test 512 and save sig for later to ensure we can't validate once min
	// decryption version is set
	req.Data["algorithm"] = "sha2-512"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	v1sig := sig

	// Test bad algorithm
	req.Data["algorithm"] = "foobar"
	signRequest(req, true, "")

	// Test bad input
	req.Data["algorithm"] = "sha2-256"
	req.Data["input"] = "foobar"
	signRequest(req, true, "")

	// Rotate and set min decryption version
	err = p.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}

	p.MinDecryptionVersion = 2
	if err = p.Persist(storage); err != nil {
		t.Fatal(err)
	}

	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.Data["algorithm"] = "sha2-256"
	// Make sure signing still works fine
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)
	// Now try the v1
	verifyRequest(req, true, "", v1sig)
}

func TestTransit_SignVerify_ED25519(t *testing.T) {
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
		Data: map[string]interface{}{
			"type": "ed25519",
		},
	}
	_, err := b.HandleRequest(req)
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
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	// Get the keys for later
	fooP, lock, err := b.lm.GetPolicyShared(storage, "foo")
	if err != nil {
		t.Fatal(err)
	}
	// We don't care as we're the only one using this
	lock.RUnlock()

	barP, lock, err := b.lm.GetPolicyShared(storage, "bar")
	if err != nil {
		t.Fatal(err)
	}
	lock.RUnlock()

	signRequest := func(req *logical.Request, errExpected bool, postpath string) string {
		// Delete any key that exists in the request
		delete(req.Data, "public_key")
		req.Path = "sign/" + postpath
		resp, err := b.HandleRequest(req)
		if err != nil && !errExpected {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
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
		// memoize any pubic key
		if key, ok := resp.Data["public_key"]; ok {
			req.Data["public_key"] = key
		}
		return value.(string)
	}

	verifyRequest := func(req *logical.Request, errExpected bool, postpath, sig string) {
		req.Path = "verify/" + postpath
		req.Data["signature"] = sig
		resp, err := b.HandleRequest(req)
		if err != nil && !errExpected {
			t.Fatalf("got error: %v, sig was %v", err, sig)
		}
		if errExpected {
			if resp != nil && !resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
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

		if pubKeyRaw, ok := req.Data["public_key"]; ok {
			input, _ := base64.StdEncoding.DecodeString(req.Data["input"].(string))
			splitSig := strings.Split(sig, ":")
			signature, _ := base64.StdEncoding.DecodeString(splitSig[2])
			if !ed25519.Verify(ed25519.PublicKey(pubKeyRaw.([]byte)), input, signature) && !errExpected {
				t.Fatal("invalid signature")
			}

			keyReadReq := &logical.Request{
				Operation: logical.ReadOperation,
				Path:      "keys/" + postpath,
			}
			keyReadResp, err := b.HandleRequest(keyReadReq)
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
			keyReadResp, err = b.HandleRequest(keyReadReq)
			if err != nil {
				t.Fatal(err)
			}
			val = keyReadResp.Data["keys"].(map[string]map[string]interface{})[strings.TrimPrefix(splitSig[1], "v")]
			if err := mapstructure.Decode(val, &ak); err != nil {
				t.Fatal(err)
			}
			if ak.PublicKey != base64.StdEncoding.EncodeToString(pubKeyRaw.([]byte)) {
				t.Fatalf("got incorrect public key; got %q, expected %q\nasymKey struct is\n%#v", ak.PublicKey, pubKeyRaw, ak)
			}
		}
	}

	req.Data = map[string]interface{}{
		"input":   "dGhlIHF1aWNrIGJyb3duIGZveA==",
		"context": "abcd",
	}

	// Test defaults
	sig := signRequest(req, false, "foo")
	verifyRequest(req, false, "foo", sig)

	sig = signRequest(req, false, "bar")
	verifyRequest(req, false, "bar", sig)

	// Test a bad signature
	verifyRequest(req, true, "foo", sig[0:len(sig)-2])
	verifyRequest(req, true, "bar", sig[0:len(sig)-2])

	v1sig := sig

	// Rotate and set min decryption version
	err = fooP.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}
	err = fooP.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}
	fooP.MinDecryptionVersion = 2
	if err = fooP.Persist(storage); err != nil {
		t.Fatal(err)
	}
	err = barP.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}
	err = barP.Rotate(storage)
	if err != nil {
		t.Fatal(err)
	}
	barP.MinDecryptionVersion = 2
	if err = barP.Persist(storage); err != nil {
		t.Fatal(err)
	}

	// Make sure signing still works fine
	sig = signRequest(req, false, "foo")
	verifyRequest(req, false, "foo", sig)
	// Now try the v1
	verifyRequest(req, true, "foo", v1sig)
	// Repeat with the other key
	sig = signRequest(req, false, "bar")
	verifyRequest(req, false, "bar", sig)
	verifyRequest(req, true, "bar", v1sig)
}
