package transit

import (
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_ConfigSettings(t *testing.T) {
	var b *backend
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}

	b = Backend(&logical.BackendConfig{
		StorageView: storage,
		System:      sysView,
	})

	doReq := func(req *logical.Request) *logical.Response {
		resp, err := b.HandleRequest(req)
		if err != nil {
			t.Fatalf("got err:\n%#v\nreq:\n%#v\n", err, *req)
		}
		return resp
	}
	doErrReq := func(req *logical.Request) {
		resp, err := b.HandleRequest(req)
		if err == nil {
			if resp == nil || !resp.IsError() {
				t.Fatalf("expected error; req:\n%#v\n", *req)
			}
		}
	}

	// First create a key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/aes",
		Data: map[string]interface{}{
			"derived": true,
		},
	}
	doReq(req)

	req.Path = "keys/ed"
	req.Data["type"] = "ed25519"
	doReq(req)

	delete(req.Data, "derived")

	req.Path = "keys/p256"
	req.Data["type"] = "ecdsa-p256"
	doReq(req)

	delete(req.Data, "type")

	req.Path = "keys/aes/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/ed/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/p256/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/aes/config"
	// Too high
	req.Data["min_decryption_version"] = 7
	doErrReq(req)
	// Too low
	req.Data["min_decryption_version"] = -1
	doErrReq(req)

	delete(req.Data, "min_decryption_version")
	// Too high
	req.Data["min_encryption_version"] = 7
	doErrReq(req)
	// Too low
	req.Data["min_encryption_version"] = 7
	doErrReq(req)

	// Not allowed, cannot decrypt
	req.Data["min_decryption_version"] = 3
	req.Data["min_encryption_version"] = 2
	doErrReq(req)

	// Allowed
	req.Data["min_decryption_version"] = 2
	req.Data["min_encryption_version"] = 3
	doReq(req)
	req.Path = "keys/ed/config"
	doReq(req)
	req.Path = "keys/p256/config"
	doReq(req)

	req.Data = map[string]interface{}{
		"plaintext": "abcd",
		"context":   "abcd",
	}

	maxKeyVersion := 5
	key := "aes"

	testHMAC := func(ver int, valid bool) {
		req.Path = "hmac/" + key
		delete(req.Data, "hmac")
		if ver == maxKeyVersion {
			delete(req.Data, "key_version")
		} else {
			req.Data["key_version"] = ver
		}

		if !valid {
			doErrReq(req)
			return
		}

		resp := doReq(req)
		ct := resp.Data["hmac"].(string)
		if strings.Split(ct, ":")[1] != "v"+strconv.Itoa(ver) {
			t.Fatal("wrong hmac version")
		}

		req.Path = "verify/" + key
		delete(req.Data, "key_version")
		req.Data["hmac"] = resp.Data["hmac"]
		doReq(req)
	}

	testEncryptDecrypt := func(ver int, valid bool) {
		req.Path = "encrypt/" + key
		delete(req.Data, "ciphertext")
		if ver == maxKeyVersion {
			delete(req.Data, "key_version")
		} else {
			req.Data["key_version"] = ver
		}

		if !valid {
			doErrReq(req)
			return
		}

		resp := doReq(req)
		ct := resp.Data["ciphertext"].(string)
		if strings.Split(ct, ":")[1] != "v"+strconv.Itoa(ver) {
			t.Fatal("wrong encryption version")
		}

		req.Path = "decrypt/" + key
		delete(req.Data, "key_version")
		req.Data["ciphertext"] = resp.Data["ciphertext"]
		doReq(req)
	}
	testEncryptDecrypt(5, true)
	testEncryptDecrypt(4, true)
	testEncryptDecrypt(3, true)
	testEncryptDecrypt(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	delete(req.Data, "plaintext")
	req.Data["input"] = "abcd"
	key = "ed"
	testSignVerify := func(ver int, valid bool) {
		req.Path = "sign/" + key
		delete(req.Data, "signature")
		if ver == maxKeyVersion {
			delete(req.Data, "key_version")
		} else {
			req.Data["key_version"] = ver
		}

		if !valid {
			doErrReq(req)
			return
		}

		resp := doReq(req)
		ct := resp.Data["signature"].(string)
		if strings.Split(ct, ":")[1] != "v"+strconv.Itoa(ver) {
			t.Fatal("wrong signature version")
		}

		req.Path = "verify/" + key
		delete(req.Data, "key_version")
		req.Data["signature"] = resp.Data["signature"]
		doReq(req)
	}
	testSignVerify(5, true)
	testSignVerify(4, true)
	testSignVerify(3, true)
	testSignVerify(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	delete(req.Data, "context")
	key = "p256"
	testSignVerify(5, true)
	testSignVerify(4, true)
	testSignVerify(3, true)
	testSignVerify(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)
}
