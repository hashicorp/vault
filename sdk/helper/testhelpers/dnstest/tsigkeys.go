// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package dnstest

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// TSIGAlgorithm represents the supported TSIG algorithm types
type TSIGAlgorithm int

const (
	HmacSHA1 TSIGAlgorithm = iota
	HmacSHA224
	HmacSHA256
	HmacSHA384
	HmacSHA512
)

// String returns the string representation of the algorithm
func (a TSIGAlgorithm) String() string {
	switch a {
	case HmacSHA1:
		return "hmac-sha1"
	case HmacSHA224:
		return "hmac-sha224"
	case HmacSHA256:
		return "hmac-sha256"
	case HmacSHA384:
		return "hmac-sha384"
	case HmacSHA512:
		return "hmac-sha512"
	default:
		return "unknown"
	}
}

// Bits returns the key size in bits for the algorithm
func (a TSIGAlgorithm) Bits() int {
	switch a {
	case HmacSHA1:
		return 160
	case HmacSHA224:
		return 224
	case HmacSHA256:
		return 256
	case HmacSHA384:
		return 384
	case HmacSHA512:
		return 512
	default:
		return 0
	}
}

type TSigKey struct {
	KeyName   string
	Algorithm TSIGAlgorithm
	Secret    string
}

// GenerateTSIGKey generates a base64 std encoded TSIG key for the specified algorithm
func GenerateTSIGKey(keyName string, algorithm TSIGAlgorithm) (TSigKey, error) {
	if keyName == "" {
		return TSigKey{}, fmt.Errorf("empty key name")
	}

	bits := algorithm.Bits()
	if bits == 0 {
		return TSigKey{}, fmt.Errorf("unsupported algorithm: %v", algorithm)
	}

	// Calculate byte length from bits
	byteLength := bits / 8

	// Generate random bytes
	key := make([]byte, byteLength)
	_, err := rand.Read(key)
	if err != nil {
		return TSigKey{}, fmt.Errorf("failed to generate random key: %w", err)
	}

	// Encode to base64
	encodedKey := base64.StdEncoding.EncodeToString(key)
	return TSigKey{KeyName: keyName, Algorithm: algorithm, Secret: encodedKey}, nil
}
