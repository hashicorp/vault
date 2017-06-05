package credsutil

import (
	"crypto/rand"
	"time"

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

// RandomAlphaNumericOfLen returns a random byte slice of characters [A-Za-z0-9]
// of the provided length.
func RandomAlphaNumericOfLen(len int) ([]byte, error) {
	retBytes := make([]byte, len)
	size := 0

	for size < len {
		// Extend the len of the random byte slice to lower odds of having to
		// re-roll.
		c := len + 3
		bArr := make([]byte, c)
		_, err := rand.Read(bArr)
		if err != nil {
			return nil, err
		}

		for _, b := range bArr {
			if size == len {
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

			// The byte is any of        0-9                  A-Z                      a-z
			byteIsAllowable := (b >= 48 && b <= 57) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122)
			if byteIsAllowable {
				retBytes[size] = b
				size++
			}
		}
	}

	return retBytes, nil
}
