// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package roottoken

import (
	"encoding/base64"
	"fmt"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/xor"
)

// DecodeToken will decode the root token returned by the Vault API
// The algorithm was initially used in the generate root command
func DecodeToken(encoded, otp string, otpLength int) (string, error) {
	switch otpLength {
	case 0:
		// Backwards compat
		tokenBytes, err := xor.XORBase64(encoded, otp)
		if err != nil {
			return "", fmt.Errorf("error xoring token: %s", err)
		}

		uuidToken, err := uuid.FormatUUID(tokenBytes)
		if err != nil {
			return "", fmt.Errorf("error formatting base64 token value: %s", err)
		}
		return strings.TrimSpace(uuidToken), nil
	default:
		tokenBytes, err := base64.RawStdEncoding.DecodeString(encoded)
		if err != nil {
			return "", fmt.Errorf("error decoding base64'd token: %v", err)
		}

		tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
		if err != nil {
			return "", fmt.Errorf("error xoring token: %v", err)
		}
		return string(tokenBytes), nil
	}
}
