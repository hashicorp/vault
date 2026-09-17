// Copyright IBM Corp. 2016, 2025
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

// parseHashAlgorithmArg validates that the provided hash algorithm argument received on the api can be used.
func parseHashAlgorithmArg(keyType keysutil.KeyType, rawHt any) (keysutil.HashType, error) {
	h, ok := rawHt.(string)
	if !ok {
		return keysutil.HashTypeNone, fmt.Errorf("argument was not a string: %T", rawHt)
	}

	hashType, ok := keysutil.HashTypeMap[h]
	if !ok {
		return keysutil.HashTypeNone, fmt.Errorf("unknown hash algorithm")
	}

	switch keyType {
	case keysutil.KeyType_RSA2048, keysutil.KeyType_RSA3072, keysutil.KeyType_RSA4096, keysutil.KeyType_MANAGED_KEY:
		return hashType, nil
	default:
		return keysutil.HashTypeNone, fmt.Errorf("unsupported key type %s for hash algorithm", keyType.String())
	}
}
