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
	"strings"
)

// Secret is used to attempt to unmarshal a Vault secret
// JSON response, as a convenience
type Secret struct {
	Data map[string]interface{} `json:"data"`
}

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

	TLSUnknown TLSUsage = 0
	TLSServer  TLSUsage = 1 << iota
	TLSClient
)

// UserError represents an error generated due to invalid user input
type UserError struct {
	Err string
}

func (e UserError) Error() string {
	return e.Err
}

// InternalError represents an error generated internally,
// presumably not due to invalid user input
type InternalError struct {
	Err string
}

func (e InternalError) Error() string {
	return e.Err
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
// and a DER-encoded certificate
type ParsedCertBundle struct {
	PrivateKeyType   int
	PrivateKeyBytes  []byte
	PrivateKey       crypto.Signer
	IssuingCABytes   []byte
	IssuingCA        *x509.Certificate
	CertificateBytes []byte
	Certificate      *x509.Certificate
}

// CSRBundle contains a key type, a PEM-encoded private key,
// and a PEM-encoded CSR
type CSRBundle struct {
	PrivateKeyType string `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	CSR            string `json:"csr" structs:"csr" mapstructure:"csr"`
	PrivateKey     string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

// ParsedCSRBundle contains a key type, a DER-encoded private key,
// and a DER-encoded certificate request
type ParsedCSRBundle struct {
	PrivateKeyType  int
	PrivateKeyBytes []byte
	PrivateKey      crypto.Signer
	CSRBytes        []byte
	CSR             *x509.CertificateRequest
}

// ToParsedCertBundle converts a string-based certificate bundle
// to a byte-based raw certificate bundle
func (c *CertBundle) ToParsedCertBundle() (*ParsedCertBundle, error) {
	result := &ParsedCertBundle{}
	var err error
	var pemBlock *pem.Block

	if len(c.PrivateKey) > 0 {
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
	}

	if len(c.Certificate) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.Certificate))
		if pemBlock == nil {
			return nil, UserError{"Error decoding certificate from cert bundle"}
		}
		result.CertificateBytes = pemBlock.Bytes
		result.Certificate, err = x509.ParseCertificate(result.CertificateBytes)
		if err != nil {
			return nil, UserError{"Error encountered parsing certificate bytes from raw bundle"}
		}
	}

	if len(c.IssuingCA) > 0 {
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

	if len(c.SerialNumber) == 0 && len(c.Certificate) > 0 {
		c.SerialNumber = GetOctalFormatted(result.Certificate.SerialNumber.Bytes(), ":")
	}

	return result, nil
}

// ToCertBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (p *ParsedCertBundle) ToCertBundle() (*CertBundle, error) {
	result := &CertBundle{}
	block := pem.Block{
		Type: "CERTIFICATE",
	}

	if p.Certificate != nil {
		result.SerialNumber = strings.TrimSpace(GetOctalFormatted(p.Certificate.SerialNumber.Bytes(), ":"))
	}

	if p.CertificateBytes != nil && len(p.CertificateBytes) > 0 {
		block.Bytes = p.CertificateBytes
		result.Certificate = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	if p.IssuingCABytes != nil && len(p.IssuingCABytes) > 0 {
		block.Bytes = p.IssuingCABytes
		result.IssuingCA = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	if p.PrivateKeyBytes != nil && len(p.PrivateKeyBytes) > 0 {
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
		result.PrivateKey = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	return result, nil
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this ParsedCertBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (p *ParsedCertBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error

	if p.PrivateKeyBytes == nil || len(p.PrivateKeyBytes) == 0 {
		return nil, UserError{"Given parsed cert bundle does not have private key information"}
	}

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

// ToParsedCSRBundle converts a string-based CSR bundle
// to a byte-based raw CSR bundle
func (c *CSRBundle) ToParsedCSRBundle() (*ParsedCSRBundle, error) {
	result := &ParsedCSRBundle{}
	var err error
	var pemBlock *pem.Block

	if len(c.PrivateKey) > 0 {
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
	}

	if len(c.CSR) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.CSR))
		if pemBlock == nil {
			return nil, UserError{"Error decoding certificate from cert bundle"}
		}
		result.CSRBytes = pemBlock.Bytes
		result.CSR, err = x509.ParseCertificateRequest(result.CSRBytes)
		if err != nil {
			return nil, UserError{"Error encountered parsing certificate bytes from raw bundle"}
		}
	}

	return result, nil
}

// ToCSRBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (p *ParsedCSRBundle) ToCSRBundle() (*CSRBundle, error) {
	result := &CSRBundle{}
	block := pem.Block{
		Type: "CERTIFICATE REQUEST",
	}

	if p.CSRBytes != nil && len(p.CSRBytes) > 0 {
		block.Bytes = p.CSRBytes
		result.CSR = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	if p.PrivateKeyBytes != nil && len(p.PrivateKeyBytes) > 0 {
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
		result.PrivateKey = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	return result, nil
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this ParsedCSRBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (p *ParsedCSRBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error

	if p.PrivateKeyBytes == nil || len(p.PrivateKeyBytes) == 0 {
		return nil, UserError{"Given parsed cert bundle does not have private key information"}
	}

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
	tlsCert := tls.Certificate{
		Certificate: [][]byte{},
	}

	tlsConfig := &tls.Config{
		NextProtos: []string{"http/1.1"},
	}

	if p.Certificate != nil {
		tlsCert.Leaf = p.Certificate
	}

	if p.PrivateKey != nil {
		tlsCert.PrivateKey = p.PrivateKey
	}

	if p.CertificateBytes != nil && len(p.CertificateBytes) > 0 {
		tlsCert.Certificate = append(tlsCert.Certificate, p.CertificateBytes)
	}

	if p.IssuingCABytes != nil && len(p.IssuingCABytes) > 0 {
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

		if usage&TLSServer > 0 {
			tlsConfig.ClientCAs = caPool
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		}
		if usage&TLSClient > 0 {
			tlsConfig.RootCAs = caPool
		}
	}

	if tlsCert.Certificate != nil && len(tlsCert.Certificate) > 0 {
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsConfig.BuildNameToCertificate()
	}

	return tlsConfig, nil
}

// IssueData is a structure that is suitable for marshaling into a request;
// either via JSON, or into a map[string]interface{} via the structs package
type IssueData struct {
	TTL        string `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	CommonName string `json:"common_name" structs:"common_name" mapstructure:"common_name"`
	AltNames   string `json:"alt_names" structs:"alt_names" mapstructure:"alt_names"`
	IPSANs     string `json:"ip_sans" structs:"ip_sans" mapstructure:"ip_sans"`
}
