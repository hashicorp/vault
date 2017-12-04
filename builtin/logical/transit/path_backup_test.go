package transit

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_BackupRestore(t *testing.T) {
	// Test encryption/decryption after a restore for supported keys
	testBackupRestore(t, "aes256-gcm96", "encrypt-decrypt")
	testBackupRestore(t, "rsa-2048", "encrypt-decrypt")
	testBackupRestore(t, "rsa-4096", "encrypt-decrypt")

	// Test signing/verification after a restore for supported keys
	testBackupRestore(t, "ecdsa-p256", "sign-verify")
	testBackupRestore(t, "ed25519", "sign-verify")
	testBackupRestore(t, "rsa-2048", "sign-verify")
	testBackupRestore(t, "rsa-4096", "sign-verify")

	// Test HMAC/verification after a restore for all key types
	testBackupRestore(t, "aes256-gcm96", "hmac-verify")
	testBackupRestore(t, "ecdsa-p256", "hmac-verify")
	testBackupRestore(t, "ed25519", "hmac-verify")
	testBackupRestore(t, "rsa-2048", "hmac-verify")
	testBackupRestore(t, "rsa-4096", "hmac-verify")
}

func testBackupRestore(t *testing.T, keyType, feature string) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	// Create a key
	keyReq := &logical.Request{
		Path:      "keys/test",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"type": keyType,
		},
	}
	resp, err = b.HandleRequest(keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	configReq := &logical.Request{
		Path:      "keys/test/config",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"deletion_allowed": true,
		},
	}
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	backupReq := &logical.Request{
		Path:      "backup/test",
		Operation: logical.ReadOperation,
		Storage:   s,
	}
	resp, err = b.HandleRequest(backupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}
	backup := resp.Data["backup"]

	restoreReq := &logical.Request{
		Path:      "restore",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"backup": backup,
		},
	}
	resp, err = b.HandleRequest(restoreReq)
	if resp != nil && resp.IsError() {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}
	if err == nil {
		t.Fatalf("expected an error")
	}

	plaintextB64 := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"

	var encryptReq, signReq, hmacReq *logical.Request
	var ciphertext, signature, hmac string
	switch feature {
	case "encrypt-decrypt":
		encryptReq = &logical.Request{
			Path:      "encrypt/test",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"plaintext": plaintextB64,
			},
		}
		resp, err = b.HandleRequest(encryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}
		ciphertext = resp.Data["ciphertext"].(string)
	case "sign-verify":
		signReq = &logical.Request{
			Path:      "sign/test",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"input": plaintextB64,
			},
		}
		resp, err = b.HandleRequest(signReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}
		signature = resp.Data["signature"].(string)

	case "hmac-verify":
		hmacReq = &logical.Request{
			Path:      "hmac/test",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"input": plaintextB64,
			},
		}
		resp, err = b.HandleRequest(hmacReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}
		hmac = resp.Data["hmac"].(string)
	}

	keyReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	resp, err = b.HandleRequest(restoreReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	validationFunc := func(keyName string) {
		var decryptReq *logical.Request
		var verifyReq *logical.Request
		switch feature {
		case "encrypt-decrypt":
			decryptReq = &logical.Request{
				Path:      "decrypt/" + keyName,
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"ciphertext": ciphertext,
				},
			}
			resp, err = b.HandleRequest(decryptReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("resp: %#v\nerr: %v", resp, err)
			}

			if resp.Data["plaintext"].(string) != plaintextB64 {
				t.Fatalf("bad: plaintext; expected: %q, actual: %q", plaintextB64, resp.Data["plaintext"].(string))
			}
		case "sign-verify":
			verifyReq = &logical.Request{
				Path:      "verify/" + keyName,
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"signature": signature,
					"input":     plaintextB64,
				},
			}
			resp, err = b.HandleRequest(verifyReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("resp: %#v\nerr: %v", resp, err)
			}
			if resp.Data["valid"].(bool) != true {
				t.Fatalf("bad: signature verification failed for key type %q", keyType)
			}

		case "hmac-verify":
			verifyReq = &logical.Request{
				Path:      "verify/" + keyName,
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"hmac":  hmac,
					"input": plaintextB64,
				},
			}
			resp, err = b.HandleRequest(verifyReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("resp: %#v\nerr: %v", resp, err)
			}
			if resp.Data["valid"].(bool) != true {
				t.Fatalf("bad: HMAC verification failed for key type %q", keyType)
			}
		}
	}

	validationFunc("test")

	resp, err = b.HandleRequest(keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	restoreReq.Path = "restore/test1"
	resp, err = b.HandleRequest(restoreReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	validationFunc("test1")
}
