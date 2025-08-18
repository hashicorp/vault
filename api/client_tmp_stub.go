//go:build !linux

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"crypto/tls"
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

// loadX509KeyPairWithTPM loads a certificate and key pair, with error for TSS2 on non-Linux
func loadX509KeyPairWithTPM(certFile, keyFile, tmpDevice string) (tls.Certificate, error) {
	// Read key to check format
	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to read key file: %w", err)
	}

	// Check if it's a TSS2 key and return appropriate error
	if isTSS2Key(keyData) {
		return tls.Certificate{}, fmt.Errorf("TPM authentication with TSS2 keys is only supported on Linux systems with TPM 2.0 hardware (current OS: %s)", runtime.GOOS)
	}

	// Fall back to standard key loading
	return tls.LoadX509KeyPair(certFile, keyFile)
}