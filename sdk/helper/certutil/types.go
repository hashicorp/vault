// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/errutil"
)

const (
	PrivateKeyTypeP521 = "p521"
)

// This can be one of a few key types so the different params may or may not be filled
type ClusterKeyParams struct {
	Type string   `json:"type" structs:"type" mapstructure:"type"`
	X    *big.Int `json:"x" structs:"x" mapstructure:"x"`
	Y    *big.Int `json:"y" structs:"y" mapstructure:"y"`
	D    *big.Int `json:"d" structs:"d" mapstructure:"d"`
}

// Secret is used to attempt to unmarshal a Vault secret
// JSON response, as a convenience
type Secret struct {
	Data map[string]interface{} `json:"data"`
}

// PrivateKeyType holds a string representation of the type of private key (ec
// or rsa) referenced in CertBundle and ParsedCertBundle. This uses colloquial
// names rather than official names, to eliminate confusion
type PrivateKeyType string

// Well-known PrivateKeyTypes
const (
	UnknownPrivateKey PrivateKeyType = ""
	RSAPrivateKey     PrivateKeyType = "rsa"
	ECPrivateKey      PrivateKeyType = "ec"
	Ed25519PrivateKey PrivateKeyType = "ed25519"
	ManagedPrivateKey PrivateKeyType = "ManagedPrivateKey"
)

// TLSUsage controls whether the intended usage of a *tls.Config
// returned from ParsedCertBundle.getTLSConfig is for server use,
// client use, or both, which affects which values are set
type TLSUsage int

// Well-known TLSUsage types
const (
	TLSUnknown TLSUsage = 0
	TLSServer  TLSUsage = 1 << iota
	TLSClient
)

// BlockType indicates the serialization format of the key
type BlockType string

// Well-known formats
const (
	UnknownBlock BlockType = ""
	PKCS1Block   BlockType = "RSA PRIVATE KEY"
	PKCS8Block   BlockType = "PRIVATE KEY"
	ECBlock      BlockType = "EC PRIVATE KEY"
)

// ParsedPrivateKeyContainer allows common key setting for certs and CSRs
type ParsedPrivateKeyContainer interface {
	SetParsedPrivateKey(crypto.Signer, PrivateKeyType, []byte)
}

// CertBlock contains the DER-encoded certificate and the PEM
// block's byte array
type CertBlock struct {
	Certificate *x509.Certificate
	Bytes       []byte
}

