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

	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// Store at an older path
	err = config.StorageView.Put(&logical.StorageEntry{
		Key:   CAPrivateKeyStoragePathDeprecated,
		Value: []byte(privateKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	storedPrivateKey, err := caKey(config.StorageView, "private")
	if err != nil {
		t.Fatal(err)
	}
	if storedPrivateKey == "" {
		t.Fatalf("failed to read the stored private key")
	}

	entry, err := config.StorageView.Get(CAPrivateKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(CAPrivateKeyStoragePath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatalf("bad: expected a non-nil entry after upgrade")
	}

	// Store at an older path
	err = config.StorageView.Put(&logical.StorageEntry{
		Key:   CAPublicKeyStoragePathDeprecated,
		Value: []byte(publicKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	storedPublicKey, err := caKey(config.StorageView, "public")
	if err != nil {
		t.Fatal(err)
	}
	if storedPublicKey == "" {
		t.Fatalf("failed to read the stored public key")
	}

	entry, err = config.StorageView.Get(CAPublicKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(CAPublicKeyStoragePath)
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
