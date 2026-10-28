// Copyright IBM Corp. 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pkcs7

import (
	"crypto/x509"
)

// CustomParseX509Certificates parses one or more DER-encoded certificates from
// data using the standard library parser.
func CustomParseX509Certificates(data []byte) ([]*x509.Certificate, error) {
	return x509.ParseCertificates(data)
}
