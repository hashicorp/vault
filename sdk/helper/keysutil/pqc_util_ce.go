// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package keysutil

import (
	"errors"
	"io"

	"github.com/cloudflare/circl/sign"
)

func (p *Policy) generateMLDSAKey(_ io.Reader) (sign.PrivateKey, error) {
	return nil, errors.New("PQC key types are only available in enterprise versions of Vault")
}

func (p *Policy) signWithMLDSA(_ []byte, _ int) ([]byte, error) {
	return nil, errors.New("PQC key types are only available in enterprise versions of Vault")
}

func (p *Policy) verifyWithMLDSA(_, _ []byte, _ int) (bool, error) {
	return false, errors.New("PQC key types are only available in enterprise versions of Vault")
}
