// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

import (
	"crypto"
	"fmt"
	"io"
	"sync"

	"github.com/google/go-attestation/attest"
)

// TPMSigner implements crypto.Signer using a TPM-backed key.
// It is safe for concurrent use.
type TPMSigner struct {
	cfg         TPMConfig
	keyMaterial *TPMKeyMaterial
	mu          sync.Mutex
	public      crypto.PublicKey
}

var _ crypto.Signer = &TPMSigner{}

// NewTPMSigner creates a new TPMSigner from key material.
// If cfg.TestTPM is non-nil it is reused for every signing operation.
// Otherwise the TPM is opened and closed on each operation using cfg.DevicePath.
func NewTPMSigner(cfg TPMConfig, keyMaterial *TPMKeyMaterial) (*TPMSigner, error) {
	ts := &TPMSigner{
		cfg:         cfg,
		keyMaterial: keyMaterial,
	}

	// Load the public key once and cache it
	key, closer, err := ts.getKey()
	if err != nil {
		return nil, err
	}
	ts.public = key.Public()
	key.Close()
	closer()

	return ts, nil
}

// Public returns the public key associated with the private key in the TPM.
// This is called frequently by Go's TLS stack, so it must be cached in memory.
func (ts *TPMSigner) Public() crypto.PublicKey {
	return ts.public
}

// getKey loads the TPM key, opening the TPM if necessary.
// Returns the key, a closer function, and any error.
func (ts *TPMSigner) getKey() (*attest.Key, func(), error) {
	tpm, err := openTPM(ts.cfg.DevicePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open TPM: %w", err)
	}

	key, err := tpm.LoadKey(ts.keyMaterial.AppBlob)
	if err != nil {
		tpm.Close()
		return nil, nil, err
	}

	return key, func() { tpm.Close() }, nil
}

// Sign is called by the TLS stack during the handshake to sign data.
// It implements the crypto.Signer interface.
func (ts *TPMSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	key, closer, err := ts.getKey()
	if err != nil {
		return nil, err
	}
	defer key.Close()
	defer closer()

	// Get the crypto.Signer from TPM key
	pubKey := key.Public()
	privKey, err := key.Private(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get TPM private key: %w", err)
	}

	signer, ok := privKey.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("TPM private key does not implement crypto.Signer (got type %T)", privKey)
	}

	return signer.Sign(rand, digest, opts)
}
