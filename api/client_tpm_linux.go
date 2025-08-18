//go:build linux

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"crypto"
	"crypto/tls"
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

// loadX509KeyPairWithTPM loads a certificate and key pair, with TPM support for TSS2 keys
func loadX509KeyPairWithTPM(certFile, keyFile, tmpDevice string) (tls.Certificate, error) {
	// Read certificate
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Read key
	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to read key file: %w", err)
	}

	// Parse certificate
	certPEMBlock, _ := pem.Decode(certData)
	if certPEMBlock == nil {
		return tls.Certificate{}, fmt.Errorf("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(certPEMBlock.Bytes)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Check if key is TSS2 format
	if isTSS2Key(keyData) {
		// Load TPM key
		signer, err := createTPMSigner(keyFile, tmpDevice)
		if err != nil {
			return tls.Certificate{}, fmt.Errorf("failed to load TPM key: %w", err)
		}

		return tls.Certificate{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  signer,
			Leaf:        cert,
		}, nil
	}

	// Fall back to standard key loading
	return tls.LoadX509KeyPair(certFile, keyFile)
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