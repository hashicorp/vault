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
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
)

// The type of of the Private Key referenced in CertBundle
// and RawCertBundle. This uses colloquial names rather than
// official names, to eliminate confusion
const (
	UnknownPrivateKeyType = iota
	RSAPrivateKeyType
	ECPrivateKeyType
)

// CertBundle contains a key type, a private key,
// a certificate, and a string-encoded serial number,
// returned from a successful Issue request
type CertBundle struct {
	PrivateKeyType string `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	Certificate    string `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	IssuingCA      string `json:"issuing_ca" structs:"issuing_ca" mapstructure:"issuing_ca"`
	PrivateKey     string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	SerialNumber   string `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

// ToRawCertBundle converts a string-based certificate bundle
// to a byte-based raw certificate bundle
func (c *CertBundle) ToRawCertBundle() (*RawCertBundle, error) {
	result := &RawCertBundle{}
	switch c.PrivateKeyType {
	case "ec":
		result.PrivateKeyType = ECPrivateKeyType
	case "rsa":
		result.PrivateKeyType = RSAPrivateKeyType
	default:
		return nil, fmt.Errorf("Unknown private key type in bundle: %s", c.PrivateKeyType)
	}

	var pemBlock *pem.Block
	pemBlock, _ = pem.Decode([]byte(c.PrivateKey))
	if pemBlock == nil {
		return nil, fmt.Errorf("Error decoding private key from cert bundle")
	}
	result.PrivateKeyBytes = pemBlock.Bytes

	pemBlock, _ = pem.Decode([]byte(c.Certificate))
	if pemBlock == nil {
		return nil, fmt.Errorf("Error decoding certificate from cert bundle")
	}
	result.CertificateBytes = pemBlock.Bytes

	if len(c.IssuingCA) != 0 {
		pemBlock, _ = pem.Decode([]byte(c.IssuingCA))
		if pemBlock == nil {
			return nil, fmt.Errorf("Error decoding issuing CA from cert bundle")
		}
		result.IssuingCABytes = pemBlock.Bytes
	}

	if err := result.populateSerialNumber(); err != nil {
		return nil, err
	}

	return result, nil
}

// RawCertBundle contains a key type, a DER-encoded private key,
// a DER-encoded certificate, and a big.Int serial number
type RawCertBundle struct {
	PrivateKeyType   int
	PrivateKeyBytes  []byte
	IssuingCABytes   []byte
	CertificateBytes []byte
	SerialNumber     *big.Int
}

// ToCertBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (r *RawCertBundle) ToCertBundle() (*CertBundle, error) {
	result := &CertBundle{
		SerialNumber: GetOctalFormatted(r.SerialNumber.Bytes(), ":"),
	}

	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: r.CertificateBytes,
	}
	result.Certificate = string(pem.EncodeToMemory(&block))

	if len(r.IssuingCABytes) != 0 {
		block.Bytes = r.IssuingCABytes
		result.IssuingCA = string(pem.EncodeToMemory(&block))
	}

	block.Bytes = r.PrivateKeyBytes
	switch r.PrivateKeyType {
	case RSAPrivateKeyType:
		result.PrivateKeyType = "rsa"
		block.Type = "RSA PRIVATE KEY"
	case ECPrivateKeyType:
		result.PrivateKeyType = "ec"
		block.Type = "EC PRIVATE KEY"
	default:
		return nil, fmt.Errorf("Could not determine private key type when creating block")
	}
	result.PrivateKey = string(pem.EncodeToMemory(&block))

	return result, nil
}

func (r *RawCertBundle) populateSerialNumber() error {
	cert, err := x509.ParseCertificate(r.CertificateBytes)
	if err != nil {
		return fmt.Errorf("Error encountered parsing certificate bytes from raw bundle")
	}
	r.SerialNumber = cert.SerialNumber
	return nil
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this RawCertBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (r *RawCertBundle) GetSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error
	switch r.PrivateKeyType {
	case ECPrivateKeyType:
		signer, err = x509.ParseECPrivateKey(r.PrivateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse CA's private EC key: %s", err)
		}
	case RSAPrivateKeyType:
		signer, err = x509.ParsePKCS1PrivateKey(r.PrivateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse CA's private RSA key: %s", err)
		}
	default:
		return nil, fmt.Errorf("Unable to determine the type of private key")
	}
	return signer, nil
}

// GetSubjKeyID returns the subject key ID, e.g. the SHA1 sum
// of the marshaled public key
func (r *RawCertBundle) GetSubjKeyID() ([]byte, error) {
	privateKey, err := r.GetSigner()
	if err != nil {
		return nil, err
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("Error marshalling public key: %s", err)
	}

	subjKeyID := sha1.Sum(marshaledKey)

	return subjKeyID[:], nil
}

// IssueData is a structure that is suitable for marshaling into a request;
// either via JSON, or into a map[string]interface{} via the structs package
type IssueData struct {
	Lease      string `json:"lease" structs:"lease" mapstructure:"lease"`
	CommonName string `json:"common_name" structs:"common_name" mapstructure:"common_name"`
	AltNames   string `json:"alt_names" structs:"alt_names" mapstructure:"alt_names"`
	IPSANs     string `json:"ip_sans" structs:"ip_sans" mapstructure:"ip_sans"`
}
