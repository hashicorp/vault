// +build go1.4

package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
)

// SigningMethodRSAPSS implements the RSAPSS family of SigningMethods.
type SigningMethodRSAPSS struct {
	*SigningMethodRSA
	Options *rsa.PSSOptions
}

// Specific instances for RS/PS SigningMethods.
var (
	// SigningMethodPS256 implements PS256.
	SigningMethodPS256 = &SigningMethodRSAPSS{
		&SigningMethodRSA{
			Name: "PS256",
			Hash: crypto.SHA256,
		},
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA256,
		},
	}

	// SigningMethodPS384 implements PS384.
	SigningMethodPS384 = &SigningMethodRSAPSS{
		&SigningMethodRSA{
			Name: "PS384",
			Hash: crypto.SHA384,
		},
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA384,
		},
	}

	// SigningMethodPS512 implements PS512.
	SigningMethodPS512 = &SigningMethodRSAPSS{
		&SigningMethodRSA{
			Name: "PS512",
			Hash: crypto.SHA512,
		},
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA512,
		},
	}
)

// Verify implements the Verify method from SigningMethod.
// For this verify method, key must be an *rsa.PublicKey.
func (m *SigningMethodRSAPSS) Verify(raw []byte, signature Signature, key interface{}) error {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKey
	}
	return rsa.VerifyPSS(rsaKey, m.Hash, m.sum(raw), signature, m.Options)
}

// Sign implements the Sign method from SigningMethod.
// For this signing method, key must be an *rsa.PrivateKey.
func (m *SigningMethodRSAPSS) Sign(raw []byte, key interface{}) (Signature, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	sigBytes, err := rsa.SignPSS(rand.Reader, rsaKey, m.Hash, m.sum(raw), m.Options)
	if err != nil {
		return nil, err
	}
	return Signature(sigBytes), nil
}

func (m *SigningMethodRSAPSS) sum(b []byte) []byte {
	h := m.Hash.New()
	h.Write(b)
	return h.Sum(nil)
}

// Hasher implements the Hasher method from SigningMethod.
func (m *SigningMethodRSAPSS) Hasher() crypto.Hash { return m.Hash }

// MarshalJSON implements json.Marshaler.
// See SigningMethodECDSA.MarshalJSON() for information.
func (m *SigningMethodRSAPSS) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.Alg() + `"`), nil
}

var _ json.Marshaler = (*SigningMethodRSAPSS)(nil)
