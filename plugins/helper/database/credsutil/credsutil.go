package credsutil

import (
	"crypto/rand"
	"time"

	"fmt"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
)

// CredentialsProducer can be used as an embeded interface in the Database
// definition. It implements the methods for generating user information for a
// particular database type and is used in all the builtin database types.
type CredentialsProducer interface {
	GenerateUsername(usernameConfig dbplugin.UsernameConfig) (string, error)
	GeneratePassword() (string, error)
	GenerateExpiration(ttl time.Time) (string, error)
}

const (
	reqStr    = `A1a-`
	minStrLen = 10
)

// RandomAlphaNumeric returns a random string of characters [A-Za-z0-9-]
// of the provided length. The string generated takes up to 4 characters
// of space that are predefined and prepended to ensure password
// character requirements. It also requires a min length of 10 characters.
func RandomAlphaNumeric(length int, prependA1a bool) (string, error) {
	if length < minStrLen {
		return "", fmt.Errorf("minimum length of %d is required", minStrLen)
	}

	var size int
	var retBytes []byte
	if prependA1a {
		size = len(reqStr)
		retBytes = make([]byte, length-size)
		// Enforce alphanumeric requirements
		retBytes = append([]byte(reqStr), retBytes...)
	} else {
		retBytes = make([]byte, length)
	}

	for size < length {
		// Extend the len of the random byte slice to lower odds of having to
		// re-roll.
		c := length + len(reqStr)
		bArr := make([]byte, c)
		_, err := rand.Read(bArr)
		if err != nil {
			return "", err
		}

		for _, b := range bArr {
			if size == length {
				break
			}

			/**
			 * Each byte will be in [0, 256), but we only care about:
			 *
			 * [48, 57]     0-9
			 * [65, 90]     A-Z
			 * [97, 122]    a-z
			 *
			 * Which means that the highest bit will always be zero, since the last byte with high bit
			 * zero is 01111111 = 127 which is higher than 122. Lower our odds of having to re-roll a byte by
			 * dividing by two (right bit shift of 1).
			 */

			b = b >> 1
			// Bitwise OR to set min to 48, further reduces re-roll
			b |= 0x30

			// The byte is any of        0-9                  A-Z                      a-z
			byteIsAllowable := (b >= 48 && b <= 57) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122)
			if byteIsAllowable {
				retBytes[size] = b
				size++
			}
		}
	}

	return string(retBytes), nil
}
