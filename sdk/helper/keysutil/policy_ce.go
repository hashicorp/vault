// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package keysutil

import (
	"fmt"
	"io"

	"github.com/hashicorp/vault/sdk/helper/errutil"
)

type entKeyEntry struct{}

func (e entKeyEntry) IsEntPrivateKeyMissing() bool {
	return true
}

func entSignWithOptions(p *Policy, input, context []byte, ver int, hashAlgorithm HashType, options *SigningOptions) ([]byte, error) {
	return nil, fmt.Errorf("unsupported key type %v", p.Type)
}

func entVerifySignatureWithOptions(p *Policy, input, context []byte, sigBytes []byte, ver int, options *SigningOptions) (bool, error) {
	return false, errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
}

func entRotateInMemory(p *Policy, entry *KeyEntry, rand io.Reader) error {
	return fmt.Errorf("unsupported key type %v", p.Type)
}
