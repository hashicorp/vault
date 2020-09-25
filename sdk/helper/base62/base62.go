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

// Encode encodes bytes to base62.  This does *not* scale linearly with input as base64, so use caution
// when using on large inputs.
func Encode(src []byte) string {
	if src == nil {
		return ""
	}

	var b strings.Builder
	b.Grow(int(float32(len(src))*1.4))

	var zero big.Int
	var rem big.Int
	var x big.Int
	x.SetBytes(src)

	// for x > 0 {
	//   str = (charset[x%62]) + str
	//   x = x/62
	// }
	for x.CmpAbs(&zero) > 0 {
		x.DivMod(&x, csLenBig, &rem)
		b.WriteByte(charset[int(rem.Int64())])
	}
	return reverse(b.String())
}

// Decode decodes a base62 string into bytes. This does *not* scale linearly with input as base64, so use caution
//// when using on large inputs.
func Decode(dst []byte, src string) ([]byte, error) {
	var num big.Int
	var x big.Int

	// n = c[0]*62^0 + c[1]*62^1 + c[2]*62^2 ...
	for i, c := range src {
		if i > 0 {
			num.Mul(&num, csLenBig)
		}
		idx := strings.IndexRune(charset, c)
		if idx < 0 {
			return nil, errors.New("invalid base62 character")
		}
		x.SetUint64(uint64(idx))
		num.Add(&num, &x)
	}

	b := num.Bytes()
	if cap(dst)<len(b) {
		return num.Bytes(), nil
	}
	copy(dst[cap(dst)-len(b):], b)
	return dst, nil
}

func reverse(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}