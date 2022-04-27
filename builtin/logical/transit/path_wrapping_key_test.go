package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storagePath = "import/policy/" + WrappingKeyName
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
	keyEntry, err := storage.Get(context.Background(), storagePath)
	if err != nil {
		t.Fatalf("error retrieving wrapping key from storage: %s", err)
	}
	if keyEntry != nil {
		t.Fatal("wrapping key unexpectedly exists")
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
	if resp == nil || resp.Data == nil || resp.Data["public_key"] == nil {
		t.Fatal("expected non-nil response")
	}
	pubKey := resp.Data["public_key"]

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
	if resp == nil || resp.Data == nil || resp.Data["public_key"] == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Data["public_key"] != pubKey {
		t.Fatal("wrapping key public component changed between requests")
	}
}
