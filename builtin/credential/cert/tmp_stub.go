//go:build !linux

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"runtime"
)

// isTSS2Key checks if the key data contains a TSS2 format private key
func isTSS2Key(keyData []byte) bool {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return false
	}
	return block.Type == "TSS2 PRIVATE KEY"
}

// createTPMSigner returns an error since TPM is not supported on non-Linux platforms
func createTPMSigner(keyPath, tmpDevice string) (crypto.Signer, error) {
	return nil, fmt.Errorf("TPM authentication is only supported on Linux systems with TPM 2.0 hardware (current OS: %s)", runtime.GOOS)
}

// loadPrivateKey attempts to load a private key, with TPM error on non-Linux
func loadPrivateKey(keyPath, tmpDevice string) (crypto.Signer, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// Check if it's a TSS2 key and return appropriate error
	if isTSS2Key(keyData) {
		return nil, fmt.Errorf("TPM authentication with TSS2 keys is only supported on Linux systems with TPM 2.0 hardware (current OS: %s)", runtime.GOOS)
	}

	// Fall back to standard key parsing
	return loadStandardPrivateKey(keyData)
}

// loadStandardPrivateKey parses standard format private keys (RSA, ECC, PKCS8)
func loadStandardPrivateKey(keyData []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if signer, ok := key.(crypto.Signer); ok {
			return signer, nil
		}
		return nil, fmt.Errorf("private key does not implement crypto.Signer")
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}
}

// tmpSupported returns false since this is the non-Linux stub
func tmpSupported() bool {
	return false
}

// getTPMDevicePath returns empty string since TPM is not supported
func getTPMDevicePath(tmpDevice string) string {
	return ""
}