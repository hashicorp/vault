package jwt

import (
	"crypto"
	"crypto/hmac"
	"errors"
)

// SigningMethodHMAC implements the HMAC-SHA family of signing methods
// Expects key type of []byte for both signing and validation
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for HS256 and company
var (
	SigningMethodHS256  *SigningMethodHMAC
	SigningMethodHS384  *SigningMethodHMAC
	SigningMethodHS512  *SigningMethodHMAC
	ErrSignatureInvalid = errors.New("signature is invalid")
)

func init() {
	// HS256
	SigningMethodHS256 = &SigningMethodHMAC{"HS256", crypto.SHA256}
	RegisterSigningMethod(SigningMethodHS256.Alg(), func() SigningMethod {
		return SigningMethodHS256
	})

	// HS384
	SigningMethodHS384 = &SigningMethodHMAC{"HS384", crypto.SHA384}
	RegisterSigningMethod(SigningMethodHS384.Alg(), func() SigningMethod {
		return SigningMethodHS384
	})

	// HS512
	SigningMethodHS512 = &SigningMethodHMAC{"HS512", crypto.SHA512}
	RegisterSigningMethod(SigningMethodHS512.Alg(), func() SigningMethod {
		return SigningMethodHS512
	})
}

// Alg implements SigningMethod
func (m *SigningMethodHMAC) Alg() string {
	return m.Name
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
func (m *SigningMethodHMAC) Verify(signingString, signature string, key interface{}) error {
	// Verify the key is the right type
	keyBytes, ok := key.([]byte)
	if !ok {
		return NewInvalidKeyTypeError("[]byte", key)
	}

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// Can we use the specified hashing method?
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	// No validation errors.  Signature is good.
	return nil
}

// Sign implements the Sign method from SigningMethod
// Key must be []byte
func (m *SigningMethodHMAC) Sign(signingString string, key interface{}) (string, error) {
	keyBytes, ok := key.([]byte)
	if !ok {
		return "", NewInvalidKeyTypeError("[]byte", key)
	}

	if !m.Hash.Available() {
		return "", ErrHashUnavailable
	}

	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))

	return EncodeSegment(hasher.Sum(nil)), nil
}
