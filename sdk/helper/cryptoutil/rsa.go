// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cryptoutil

import (
	"crypto/rsa"
	"io"
	"os"

	"github.com/hashicorp/go-secure-stdlib/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

var disabled bool

func init() {
	s := os.Getenv("VAULT_DISABLE_RSA_DRBG")
	var err error
	disabled, err = parseutil.ParseBool(s)
	if err != nil {
		// Assume it's a typo and disable
		disabled = true
	}
}

// Uses go-secure-stdlib's GenerateRSAKey routine conditionally.  This exists to be able to disable the feature
// via an ENV var in a pinch
func GenerateRSAKey(randomSource io.Reader, bits int) (*rsa.PrivateKey, error) {
	if disabled {
		return rsa.GenerateKey(randomSource, bits)
	}
	return cryptoutil.GenerateRSAKey(randomSource, bits)
}