// CertBundle contains a key type, a PEM-encoded private key,
// a PEM-encoded certificate, and a string-encoded serial number,
// returned from a successful Issue request
type CertBundle struct {
	PrivateKeyType PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	Certificate    string         `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	IssuingCA      string         `json:"issuing_ca" structs:"issuing_ca" mapstructure:"issuing_ca"`
	CAChain        []string       `json:"ca_chain" structs:"ca_chain" mapstructure:"ca_chain"`
	PrivateKey     string         `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	SerialNumber   string         `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

// ParsedCertBundle contains a key type, a DER-encoded private key,
// and a DER-encoded certificate
type ParsedCertBundle struct {
	PrivateKeyType   PrivateKeyType
	PrivateKeyFormat BlockType
	PrivateKeyBytes  []byte
	PrivateKey       crypto.Signer
	CertificateBytes []byte
	Certificate      *x509.Certificate
	CAChain          []*CertBlock
}

// CSRBundle contains a key type, a PEM-encoded private key,
// and a PEM-encoded CSR
type CSRBundle struct {
	PrivateKeyType PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	CSR            string         `json:"csr" structs:"csr" mapstructure:"csr"`
	PrivateKey     string         `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

// ParsedCSRBundle contains a key type, a DER-encoded private key,
// and a DER-encoded certificate request
type ParsedCSRBundle struct {
	PrivateKeyType  PrivateKeyType
	PrivateKeyBytes []byte
	PrivateKey      crypto.Signer
	CSRBytes        []byte
	CSR             *x509.CertificateRequest
}

type KeyBundle struct {
	PrivateKeyType  PrivateKeyType
	PrivateKeyBytes []byte
	PrivateKey      crypto.Signer
}

func GetPrivateKeyTypeFromSigner(signer crypto.Signer) PrivateKeyType {
	// We look at the public key types to work-around limitations/typing of managed keys.
	switch signer.Public().(type) {
	case *rsa.PublicKey:
		return RSAPrivateKey
	case *ecdsa.PublicKey:
		return ECPrivateKey
	case ed25519.PublicKey:
		return Ed25519PrivateKey
	}
	return UnknownPrivateKey
}

// GetPrivateKeyTypeFromPublicKey based on the public key, return the PrivateKeyType
// that would be associated with it, returning UnknownPrivateKey for unsupported types
func GetPrivateKeyTypeFromPublicKey(pubKey crypto.PublicKey) PrivateKeyType {
	switch pubKey.(type) {
	case *rsa.PublicKey:
		return RSAPrivateKey
	case *ecdsa.PublicKey:
		return ECPrivateKey
	case ed25519.PublicKey:
		return Ed25519PrivateKey
	default:
		return UnknownPrivateKey
	}
}

// ToPEMBundle converts a string-based certificate bundle
// to a PEM-based string certificate bundle in trust path
// order, leaf certificate first
func (c *CertBundle) ToPEMBundle() string {
	var result []string

	if len(c.PrivateKey) > 0 {
		result = append(result, c.PrivateKey)
	}
	if len(c.Certificate) > 0 {
		result = append(result, c.Certificate)
	}
	if len(c.CAChain) > 0 {
		result = append(result, c.CAChain...)
	}

	return strings.Join(result, "\n")
}

// ToParsedCertBundle converts a string-based certificate bundle
// to a byte-based raw certificate bundle
func (c *CertBundle) ToParsedCertBundle() (*ParsedCertBundle, error) {
	return c.ToParsedCertBundleWithExtractor(extractAndSetPrivateKey)
}

// PrivateKeyExtractor extract out a private key from the passed in
// CertBundle and set the appropriate bits within the ParsedCertBundle.
type PrivateKeyExtractor func(c *CertBundle, parsedBundle *ParsedCertBundle) error

func (c *CertBundle) ToParsedCertBundleWithExtractor(privateKeyExtractor PrivateKeyExtractor) (*ParsedCertBundle, error) {
	var err error
	var pemBlock *pem.Block
	result := &ParsedCertBundle{}

	err = privateKeyExtractor(c, result)
	if err != nil {
		return nil, err
	}

	if len(c.Certificate) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.Certificate))
		if pemBlock == nil {
			return nil, errutil.UserError{Err: "Error decoding certificate from cert bundle"}
		}
		result.CertificateBytes = pemBlock.Bytes
		result.Certificate, err = x509.ParseCertificate(result.CertificateBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Error encountered parsing certificate bytes from raw bundle: %v", err)}
		}
	}
	switch {
	case len(c.CAChain) > 0:
		for _, cert := range c.CAChain {
			pemBlock, _ := pem.Decode([]byte(cert))
			if pemBlock == nil {
				return nil, errutil.UserError{Err: "Error decoding certificate from cert bundle"}
			}

			parsedCert, err := x509.ParseCertificate(pemBlock.Bytes)
			if err != nil {
				return nil, errutil.UserError{Err: fmt.Sprintf("Error encountered parsing certificate bytes from raw bundle via CA chain: %v", err)}
			}

			certBlock := &CertBlock{
				Bytes:       pemBlock.Bytes,
				Certificate: parsedCert,
			}
			result.CAChain = append(result.CAChain, certBlock)
		}

	// For backwards compatibility
	case len(c.IssuingCA) > 0:
		pemBlock, _ = pem.Decode([]byte(c.IssuingCA))
		if pemBlock == nil {
			return nil, errutil.UserError{Err: "Error decoding ca certificate from cert bundle"}
		}

		parsedCert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Error encountered parsing certificate bytes from raw bundle via issuing CA: %v", err)}
		}

		certBlock := &CertBlock{
			Bytes:       pemBlock.Bytes,
			Certificate: parsedCert,
		}
		result.CAChain = append(result.CAChain, certBlock)
	}

	// Populate if it isn't there already
	if len(c.SerialNumber) == 0 && len(c.Certificate) > 0 {
		c.SerialNumber = GetHexFormatted(result.Certificate.SerialNumber.Bytes(), ":")
	}

	return result, nil
}

