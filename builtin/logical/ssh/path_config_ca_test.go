package ssh

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestSSH_ConfigCAStorageUpgrade(t *testing.T) {
	var err error

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// Store at an older path
	err = config.StorageView.Put(&logical.StorageEntry{
		Key:   caPrivateKeyStoragePathDeprecated,
		Value: []byte(privateKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	privateKeyEntry, err := caKey(config.StorageView, caPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	if privateKeyEntry == nil || privateKeyEntry.Key == "" {
		t.Fatalf("failed to read the stored private key")
	}

	entry, err := config.StorageView.Get(caPrivateKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(caPrivateKeyStoragePath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatalf("bad: expected a non-nil entry after upgrade")
	}

	// Store at an older path
	err = config.StorageView.Put(&logical.StorageEntry{
		Key:   caPublicKeyStoragePathDeprecated,
		Value: []byte(publicKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	publicKeyEntry, err := caKey(config.StorageView, caPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if publicKeyEntry == nil || publicKeyEntry.Key == "" {
		t.Fatalf("failed to read the stored public key")
	}

	entry, err = config.StorageView.Get(caPublicKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(caPublicKeyStoragePath)
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

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	caReq := &logical.Request{
		Path:      "config/ca",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
	}

	// Auto-generate the keys
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(caReq)
	if err == nil {
		t.Fatalf("expected an error")
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = map[string]interface{}{
		"public_key":  publicKey,
		"private_key": privateKey,
	}

	// Successfully create a new one
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(caReq)
	if err == nil {
		t.Fatalf("expected an error")
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = nil

	// Successfully create a new one
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}
}
