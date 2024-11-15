// Package base62 provides utilities for working with base62 strings.
// base62 strings will only contain characters: 0-9, a-z, A-Z
package base62

import (
	"crypto/rand"
	"io"

	uuid "github.com/hashicorp/go-uuid"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	csLen   = byte(len(charset))
)

// MustRandom generates a random string using base-62 characters. Resulting
// entropy is ~5.95 bits/character. If an error is encountered, MustRandom
// panics.
func MustRandom(length int) string {
	out, err := RandomWithReader(length, rand.Reader)
	if err != nil {
		panic(err)
	}
	return out
}

// Random generates a random string using base-62 characters.
// Resulting entropy is ~5.95 bits/character.
func Random(length int) (string, error) {
	return RandomWithReader(length, rand.Reader)
}

// RandomWithReader generates a random string using base-62 characters and a given reader.
// Resulting entropy is ~5.95 bits/character.
func RandomWithReader(length int, reader io.Reader) (string, error) {
	if length == 0 {
		return "", nil
	}
	output := make([]byte, 0, length)

	// Request a bit more than length to reduce the chance
	// of needing more than one batch of random bytes
	batchSize := length + length/4

	for {
		buf, err := uuid.GenerateRandomBytesWithReader(batchSize, reader)
		if err != nil {
			return "", err
		}

		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return string(output), nil
				}
			}
		}
	}
}
