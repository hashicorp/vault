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
var lookup [256]*big.Int

func init() {
	for i, c := range charset {
		lookup[c]= big.NewInt(int64(i))
	}
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

// Encode encodes bytes to base62.  This does *not* scale linearly with input as base64, so use caution
// when using on large inputs.
func Encode(src []byte) string {
	if src == nil {
		return ""
	}

	var b strings.Builder
	b.Grow(int(float32(len(src))*1.4))

	//var zero big.Int
	var rem big.Int
	var x big.Int
	x.SetBytes(src)

	// for x > 0 {
	//   str = (charset[x%62]) + str
	//   x = x/62
	// }
	for x.BitLen() > 0 {
		x.DivMod(&x, csLenBig, &rem)
		b.WriteByte(charset[int(rem.Int64())])
	}

	for i:=0; i<len(src)-1 && src[i]==0; i++ {
		b.WriteByte(0)
	}
	return reverse(b.String())
}

var errInvalidBase62Char = errors.New("invalid base62 character")

// Decode decodes a base62 string into bytes. This does *not* scale linearly with input as base64, so use caution
//// when using on large inputs.
func DecodeString(src string) ([]byte, error) {
	if src=="" {
		return nil, nil
	}

	// n = c[0]
	a := lookup[src[0]]
	if a==nil {
		return nil, errInvalidBase62Char
	}

	var num big.Int
	num.Set(a)

	// n = n * 62 + c[1] ...
	for _, c := range src[1:] {
		num.Mul(&num, csLenBig)
		a = lookup[c]
		if a == nil {
			return nil, errInvalidBase62Char
		}
		num.Add(&num, a)
	}

	return num.Bytes(), nil
}

func reverse(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}