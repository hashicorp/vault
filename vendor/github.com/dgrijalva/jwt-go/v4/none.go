package jwt

// SigningMethodNone implements the none signing method.  This is required by the spec
// but you probably should never use it.
var SigningMethodNone *signingMethodNone

// UnsafeAllowNoneSignatureType must be returned from Keyfunc in order for the
// none signing method to be allowed. This is intended to make is possible to use
// this signing method, but not by accident
const UnsafeAllowNoneSignatureType unsafeNoneMagicConstant = "none signing method allowed"

// NoneSignatureTypeDisallowedError is the error value returned when the none signing method
// is used without UnsafeAllowNoneSignatureType
var NoneSignatureTypeDisallowedError error

type signingMethodNone struct{}
type unsafeNoneMagicConstant string

func init() {
	SigningMethodNone = &signingMethodNone{}
	NoneSignatureTypeDisallowedError = &InvalidSignatureError{Message: "'none' signature type is not allowed"}

	RegisterSigningMethod(SigningMethodNone.Alg(), func() SigningMethod {
		return SigningMethodNone
	})
}

func (m *signingMethodNone) Alg() string {
	return "none"
}

// Only allow 'none' alg type if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Verify(signingString, signature string, key interface{}) (err error) {
	// Key must be UnsafeAllowNoneSignatureType to prevent accidentally
	// accepting 'none' signing method
	if _, ok := key.(unsafeNoneMagicConstant); !ok {
		return NoneSignatureTypeDisallowedError
	}
	// If signing method is none, signature must be an empty string
	if signature != "" {
		return &InvalidSignatureError{Message: "'none' signing method with non-empty signature"}
	}

	// Accept 'none' signing method.
	return nil
}

// Only allow 'none' signing if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Sign(signingString string, key interface{}) (string, error) {
	if _, ok := key.(unsafeNoneMagicConstant); ok {
		return "", nil
	}
	return "", NoneSignatureTypeDisallowedError
}
