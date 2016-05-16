// This package is used to implement Key Derivation Functions (KDF)
// based on the recommendations of NIST SP 800-108. These are useful
// for generating unique-per-transaction keys, or situations in which
// a key hierarchy may be useful.
package kdf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// PRF is a pseudo-random function that takes a key or seed,
// as well as additional binary data and generates output that is
// indistinguishable from random. Examples are cryptographic hash
// functions or block ciphers.
type PRF func([]byte, []byte) ([]byte, error)

// CounterMode implements the counter mode KDF that uses a psuedo-random-function (PRF)
// along with a counter to generate derived keys. The KDF takes a base key
// a derivation context, and the required number of output bits.
func CounterMode(prf PRF, prfLen uint32, key []byte, context []byte, bits uint32) ([]byte, error) {
	// Ensure the PRF is byte aligned
	if prfLen%8 != 0 {
		return nil, fmt.Errorf("PRF must be byte aligned")
	}

	// Ensure the bits required are byte aligned
	if bits%8 != 0 {
		return nil, fmt.Errorf("bits required must be byte aligned")
	}

	// Determine the number of rounds required
	rounds := bits / prfLen
	if bits%prfLen != 0 {
		rounds++
	}

	// Allocate and setup the input
	input := make([]byte, 4+len(context)+4)
	copy(input[4:], context)
	binary.BigEndian.PutUint32(input[4+len(context):], bits)

	// Iteratively generate more key material
	var out []byte
	var i uint32
	for i = 0; i < rounds; i++ {
		// Update the counter in the input string
		binary.BigEndian.PutUint32(input[:4], i)

		// Compute more key material
		part, err := prf(key, input)
		if err != nil {
			return nil, err
		}
		if uint32(len(part)*8) != prfLen {
			return nil, fmt.Errorf("PRF length mis-match (%d vs %d)", len(part)*8, prfLen)
		}
		out = append(out, part...)
	}

	// Return the desired number of output bytes
	return out[:bits/8], nil
}

const (
	// HMACSHA256PRFLen is the length of output from HMACSHA256PRF
	HMACSHA256PRFLen uint32 = 256
)

// HMACSHA256PRF is a pseudo-random-function (PRF) that uses an HMAC-SHA256
func HMACSHA256PRF(key []byte, data []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil), nil
}
