//go:build linux

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	tpmkeyfiles "github.com/foxboron/go-tpm-keyfiles"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
)

const defaultTPMDevice = "/dev/tpmrm0"

// isTSS2Key checks if the key data contains a TSS2 format private key
func isTSS2Key(keyData []byte) bool {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return false
	}
	return block.Type == "TSS2 PRIVATE KEY"
}

// createTPMSigner creates a crypto.Signer from a TSS2 format key file
func createTPMSigner(keyPath, tmpDevice string) (crypto.Signer, error) {
	if tmpDevice == "" {
		tmpDevice = defaultTPMDevice
	}

	// Check if TPM device exists
	if _, err := os.Stat(tmpDevice); os.IsNotExist(err) {
		return nil, fmt.Errorf("TPM device not found: %s", tmpDevice)
	}

	// Read the key file
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// Verify it's a TSS2 key
	if !isTSS2Key(keyData) {
		return nil, fmt.Errorf("key file is not in TSS2 format")
	}

	// Open TPM device
	tpmTransport, err := linuxtpm.Open(tmpDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to open TPM device %s: %w", tmpDevice, err)
	}

	// Parse the TSS2 key
	key, err := tpmkeyfiles.Decode(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TSS2 key: %w", err)
	}

	// Create the signer
	signer, err := key.Signer(tpmTransport)
	if err != nil {
		return nil, fmt.Errorf("failed to create TPM signer: %w", err)
	}

	return signer, nil
}

// loadPrivateKey attempts to load a private key, with automatic TPM detection
func loadPrivateKey(keyPath, tmpDevice string) (crypto.Signer, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// Check if it's a TSS2 key first
	if isTSS2Key(keyData) {
		return createTPMSigner(keyPath, tmpDevice)
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

// tmpSupported returns true since this is the Linux implementation
func tmpSupported() bool {
	return true
}

// getTPMDevicePath returns the TPM device path, using default if empty
func getTPMDevicePath(tmpDevice string) string {
	if tmpDevice == "" {
		return defaultTPMDevice
	}
	return tmpDevice
}
