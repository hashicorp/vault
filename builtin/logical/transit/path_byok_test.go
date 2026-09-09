// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestTransit_BYOKExportImport(t *testing.T) {
	// Test encryption/decryption after a restore for supported keys
	testBYOKExportImport(t, "aes128-gcm96", "encrypt-decrypt")
	testBYOKExportImport(t, "aes256-gcm96", "encrypt-decrypt")
	testBYOKExportImport(t, "chacha20-poly1305", "encrypt-decrypt")
	testBYOKExportImport(t, "rsa-2048", "encrypt-decrypt")
	testBYOKExportImport(t, "rsa-3072", "encrypt-decrypt")
	testBYOKExportImport(t, "rsa-4096", "encrypt-decrypt")

	// Test signing/verification after a restore for supported keys
	testBYOKExportImport(t, "ecdsa-p256", "sign-verify")
	testBYOKExportImport(t, "ecdsa-p384", "sign-verify")
	testBYOKExportImport(t, "ecdsa-p521", "sign-verify")
	testBYOKExportImport(t, "ed25519", "sign-verify")
	testBYOKExportImport(t, "rsa-2048", "sign-verify")
	testBYOKExportImport(t, "rsa-3072", "sign-verify")
	testBYOKExportImport(t, "rsa-4096", "sign-verify")

	// Test HMAC sign/verify after a restore for supported keys.
	testBYOKExportImport(t, "hmac", "hmac-verify")
}

func testBYOKExportImport(t *testing.T, keyType, feature string) {
	var resp *logical.Response
	var err error

	b, s, obsRecorder := createBackendWithObservationRecorder(t)

	// Create a key
	keyReq := &logical.Request{
		Path:      "keys/test-source",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"type":       keyType,
			"exportable": true,
		},
	}
	if keyType == "hmac" {
		keyReq.Data["key_size"] = 32
	}
	resp, err = b.HandleRequest(context.Background(), keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	// Read the wrapping key.
	wrapKeyReq := &logical.Request{
		Path:      "wrapping_key",
		Operation: logical.ReadOperation,
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), wrapKeyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	// Import the wrapping key.
	wrapKeyImportReq := &logical.Request{
		Path:      "keys/wrapper/import",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"public_key": resp.Data["public_key"],
			"type":       "rsa-4096",
		},
	}
	resp, err = b.HandleRequest(context.Background(), wrapKeyImportReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	// Export the key
	backupReq := &logical.Request{
		Path:      "byok-export/wrapper/test-source",
		Operation: logical.ReadOperation,
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), backupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyExportBYOK)
	require.NotNil(t, obs)
	require.Equal(t, "test-source", obs.Data["key_name"])
	require.Equal(t, "wrapper", obs.Data["destination_key"])
	keys := resp.Data["keys"].(map[string]string)

	// Import the key to a new name.
	restoreReq := &logical.Request{
		Path:      "keys/test/import",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"ciphertext": keys["1"],
			"type":       keyType,
		},
	}
	resp, err = b.HandleRequest(context.Background(), restoreReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	plaintextB64 := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"
	// Perform encryption, signing or hmac-ing based on the set 'feature'
	var encryptReq, signReq, hmacReq *logical.Request
	var ciphertext, signature, hmac string
	switch feature {
	case "encrypt-decrypt":
		encryptReq = &logical.Request{
			Path:      "encrypt/test-source",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"plaintext": plaintextB64,
			},
		}
		resp, err = b.HandleRequest(context.Background(), encryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}
		ciphertext = resp.Data["ciphertext"].(string)

	case "sign-verify":
		signReq = &logical.Request{
			Path:      "sign/test-source",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"input": plaintextB64,
			},
		}
		resp, err = b.HandleRequest(context.Background(), signReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}
		signature = resp.Data["signature"].(string)

	case "hmac-verify":
		hmacReq = &logical.Request{
			Path:      "hmac/test-source",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"input": plaintextB64,
			},
		}
		resp, err = b.HandleRequest(context.Background(), hmacReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}
		hmac = resp.Data["hmac"].(string)
	}

	// validationFunc verifies the ciphertext, signature or hmac based on the
	// set 'feature'
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
			resp, err = b.HandleRequest(context.Background(), decryptReq)
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
			resp, err = b.HandleRequest(context.Background(), verifyReq)
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
			resp, err = b.HandleRequest(context.Background(), verifyReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("resp: %#v\nerr: %v", resp, err)
			}
			if resp.Data["valid"].(bool) != true {
				t.Fatalf("bad: HMAC verification failed for key type %q", keyType)
			}
		}
	}

	// Ensure that the restored key is functional
	validationFunc("test")

	// Ensure the original key is functional
	validationFunc("test-source")
}
