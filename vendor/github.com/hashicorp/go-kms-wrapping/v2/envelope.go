// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wrapping

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	fmt "fmt"

	uuid "github.com/hashicorp/go-uuid"
)

// EnvelopeEncrypt takes in plaintext and envelope encrypts it, generating an
// EnvelopeInfo value.  An empty plaintext is a valid parameter and will not cause
// an error.  Also note: if you provide a plaintext of []byte(""),
// EnvelopeDecrypt will return []byte(nil).
//
// Supported options:
//
// * wrapping.WithAad: Additional authenticated data that should be sourced from
// a separate location, and must also be provided during envelope decryption
func EnvelopeEncrypt(plaintext []byte, opt ...Option) (*EnvelopeInfo, error) {
	opts, err := GetOpts(opt...)
	if err != nil {
		return nil, err
	}

	// Generate DEK
	key, err := uuid.GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}

	var iv []byte
	if opts.WithIv != nil {
		if len(opts.WithIv) != 12 {
			return nil, fmt.Errorf("invalid IV provided: expected 12 bytes, got %d", len(opts.WithIv))
		}
		iv = opts.WithIv
	} else {
		iv, err = uuid.GenerateRandomBytes(12)
		if err != nil {
			return nil, err
		}
	}

	aead, err := aeadEncrypter(key)
	if err != nil {
		return nil, err
	}

	return &EnvelopeInfo{
		Ciphertext: aead.Seal(nil, iv, plaintext, opts.WithAad),
		Key:        key,
		Iv:         iv,
	}, nil
}

// EnvelopeDecrypt takes in EnvelopeInfo and potentially additional options and
// decrypts.  Also note: if you provided a plaintext of []byte("") to
// EnvelopeEncrypt, then this function will return []byte(nil).
//
// Supported options:
//
// * wrapping.WithAad: Additional authenticated data that should be sourced from
// a separate location, and must match what was provided during envelope
// encryption.
func EnvelopeDecrypt(data *EnvelopeInfo, opt ...Option) ([]byte, error) {
	// need to check data or we could panic when trying to access data.Key
	if data == nil {
		return nil, fmt.Errorf("missing envelope info: %w", ErrInvalidParameter)
	}
	opts, err := GetOpts(opt...)
	if err != nil {
		return nil, err
	}

	aead, err := aeadEncrypter(data.Key)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, data.Iv, data.Ciphertext, opts.WithAad)
}

func aeadEncrypter(key []byte) (cipher.AEAD, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create the GCM mode AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, errors.New("failed to initialize GCM mode")
	}

	return gcm, nil
}
