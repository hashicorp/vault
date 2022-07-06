package transit

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
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

func TestTransit_SignVerify_ECDSA(t *testing.T) {
	t.Run("256", func(t *testing.T) {
		testTransit_SignVerify_ECDSA(t, 256)
	})
	t.Run("384", func(t *testing.T) {
		testTransit_SignVerify_ECDSA(t, 384)
	})
	t.Run("521", func(t *testing.T) {
		testTransit_SignVerify_ECDSA(t, 521)
	})
}

func testTransit_SignVerify_ECDSA(t *testing.T, bits int) {
	b, storage := createBackendWithSysView(t)

	// First create a key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo",
		Data: map[string]interface{}{
			"type": fmt.Sprintf("ecdsa-p%d", bits),
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Now, change the key value to something we control
	p, _, err := b.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "foo",
	}, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}

	// Useful code to output a key for openssl verification
	/*
		if bits == 384 {
			var curve elliptic.Curve
			switch bits {
			case 521:
				curve = elliptic.P521()
			case 384:
				curve = elliptic.P384()
			default:
				curve = elliptic.P256()
			}
			key := p.Keys[strconv.Itoa(p.LatestVersion)]
			keyBytes, _ := x509.MarshalECPrivateKey(&ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: curve,
					X:     key.EC_X,
					Y:     key.EC_Y,
				},
				D: key.EC_D,
			})
			pemBlock := &pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: keyBytes,
			}
			pemBytes := pem.EncodeToMemory(pemBlock)
			t.Fatalf("X: %s, Y: %s, D: %s, marshaled: %s", key.EC_X.Text(16), key.EC_Y.Text(16), key.EC_D.Text(16), string(pemBytes))
		}
	*/

	var xString, yString, dString string
	switch bits {
	case 384:
		xString = "703457a84e48bfcb037cfb509f1870d2aa5b74c109c2f24624ab21444492575229f8711453e5c656dab596b4e26db30e"
		yString = "411c5b7092a893dc8b7af39de3d21d1c26f45b27616baeac4c479ef3c9f21c194b5ac501dee47ba2b2cb243a54256524"
		dString = "3de3e4fd2ecbc490e956f41f5003a1e57a84763cec7b722fa3427cf461a1148ea4d5206023bcce0422289f6633730759"
		/*
			-----BEGIN EC PRIVATE KEY-----
			MIGkAgEBBDA94+T9LsvEkOlW9B9QA6HleoR2POx7ci+jQnz0YaEUjqTVIGAjvM4E
			IiifZjNzB1mgBwYFK4EEACKhZANiAARwNFeoTki/ywN8+1CfGHDSqlt0wQnC8kYk
			qyFERJJXUin4cRRT5cZW2rWWtOJtsw5BHFtwkqiT3It6853j0h0cJvRbJ2FrrqxM
			R57zyfIcGUtaxQHe5HuissskOlQlZSQ=
			-----END EC PRIVATE KEY-----
		*/
	case 521:
		xString = "1913f75fc044fe5d1f871c2629a377462fd819b174a41d3ec7d04ebd5ae35475ff8de544f4e19a9aa6b16a8f67af479be6884e00ca3147dc24d5924d66ac395e04b"
		yString = "4919406b90d8323fdb5c9c4f48259c56ebcea37b40ad1a82bbbfad62a9b9c2dce515772274b84725471c7d0b7c62e10c23296b1a9d2b2586ada67735ff5d9fffc4"
		dString = "1867d0fcd9bac4c5821b70a6b13117499438f8c274579c0aba254fbd85fa98892c3608576197d5534366a9aab0f904155bec46d800d23a57f7f053d91526568b09"
		/*
			-----BEGIN EC PRIVATE KEY-----
			MIHcAgEBBEIAGGfQ/Nm6xMWCG3CmsTEXSZQ4+MJ0V5wKuiVPvYX6mIksNghXYZfV
			U0Nmqaqw+QQVW+xG2ADSOlf38FPZFSZWiwmgBwYFK4EEACOhgYkDgYYABAGRP3X8
			BE/l0fhxwmKaN3Ri/YGbF0pB0+x9BOvVrjVHX/jeVE9OGamqaxao9nr0eb5ohOAM
			oxR9wk1ZJNZqw5XgSwBJGUBrkNgyP9tcnE9IJZxW686je0CtGoK7v61iqbnC3OUV
			dyJ0uEclRxx9C3xi4QwjKWsanSslhq2mdzX/XZ//xA==
			-----END EC PRIVATE KEY-----
		*/
	default:
		xString = "7336010a6da5935113d26d9ea4bb61b3b8d102c9a8083ed432f9b58fd7e80686"
		yString = "4040aa31864691a8a9e7e3ec9250e85425b797ad7be34ba8df62bfbad45ebb0e"
		dString = "99e5569be8683a2691dfc560ca9dfa71e887867a3af60635a08a3e3655aba3ef"
	}

	keyEntry := p.Keys[strconv.Itoa(p.LatestVersion)]
	_, ok := keyEntry.EC_X.SetString(xString, 16)
	if !ok {
		t.Fatal("could not set X")
	}
	_, ok = keyEntry.EC_Y.SetString(yString, 16)
	if !ok {
		t.Fatal("could not set Y")
	}
	_, ok = keyEntry.EC_D.SetString(dString, 16)
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
		if err != nil {
			if errExpected {
				return
			}
			t.Fatalf("got error: %v, sig was %v", err, sig)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.IsError() {
			if errExpected {
				return
			}
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		value, ok := resp.Data["valid"]
		if !ok {
			t.Fatalf("no valid key found in returned data, got resp data %#v", resp.Data)
		}
		if !value.(bool) && !errExpected {
			t.Fatalf("verification failed; req was %#v, resp is %#v", *req, *resp)
		} else if value.(bool) && errExpected {
			t.Fatalf("expected error and didn't get one; req was %#v, resp is %#v", *req, *resp)
		}
	}

	// Comparisons are against values generated via openssl

	// Test defaults -- sha2-256
	sig := signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	// Test a bad signature
	verifyRequest(req, true, "", sig[0:len(sig)-2])

	// Test a signature generated with the same key by openssl
	switch bits {
	case 384:
		sig = `vault:v1:MGUCMHHZLRN/3ehWuWACfSCMLtFtNEAdx6Rkwon2Lx6FWCyXCXqH6A8Pz8er0Qkgvm2ElQIxAO922LmUeYzHmDSfC5is/TjFu3b4Fb+1XtoBXncc2u4t2vSuTAxEv7WMh2D2YDdxeA==`
	case 521:
		sig = `vault:v1:MIGIAkIBYhspOgSs/K/NUWtlBN+CfYe1IVFpUbQNSqdjT7s+QKcr6GKmdGLIQAXw0q6K0elBgzi1wgLjxwdscwMeW7tm/QQCQgDzdITGlUEd9Z7DOfLCnDP4X8pGsfO60Tvsh/BN44drZsHLtXYBXLczB/XZfIWAsPMuI5F7ExwVNbmQP0FBVri/QQ==`
	default:
		sig = `vault:v1:MEUCIAgnEl9V8P305EBAlz68Nq4jZng5fE8k6MactcnlUw9dAiEAvJVePg3dazW6MaW7lRAVtEz82QJDVmR98tXCl8Pc7DA=`
	}
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

	req.Data["hash_algorithm"] = "sha2-512"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["hash_algorithm"] = "sha3-224"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["hash_algorithm"] = "sha3-256"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["hash_algorithm"] = "sha3-384"
	sig = signRequest(req, false, "")
	verifyRequest(req, false, "", sig)

	req.Data["hash_algorithm"] = "sha3-512"
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
	err = p.Rotate(context.Background(), storage, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}
	err = p.Rotate(context.Background(), storage, b.GetRandomReader())
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
	fooP, _, err := b.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "foo",
	}, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}

	barP, _, err := b.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "bar",
	}, b.GetRandomReader())
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
	err = fooP.Rotate(context.Background(), storage, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}
	err = fooP.Rotate(context.Background(), storage, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}
	fooP.MinDecryptionVersion = 2
	if err = fooP.Persist(context.Background(), storage); err != nil {
		t.Fatal(err)
	}
	err = barP.Rotate(context.Background(), storage, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}
	err = barP.Rotate(context.Background(), storage, b.GetRandomReader())
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

