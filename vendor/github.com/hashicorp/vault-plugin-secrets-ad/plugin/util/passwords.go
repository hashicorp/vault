package util

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/go-uuid"
)

var (
	// Per https://en.wikipedia.org/wiki/Password_strength#Guidelines_for_strong_passwords
	minimumLengthOfComplexString = 8

	PasswordComplexityPrefix = "?@09AZ"
	MinimumPasswordLength    = len(PasswordComplexityPrefix) + minimumLengthOfComplexString
)

func GeneratePassword(desiredLength int) (string, error) {
	if desiredLength < MinimumPasswordLength {
		return "", fmt.Errorf("it's not possible to generate a _secure_ password of length %d, please boost length to %d, though Vault recommends higher", desiredLength, MinimumPasswordLength)
	}

	b, err := uuid.GenerateRandomBytes(desiredLength)
	if err != nil {
		return "", err
	}

	result := ""
	// Though the result should immediately be longer than the desiredLength,
	// do this in a loop to ensure there's absolutely no risk of a panic when slicing it down later.
	for len(result) <= desiredLength {
		// Encode to base64 because it's more complex.
		result += base64.StdEncoding.EncodeToString(b)
	}

	result = PasswordComplexityPrefix + result
	return result[:desiredLength], nil
}
