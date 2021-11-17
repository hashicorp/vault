package crypto

import "crypto"

// SigningMethod is an interface that provides a way to sign JWS tokens.
type SigningMethod interface {
	// Alg describes the signing algorithm, and is used to uniquely
	// describe the specific crypto.SigningMethod.
	Alg() string

	// Verify accepts the raw content, the signature, and the key used
	// to sign the raw content, and returns any errors found while validating
	// the signature and content.
	Verify(raw []byte, sig Signature, key interface{}) error

	// Sign returns a Signature for the raw bytes, as well as any errors
	// that occurred during the signing.
	Sign(raw []byte, key interface{}) (Signature, error)

	// Used to cause quick panics when a crypto.SigningMethod whose form of hashing
	// isn't linked in the binary when you register a crypto.SigningMethod.
	// To spoof this, see "crypto.SigningMethodNone".
	Hasher() crypto.Hash
}
