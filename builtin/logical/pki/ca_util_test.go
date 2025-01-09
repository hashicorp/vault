// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
)

func TestGetKeyTypeAndBitsFromPublicKeyForRole(t *testing.T) {
	rsaKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error generating rsa key: %s", err)
	}

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("error generating ecdsa key: %s", err)
	}

	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("error generating ed25519 key: %s", err)
	}

	testCases := map[string]struct {
		publicKey       crypto.PublicKey
		expectedKeyType certutil.PrivateKeyType
		expectedKeyBits int
		expectError     bool
	}{
		"rsa": {
			publicKey:       rsaKey.Public(),
			expectedKeyType: certutil.RSAPrivateKey,
			expectedKeyBits: 2048,
		},
		"ecdsa": {
			publicKey:       ecdsaKey.Public(),
			expectedKeyType: certutil.ECPrivateKey,
			expectedKeyBits: 0,
		},
		"ed25519": {
			publicKey:       publicKey,
			expectedKeyType: certutil.Ed25519PrivateKey,
			expectedKeyBits: 0,
		},
		"bad key type": {
			publicKey:       []byte{},
			expectedKeyType: certutil.UnknownPrivateKey,
			expectedKeyBits: 0,
			expectError:     true,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			keyType, keyBits, err := getKeyTypeAndBitsFromPublicKeyForRole(tt.publicKey)
			if err != nil && !tt.expectError {
				t.Fatalf("unexpected error: %s", err)
			}
			if err == nil && tt.expectError {
				t.Fatal("expected error, got nil")
			}

			if keyType != tt.expectedKeyType {
				t.Fatalf("key type mismatch: expected %s, got %s", tt.expectedKeyType, keyType)
			}

			if keyBits != tt.expectedKeyBits {
				t.Fatalf("key bits mismatch: expected %d, got %d", tt.expectedKeyBits, keyBits)
			}
		})
	}
}
