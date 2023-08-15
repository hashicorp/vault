// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package keysutil

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/logical"
)

type ManagedKeyParameters struct {
	ManagedKeySystemView logical.ManagedKeySystemView
	BackendUUID          string
	Context              context.Context
}

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func (p *Policy) decryptWithManagedKey(params ManagedKeyParameters, keyEntry KeyEntry, ciphertext []byte, nonce []byte, aad []byte) (plaintext []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) encryptWithManagedKey(params ManagedKeyParameters, keyEntry KeyEntry, plaintext []byte, nonce []byte, aad []byte) (ciphertext []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) signWithManagedKey(options *SigningOptions, keyEntry KeyEntry, input []byte) (sig []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) verifyWithManagedKey(options *SigningOptions, keyEntry KeyEntry, input, sig []byte) (verified bool, err error) {
	return false, errEntOnly
}

func (p *Policy) HMACWithManagedKey(ctx context.Context, ver int, managedKeySystemView logical.ManagedKeySystemView, backendUUID string, algorithm string, data []byte) (hmacBytes []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) RotateManagedKey(ctx context.Context, storage logical.Storage, managedKeyUUID string) error {
	return errEntOnly
}