func extractAndSetPrivateKey(c *CertBundle, parsedBundle *ParsedCertBundle) error {
	if len(c.PrivateKey) == 0 {
		return nil
	}

	pemBlock, _ := pem.Decode([]byte(c.PrivateKey))
	if pemBlock == nil {
		return errutil.UserError{Err: "Error decoding private key from cert bundle"}
	}

	parsedBundle.PrivateKeyBytes = pemBlock.Bytes
	parsedBundle.PrivateKeyFormat = BlockType(strings.TrimSpace(pemBlock.Type))

	switch parsedBundle.PrivateKeyFormat {
	case ECBlock:
		parsedBundle.PrivateKeyType, c.PrivateKeyType = ECPrivateKey, ECPrivateKey
	case PKCS1Block:
		c.PrivateKeyType, parsedBundle.PrivateKeyType = RSAPrivateKey, RSAPrivateKey
	case PKCS8Block:
		t, err := getPKCS8Type(pemBlock.Bytes)
		if err != nil {
			return errutil.UserError{Err: fmt.Sprintf("Error getting key type from pkcs#8: %v", err)}
		}
		parsedBundle.PrivateKeyType = t
		switch t {
		case ECPrivateKey:
			c.PrivateKeyType = ECPrivateKey
		case RSAPrivateKey:
			c.PrivateKeyType = RSAPrivateKey
		case Ed25519PrivateKey:
			c.PrivateKeyType = Ed25519PrivateKey
		case ManagedPrivateKey:
			c.PrivateKeyType = ManagedPrivateKey
		}
	default:
		return errutil.UserError{Err: fmt.Sprintf("Unsupported key block type: %s", pemBlock.Type)}
	}

	var err error
	parsedBundle.PrivateKey, err = parsedBundle.getSigner()
	if err != nil {
		return errutil.UserError{Err: fmt.Sprintf("Error getting signer: %s", err)}
	}
	return nil
}

// ToCertBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (p *ParsedCertBundle) ToCertBundle() (*CertBundle, error) {
	result := &CertBundle{}
	block := pem.Block{
		Type: "CERTIFICATE",
	}

	if p.Certificate != nil {
		result.SerialNumber = strings.TrimSpace(GetHexFormatted(p.Certificate.SerialNumber.Bytes(), ":"))
	}

	if p.CertificateBytes != nil && len(p.CertificateBytes) > 0 {
		block.Bytes = p.CertificateBytes
		result.Certificate = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	for _, caCert := range p.CAChain {
		block.Bytes = caCert.Bytes
		certificate := strings.TrimSpace(string(pem.EncodeToMemory(&block)))

		result.CAChain = append(result.CAChain, certificate)
	}

	if p.PrivateKeyBytes != nil && len(p.PrivateKeyBytes) > 0 {
		block.Type = string(p.PrivateKeyFormat)
		block.Bytes = p.PrivateKeyBytes
		result.PrivateKeyType = p.PrivateKeyType

		// Handle bundle not parsed by us
		if block.Type == "" {
			switch p.PrivateKeyType {
			case ECPrivateKey:
				block.Type = string(ECBlock)
			case RSAPrivateKey:
				block.Type = string(PKCS1Block)
			case Ed25519PrivateKey:
				block.Type = string(PKCS8Block)
			}
		}

		result.PrivateKey = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	return result, nil
}

// Verify checks if the parsed bundle is valid.  It validates the public
// key of the certificate to the private key and checks the certificate trust
// chain for path issues.
func (p *ParsedCertBundle) Verify() error {
	// If private key exists, check if it matches the public key of cert
	if p.PrivateKey != nil && p.Certificate != nil {
		equal, err := ComparePublicKeys(p.Certificate.PublicKey, p.PrivateKey.Public())
		if err != nil {
			return errwrap.Wrapf("could not compare public and private keys: {{err}}", err)
		}
		if !equal {
			return fmt.Errorf("public key of certificate does not match private key")
		}
	}

	certPath := p.GetCertificatePath()
	if len(certPath) > 1 {
		for i, caCert := range certPath[1:] {
			if !caCert.Certificate.IsCA {
				return fmt.Errorf("certificate %d of certificate chain is not a certificate authority", i+1)
			}
			if !bytes.Equal(certPath[i].Certificate.AuthorityKeyId, caCert.Certificate.SubjectKeyId) {
				return fmt.Errorf("certificate %d of certificate chain ca trust path is incorrect (%q/%q) (%X/%X)",
					i+1,
					certPath[i].Certificate.Subject.CommonName, caCert.Certificate.Subject.CommonName,
					certPath[i].Certificate.AuthorityKeyId, caCert.Certificate.SubjectKeyId)
			}
		}
	}

	return nil
}

// GetCertificatePath returns a slice of certificates making up a path, pulled
// from the parsed cert bundle
func (p *ParsedCertBundle) GetCertificatePath() []*CertBlock {
	var certPath []*CertBlock

	certPath = append(certPath, &CertBlock{
		Certificate: p.Certificate,
		Bytes:       p.CertificateBytes,
	})

	if len(p.CAChain) > 0 {
		// Root CA puts itself in the chain
		if p.CAChain[0].Certificate.SerialNumber != p.Certificate.SerialNumber {
			certPath = append(certPath, p.CAChain...)
		}
	}

	return certPath
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this ParsedCertBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (p *ParsedCertBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error

	if p.PrivateKeyBytes == nil || len(p.PrivateKeyBytes) == 0 {
		return nil, errutil.UserError{Err: "Given parsed cert bundle does not have private key information"}
	}

	switch p.PrivateKeyFormat {
	case ECBlock:
		signer, err = x509.ParseECPrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Unable to parse CA's private EC key: %s", err)}
		}

	case PKCS1Block:
		signer, err = x509.ParsePKCS1PrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Unable to parse CA's private RSA key: %s", err)}
		}

	case PKCS8Block:
		if k, err := x509.ParsePKCS8PrivateKey(p.PrivateKeyBytes); err == nil {
			switch k := k.(type) {
			case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
				return k.(crypto.Signer), nil
			default:
				return nil, errutil.UserError{Err: "Found unknown private key type in pkcs#8 wrapping"}
			}
		}
		return nil, errutil.UserError{Err: fmt.Sprintf("Failed to parse pkcs#8 key: %v", err)}
	default:
		return nil, errutil.UserError{Err: "Unable to determine type of private key; only RSA and EC are supported"}
	}
	return signer, nil
}

// SetParsedPrivateKey sets the private key parameters on the bundle
func (p *ParsedCertBundle) SetParsedPrivateKey(privateKey crypto.Signer, privateKeyType PrivateKeyType, privateKeyBytes []byte) {
	p.PrivateKey = privateKey
	p.PrivateKeyType = privateKeyType
	p.PrivateKeyBytes = privateKeyBytes
}

