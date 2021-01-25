package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"fmt"
	"math/big"
)

// SigningMethodECDSA implements the ECDSA family of signing methods signing methods
// Expects *ecdsa.PrivateKey for signing and *ecdsa.PublicKey for verification
type SigningMethodECDSA struct {
	Name      string
	Hash      crypto.Hash
	KeySize   int
	CurveBits int
}

// Mirrors the struct from crypto/ecdsa, we expect ecdsa.PrivateKey.Sign function to return this struct asn1 encoded
type ecdsaSignature struct {
	R, S *big.Int
}

// Specific instances for EC256 and company
var (
	SigningMethodES256 *SigningMethodECDSA
	SigningMethodES384 *SigningMethodECDSA
	SigningMethodES512 *SigningMethodECDSA
)

func init() {
	// ES256
	SigningMethodES256 = &SigningMethodECDSA{"ES256", crypto.SHA256, 32, 256}
	RegisterSigningMethod(SigningMethodES256.Alg(), func() SigningMethod {
		return SigningMethodES256
	})

	// ES384
	SigningMethodES384 = &SigningMethodECDSA{"ES384", crypto.SHA384, 48, 384}
	RegisterSigningMethod(SigningMethodES384.Alg(), func() SigningMethod {
		return SigningMethodES384
	})

	// ES512
	SigningMethodES512 = &SigningMethodECDSA{"ES512", crypto.SHA512, 66, 521}
	RegisterSigningMethod(SigningMethodES512.Alg(), func() SigningMethod {
		return SigningMethodES512
	})
}

// Alg implements SigningMethod
func (m *SigningMethodECDSA) Alg() string {
	return m.Name
}

// Verify implements the Verify method from SigningMethod
// For this verify method, key must be an ecdsa.PublicKey struct
func (m *SigningMethodECDSA) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = DecodeSegment(signature); err != nil {
		return err
	}

	// Get the key
	var ecdsaKey *ecdsa.PublicKey
	var ok bool

	switch k := key.(type) {
	case *ecdsa.PublicKey:
		ecdsaKey = k
	case crypto.Signer:
		pub := k.Public()
		if ecdsaKey, ok = pub.(*ecdsa.PublicKey); !ok {
			return &InvalidKeyError{Message: fmt.Sprintf("crypto.Signer returned an unexpected public key type: %T", pub)}
		}
	default:
		return NewInvalidKeyTypeError("*ecdsa.PublicKey or crypto.Signer", key)
	}

	if len(sig) != 2*m.KeySize {
		return &UnverfiableTokenError{Message: "signature length is invalid"}
	}

	r := big.NewInt(0).SetBytes(sig[:m.KeySize])
	s := big.NewInt(0).SetBytes(sig[m.KeySize:])

	// Create hasher
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Verify the signature
	if verifystatus := ecdsa.Verify(ecdsaKey, hasher.Sum(nil), r, s); verifystatus == true {
		return nil
	}
	return new(InvalidSignatureError)
}

// Sign implements the Sign method from SigningMethod
// For this signing method, key must be an ecdsa.PrivateKey struct
func (m *SigningMethodECDSA) Sign(signingString string, key interface{}) (string, error) {
	var signer crypto.Signer
	var pub *ecdsa.PublicKey
	var ok bool

	if signer, ok = key.(crypto.Signer); !ok {
		return "", NewInvalidKeyTypeError("*ecdsa.PrivateKey or crypto.Signer", key)
	}

	//sanity check that the signer is an ecdsa signer
	if pub, ok = signer.Public().(*ecdsa.PublicKey); !ok {
		return "", &InvalidKeyError{Message: fmt.Sprintf("signer returned unexpected public key type: %T", pub)}
	}

	// Create the hasher
	if !m.Hash.Available() {
		return "", ErrHashUnavailable
	}

	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return r, s
	asn1Sig, err := signer.Sign(rand.Reader, hasher.Sum(nil), m.Hash)
	if err != nil {
		return "", err
	}

	//the ecdsa.PrivateKey Sign function returns an asn1 encoded signature which is not what we want
	// so we unmarshal it to get r and s to encode as described in rfc7518 section-3.4
	var ecdsaSig ecdsaSignature
	rest, err := asn1.Unmarshal(asn1Sig, &ecdsaSig)
	if err != nil {
		return "", err
	}

	if len(rest) != 0 {
		return "", &UnverfiableTokenError{Message: "unexpected extra bytes in ecda signature"}
	}

	curveBits := pub.Curve.Params().BitSize

	if m.CurveBits != curveBits {
		return "", &InvalidKeyError{Message: "CurveBits in public key don't match those in signing method"}
	}

	keyBytes := curveBits / 8
	if curveBits%8 > 0 {
		keyBytes++
	}

	// We serialize the output (r and s) into big-endian byte arrays and pad
	// them with zeros on the left to make sure the sizes work out. Both arrays
	// must be keyBytes long, and the output must be 2*keyBytes long.
	rBytes := ecdsaSig.R.Bytes()
	rBytesPadded := make([]byte, keyBytes)
	copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

	sBytes := ecdsaSig.S.Bytes()
	sBytesPadded := make([]byte, keyBytes)
	copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

	out := append(rBytesPadded, sBytesPadded...)

	return EncodeSegment(out), nil
}
