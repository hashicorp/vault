// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storagePath = "policy/import/" + WrappingKeyName
)

func TestTransit_WrappingKey(t *testing.T) {
	// Set up shared backend for subtests
	b, s := createBackendWithStorage(t)

	// Ensure the key does not exist before requesting it.
	keyEntry, err := s.Get(context.Background(), storagePath)
	if err != nil {
		t.Fatalf("error retrieving wrapping key from storage: %s", err)
	}
	if keyEntry != nil {
		t.Fatal("wrapping key unexpectedly exists")
	}

	// Generate the key pair by requesting the public key.
	req := &logical.Request{
		Storage:   s,
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
	pubKeyPEM := resp.Data["public_key"]

	// Ensure the returned key is a 4096-bit RSA key.
	pubKeyBlock, _ := pem.Decode([]byte(pubKeyPEM.(string)))
	rawPubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse public wrapping key: %s", err)
	}
	wrappingKey, ok := rawPubKey.(*rsa.PublicKey)
	if !ok || wrappingKey.Size() != 512 {
		t.Fatal("public wrapping key is not a 4096-bit RSA key")
	}

	// Request the wrapping key again to ensure it isn't regenerated.
	req = &logical.Request{
		Storage:   s,
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
	if resp.Data["public_key"] != pubKeyPEM {
		t.Fatal("wrapping key public component changed between requests")
	}
}
