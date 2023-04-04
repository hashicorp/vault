// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package tlsutil

import (
	"crypto/tls"

	exttlsutil "github.com/hashicorp/go-secure-stdlib/tlsutil"
)

var ErrInvalidCertParams = exttlsutil.ErrInvalidCertParams

var TLSLookup = exttlsutil.TLSLookup

func ParseCiphers(cipherStr string) ([]uint16, error) {
	return exttlsutil.ParseCiphers(cipherStr)
}

func GetCipherName(cipher uint16) (string, error) {
	return exttlsutil.GetCipherName(cipher)
}

func ClientTLSConfig(caCert []byte, clientCert []byte, clientKey []byte) (*tls.Config, error) {
	return exttlsutil.ClientTLSConfig(caCert, clientCert, clientKey)
}

func LoadClientTLSConfig(caCert, clientCert, clientKey string) (*tls.Config, error) {
	return exttlsutil.LoadClientTLSConfig(caCert, clientCert, clientKey)
}

func SetupTLSConfig(conf map[string]string, address string) (*tls.Config, error) {
	return exttlsutil.SetupTLSConfig(conf, address)
}
