package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

// TODO: Replace these with the real values. - schultz
const (
	pubKeyStoragePath  = "/path/to/public/wrapping/key"
	privKeyStoragePath = "/path/to/private/wrapping/key"
)

func TestTransit_WrappingKey(t *testing.T) {
	// Set up shared backend for subtests
	var b *backend
	storage := &logical.InmemStorage{}
	sysView := logical.TestSystemView()
	b, _ = Backend(
		context.Background(),
		&logical.BackendConfig{
			StorageView: storage,
			System:      sysView,
		},
	)

	// Ensure the key does not exist before requesting it.
	pubKeyEntry, err := storage.Get(context.Background(), pubKeyStoragePath)
	if err != nil {
		t.Fatalf("error retrieving public wrapping key from storage: %s", err)
	}
	if pubKeyEntry != nil {
		t.Fatal("public wrapping key unexpectedly exists")
	}

	privKeyEntry, err := storage.Get(context.Background(), privKeyStoragePath)
	if err != nil {
		t.Fatalf("error retrieving private wrapping key from storage: %s", err)
	}
	if privKeyEntry != nil {
		t.Fatal("private wrapping key unexpectedly exists")
	}

	// Generate the key pair by requesting the public key.
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "wrapping_key",
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected request error: %s", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	// TODO: Save public key from response. - schultz
	pubKey := ""

	// TODO: Ensure only public key is returned in response. - schultz

	// Ensure the key exists in storage and that the public key matches what was returned.
	pubKeyEntry, err = storage.Get(context.Background(), pubKeyStoragePath)
	if err != nil {
		t.Fatalf("error retrieving public wrapping key from storage: %s", err)
	}
	if pubKeyEntry == nil {
		t.Fatal("expected non-nil public wrapping key")
	}
	if string(pubKeyEntry.Value) != pubKey {
		t.Fatal("returned public wrapping key does not match value in storage")
	}

	privKeyEntry, err = storage.Get(context.Background(), privKeyStoragePath)
	if err != nil {
		t.Fatalf("error retrieving private wrapping key from storage: %s", err)
	}
	if privKeyEntry == nil {
		t.Fatal("expected non-nil private wrapping key")
	}

	// Request the wrapping key again to ensure it isn't regenerated.
	req = &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "wrapping_key",
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected request error: %s", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	// TODO: Compare response body to previously returned public key. - schultz
}
