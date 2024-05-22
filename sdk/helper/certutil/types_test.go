// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestGetPrivateKeyTypeFromPublicKey(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
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
		expectedKeyType PrivateKeyType
	}{
		"rsa": {
			publicKey:       rsaKey.Public(),
			expectedKeyType: RSAPrivateKey,
		},
		"ecdsa": {
			publicKey:       ecdsaKey.Public(),
			expectedKeyType: ECPrivateKey,
		},
		"ed25519": {
			publicKey:       publicKey,
			expectedKeyType: Ed25519PrivateKey,
		},
		"bad key type": {
			publicKey:       []byte{},
			expectedKeyType: UnknownPrivateKey,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			keyType := GetPrivateKeyTypeFromPublicKey(tt.publicKey)

			if keyType != tt.expectedKeyType {
				t.Fatalf("key type mismatch: expected %s, got %s", tt.expectedKeyType, keyType)
			}
		})
	}
}
