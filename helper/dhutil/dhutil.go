// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package dhutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

type PublicKeyInfo struct {
	Curve25519PublicKey []byte `json:"curve25519_public_key"`
}

type Envelope struct {
	Curve25519PublicKey []byte `json:"curve25519_public_key"`
	Nonce               []byte `json:"nonce"`
	EncryptedPayload    []byte `json:"encrypted_payload"`
}

// generatePublicPrivateKey uses curve25519 to generate a public and private key
// pair.
func GeneratePublicPrivateKey() ([]byte, []byte, error) {
	var scalar, public [32]byte

	if _, err := io.ReadFull(rand.Reader, scalar[:]); err != nil {
		return nil, nil, err
	}

	curve25519.ScalarBaseMult(&public, &scalar)
	return public[:], scalar[:], nil
}

// GenerateSharedSecret uses the private key and the other party's public key to
// generate the shared secret.
func GenerateSharedSecret(ourPrivate, theirPublic []byte) ([]byte, error) {
	if len(ourPrivate) != 32 {
		return nil, fmt.Errorf("invalid private key length: %d", len(ourPrivate))
	}
	if len(theirPublic) != 32 {
		return nil, fmt.Errorf("invalid public key length: %d", len(theirPublic))
	}

	return curve25519.X25519(ourPrivate, theirPublic)
}

// DeriveSharedKey uses HKDF to derive a key from a shared secret and public keys
func DeriveSharedKey(secret, ourPublic, theirPublic []byte) ([]byte, error) {
	// Derive the final key from the HKDF of the secret and public keys.

	/*
		Internally, HKDF hashes the secret and two public keys. If Alice and Bob are doing DH key exchange, Alice calculates:

		HKDF(secret, A, B) since ourPublic is A.

			Bob calculates HKDF(secret, B, A), since Bob's ours is B. That produces a different value. Now we only care
		    that both public keys participate in the derivation, so simply sorting them so they are in a consistent
		    numerical order (either one would do) arrives at an agreed value.
	*/

	var pub1 []byte
	var pub2 []byte
	switch bytes.Compare(ourPublic, theirPublic) {
	case 0:
		return nil, errors.New("same public key supplied for both participants")
	case -1:
		pub1 = ourPublic
		pub2 = theirPublic
	case 1:
		pub1 = theirPublic
		pub2 = ourPublic
	}

	kio := hkdf.New(sha256.New, secret, pub1, pub2)

	var key [32]byte
	n, err := io.ReadFull(kio, key[:])
	if err != nil {
		// Don't return the key along with the error to prevent misuse
		return nil, err
	}
	if n != 32 {
		return nil, errors.New("short read from hkdf")
	}

	return key[:], nil
}

// Use AES256-GCM to encrypt some plaintext with a provided key. The returned values are
// the ciphertext, the nonce, and error respectively.
func EncryptAES(key, plaintext, aad []byte) ([]byte, []byte, error) {
	// We enforce AES-256, so check explicitly for 32 bytes on the key
	if len(key) != 32 {
		return nil, nil, fmt.Errorf("invalid key length: %d", len(key))
	}

	if len(plaintext) == 0 {
		return nil, nil, errors.New("empty plaintext provided")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, aad)

	return ciphertext, nonce, nil
}

// Use AES256-GCM to decrypt some ciphertext with a provided key and nonce. The
// returned values are the plaintext and error respectively.
func DecryptAES(key, ciphertext, nonce, aad []byte) ([]byte, error) {
	// We enforce AES-256, so check explicitly for 32 bytes on the key
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(key))
	}

	if len(ciphertext) == 0 {
		return nil, errors.New("empty ciphertext provided")
	}

	if len(nonce) == 0 {
		return nil, errors.New("empty nonce provided")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
