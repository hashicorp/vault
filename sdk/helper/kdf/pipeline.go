// See note about this package in kdf.go; these implementations aim to provide
// a PKCS#11 v3.0 CKM_SP800_108_*_KDF compatible interface.
package kdf

import (
	"fmt"
	"hash"
	"io"
)

// Creates a new SP800-108 Double Pipeline Mode KDF instance. This function
// takes a PRF instance (preferably HMAC or CMAC) already initialized with a
// key, an ordered list of KDF parameters to use (see params.go), an
// initialization vector (for the first PRF invocation), and a list of bit
// lengths of keys to derive (in order). When successfully constructed, the
// return is an io.Reader instance, which when called with byte arrays sized
// according to the specified originally specified keyLengths, is guaranteed
// to successfully return a key.
//
// See also NIST SP800-108 and PKCS#11 v3.0 for security concerns and
// more information about KBKDFParameters.
func NewPipeline(prf hash.Hash, params []KBKDFParameter, keyLengths []int) (io.Reader, error) {
	// Validate all parameters have good values.
	var foundChaining = false
	for index, parameter := range params {
		if err := parameter.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate parameter at index %d: %v", index, err)
		}

		if _, ok := parameter.(ChainingVariable); ok {
			foundChaining = true
		}
	}

	if !foundChaining {
		return nil, fmt.Errorf("missing required parameter of type ChainingVariable for Double Pipeline Mode KDFs")
	}

	// Validate we have a correct number of key lengths; need at least one.
	//
	// Technically we could avoid this restriction IF we guarantee that the
	// DKMLength parameter isn't specified. However, it makes sense to keep
	// this restriction unconditionally in my view.
	if len(keyLengths) == 0 {
		return nil, fmt.Errorf("missing required parameter keyLengths")
	}

	// Validate that each specified keyLength is in bits, not bytes.
	for index, keyLength := range keyLengths {
		if keyLength <= 0 {
			return nil, fmt.Errorf("key length at index %d was not positive; size must be specified in bits; got %d", index, keyLength)
		}

		if keyLength%8 != 0 {
			return nil, fmt.Errorf("key length at index %d wasn't a multiple of 8; size must be specified in bits; got %d", index, keyLength)
		}
	}

	// Construct the PRF. Note that counter needs to start at one!
	return &pipelineKDF{prf, params, keyLengths, 1, nil, 0}, nil
}
