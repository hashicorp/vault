package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// ECDSA parsing errors.
var (
	ErrNotECPublicKey  = errors.New("Key is not a valid ECDSA public key")
	ErrNotECPrivateKey = errors.New("Key is not a valid ECDSA private key")
)

// ParseECPrivateKeyFromPEM will parse a PEM encoded EC Private
// Key Structure.
func ParseECPrivateKeyFromPEM(key []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

// ParseECPublicKeyFromPEM will parse a PEM encoded PKCS1 or PKCS8 public key
func ParseECPublicKeyFromPEM(key []byte) (*ecdsa.PublicKey, error) {

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		parsedKey = cert.PublicKey
	}

	pkey, ok := parsedKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrNotECPublicKey
	}
	return pkey, nil
}
