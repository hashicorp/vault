// See note about this package in kdf.go; these implementations aim to provide
// a PKCS#11 v3.0 CKM_SP800_108_*_KDF compatible interface.
package kdf

import (
	"encoding/binary"
	"fmt"
	"hash"
)

const nistMaxInvocations uint64 = (1 << 32) - 1

// Generic interface all KBKDF Parameter types (below) must implement.
type KBKDFParameter interface {
	Validate() error
}

// In the below, we deviate from PKCS#11 v3.0 by necessity of a cleaner Go
// interface. However, this deviation is reversible and the two APIs are
// equivalent. We specify CounterVariable to mean an incrementing counter,
// regardless of KBKDF type, and ChainingVariable to strictly mean the
// previous PRF invocation type (as used in Feedback and Double Pipeline
// modes).
//
// In PKCS#11 v3.0, the counter is split between two parameter types,
// CK_SP800_108_ITERATION_VARIABLE, a mandatory field, and
// CK_SP800_108_COUNTER, an optional field ironically forbidden from Counter
// mode KBKDFs. In a Counter Mode PRF, ITERATION_VARIABLE takes a
// CounterVariable type, whereas in Feedback or Double Pipeline Modes,
// ITERATION_VARIABLE takes a NULL value (to signal the previous PRF output,
// like ChainingVariable does). CK_SP800_108_COUNTER, when allowed, always
// takes a CounterVariable.
//
// Not strictly following the PKCS#11 v3.0 API allows us to define concrete
// parameter structures holding actual data, rather than using the
// type+ptr+length indications that PKCS#11 v3.0 prefers.

// Defines an incrementing integer variable; the required incrementing counter
// variable in Counter mode and an optional variable in Feedback and Double
// Pipeline modes.
type CounterVariable struct {
	// Whether or not to encode this value in little endian or big endian.
	LittleEndian bool

	// Bit width of the specified variable; must be large enough to hold
	// the value to encode; must be divisible by 8 with a maximum value of
	// 32, per section 5.0 of NIST SP800-108.
	Width uint8
}

// Defines a placeholder value for the previous PRF output (as used in
// Feedback and Double Pipeline modes). Not permitted in Counter mode.
type ChainingVariable struct {
	// No fields.
}

// The DKMLength field type encodes the length of the Derived Key Material.
// PKCS#11 v3.0 stipulates that no two keys can be produced from the same
// PRF output segment. This presents a problem for calculating the derived
// key material's length: is it the sum of the length of all requested
// keys or is it the sum of the length of all segments used to derive keys?
// This constant (provided in the DKMLength parameter) decides this.
//
// NIST SP800-108's three KBKDFs all use SumOfKeys method with a single key.
type DKMLengthMethod int

const (
	// Sum the exact size of all derived keys.
	SumOfKeys DKMLengthMethod = iota + 1

	// Sum the length of the PRF segments used.
	SumOfSegments
)

// Defines the derived key material length field.
type DKMLength struct {
	// Which method for calculating DKM length should be used.
	Method DKMLengthMethod

	// Whether to encode the result in little endian or big endian.
	LittleEndian bool

	// Bit width of the specified variable; must be large enough to hold
	// the value to encode; must be divisible by 8 with a maximum value of
	// 64. Note that this is longer than CounterVariable.Width as it is not
	// strictly specified by NIST. With 2^32 - 1 PRF invocations (max allowed
	// by NIST), each generating 1024 bits (and generating 1024-bit keys),
	// this could theoretically require a 42-bit counter, hence 64 seems a
	// sane upper bound here.
	Width uint8
}

// Defines a fixed byte array value (such as the info or context fields from
// NIST SP800-108). No additional manipulation are performed on these fields.
type ByteArray []byte

// Encodes an integer value (as a uint64) into the specified byte width,
// using the desired endianness.
//
// Note that all parameters above, like the PKCS#11 v3.0 standard, use _bits_
// to specify the width; here we take _bytes_ to make the function easier to
// read.
func encodeInteger(value uint64, bytes uint8, little bool) []byte {
	// Always allocate 8 bytes == 64 bits, our maximum possible size; we'll
	// truncate it later as desired.
	var result []byte = make([]byte, 8)

	if little {
		binary.LittleEndian.PutUint64(result, value)

		// In little endian encoding, the value is left aligned; this means,
		// to truncate to the specified width, we should remove the values
		// thereafter.
		result = result[:bytes]
	} else {
		binary.BigEndian.PutUint64(result, value)

		// In big endian encoding, the value is right aligned; this means,
		// to truncate to the specified width, we need to remove the leading
		// values. To do so, we skip len(result) - $width bytes at the front,
		// to leave $width bytes at the end.
		result = result[8-bytes : 8]
	}

	return result
}

// Encodes, using the CounterVariable template, a given integer value.
func (cv CounterVariable) Encode(value uint64) []byte {
	return encodeInteger(value, cv.Width/8, cv.LittleEndian)
}

// Validates this parameter.
func (cv CounterVariable) Validate() error {
	if cv.Width == 0 {
		return fmt.Errorf("potentially uninitialized CounterVariable: expected Width to be greater than zero; was 0")
	}

	if cv.Width%8 != 0 {
		return fmt.Errorf("invalid value for CounterVariable: expected Width to be a multiple of 8; got %d", cv.Width)
	}

	if cv.Width > 32 {
		return fmt.Errorf("invalid value for CounterVariable: expected Width to not exceed 32; got %d", cv.Width)
	}

	return nil
}

// Validates this parameter
func (cv ChainingVariable) Validate() error {
	return nil
}

func (kl DKMLength) CalculateDKMLength(prfSize int, keys []int) uint64 {
	// Assumption: Validate has already passed, meaning there is exactly one
	// of two values for kl.Method.
	var result uint64

	for _, keyBits := range keys {
		if kl.Method == SumOfKeys {
			result += uint64(keyBits)
		} else {
			// Because PKCS#11 v3.0 guarantees no two keys will be in the same
			// PRF output segment, we need to compute:
			//     ceil(keyBits/prfSize)
			// to get the number of segments, and then multiply by the segment
			// length (prfSize).
			var keySegments = (keyBits + prfSize - 1) / prfSize
			result += uint64(keySegments) * uint64(prfSize)
		}
	}

	return result
}

// Encode the DKM Length to a byte array.
func (kl DKMLength) Encode(prf hash.Hash, keys []int) []byte {
	// Need the PRF size in bits, not bytes; multiply by 8.
	length := kl.CalculateDKMLength(prf.Size()*8, keys)
	return encodeInteger(length, kl.Width/8, kl.LittleEndian)
}

func (kl DKMLength) Validate() error {
	if kl.Method != SumOfKeys && kl.Method != SumOfSegments {
		return fmt.Errorf("invalid value for DKMLength: expected Method to be either SumOfKeys or SumOfSegments; got %d", kl.Method)
	}

	if kl.Width == 0 {
		return fmt.Errorf("potentially uninitialized DKMLength: expected Width to be greater than zero; was 0")
	}

	if kl.Width%8 != 0 {
		return fmt.Errorf("invalid value for DKMLength: expected Width to be a multiple of 8; got %d", kl.Width)
	}

	if kl.Width > 64 {
		return fmt.Errorf("invalid value for DKMLength: expected Width to not exceed 64; got %d", kl.Width)
	}

	return nil
}

func (ba ByteArray) Validate() error {
	return nil
}
