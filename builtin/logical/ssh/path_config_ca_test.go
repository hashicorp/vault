package ssh

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestSSH_ConfigCAStorageUpgrade(t *testing.T) {
	var err error

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Store at an older path
	err = config.StorageView.Put(context.Background(), &logical.StorageEntry{
		Key:   caPrivateKeyStoragePathDeprecated,
		Value: []byte(testCAPrivateKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	privateKeyEntry, err := caKey(context.Background(), config.StorageView, caPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	if privateKeyEntry == nil || privateKeyEntry.Key == "" {
		t.Fatalf("failed to read the stored private key")
	}

	entry, err := config.StorageView.Get(context.Background(), caPrivateKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(context.Background(), caPrivateKeyStoragePath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatalf("bad: expected a non-nil entry after upgrade")
	}

	// Store at an older path
	err = config.StorageView.Put(context.Background(), &logical.StorageEntry{
		Key:   caPublicKeyStoragePathDeprecated,
		Value: []byte(testCAPublicKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	publicKeyEntry, err := caKey(context.Background(), config.StorageView, caPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if publicKeyEntry == nil || publicKeyEntry.Key == "" {
		t.Fatalf("failed to read the stored public key")
	}

	entry, err = config.StorageView.Get(context.Background(), caPublicKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(context.Background(), caPublicKeyStoragePath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatalf("bad: expected a non-nil entry after upgrade")
	}
}

func TestSSH_ConfigCAUpdateDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	caReq := &logical.Request{
		Path:      "config/ca",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
	}

	// Auto-generate the keys
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error, got %#v", *resp)
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = map[string]interface{}{
		"public_key":  testCAPublicKey,
		"private_key": testCAPrivateKey,
	}

	// Successfully create a new one
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error, got %#v", *resp)
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = nil

	// Successfully create a new one
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Delete the configured keys
	caReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}
}

func createDeleteHelper(t *testing.T, b logical.Backend, config *logical.BackendConfig, index int, keyType string, keyBits int) {
	// Check that we can create a new key of the specified type
	caReq := &logical.Request{
		Path:      "config/ca",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
	}
	caReq.Data = map[string]interface{}{
		"generate_signing_key": true,
		"key_type":             keyType,
		"key_bits":             keyBits,
	}
	resp, err := b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad case %v: err: %v, resp: %v", index, err, resp)
	}
	if !strings.Contains(resp.Data["public_key"].(string), caReq.Data["key_type"].(string)) {
		t.Fatalf("bad case %v: expected public key of type %v but was %v", index, caReq.Data["key_type"], resp.Data["public_key"])
	}

	issueOptions := map[string]interface{}{
		"public_key": testCAPublicKeyEd25519,
	}
	issueReq := &logical.Request{
		Path:      "sign/ca-issuance",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      issueOptions,
	}
	resp, err = b.HandleRequest(context.Background(), issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad case %v: err: %v, resp: %v", index, err, resp)
	}

	// Delete the configured keys
	caReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad case %v: err: %v, resp: %v", index, err, resp)
	}
}

func TestSSH_ConfigCAKeyTypes(t *testing.T) {
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	cases := []struct {
		keyType string
		keyBits int
	}{
		{"ssh-rsa", 2048},
		{"ssh-rsa", 4096},
		{"ssh-rsa", 0},
		{"rsa", 2048},
		{"rsa", 4096},
		{"ecdsa-sha2-nistp256", 0},
		{"ecdsa-sha2-nistp384", 0},
		{"ecdsa-sha2-nistp521", 0},
		{"ec", 256},
		{"ec", 384},
		{"ec", 521},
		{"ec", 0},
		{"ssh-ed25519", 0},
		{"ed25519", 0},
	}

	// Create a role for ssh signing.
	roleOptions := map[string]interface{}{
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"key_type":                "ca",
		"ttl":                     "30s",
		"not_before_duration":     "2h",
	}
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/ca-issuance",
		Data:      roleOptions,
		Storage:   config.StorageView,
	}
	_, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil {
		t.Fatalf("Cannot create role to issue against: %s", err)
	}

	for index, scenario := range cases {
		createDeleteHelper(t, b, config, index, scenario.keyType, scenario.keyBits)
	}
}
