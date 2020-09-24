// Package base62 provides utilities for working with base62 strings.
// base62 strings will only contain characters: 0-9, a-z, A-Z
package base62

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const csLen = byte(len(charset))
var csLenBig = big.NewInt(int64(len(charset)))

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

// Encode encodes bytes to base62.
func Encode(src []byte) string {
	if src == nil {
		return ""
	}

	var b string

	var zero big.Int
	var rem big.Int
	var x big.Int
	x.SetBytes(src)

	for x.Cmp(&zero) > 0 {
		x.DivMod(&x, csLenBig, &rem)
		b = string(charset[int(rem.Int64())]) + b
	}
	return b
}

// Decode decodes a base62 string into bytes
func DecodeString(src string) ([]byte, error) {
	var num big.Int
	var x big.Int
	var y big.Int
	var e big.Int

	strlen := len(src)
	for i, c := range src {
		idx := strings.IndexRune(charset, c)
		if idx < 0 {
			return nil, errors.New("invalid base62 character")
		}
		y.SetInt64(int64(strlen - (i + 1)))
		e.Exp(csLenBig, &y, nil)
		x.SetInt64(int64(idx))
		x.Mul(&x, &e)
		num.Add(&num, &x)
	}
	return num.Bytes(), nil
}