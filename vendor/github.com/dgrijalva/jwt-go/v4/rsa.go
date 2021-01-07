package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// SigningMethodRSA implements the RSA family of signing methods signing methods
// Expects *rsa.PrivateKey for signing and *rsa.PublicKey for validation
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for RS256 and company
var (
	SigningMethodRS256 *SigningMethodRSA
	SigningMethodRS384 *SigningMethodRSA
	SigningMethodRS512 *SigningMethodRSA
)

func init() {
	// RS256
	SigningMethodRS256 = &SigningMethodRSA{"RS256", crypto.SHA256}
	RegisterSigningMethod(SigningMethodRS256.Alg(), func() SigningMethod {
		return SigningMethodRS256
	})

	// RS384
	SigningMethodRS384 = &SigningMethodRSA{"RS384", crypto.SHA384}
	RegisterSigningMethod(SigningMethodRS384.Alg(), func() SigningMethod {
		return SigningMethodRS384
	})

	// RS512
	SigningMethodRS512 = &SigningMethodRSA{"RS512", crypto.SHA512}
	RegisterSigningMethod(SigningMethodRS512.Alg(), func() SigningMethod {
		return SigningMethodRS512
	})
}

// Alg implements the Alg method from SigningMethod
func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

// Verify implements the Verify method from SigningMethod
// For this signing method, must be an *rsa.PublicKey structure.
func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error {
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

	// Verify the signature
	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig)
}

// Sign implements the Sign method from SigningMethod
// For this signing method, must be an *rsa.PrivateKey structure.
func (m *SigningMethodRSA) Sign(signingString string, key interface{}) (string, error) {
	var signer crypto.Signer
	var ok bool

	if signer, ok = key.(crypto.Signer); !ok {
		return "", NewInvalidKeyTypeError("*rsa.PublicKey or crypto.Signer", key)
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
	sigBytes, err := signer.Sign(rand.Reader, hasher.Sum(nil), m.Hash)
	if err != nil {
		return "", err
	}
	return EncodeSegment(sigBytes), nil
}
