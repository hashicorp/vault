package crypto

import (
	"crypto"
	"crypto/hmac"
	"encoding/json"
	"errors"
)

// SigningMethodHMAC implements the HMAC-SHA family of SigningMethods.
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
	_    struct{}
}

// Specific instances of HMAC-SHA SigningMethods.
var (
	// SigningMethodHS256 implements HS256.
	SigningMethodHS256 = &SigningMethodHMAC{
		Name: "HS256",
		Hash: crypto.SHA256,
	}

	// SigningMethodHS384 implements HS384.
	SigningMethodHS384 = &SigningMethodHMAC{
		Name: "HS384",
		Hash: crypto.SHA384,
	}

	// SigningMethodHS512 implements HS512.
	SigningMethodHS512 = &SigningMethodHMAC{
		Name: "HS512",
		Hash: crypto.SHA512,
	}

	// ErrSignatureInvalid is returned when the provided signature is found
	// to be invalid.
	ErrSignatureInvalid = errors.New("signature is invalid")
)

// Alg implements the SigningMethod interface.
func (m *SigningMethodHMAC) Alg() string { return m.Name }

// Verify implements the Verify method from SigningMethod.
// For this signing method, must be a []byte.
func (m *SigningMethodHMAC) Verify(raw []byte, signature Signature, key interface{}) error {
	keyBytes, ok := key.([]byte)
	if !ok {
		return ErrInvalidKey
	}
	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write(raw)
	if hmac.Equal(signature, hasher.Sum(nil)) {
		return nil
	}
	return ErrSignatureInvalid
}

// Sign implements the Sign method from SigningMethod for this signing method.
// Key must be a []byte.
func (m *SigningMethodHMAC) Sign(data []byte, key interface{}) (Signature, error) {
	keyBytes, ok := key.([]byte)
	if !ok {
		return nil, ErrInvalidKey
	}
	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write(data)
	return Signature(hasher.Sum(nil)), nil
}

// Hasher implements the SigningMethod interface.
func (m *SigningMethodHMAC) Hasher() crypto.Hash { return m.Hash }

// MarshalJSON implements json.Marshaler.
// See SigningMethodECDSA.MarshalJSON() for information.
func (m *SigningMethodHMAC) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.Alg() + `"`), nil
}

var _ json.Marshaler = (*SigningMethodHMAC)(nil)