func TestTransit_SignVerify_RSAPSS(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	//create RSA Key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo",
		Data: map[string]interface{}{
			"type": "rsa-3072",
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := b.GetPolicy(context.Background(), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "foo",
	}, b.GetRandomReader())
	if err != nil {
		t.Fatal(err)
	}

	// Set the RSA key to a known value
	keyEntry := p.Keys[strconv.Itoa(p.LatestVersion)]
	d := &big.Int{}
	d.UnmarshalText([]byte("1202968648230909168045620068417760275295456979639745012353765236549290010690691905232434050345715329321148006148161911242113887010813816466646948037232556705436709826907988822128130749416941720054147815158540411664237526359829784722708437312431113835967495963181332704211806465233209363666530087956398936068027950806184656645868580282142708603503463009125869540065719450773247999180553084327731820344258433172141828990096131150165738737489869763947541855698613277701043218701250253493792267496965291335622150101989283719195285427737832201644492380250469887304294943179455031543357697328739719758686609441455951441833631993398836568668540418920917619195944820633703763280429714796664544305449572971742629712279365509653320628674969017623395860862203981513701363328423779752796189206869111835384413950280733911285492340511256734352674630934319556427684747862695675766069427830877256824237554222845196827379299919432194942144473"))
	if err != nil {
		t.Fatal(err)
	}
	prime1 := &big.Int{}
	err = prime1.UnmarshalText([]byte("2055662681415640757527037964786977991821510354423402187146841519469586519238092913775083064402486556051237843029367172508018550980898416475664728551997413582059604037474338986870339870718709679882961395993989214128499165904389954550732462092150776361615542970227027371350074370268495466991086366718529590405085824584855391251864782854458716354893872754315376044606226143266734976714763503600335290996275953717131880996301086793047310442708072514030450850488551751"))
	if err != nil {
		t.Fatal(err)
	}
	prime2 := &big.Int{}
	err = prime2.UnmarshalText([]byte("2155338210536045375345461329945095769089640047457940302875938444864355576790859924695523433820262936343393526707386360465750814035695569282097499936405507048127441171817791800992189812299302833657366046017230614262203396538000510758997141512210577333647931181430496330608397833397098973276033750722515264792258002587897982676420464615298870584510171052786328379200394276816247444638334698383524148609230657368382216986985533157927862442013874101932275115356587027"))
	if err != nil {
		t.Fatal(err)
	}

	N := &big.Int{}
	err = N.UnmarshalText([]byte("4430648325228115889974474678200222275038685178973247660707469501389840363641445172148928310526421576807917100086101336241115983647729857804576881730701925862886683821853931855446291161320563645452887791505297457527207753458703192164333081721074401903327064344105597473076832657749176242925333391839862710637875003483473296762970054285196509708205375701308424864970611197332041931117000533190207952675152575857179894713045416836484883648357513471947288445370350593212533398640673311877931678397838305041906210088270200886267736677296839706550629699361980127001576113365321878051556656462048787180983865155995976862145691440916151868326704094676484075067961203015485370448193132673395545691929211260967925112935845280679371059652028613973186835951092078768057547173601361986909326965500543006815882540522111841040396486895287727002078916720900659021612983110731417503963681431384865409248735458853095156697497092310531424734277"))
	if err != nil {
		t.Fatal(err)
	}
	rsaKeyEntry := rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: N,
			E: 65537,
		},
		Primes:      []*big.Int{prime1, prime2},
		D:           d,
		Precomputed: rsa.PrecomputedValues{},
	}
	rsaKeyEntry.Precompute()
	keyEntry.RSAKey = &rsaKeyEntry
	keyEntry.RSAKey.Precompute()
	p.Keys[strconv.Itoa(p.LatestVersion)] = keyEntry

	if err = p.Persist(context.Background(), storage); err != nil {
		t.Fatal(err)
	}

	//"the quick brown fox"
	inputDataB64 := "dGhlIHF1aWNrIGJyb3duIGZveA=="

	req.Data = map[string]interface{}{
		"input": inputDataB64,
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
		if err != nil {
			if errExpected {
				return
			}
			t.Fatalf("got error: %v, sig was %v", err, sig)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.IsError() {
			if errExpected {
				return
			}
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		value, ok := resp.Data["valid"]
		if !ok {
			t.Fatalf("no valid key found in returned data, got resp data %#v", resp.Data)
		}
		if !value.(bool) && !errExpected {
			t.Fatalf("verification failed; req was %#v, resp is %#v", *req, *resp)
		} else if value.(bool) && errExpected {
			t.Fatalf("expected error and didn't get one; req was %#v, resp is %#v", *req, *resp)
		}
	}
	//get a signature with default settings
	sig := signRequest(req, false, "/sha2-256")
	verifyRequest(req, false, "", sig)

	// Test a bad signature
	verifyRequest(req, true, "/sha2-256", sig[0:len(sig)-2])

	//check if a new signature is different (because of the salt)
	oldSig := sig
	sig = signRequest(req, false, "/sha2-256")

	if oldSig == sig {
		t.Errorf("Two PSS signatures should be different\n")
	}

	//now do some pkcs1v15 signatures
	req.Data["signature_algorithm"] = "pkcs1v15"

	sig = signRequest(req, false, "/sha2-256")
	verifyRequest(req, false, "", sig)
	oldSig = sig
	if oldSig != sig {
		t.Errorf("Two RSA PKCS1#1.5 signatures should be equal\n")
	}
	//verify PSS locally - focus on salt length
	req.Data["signature_algorithm"] = "pss"
	req.Data["pss_salt_length"] = "0" //default but explicitly set here

	sig = signRequest(req, false, "/sha2-256")

	inputRaw, err := base64.StdEncoding.DecodeString(inputDataB64)
	if err != nil {
		t.Error(err)
	}
	inputHash := sha256.Sum256(inputRaw)

	rawSig := getRawSig(sig)

	//verify with PSSOptions == nil
	err = rsa.VerifyPSS(&rsaKeyEntry.PublicKey, crypto.SHA256, inputHash[:], rawSig, nil)
	if err != nil {
		t.Errorf("PSS verification should succeed but error %v\n", err)
	}

	//verify with a salt length of "0"
	err = rsa.VerifyPSS(&rsaKeyEntry.PublicKey, crypto.SHA256, inputHash[:], rawSig, &rsa.PSSOptions{SaltLength: 0})
	if err != nil {
		t.Errorf("PSS verification should succeed but error %v\n", err)
	}
	//verify with a salt length of "350" which should be calculated by go stdlib
	err = rsa.VerifyPSS(&rsaKeyEntry.PublicKey, crypto.SHA256, inputHash[:], rawSig, &rsa.PSSOptions{SaltLength: 350})
	if err != nil {
		t.Errorf("PSS verification should succeed but error %v\n", err)
	}
	//verify with a salt length of "200" which should be wrong and error
	err = rsa.VerifyPSS(&rsaKeyEntry.PublicKey, crypto.SHA256, inputHash[:], rawSig, &rsa.PSSOptions{SaltLength: 200})
	if err == nil {
		t.Errorf("PSS verification should not succeed, but did\n")
	}

	// This does not work because it returns a "nil" response with an error attached
	// The underlying cause is that the crypto/rsa Package just returns a
	// errors.New("crypto/rsa: key size too small for PSS signature")
	// for a salt that's to big. That is not really a "testable" error
	// Alternatives: Salt Size Check in helpers/keysutil/policy.go
	// Direct check for the erorr string "crypto/rsa: key size too small for PSS signature"
	// in builtin/logical/transit/path_sign_verify.go

	//sign with an explicit salt length of 1000 which should not work (Salt is too large)
	//req.Data["pss_salt_length"] = "1000"
	//_ = signRequest(req, true, "/sha2-256")

	// common value used, Salt Length == Hash Length, see
	// i.e. TLS 1.3 specification
	// https://datatracker.ietf.org/doc/html/rfc8446 page 43
	req.Data["pss_salt_length"] = "32"

	sig = signRequest(req, false, "/sha2-256")
	verifyRequest(req, false, "/sha2-256", sig)

	//let crypto/rsa verify
	rawSig = getRawSig(sig)
	//verify with PSSOptions == nil - should error
	err = rsa.VerifyPSS(&rsaKeyEntry.PublicKey, crypto.SHA256, inputHash[:], rawSig, &rsa.PSSOptions{SaltLength: 34})
	if err == nil {
		t.Errorf("PSS verification should error beause of salt length mismatch\n")
	}
	//verify with saltLength = 32
	err = rsa.VerifyPSS(&rsaKeyEntry.PublicKey, crypto.SHA256, inputHash[:], rawSig, &rsa.PSSOptions{SaltLength: 32})
	if err != nil {
		t.Errorf("PSS verification should succeed but error %v\n", err)
	}

	//let vault verify
	verifyRequest(req, false, "/sha2-256", sig)
	delete(req.Data, "pss_salt_length")
	verifyRequest(req, false, "/sha2-256", sig)

}

func getRawSig(vaultSig string) []byte {
	raw, err := base64.StdEncoding.DecodeString(strings.Split(vaultSig, ":")[2])
	if err != nil {
		panic(err)
	}
	return raw
}
