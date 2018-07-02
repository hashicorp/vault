package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
)

// SigningMethodRSA implements the RSA family of SigningMethods.
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
	_    struct{}
}

// Specific instances of RSA SigningMethods.
var (
	// SigningMethodRS256 implements RS256.
	SigningMethodRS256 = &SigningMethodRSA{
		Name: "RS256",
		Hash: crypto.SHA256,
	}

	// SigningMethodRS384 implements RS384.
	SigningMethodRS384 = &SigningMethodRSA{
		Name: "RS384",
		Hash: crypto.SHA384,
	}

	// SigningMethodRS512 implements RS512.
	SigningMethodRS512 = &SigningMethodRSA{
		Name: "RS512",
		Hash: crypto.SHA512,
	}
)

// Alg implements the SigningMethod interface.
func (m *SigningMethodRSA) Alg() string { return m.Name }

// Verify implements the Verify method from SigningMethod.
// For this signing method, must be an *rsa.PublicKey.
func (m *SigningMethodRSA) Verify(raw []byte, sig Signature, key interface{}) error {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKey
	}
	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, m.sum(raw), sig)
}

// Sign implements the Sign method from SigningMethod.
// For this signing method, must be an *rsa.PrivateKey structure.
func (m *SigningMethodRSA) Sign(data []byte, key interface{}) (Signature, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, m.sum(data))
	if err != nil {
		return nil, err
	}
	return Signature(sigBytes), nil
}

func (m *SigningMethodRSA) sum(b []byte) []byte {
	h := m.Hash.New()
	h.Write(b)
	return h.Sum(nil)
}

// Hasher implements the SigningMethod interface.
func (m *SigningMethodRSA) Hasher() crypto.Hash { return m.Hash }

// MarshalJSON implements json.Marshaler.
// See SigningMethodECDSA.MarshalJSON() for information.
func (m *SigningMethodRSA) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.Alg() + `"`), nil
}

var _ json.Marshaler = (*SigningMethodRSA)(nil)
