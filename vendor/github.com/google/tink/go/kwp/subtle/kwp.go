// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

// Package subtle implements the key wrapping primitive KWP defined in
// NIST SP 800 38f.
//
// The same encryption mode is also defined in RFC 5649. The NIST document is
// used here as a primary reference, since it contains a security analysis and
// further recommendations. In particular, Section 8 of NIST SP 800 38f
// suggests that the allowed key sizes may be restricted. The implementation in
// this package requires that the key sizes are in the range MinWrapSize and
// MaxWrapSize.
//
// The minimum of 16 bytes has been chosen, because 128 bit keys are the
// smallest key sizes used in tink. Additionally, wrapping short keys with KWP
// does not use the function W and hence prevents using security arguments
// based on the assumption that W is a strong pseudorandom. One consequence of
// using a strong pseudorandom permutation as an underlying function is that
// leaking partial information about decrypted bytes is not useful for an
// attack.
//
// The upper bound for the key size is somewhat arbitrary. Setting an upper
// bound is motivated by the analysis in section A.4 of NIST SP 800 38f:
// forgery of long messages is simpler than forgery of short messages.
package subtle

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"math"

	// Placeholder for internal crypto/cipher allowlist, please ignore.
)

const (
	// MinWrapSize is the smallest key byte length that may be wrapped.
	MinWrapSize = 16
	// MaxWrapSize is the largest key byte length that may be wrapped.
	MaxWrapSize = 8192

	roundCount = 6
	ivPrefix   = uint32(0xA65959A6)
)

// KWP is an implementation of an AES-KWP key wrapping cipher.
type KWP struct {
	block cipher.Block
}

// NewKWP returns a KWP instance.
//
// The key argument should be the AES wrapping key, either 16 or 32 bytes
// to select AES-128 or AES-256.
func NewKWP(wrappingKey []byte) (*KWP, error) {
	switch len(wrappingKey) {
	default:
		return nil, fmt.Errorf("kwp: invalid AES key size; want 16 or 32, got %d", len(wrappingKey))
	case 16, 32:
		block, err := aes.NewCipher(wrappingKey)
		if err != nil {
			return nil, fmt.Errorf("kwp: error building AES cipher: %v", err)
		}
		return &KWP{block: block}, nil
	}
}

// wrappingSize computes the byte length of the ciphertext output for the
// provided plaintext input.
func wrappingSize(inputSize int) int {
	paddingSize := 7 - (inputSize+7)%8
	return inputSize + paddingSize + 8
}

// computeW computes the pseudorandom permutation W over the IV concatenated
// with zero-padded key material.
func (kwp *KWP) computeW(iv, key []byte) ([]byte, error) {
	// Checks the parameter sizes for which W is defined.
	// Note that the caller ensures stricter limits.
	if len(key) <= 8 || len(key) > math.MaxInt32-16 || len(iv) != 8 {
		return nil, fmt.Errorf("kwp: computeW called with invalid parameters")
	}

	data := make([]byte, wrappingSize(len(key)))
	copy(data, iv)
	copy(data[8:], key)
	blockCount := len(data)/8 - 1

	buf := make([]byte, 16)
	copy(buf, data[:8])

	for i := 0; i < roundCount; i++ {
		for j := 0; j < blockCount; j++ {

			copy(buf[8:], data[8*(j+1):])
			kwp.block.Encrypt(buf, buf)

			// xor the round constant in big endian order
			// to the left half of the buffer
			roundConst := uint(i*blockCount + j + 1)
			for b := 0; b < 4; b++ {
				buf[7-b] ^= byte(roundConst & 0xFF)
				roundConst >>= 8
			}

			copy(data[8*(j+1):], buf[8:])
		}
	}
	copy(data[:8], buf)
	return data, nil
}

// invertW computes the inverse of the pseudorandom permutation W. Note that
// invertW does not perform an integrity check.
func (kwp *KWP) invertW(wrapped []byte) ([]byte, error) {
	// Checks the input size for which invertW is defined.
	// Note that the caller ensures stricter limits.
	if len(wrapped) < 24 || len(wrapped)%8 != 0 {
		return nil, fmt.Errorf("kwp: incorrect data size")
	}

	data := make([]byte, len(wrapped))
	copy(data, wrapped)

	blockCount := len(data)/8 - 1

	buf := make([]byte, 16)
	copy(buf, data[:8])

	for i := roundCount - 1; i >= 0; i-- {
		for j := blockCount - 1; j >= 0; j-- {
			copy(buf[8:], data[8*(j+1):])

			// xor the round constant in big endian order
			// to the left half of the buffer
			roundConst := uint(i*blockCount + j + 1)
			for b := 0; b < 4; b++ {
				buf[7-b] ^= byte(roundConst & 0xFF)
				roundConst >>= 8
			}

			kwp.block.Decrypt(buf, buf)
			copy(data[8*(j+1):], buf[8:])
		}
	}

	copy(data, buf[:8])
	return data, nil
}

// Wrap wraps the provided key material.
func (kwp *KWP) Wrap(data []byte) ([]byte, error) {
	if len(data) < MinWrapSize {
		return nil, fmt.Errorf("kwp: key size to wrap too small")
	}
	if len(data) > MaxWrapSize {
		return nil, fmt.Errorf("kwp: key size to wrap too large")
	}

	iv := make([]byte, 8)
	binary.BigEndian.PutUint32(iv, ivPrefix)
	binary.BigEndian.PutUint32(iv[4:], uint32(len(data)))

	return kwp.computeW(iv, data)
}

var errIntegrity = fmt.Errorf("kwp: unwrap failed integrity check")

// Unwrap unwraps a wrapped key.
func (kwp *KWP) Unwrap(data []byte) ([]byte, error) {
	if len(data) < wrappingSize(MinWrapSize) {
		return nil, fmt.Errorf("kwp: wrapped key size too small")
	}
	if len(data) > wrappingSize(MaxWrapSize) {
		return nil, fmt.Errorf("kwp: wrapped key size too large")
	}
	if len(data)%8 != 0 {
		return nil, fmt.Errorf("kwp: wrapped key size must be a multiple of 8 bytes")
	}

	unwrapped, err := kwp.invertW(data)
	if err != nil {
		return nil, err
	}

	// Check the IV and padding.
	// W has been designed to be strong pseudorandom permutation, so
	// leaking information about improperly padded keys would not be a
	// vulnerability. This means we don't have to go to extra lengths to
	// ensure that the integrity checks run in constant time.

	if binary.BigEndian.Uint32(unwrapped) != ivPrefix {
		return nil, errIntegrity
	}

	encodedSize := int(binary.BigEndian.Uint32(unwrapped[4:]))
	if encodedSize < 0 || wrappingSize(encodedSize) != len(unwrapped) {
		return nil, errIntegrity
	}

	for i := 8 + encodedSize; i < len(unwrapped); i++ {
		if unwrapped[i] != 0 {
			return nil, errIntegrity
		}
	}

	return unwrapped[8 : 8+encodedSize], nil
}
