package auth

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CertValidationError is returned in case of X09 certificate validation errors.
type CertValidationError struct {
	t    time.Time
	cert *x509.Certificate
}

func (e CertValidationError) Error() string {
	return fmt.Sprintf("certificate issuer %s subject %s not in validity period from %s to %s - now %s",
		e.cert.Issuer.ToRDNSequence().String(),
		e.cert.Subject.ToRDNSequence().String(),
		e.cert.NotBefore,
		e.cert.NotAfter,
		e.t,
	)
}

// CertKey represents a X509 certificate and key.
type CertKey struct {
	cert, key  string // define as string for being immutable
	certBlocks []*pem.Block
	certs      []*x509.Certificate
	keyBlock   *pem.Block
}

// NewCertKey returns a new certificate and key instance.
func NewCertKey(cert, key []byte) (*CertKey, error) {
	certBlocks, err := decodeClientCert(cert)
	if err != nil {
		return nil, err
	}
	certs, err := parseCerts(certBlocks)
	if err != nil {
		return nil, err
	}
	keyBlock, err := decodeClientKey(key)
	if err != nil {
		return nil, err
	}
	return &CertKey{cert: string(cert), key: string(key), certBlocks: certBlocks, certs: certs, keyBlock: keyBlock}, nil
}

func (ck *CertKey) String() string { return fmt.Sprintf("cert %s key %s", ck.cert, ck.key) }

// Equal returns true if the certificate and key equals the instance data, false otherwise.
func (ck *CertKey) Equal(cert, key []byte) bool {
	return string(cert) == ck.cert && string(key) == ck.key
}

// Cert returns the certificate.
func (ck *CertKey) Cert() []byte { return []byte(ck.cert) }

// Key returns the key.
func (ck *CertKey) Key() []byte { return []byte(ck.key) }

// validate validates the certificate (currently validity period only).
func (ck *CertKey) validate(t time.Time) error {
	t = t.UTC() // cert.NotBefore and cert.NotAfter in UTC as well
	for _, cert := range ck.certs {
		// checks
		// .check validity period
		if t.Before(cert.NotBefore) || t.After(cert.NotAfter) {
			return &CertValidationError{t: t, cert: cert}
		}
	}
	return nil
}

// signer returns the cryptographic signer of the key.
func (ck *CertKey) signer() (crypto.Signer, error) {
	switch ck.keyBlock.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(ck.keyBlock.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(ck.keyBlock.Bytes)
		if err != nil {
			return nil, err
		}
		signer, ok := key.(crypto.Signer)
		if !ok {
			return nil, errors.New("internal error: parsed PKCS8 private key is not a crypto.Signer")
		}
		return signer, nil
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(ck.keyBlock.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type %q", ck.keyBlock.Type)
	}
}

func (ck *CertKey) sign(message *bytes.Buffer) ([]byte, error) {
	signer, err := ck.signer()
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(message.Bytes())
	return signer.Sign(rand.Reader, hashed[:], crypto.SHA256)
}

func decodePEM(data []byte) []*pem.Block {
	var blocks []*pem.Block
	block, rest := pem.Decode(data)
	for block != nil {
		blocks = append(blocks, block)
		block, rest = pem.Decode(rest)
	}
	return blocks
}

func decodeClientCert(data []byte) ([]*pem.Block, error) {
	blocks := decodePEM(data)
	switch {
	case blocks == nil:
		return nil, errors.New("invalid client certificate")
	case len(blocks) < 1:
		return nil, fmt.Errorf("invalid number of blocks in certificate file %d - expected min 1", len(blocks))
	}
	return blocks, nil
}

func parseCerts(blocks []*pem.Block) (certs []*x509.Certificate, err error) {
	for _, block := range blocks {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// encryptedBlock tells whether a private key is
// encrypted by examining its Proc-Type header
// for a mention of ENCRYPTED
// according to RFC 1421 Section 4.6.1.1.
func encryptedBlock(block *pem.Block) bool {
	return strings.Contains(block.Headers["Proc-Type"], "ENCRYPTED")
}

func decodeClientKey(data []byte) (*pem.Block, error) {
	blocks := decodePEM(data)
	switch {
	case blocks == nil:
		return nil, errors.New("invalid client key")
	case len(blocks) != 1:
		return nil, fmt.Errorf("invalid number of blocks in key file %d - expected 1", len(blocks))
	}
	block := blocks[0]
	if encryptedBlock(block) {
		return nil, errors.New("client key is password encrypted")
	}
	return block, nil
}
