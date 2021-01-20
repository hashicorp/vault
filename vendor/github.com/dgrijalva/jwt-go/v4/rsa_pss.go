// +build go1.4

package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// SigningMethodRSAPSS implements the RSAPSS family of signing methods
type SigningMethodRSAPSS struct {
	*SigningMethodRSA
	Options *rsa.PSSOptions
}

// Specific instances for RS/PS and company
var (
	SigningMethodPS256 *SigningMethodRSAPSS
	SigningMethodPS384 *SigningMethodRSAPSS
	SigningMethodPS512 *SigningMethodRSAPSS
)

func init() {
	// PS256
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
	RegisterSigningMethod(SigningMethodPS256.Alg(), func() SigningMethod {
		return SigningMethodPS256
	})

	// PS384
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
	RegisterSigningMethod(SigningMethodPS384.Alg(), func() SigningMethod {
		return SigningMethodPS384
	})

	// PS512
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
	RegisterSigningMethod(SigningMethodPS512.Alg(), func() SigningMethod {
		return SigningMethodPS512
	})
}

// Verify implements the Verify method from SigningMethod
// For this verify method, key must be an rsa.PublicKey struct
func (m *SigningMethodRSAPSS) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = DecodeSegment(signature); err != nil {
		return err
	}

	var rsaKey *rsa.PublicKey
	var ok bool

	switch k := key.(type) {
	case *rsa.PublicKey:
		rsaKey = k
	case crypto.Signer:
		pub := k.Public()
		if rsaKey, ok = pub.(*rsa.PublicKey); !ok {
			return &InvalidKeyError{Message: fmt.Sprintf("signer returned unexpected public key type: %T", pub)}
		}
	default:
		return NewInvalidKeyTypeError("*rsa.PublicKey or crypto.Signer", key)
	}

	// Create hasher
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	return rsa.VerifyPSS(rsaKey, m.Hash, hasher.Sum(nil), sig, m.Options)
}

// Sign implements the Sign method from SigningMethod
// For this signing method, key must be an rsa.PrivateKey struct
func (m *SigningMethodRSAPSS) Sign(signingString string, key interface{}) (string, error) {
	var signer crypto.Signer
	var ok bool

	if signer, ok = key.(crypto.Signer); !ok {
		return "", NewInvalidKeyTypeError("*rsa.PrivateKey or crypto.Signer", key)
	}

	//sanity check that the signer is an rsa signer
	if pub, ok := signer.Public().(*rsa.PublicKey); !ok {
		return "", &InvalidKeyError{Message: fmt.Sprintf("signer returned unexpected public key type: %T", pub)}
	}

	// Create the hasher
	if !m.Hash.Available() {
		return "", ErrHashUnavailable
	}

	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	sigBytes, err := signer.Sign(rand.Reader, hasher.Sum(nil), m.Options)
	if err != nil {
		return "", err
	}
	return EncodeSegment(sigBytes), nil

}
