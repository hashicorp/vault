// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/go-rootcerts"
)

// BuildTLSConfig creates a tls.Config from an attestation result.
// The resulting config uses TPM-backed certificates for client authentication.
// The TPM will be opened on each signing operation.
func BuildTLSConfig(result *AttestationResult, caCert, caPath string, cfg TPMConfig) (*tls.Config, error) {
	if result == nil {
		return nil, fmt.Errorf("attestation result is nil")
	}
	if result.KeyMaterial == nil {
		return nil, fmt.Errorf("key material is nil")
	}

	// Parse certificate
	certBlock, _ := pem.Decode(result.CertPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Create TPM signer with key material
	signer, err := NewTPMSigner(cfg, result.KeyMaterial)
	if err != nil {
		return nil, fmt.Errorf("failed to create TPM signer: %w", err)
	}

	tlsCert := &tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  signer,
		Leaf:        cert,
	}

	tlsConfig := &tls.Config{
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return tlsCert, nil
		},
		MinVersion: tls.VersionTLS12,
	}

	// Configure root CAs
	rootConfig := &rootcerts.Config{
		CAFile: caCert,
		CAPath: caPath,
	}
	if err := rootcerts.ConfigureTLS(tlsConfig, rootConfig); err != nil {
		return nil, fmt.Errorf("failed to configure root CAs: %w", err)
	}

	return tlsConfig, nil
}
