// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
)

// parsePaddingSchemeArg validate that the provided padding scheme argument received on the api can be used.
func parsePaddingSchemeArg(keyType keysutil.KeyType, rawPs any) (keysutil.PaddingScheme, error) {
	ps, ok := rawPs.(string)
	if !ok {
		return "", fmt.Errorf("argument was not a string: %T", rawPs)
	}

	paddingScheme, err := keysutil.ParsePaddingScheme(ps)
	if err != nil {
		return "", err
	}

	if !keyType.PaddingSchemesSupported() {
		return "", fmt.Errorf("unsupported key type %s for padding scheme", keyType.String())
	}

	return paddingScheme, nil
}
