package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var (
	// per https://en.wikipedia.org/wiki/Password_strength#Guidelines_for_strong_passwords
	minimumLengthOfComplexString = 8
	complexityPrefix             = "?@09AZ"

	MinimumPasswordLength = len(complexityPrefix) + minimumLengthOfComplexString
)

func GeneratePassword(desiredLength int) (string, error) {

	if desiredLength <= 0 {
		return "", fmt.Errorf("it's not possible to generate a password of password_length %d", desiredLength)
	}
	if desiredLength < MinimumPasswordLength {
		return "", fmt.Errorf("it's not possible to generate a _secure_ password of length %d, please boost password_length to %d, though Vault recommends higher", desiredLength, MinimumPasswordLength)
	}

	// First, get some cryptographically secure pseudorandom bytes.
	b := make([]byte, desiredLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	result := ""
	// Though the result should immediately be longer than the desiredLength,
	// do this in a loop to ensure there's absolutely no risk of a panic when slicing it down later.
	for len(result) <= desiredLength {
		// Encode to base64 because it's more complex and performant than base62.
		result += base64.StdEncoding.EncodeToString(b)
	}

	result = complexityPrefix + result
	return result[:desiredLength], nil
}
