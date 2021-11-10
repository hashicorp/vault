// See note about this package in kdf.go; these implementations aim to provide
// a PKCS#11 v3.0 CKM_SP800_108_*_KDF compatible interface.
package kdf

import (
	"fmt"
	"hash"
	"io"
)

// Creates a new SP800-108 Counter Mode KDF instance. This function takes
// a PRF instance (preferably HMAC or CMAC) already initialized with a key,
// an ordered list of KDF parameters to use (see params.go), and a list of
// bit lengths of keys to derive (in order). When successfully constructed,
// the return is an io.Reader instance, which when called with byte arrays
// sized according to the originally specified keyLengths, is guaranteed to
// successfully return a key.
//
// This method is fairly low-level and generic; preference should generally
// be given to a higher-level function like CounterDerive(...) or using this
// function with specific parameter-yielding functions like
// SP800CounterParams(...).
//
// See also NIST SP800-108 and PKCS#11 v3.0 for security concerns and
// more information about KBKDFParameters.
func NewCounter(prf hash.Hash, params []KBKDFParameter, keyLengths []int) (io.Reader, error) {
	// Validate all parameters have good values. This validates that the
	// required CounterVariable is present (in PKCS#11 v3.0 language, a
	// CK_SP800_108_ITERATION_VARIABLE) and that we don't have a
	// ChainingVariable instance lurking somewhere.
	foundCounter := false
	for index, parameter := range params {
		if err := parameter.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate parameter at index %d: %v", index, err)
		}

		if _, ok := parameter.(ChainingVariable); ok {
			return nil, fmt.Errorf("parameter at index %d is of type ChainingVariable; not allowed in Counter Mode KDFs", index)
		}

		if _, ok := parameter.(CounterVariable); ok {
			foundCounter = true
		}
	}

	if !foundCounter {
		return nil, fmt.Errorf("missing required parameter of type CounterVariable for Counter Mode KDFs")
	}

	// Validate we have a correct number of key lengths; need at least one.
	//
	// Technically we could avoid this restriction IF we guarantee that the
	// DKMLength parameter isn't specified. However, it makes sense to keep
	// this restriction unconditionally in my view.
	if len(keyLengths) == 0 {
		return nil, fmt.Errorf("missing required parameter keyLengths")
	}

	// In Section 5.0 of SP800-108, NIST restricts all KBKDFs to at most
	// 2^32 - 1 loops. Use a temporary DKMLength variable (with SumOfSegments
	// to make the math correct) to calculate whether we'll have more than
	// 2^32 - 1 loops.
	helper := DKMLength{SumOfSegments, false, 64}
	prfBitLen := prf.Size() * 8
	length := helper.CalculateDKMLength(prfBitLen, keyLengths) // in bits
	invocations := length/uint64(prfBitLen)
	if invocations > nistMaxInvocations {
		return nil, fmt.Errorf("too much key material requested; max of %d but calculated %d invocations of the PRF needed", nistMaxInvocations, invocations)
	}

	// Validate that each specified keyLength is in bits, not bytes.
	for index, keyLength := range keyLengths {
		if keyLength%8 != 0 {
			return nil, fmt.Errorf("key length at index %d wasn't a multiple of 8; size must be specified in bits; got %d", index, keyLength)
		}
	}

	// Construct the PRF.
	return newCounter(prf, params, keyLengths)
}

// SP800-108, Section 5.1 Counter Mode, uses the parameter structure returned
// by this function. Counter and Length variables are left to be inferred by
// the reader, but generally agreement is that it is a big-endian, 4-byte
// value (from NIST CAVP test cases and OpenSSL's KBKDF implementations).
// Additionally, NIST under-specifies multiple derived keys from the same
// KBKDF invocation (though, seems to allow them), so we default to SumOfKeys
// DKM method.
func SP800CounterParams(label []byte, context []byte) []KBKDFParameter {
	return []KBKDFParameter{
		CounterVariable{false, 32},
		ByteArray(label),
		ByteArray([]byte{0}),
		ByteArray(context),
		DKMLength{SumOfKeys, false, 32},
	}
}

// One-shot Key derivation using SP800-108 Counter-Mode KDF with fixed
// default parameters (aligning with the existing CounterMode(...) function).
// Note that the output is different than that of CounterMode(...) because
// that function starts the counter at zero rather than 1.
func CounterDerive(prf hash.Hash, context []byte, keyLength int) ([]byte, error) {
	params := []KBKDFParameter{
		CounterVariable{false, 32},
		ByteArray(context),
		DKMLength{SumOfKeys, false, 32},
	}

	kdf, err := NewCounter(prf, params, []int{keyLength})
	if err != nil {
		return nil, err
	}

	result := make([]byte, keyLength)
	_, err = kdf.Read(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
