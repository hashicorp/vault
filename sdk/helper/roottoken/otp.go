// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package roottoken

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/base62"
)

// DefaultBase64EncodedOTPLength is the number of characters that will be randomly generated
// before the Base64 encoding process takes place.
const defaultBase64EncodedOTPLength = 16

// GenerateOTP generates a random token and encodes it as a Base64 or as a Base62 encoded string.
// Returns 0 if the generation completed without any error, 2 otherwise, along with the error.
func GenerateOTP(otpLength int) (string, error) {
	switch otpLength {
	case 0:
		// This is the fallback case
		buf := make([]byte, defaultBase64EncodedOTPLength)
		readLen, err := rand.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error reading random bytes: %s", err)
		}

		if readLen != defaultBase64EncodedOTPLength {
			return "", fmt.Errorf("read %d bytes when we should have read 16", readLen)
		}

		return base64.StdEncoding.EncodeToString(buf), nil
	default:
		otp, err := base62.Random(otpLength)
		if err != nil {
			return "", fmt.Errorf("error reading random bytes: %w", err)
		}

		return otp, nil
	}
}
