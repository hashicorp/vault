package transit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_BackupRestore(t *testing.T) {
	// Test encryption/decryption after a restore for supported keys
	testEncryptDecryptBackupRestore(t, "aes256-gcm96")
	testEncryptDecryptBackupRestore(t, "rsa-2048")
	testEncryptDecryptBackupRestore(t, "rsa-4096")

	// Test signing/verification after a restore for supported keys
	testSignVerifyBackupRestore(t, "ecdsa-p256")
	testSignVerifyBackupRestore(t, "ed25519")
	testSignVerifyBackupRestore(t, "rsa-2048")
	testSignVerifyBackupRestore(t, "rsa-4096")

	// Test HMAC/verification after a restore for all key types
	testHMACBackupRestore(t, "aes256-gcm96")
	testHMACBackupRestore(t, "ecdsa-p256")
	testHMACBackupRestore(t, "ed25519")
	testHMACBackupRestore(t, "rsa-2048")
	testHMACBackupRestore(t, "rsa-4096")
}

func testEncryptDecryptBackupRestore(t *testing.T, keyType string) {
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

	encryptReq := &logical.Request{
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
	ciphertext := resp.Data["ciphertext"]

	keyReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	resp, err = b.HandleRequest(restoreReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	decryptReq := &logical.Request{
		Path:      "decrypt/test",
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

	resp, err = b.HandleRequest(keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	restoreReq.Path = "restore/test1"
	resp, err = b.HandleRequest(restoreReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	decryptReq.Path = "decrypt/test1"
	resp, err = b.HandleRequest(decryptReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	if resp.Data["plaintext"].(string) != plaintextB64 {
		t.Fatalf("bad: plaintext; expected: %q, actual: %q", plaintextB64, resp.Data["plaintext"].(string))
	}
}

func testSignVerifyBackupRestore(t *testing.T, keyType string) {
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

	signReq := &logical.Request{
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

	signature := resp.Data["signature"]

	verifyReq := &logical.Request{
		Path:      "verify/test",
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
}

func testHMACBackupRestore(t *testing.T, keyType string) {
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

	hmacReq := &logical.Request{
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

	hmac := resp.Data["hmac"]

	verifyReq := &logical.Request{
		Path:      "verify/test",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"hmac":  hmac,
			"input": plaintextB64,
		},
	}

	fmt.Printf("verifyReq.Data: %#v\n", verifyReq.Data)

	resp, err = b.HandleRequest(verifyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	if resp.Data["valid"].(bool) != true {
		t.Fatalf("bad: HMAC verification failed for key type %q", keyType)
	}
}
