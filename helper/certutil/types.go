// Package certutil contains helper functions that are mostly used
// with the PKI backend but can be generally useful. Functionality
// includes helpers for converting a certificate/private key bundle
// between DER and PEM, printing certificate serial numbers, and more.
//
// Functionality specific to the PKI backend includes some types
// and helper methods to make requesting certificates from the
// backend easy.
package certutil

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// TLSUsage controls whether the intended usage of a *tls.Config
// returned from ParsedCertBundle.GetTLSConfig is for server use,
// client use, or both, which affects which values are set
type TLSUsage int

// The type of of the Private Key referenced in CertBundle
// and ParsedCertBundle. This uses colloquial names rather than
// official names, to eliminate confusion
const (
	UnknownPrivateKey = iota
	RSAPrivateKey
	ECPrivateKey

	TLSServer TLSUsage = 1 << iota
	TLSClient
)

// UserError represents an error generated due to invalid user input
type UserError struct {
	s string
}

func (e UserError) Error() string {
	return e.s
}

// InternalError represents an error generated internally,
// presumably not due to invalid user input
type InternalError struct {
	s string
}

func (e InternalError) Error() string {
	return e.s
}

// CertBundle contains a key type, a PEM-encoded private key,
// a PEM-encoded certificate, and a string-encoded serial number,
// returned from a successful Issue request
type CertBundle struct {
	PrivateKeyType string `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	Certificate    string `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	IssuingCA      string `json:"issuing_ca" structs:"issuing_ca" mapstructure:"issuing_ca"`
	PrivateKey     string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	SerialNumber   string `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

// ParsedCertBundle contains a key type, a DER-encoded private key,
// a DER-encoded certificate, and a big.Int serial number
type ParsedCertBundle struct {
	PrivateKeyType   int
	PrivateKeyBytes  []byte
	PrivateKey       crypto.Signer
	IssuingCABytes   []byte
	IssuingCA        *x509.Certificate
	CertificateBytes []byte
	Certificate      *x509.Certificate
}

// ToParsedCertBundle converts a string-based certificate bundle
// to a byte-based raw certificate bundle
func (c *CertBundle) ToParsedCertBundle() (*ParsedCertBundle, error) {
	result := &ParsedCertBundle{}
	var err error
	var pemBlock *pem.Block

	pemBlock, _ = pem.Decode([]byte(c.PrivateKey))
	if pemBlock == nil {
		return nil, UserError{"Error decoding private key from cert bundle"}
	}
	result.PrivateKeyBytes = pemBlock.Bytes

	switch c.PrivateKeyType {
	case "ec":
		result.PrivateKeyType = ECPrivateKey
	case "rsa":
		result.PrivateKeyType = RSAPrivateKey
	default:
		// Try to figure it out and correct
		if _, err := x509.ParseECPrivateKey(pemBlock.Bytes); err == nil {
			result.PrivateKeyType = ECPrivateKey
			c.PrivateKeyType = "ec"
		} else if _, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err == nil {
			result.PrivateKeyType = RSAPrivateKey
			c.PrivateKeyType = "rsa"
		} else {
			return nil, UserError{fmt.Sprintf("Unknown private key type in bundle: %s", c.PrivateKeyType)}
		}
	}

	result.PrivateKey, err = result.getSigner()
	if err != nil {
		return nil, UserError{fmt.Sprintf("Error getting signer: %s", err)}
	}

	pemBlock, _ = pem.Decode([]byte(c.Certificate))
	if pemBlock == nil {
		return nil, UserError{"Error decoding certificate from cert bundle"}
	}
	result.CertificateBytes = pemBlock.Bytes
	result.Certificate, err = x509.ParseCertificate(result.CertificateBytes)
	if err != nil {
		return nil, UserError{"Error encountered parsing certificate bytes from raw bundle"}
	}

	if len(c.IssuingCA) != 0 {
		pemBlock, _ = pem.Decode([]byte(c.IssuingCA))
		if pemBlock == nil {
			return nil, UserError{"Error decoding issuing CA from cert bundle"}
		}
		result.IssuingCABytes = pemBlock.Bytes
		result.IssuingCA, err = x509.ParseCertificate(result.IssuingCABytes)
		if err != nil {
			return nil, UserError{fmt.Sprintf("Error parsing CA certificate: %s", err)}
		}
	}

	if len(c.SerialNumber) == 0 {
		c.SerialNumber = GetOctalFormatted(result.Certificate.SerialNumber.Bytes(), ":")
	}

	return result, nil
}

// ToCertBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (p *ParsedCertBundle) ToCertBundle() (*CertBundle, error) {
	result := &CertBundle{
		SerialNumber: GetOctalFormatted(p.Certificate.SerialNumber.Bytes(), ":"),
	}

	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: p.CertificateBytes,
	}
	result.Certificate = string(pem.EncodeToMemory(&block))

	if len(p.IssuingCABytes) != 0 {
		block.Bytes = p.IssuingCABytes
		result.IssuingCA = string(pem.EncodeToMemory(&block))
	}

	block.Bytes = p.PrivateKeyBytes
	switch p.PrivateKeyType {
	case RSAPrivateKey:
		result.PrivateKeyType = "rsa"
		block.Type = "RSA PRIVATE KEY"
	case ECPrivateKey:
		result.PrivateKeyType = "ec"
		block.Type = "EC PRIVATE KEY"
	default:
		return nil, InternalError{"Could not determine private key type when creating block"}
	}
	result.PrivateKey = string(pem.EncodeToMemory(&block))

	return result, nil
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this ParsedCertBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (p *ParsedCertBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error
	switch p.PrivateKeyType {
	case ECPrivateKey:
		signer, err = x509.ParseECPrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, UserError{fmt.Sprintf("Unable to parse CA's private EC key: %s", err)}
		}
	case RSAPrivateKey:
		signer, err = x509.ParsePKCS1PrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, UserError{fmt.Sprintf("Unable to parse CA's private RSA key: %s", err)}
		}
	default:
		return nil, UserError{"Unable to determine type of private key; only RSA and EC are supported"}
	}
	return signer, nil
}

// GetTLSConfig returns a TLS config generally suitable for client
// authentiation. The returned TLS config can be modified slightly
// to be made suitable for a server requiring client authentication;
// specifically, you should set the value of ClientAuth in the returned
// config to match your needs.
func (p *ParsedCertBundle) GetTLSConfig(usage TLSUsage) (*tls.Config, error) {
	tlsCert := &tls.Certificate{
		Certificate: [][]byte{
			p.CertificateBytes,
		},
		PrivateKey: p.PrivateKey,
		Leaf:       p.Certificate,
	}

	tlsConfig := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{*tlsCert},
	}

	if len(p.IssuingCABytes) > 0 {
		tlsCert.Certificate = append(tlsCert.Certificate, p.IssuingCABytes)

		// Technically we only need one cert, but this doesn't duplicate code
		certBundle, err := p.ToCertBundle()
		if err != nil {
			return nil, fmt.Errorf("Error converting parsed bundle to string bundle when getting TLS config: %s", err)
		}

		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM([]byte(certBundle.IssuingCA))
		if !ok {
			return nil, fmt.Errorf("Could not append CA certificate")
		}

		if usage&TLSServer != 0 {
			tlsConfig.ClientCAs = caPool
		}
		if usage&TLSClient != 0 {
			tlsConfig.RootCAs = caPool
		}
	}

	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

// IssueData is a structure that is suitable for marshaling into a request;
// either via JSON, or into a map[string]interface{} via the structs package
type IssueData struct {
	Lease      string `json:"lease" structs:"lease" mapstructure:"lease"`
	CommonName string `json:"common_name" structs:"common_name" mapstructure:"common_name"`
	AltNames   string `json:"alt_names" structs:"alt_names" mapstructure:"alt_names"`
	IPSANs     string `json:"ip_sans" structs:"ip_sans" mapstructure:"ip_sans"`
}
