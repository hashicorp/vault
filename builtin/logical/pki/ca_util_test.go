// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/mldsa"
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

	mldsa44Key, err := mldsa.GenerateKey(mldsa.MLDSA44())
	if err != nil {
		t.Fatalf("error generating mldsa 44 key: %s", err)
	}

	mldsa65Key, err := mldsa.GenerateKey(mldsa.MLDSA65())
	if err != nil {
		t.Fatalf("error generating mldsa 65 key: %s", err)
	}

	mldsa87Key, err := mldsa.GenerateKey(mldsa.MLDSA87())
	if err != nil {
		t.Fatalf("error generating mldsa 87 key: %s", err)
	}

	testCases := map[string]struct {
		publicKey            crypto.PublicKey
		expectedKeyType      certutil.PrivateKeyType
		expectedKeyBits      int
		expectedParameterSet certutil.ParameterSet
		expectError          bool
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
		"ml-dsa-44": {
			publicKey:            mldsa44Key.Public(),
			expectedKeyType:      certutil.MLDSAPrivateKey,
			expectedKeyBits:      0,
			expectedParameterSet: certutil.MLDSA44,
		},
		"ml-dsa-65": {
			publicKey:            mldsa65Key.Public(),
			expectedKeyType:      certutil.MLDSAPrivateKey,
			expectedKeyBits:      0,
			expectedParameterSet: certutil.MLDSA65,
		},
		"ml-dsa-87": {
			publicKey:            mldsa87Key.Public(),
			expectedKeyType:      certutil.MLDSAPrivateKey,
			expectedKeyBits:      0,
			expectedParameterSet: certutil.MLDSA87,
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
			keyType, keyBits, parameterSet, err := getKeyTypeAndBitsFromPublicKeyForRole(tt.publicKey)
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

			if parameterSet != tt.expectedParameterSet {
				t.Fatalf("parameter set mismatch: expected %#v, got %#v", tt.expectedParameterSet, parameterSet)
			}
		})
	}
}