func getPKCS8Type(bs []byte) (PrivateKeyType, error) {
	k, err := x509.ParsePKCS8PrivateKey(bs)
	if err != nil {
		return UnknownPrivateKey, errutil.UserError{Err: fmt.Sprintf("Failed to parse pkcs#8 key: %v", err)}
	}

	switch k.(type) {
	case *ecdsa.PrivateKey:
		return ECPrivateKey, nil
	case *rsa.PrivateKey:
		return RSAPrivateKey, nil
	case ed25519.PrivateKey:
		return Ed25519PrivateKey, nil
	default:
		return UnknownPrivateKey, errutil.UserError{Err: "Found unknown private key type in pkcs#8 wrapping"}
	}
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
			return nil, errutil.UserError{Err: "Error decoding private key from cert bundle"}
		}
		result.PrivateKeyBytes = pemBlock.Bytes

		switch BlockType(pemBlock.Type) {
		case ECBlock:
			result.PrivateKeyType = ECPrivateKey
		case PKCS1Block:
			result.PrivateKeyType = RSAPrivateKey
		default:
			// Try to figure it out and correct
			if _, err := x509.ParseECPrivateKey(pemBlock.Bytes); err == nil {
				result.PrivateKeyType = ECPrivateKey
				c.PrivateKeyType = "ec"
			} else if _, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err == nil {
				result.PrivateKeyType = RSAPrivateKey
				c.PrivateKeyType = "rsa"
			} else if _, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes); err == nil {
				result.PrivateKeyType = Ed25519PrivateKey
				c.PrivateKeyType = "ed25519"
			} else {
				return nil, errutil.UserError{Err: fmt.Sprintf("Unknown private key type in bundle: %s", c.PrivateKeyType)}
			}
		}

		result.PrivateKey, err = result.getSigner()
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Error getting signer: %s", err)}
		}
	}

	if len(c.CSR) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.CSR))
		if pemBlock == nil {
			return nil, errutil.UserError{Err: "Error decoding certificate from cert bundle"}
		}
		result.CSRBytes = pemBlock.Bytes
		result.CSR, err = x509.ParseCertificateRequest(result.CSRBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Error encountered parsing certificate bytes from raw bundle via CSR: %v", err)}
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
		case Ed25519PrivateKey:
			result.PrivateKeyType = "ed25519"
			block.Type = "PRIVATE KEY"
		case ManagedPrivateKey:
			result.PrivateKeyType = ManagedPrivateKey
			block.Type = "PRIVATE KEY"
		default:
			return nil, errutil.InternalError{Err: "Could not determine private key type when creating block"}
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
		return nil, errutil.UserError{Err: "Given parsed cert bundle does not have private key information"}
	}

	switch p.PrivateKeyType {
	case ECPrivateKey:
		signer, err = x509.ParseECPrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Unable to parse CA's private EC key: %s", err)}
		}

	case RSAPrivateKey:
		signer, err = x509.ParsePKCS1PrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Unable to parse CA's private RSA key: %s", err)}
		}

	case Ed25519PrivateKey:
		signerd, err := x509.ParsePKCS8PrivateKey(p.PrivateKeyBytes)
		signer = signerd.(ed25519.PrivateKey)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("Unable to parse CA's private Ed25519 key: %s", err)}
		}

	default:
		return nil, errutil.UserError{Err: "Unable to determine type of private key; only RSA, Ed25519 and EC are supported"}
	}
	return signer, nil
}

// SetParsedPrivateKey sets the private key parameters on the bundle
func (p *ParsedCSRBundle) SetParsedPrivateKey(privateKey crypto.Signer, privateKeyType PrivateKeyType, privateKeyBytes []byte) {
	p.PrivateKey = privateKey
	p.PrivateKeyType = privateKeyType
	p.PrivateKeyBytes = privateKeyBytes
}

