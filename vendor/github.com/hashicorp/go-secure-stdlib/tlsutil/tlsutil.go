// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
)

var ErrInvalidCertParams = errors.New("invalid certificate parameters")

// TLSLookup maps the tls_min_version configuration to the internal value
var TLSLookup = map[string]uint16{
	"tls10": tls.VersionTLS10,
	"tls11": tls.VersionTLS11,
	"tls12": tls.VersionTLS12,
	"tls13": tls.VersionTLS13,
}

// cipherMap maps the cipher suite names to the internal cipher suite code.
var cipherMap = map[string]uint16{
	"TLS_RSA_WITH_RC4_128_SHA":                      tls.TLS_RSA_WITH_RC4_128_SHA,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":                 tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA":                  tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":                  tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":               tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":               tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":               tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":              tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":                tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":          tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_AES_128_GCM_SHA256":                        tls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":                        tls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256":                  tls.TLS_CHACHA20_POLY1305_SHA256,
}

// ParseCiphers parse ciphersuites from the comma-separated string into recognized slice
func ParseCiphers(cipherStr string) ([]uint16, error) {
	suites := []uint16{}
	ciphers := strutil.ParseStringSlice(cipherStr, ",")
	for _, cipher := range ciphers {
		if v, ok := cipherMap[cipher]; ok {
			suites = append(suites, v)
		} else {
			return suites, fmt.Errorf("unsupported cipher %q", cipher)
		}
	}

	return suites, nil
}

// GetCipherName returns the name of a given cipher suite code or an error if the
// given cipher is unsupported.
func GetCipherName(cipher uint16) (string, error) {
	for cipherStr, cipherCode := range cipherMap {
		if cipherCode == cipher {
			return cipherStr, nil
		}
	}
	return "", fmt.Errorf("unsupported cipher %d", cipher)
}

// ClientTLSConfig parses the CA certificate, and optionally a public/private
// client certificate key pair. The certificates must be in PEM encoded format.
func ClientTLSConfig(caCert []byte, clientCert []byte, clientKey []byte) (*tls.Config, error) {
	var tlsConfig *tls.Config
	var pool *x509.CertPool

	switch {
	case len(caCert) != 0:
		// Valid
	case len(clientCert) != 0 && len(clientKey) != 0:
		// Valid
	default:
		return nil, ErrInvalidCertParams
	}

	if len(caCert) != 0 {
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
	}

	tlsConfig = &tls.Config{
		RootCAs:    pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}

	var cert tls.Certificate
	var err error
	if len(clientCert) != 0 && len(clientKey) != 0 {
		cert, err = tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

// LoadClientTLSConfig loads and parse the CA certificate, and optionally a
// public/private client certificate key pair. The certificates must be in PEM
// encoded format.
func LoadClientTLSConfig(caCert, clientCert, clientKey string) (*tls.Config, error) {
	var tlsConfig *tls.Config
	var pool *x509.CertPool

	switch {
	case len(caCert) != 0:
		// Valid
	case len(clientCert) != 0 && len(clientKey) != 0:
		// Valid
	default:
		return nil, ErrInvalidCertParams
	}

	if len(caCert) != 0 {
		pool = x509.NewCertPool()

		data, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}

		if !pool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
	}

	tlsConfig = &tls.Config{
		RootCAs:    pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}

	var cert tls.Certificate
	var err error
	if len(clientCert) != 0 && len(clientKey) != 0 {
		cert, err = tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

func SetupTLSConfig(conf map[string]string, address string) (*tls.Config, error) {
	serverName, _, err := net.SplitHostPort(address)
	switch {
	case err == nil:
	case strings.Contains(err.Error(), "missing port"):
		serverName = conf["address"]
	default:
		return nil, err
	}

	insecureSkipVerify := false
	tlsSkipVerify := conf["tls_skip_verify"]

	if tlsSkipVerify != "" {
		b, err := parseutil.ParseBool(tlsSkipVerify)
		if err != nil {
			return nil, fmt.Errorf("failed parsing tls_skip_verify parameter: %w", err)
		}
		insecureSkipVerify = b
	}

	tlsMinVersionStr, ok := conf["tls_min_version"]
	if !ok {
		// Set the default value
		tlsMinVersionStr = "tls12"
	}

	tlsMinVersion, ok := TLSLookup[tlsMinVersionStr]
	if !ok {
		return nil, fmt.Errorf("invalid 'tls_min_version'")
	}

	tlsClientConfig := &tls.Config{
		MinVersion:         tlsMinVersion,
		InsecureSkipVerify: insecureSkipVerify,
		ServerName:         serverName,
	}

	_, okCert := conf["tls_cert_file"]
	_, okKey := conf["tls_key_file"]

	if okCert && okKey {
		tlsCert, err := tls.LoadX509KeyPair(conf["tls_cert_file"], conf["tls_key_file"])
		if err != nil {
			return nil, fmt.Errorf("client tls setup failed: %w", err)
		}

		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	} else if okCert || okKey {
		return nil, fmt.Errorf("both tls_cert_file and tls_key_file must be provided")
	}

	if tlsCaFile, ok := conf["tls_ca_file"]; ok {
		caPool := x509.NewCertPool()

		data, err := ioutil.ReadFile(tlsCaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}

		if !caPool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsClientConfig.RootCAs = caPool
	}
	return tlsClientConfig, nil
}
