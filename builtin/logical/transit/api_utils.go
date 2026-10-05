// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"fmt"
	"strconv"
	"strings"

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

// getVersion returns the key version for an input signature or ciphertext
func getVersion(input string) (int, error) {
	if !strings.HasPrefix(input, "vault:v") {
		return 0, fmt.Errorf("invalid ciphertext: no prefix")
	}

	splitVerification := strings.SplitN(strings.TrimPrefix(input, "vault:v"), ":", 2)
	if len(splitVerification) != 2 {
		return 0, fmt.Errorf("wrong number of fields delimited by ':', got %d expected 2", len(splitVerification))
	}

	ver, err := strconv.Atoi(splitVerification[0])
	if err != nil {
		return 0, fmt.Errorf("key version number %s count not be decoded", splitVerification[0])
	}

	if ver < 0 {
		return 0, fmt.Errorf("key version cannot be negative: %d", ver)
	}

	return ver, nil
}