// getTLSConfig returns a TLS config generally suitable for client
// authentication. The returned TLS config can be modified slightly
// to be made suitable for a server requiring client authentication;
// specifically, you should set the value of ClientAuth in the returned
// config to match your needs.
func (p *ParsedCertBundle) GetTLSConfig(usage TLSUsage) (*tls.Config, error) {
	tlsCert := tls.Certificate{
		Certificate: [][]byte{},
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
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

	if len(p.CAChain) > 0 {
		for _, cert := range p.CAChain {
			tlsCert.Certificate = append(tlsCert.Certificate, cert.Bytes)
		}

		// Technically we only need one cert, but this doesn't duplicate code
		certBundle, err := p.ToCertBundle()
		if err != nil {
			return nil, errwrap.Wrapf("error converting parsed bundle to string bundle when getting TLS config: {{err}}", err)
		}

		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM([]byte(certBundle.CAChain[0]))
		if !ok {
			return nil, fmt.Errorf("could not append CA certificate")
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
	}

	return tlsConfig, nil
}

// IssueData is a structure that is suitable for marshaling into a request;
// either via JSON, or into a map[string]interface{} via the structs package
type IssueData struct {
	TTL        string `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	CommonName string `json:"common_name" structs:"common_name" mapstructure:"common_name"`
	OU         string `json:"ou" structs:"ou" mapstructure:"ou"`
	AltNames   string `json:"alt_names" structs:"alt_names" mapstructure:"alt_names"`
	IPSANs     string `json:"ip_sans" structs:"ip_sans" mapstructure:"ip_sans"`
	CSR        string `json:"csr" structs:"csr" mapstructure:"csr"`
	OtherSANs  string `json:"other_sans" structs:"other_sans" mapstructure:"other_sans"`
}

type URLEntries struct {
	IssuingCertificates   []string `json:"issuing_certificates" structs:"issuing_certificates" mapstructure:"issuing_certificates"`
	CRLDistributionPoints []string `json:"crl_distribution_points" structs:"crl_distribution_points" mapstructure:"crl_distribution_points"`
	OCSPServers           []string `json:"ocsp_servers" structs:"ocsp_servers" mapstructure:"ocsp_servers"`
}

type NotAfterBehavior int

const (
	ErrNotAfterBehavior NotAfterBehavior = iota
	TruncateNotAfterBehavior
	PermitNotAfterBehavior
	AlwaysEnforceErr
)

var notAfterBehaviorNames = map[NotAfterBehavior]string{
	AlwaysEnforceErr:         "always_enforce_err",
	ErrNotAfterBehavior:      "err",
	TruncateNotAfterBehavior: "truncate",
	PermitNotAfterBehavior:   "permit",
}

func (n NotAfterBehavior) String() string {
	if name, ok := notAfterBehaviorNames[n]; ok && len(name) > 0 {
		return name
	}

	return "unknown"
}

type CAInfoBundle struct {
	ParsedCertBundle
	URLs                 *URLEntries
	LeafNotAfterBehavior NotAfterBehavior
	RevocationSigAlg     x509.SignatureAlgorithm
}

func (b *CAInfoBundle) GetCAChain() []*CertBlock {
	chain := []*CertBlock{}

	// Include issuing CA in Chain, not including Root Authority
	if (len(b.Certificate.AuthorityKeyId) > 0 &&
		!bytes.Equal(b.Certificate.AuthorityKeyId, b.Certificate.SubjectKeyId)) ||
		(len(b.Certificate.AuthorityKeyId) == 0 &&
			!bytes.Equal(b.Certificate.RawIssuer, b.Certificate.RawSubject)) {

		chain = b.GetFullChain()
	}

	return chain
}

func (b *CAInfoBundle) GetFullChain() []*CertBlock {
	var chain []*CertBlock

	// Some bundles already include the root included in the chain,
	// so don't include it twice.
	if len(b.CAChain) == 0 || !bytes.Equal(b.CAChain[0].Bytes, b.CertificateBytes) {
		chain = append(chain, &CertBlock{
			Certificate: b.Certificate,
			Bytes:       b.CertificateBytes,
		})
	}

	if len(b.CAChain) > 0 {
		chain = append(chain, b.CAChain...)
	}

	return chain
}

type CertExtKeyUsage int

const (
	AnyExtKeyUsage CertExtKeyUsage = 1 << iota
	ServerAuthExtKeyUsage
	ClientAuthExtKeyUsage
	CodeSigningExtKeyUsage
	EmailProtectionExtKeyUsage
	IpsecEndSystemExtKeyUsage
	IpsecTunnelExtKeyUsage
	IpsecUserExtKeyUsage
	TimeStampingExtKeyUsage
	OcspSigningExtKeyUsage
	MicrosoftServerGatedCryptoExtKeyUsage
	NetscapeServerGatedCryptoExtKeyUsage
	MicrosoftCommercialCodeSigningExtKeyUsage
	MicrosoftKernelCodeSigningExtKeyUsage
)

type CreationParameters struct {
	Subject                       pkix.Name
	DNSNames                      []string
	EmailAddresses                []string
	IPAddresses                   []net.IP
	URIs                          []*url.URL
	OtherSANs                     map[string][]string
	IsCA                          bool
	KeyType                       string
	KeyBits                       int
	NotAfter                      time.Time
	KeyUsage                      x509.KeyUsage
	ExtKeyUsage                   CertExtKeyUsage
	ExtKeyUsageOIDs               []string
	PolicyIdentifiers             []string
	BasicConstraintsValidForNonCA bool
	SignatureBits                 int
	UsePSS                        bool
	ForceAppendCaChain            bool

	// Only used when signing a CA cert
	UseCSRValues            bool
	PermittedDNSDomains     []string
	ExcludedDNSDomains      []string
	PermittedIPRanges       []*net.IPNet
	ExcludedIPRanges        []*net.IPNet
	PermittedEmailAddresses []string
	ExcludedEmailAddresses  []string
	PermittedURIDomains     []string
	ExcludedURIDomains      []string

	// URLs to encode into the certificate
	URLs *URLEntries

	// The maximum path length to encode
	MaxPathLength int

	// The duration the certificate will use NotBefore
	NotBeforeDuration time.Duration

	// The explicit SKID to use; especially useful for cross-signing.
	SKID []byte

	// Ignore validating the CSR's signature. This should only be enabled if the
	// sender of the CSR has proven proof of possession of the associated
	// private key by some other means, otherwise keep this set to false.
	IgnoreCSRSignature bool
}

type CreationBundle struct {
	Params        *CreationParameters
	SigningBundle *CAInfoBundle
	CSR           *x509.CertificateRequest
}

// addKeyUsages adds appropriate key usages to the template given the creation
// information
func AddKeyUsages(data *CreationBundle, certTemplate *x509.Certificate) {
	if data.Params.IsCA {
		certTemplate.KeyUsage = x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign)
		return
	}

	certTemplate.KeyUsage = data.Params.KeyUsage

	if data.Params.ExtKeyUsage&AnyExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageAny)
	}

	if data.Params.ExtKeyUsage&ServerAuthExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}

	if data.Params.ExtKeyUsage&ClientAuthExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}

	if data.Params.ExtKeyUsage&CodeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageCodeSigning)
	}

	if data.Params.ExtKeyUsage&EmailProtectionExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}

	if data.Params.ExtKeyUsage&IpsecEndSystemExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageIPSECEndSystem)
	}

	if data.Params.ExtKeyUsage&IpsecTunnelExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageIPSECTunnel)
	}

	if data.Params.ExtKeyUsage&IpsecUserExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageIPSECUser)
	}

	if data.Params.ExtKeyUsage&TimeStampingExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageTimeStamping)
	}

	if data.Params.ExtKeyUsage&OcspSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageOCSPSigning)
	}

	if data.Params.ExtKeyUsage&MicrosoftServerGatedCryptoExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageMicrosoftServerGatedCrypto)
	}

	if data.Params.ExtKeyUsage&NetscapeServerGatedCryptoExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageNetscapeServerGatedCrypto)
	}

	if data.Params.ExtKeyUsage&MicrosoftCommercialCodeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageMicrosoftCommercialCodeSigning)
	}

	if data.Params.ExtKeyUsage&MicrosoftKernelCodeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageMicrosoftKernelCodeSigning)
	}
}

// SetParsedPrivateKey sets the private key parameters on the bundle
func (p *KeyBundle) SetParsedPrivateKey(privateKey crypto.Signer, privateKeyType PrivateKeyType, privateKeyBytes []byte) {
	p.PrivateKey = privateKey
	p.PrivateKeyType = privateKeyType
	p.PrivateKeyBytes = privateKeyBytes
}

func (p *KeyBundle) ToPrivateKeyPemString() (string, error) {
	block := pem.Block{}

	if p.PrivateKeyBytes != nil && len(p.PrivateKeyBytes) > 0 {
		block.Bytes = p.PrivateKeyBytes
		switch p.PrivateKeyType {
		case RSAPrivateKey:
			block.Type = "RSA PRIVATE KEY"
		case ECPrivateKey:
			block.Type = "EC PRIVATE KEY"
		default:
			block.Type = "PRIVATE KEY"
		}
		privateKeyPemString := strings.TrimSpace(string(pem.EncodeToMemory(&block)))
		return privateKeyPemString, nil
	}

	return "", errutil.InternalError{Err: "No Private Key Bytes to Wrap"}
}

// PolicyIdentifierWithQualifierEntry Structure for Internal Storage
type PolicyIdentifierWithQualifierEntry struct {
	PolicyIdentifierOid string `json:"oid",mapstructure:"oid"`
	CPS                 string `json:"cps,omitempty",mapstructure:"cps"`
	Notice              string `json:"notice,omitempty",mapstructure:"notice"`
}

// GetPolicyIdentifierFromString parses out the internal structure of a Policy Identifier
func GetPolicyIdentifierFromString(policyIdentifier string) (*PolicyIdentifierWithQualifierEntry, error) {
	if policyIdentifier == "" {
		return nil, nil
	}
	entry := &PolicyIdentifierWithQualifierEntry{}
	// Either a OID, or a JSON Entry: First check OID:
	_, err := StringToOid(policyIdentifier)
	if err == nil {
		entry.PolicyIdentifierOid = policyIdentifier
		return entry, nil
	}
	// Now Check If JSON Entry
	jsonErr := json.Unmarshal([]byte(policyIdentifier), &entry)
	if jsonErr != nil { // Neither, if we got here
		return entry, errors.New(fmt.Sprintf("Policy Identifier %q is neither a valid OID: %s, Nor JSON Policy Identifier: %s", policyIdentifier, err.Error(), jsonErr.Error()))
	}
	return entry, nil
}

// Policy Identifier with Qualifier Structure for ASN Marshalling:

var policyInformationOid = asn1.ObjectIdentifier{2, 5, 29, 32}

type policyInformation struct {
	PolicyIdentifier asn1.ObjectIdentifier
	Qualifiers       []interface{} `asn1:"tag:optional,omitempty"`
}

var cpsPolicyQualifierID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 2, 1}

type cpsUrlPolicyQualifier struct {
	PolicyQualifierID asn1.ObjectIdentifier
	Qualifier         string `asn1:"tag:optional,ia5"`
}

var userNoticePolicyQualifierID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 2, 2}

type userNoticePolicyQualifier struct {
	PolicyQualifierID asn1.ObjectIdentifier
	Qualifier         userNotice
}

type userNotice struct {
	ExplicitText string `asn1:"tag:optional,utf8"`
}

func createPolicyIdentifierWithQualifier(entry PolicyIdentifierWithQualifierEntry) (*policyInformation, error) {
	// Each Policy is Identified by a Unique ID, as designated here:
	policyOid, err := StringToOid(entry.PolicyIdentifierOid)
	if err != nil {
		return nil, err
	}
	pi := policyInformation{
		PolicyIdentifier: policyOid,
	}
	if entry.CPS != "" {
		qualifier := cpsUrlPolicyQualifier{
			PolicyQualifierID: cpsPolicyQualifierID,
			Qualifier:         entry.CPS,
		}
		pi.Qualifiers = append(pi.Qualifiers, qualifier)
	}
	if entry.Notice != "" {
		qualifier := userNoticePolicyQualifier{
			PolicyQualifierID: userNoticePolicyQualifierID,
			Qualifier: userNotice{
				ExplicitText: entry.Notice,
			},
		}
		pi.Qualifiers = append(pi.Qualifiers, qualifier)
	}
	return &pi, nil
}

// CreatePolicyInformationExtensionFromStorageStrings parses the stored policyIdentifiers, which might be JSON Policy
// Identifier with Qualifier Entries or String OIDs, and returns an extension if everything parsed correctly, and an
// error if constructing
func CreatePolicyInformationExtensionFromStorageStrings(policyIdentifiers []string) (*pkix.Extension, error) {
	var policyInformationList []policyInformation
	for _, policyIdentifierStr := range policyIdentifiers {
		policyIdentifierEntry, err := GetPolicyIdentifierFromString(policyIdentifierStr)
		if err != nil {
			return nil, err
		}
		if policyIdentifierEntry != nil { // Okay to skip empty entries if there is no error
			policyInformationStruct, err := createPolicyIdentifierWithQualifier(*policyIdentifierEntry)
			if err != nil {
				return nil, err
			}
			policyInformationList = append(policyInformationList, *policyInformationStruct)
		}
	}
	asn1Bytes, err := asn1.Marshal(policyInformationList)
	if err != nil {
		return nil, err
	}
	return &pkix.Extension{
		Id:       policyInformationOid,
		Critical: false,
		Value:    asn1Bytes,
	}, nil
}

// Subject Attribute OIDs
var SubjectPilotUserIDAttributeOID = asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 1}
