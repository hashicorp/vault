// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package parsing

import (
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

func SerialFromCert(cert *x509.Certificate) string {
	return SerialFromBigInt(cert.SerialNumber)
}

func SerialFromBigInt(serial *big.Int) string {
	return strings.TrimSpace(certutil.GetHexFormatted(serial.Bytes(), ":"))
}

// NormalizeSerialForStorageFromBigInt given a serial number, format it as a string
// that is safe to store within a filesystem
func NormalizeSerialForStorageFromBigInt(serial *big.Int) string {
	return strings.TrimSpace(certutil.GetHexFormatted(serial.Bytes(), "-"))
}

// NormalizeSerialForStorage given a serial number with ':' characters, convert
// them to '-' which is safe to store within filesystems
func NormalizeSerialForStorage(serial string) string {
	return strings.ReplaceAll(strings.ToLower(serial), ":", "-")
}

func ParseCertificateFromString(pemCert string) (*x509.Certificate, error) {
	return ParseCertificateFromBytes([]byte(pemCert))
}

func ParseCertificateFromBytes(certBytes []byte) (*x509.Certificate, error) {
	block, err := DecodePem(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	return cert, nil
}

func ParseCertificatesFromString(pemCerts string) ([]*x509.Certificate, error) {
	return ParseCertificatesFromBytes([]byte(pemCerts))
}

func ParseCertificatesFromBytes(certBytes []byte) ([]*x509.Certificate, error) {
	block, err := DecodePem(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	cert, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	return cert, nil
}

func ParseKeyUsages(input []string) int {
	var parsedKeyUsages x509.KeyUsage
	for _, k := range input {
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "digitalsignature":
			parsedKeyUsages |= x509.KeyUsageDigitalSignature
		case "contentcommitment":
			parsedKeyUsages |= x509.KeyUsageContentCommitment
		case "keyencipherment":
			parsedKeyUsages |= x509.KeyUsageKeyEncipherment
		case "dataencipherment":
			parsedKeyUsages |= x509.KeyUsageDataEncipherment
		case "keyagreement":
			parsedKeyUsages |= x509.KeyUsageKeyAgreement
		case "certsign":
			parsedKeyUsages |= x509.KeyUsageCertSign
		case "crlsign":
			parsedKeyUsages |= x509.KeyUsageCRLSign
		case "encipheronly":
			parsedKeyUsages |= x509.KeyUsageEncipherOnly
		case "decipheronly":
			parsedKeyUsages |= x509.KeyUsageDecipherOnly
		}
	}

	return int(parsedKeyUsages)
}
