package util

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
)

var (
	// Per https://en.wikipedia.org/wiki/Password_strength#Guidelines_for_strong_passwords
	minimumLengthOfComplexString = 8

	PasswordComplexityPrefix = "?@09AZ"
	PwdFieldTmpl             = "{{PASSWORD}}"
)

func GeneratePassword(formatter string, totalLength int) (string, error) {
	if err := ValidatePwdSettings(formatter, totalLength); err != nil {
		return "", err
	}
	pwd, err := generatePassword(totalLength)
	if err != nil {
		return "", err
	}
	if formatter == "" {
		pwd = PasswordComplexityPrefix + pwd
		return pwd[:totalLength], nil
	}
	return strings.Replace(formatter, PwdFieldTmpl, pwd[:lengthOfPassword(formatter, totalLength)], 1), nil
}

func ValidatePwdSettings(formatter string, totalLength int) error {
	// Check for if there's no formatter.
	if formatter == "" {
		if totalLength < len(PasswordComplexityPrefix)+minimumLengthOfComplexString {
			return fmt.Errorf("it's not possible to generate a _secure_ password of length %d, please boost length to %d, though Vault recommends higher", totalLength, minimumLengthOfComplexString+len(PasswordComplexityPrefix))
		}
		return nil
	}

	// Check for if there is a formatter.
	if lengthOfPassword(formatter, totalLength) < minimumLengthOfComplexString {
		return fmt.Errorf("since the desired length is %d, it isn't possible to generate a sufficiently complex password - please increase desired length or remove characters from the formatter", totalLength)
	}
	numPwdFields := strings.Count(formatter, PwdFieldTmpl)
	if numPwdFields == 0 {
		return fmt.Errorf("%s must contain password replacement field of %s", formatter, PwdFieldTmpl)
	}
	if numPwdFields > 1 {
		return fmt.Errorf("%s must contain ONE password replacement field of %s", formatter, PwdFieldTmpl)
	}
	return nil
}

func lengthOfPassword(formatter string, totalLength int) int {
	lengthOfText := len(formatter) - len(PwdFieldTmpl)
	return totalLength - lengthOfText
}

// generatePassword returns a password of a length AT LEAST as long as the desired length,
// it may be longer.
func generatePassword(desiredLength int) (string, error) {
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
	return result, nil
}
