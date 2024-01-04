// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package parsing

import (
	"crypto/x509"
	"fmt"
)

func ParseCertificateRequestFromString(pemCert string) (*x509.CertificateRequest, error) {
	return ParseCertificateRequestFromBytes([]byte(pemCert))
}

func ParseCertificateRequestFromBytes(certBytes []byte) (*x509.CertificateRequest, error) {
	block, err := DecodePem(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate request: %w", err)
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate request: %w", err)
	}

	return csr, nil
}
