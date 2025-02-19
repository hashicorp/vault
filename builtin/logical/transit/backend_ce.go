// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package transit

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type entBackend struct{}

func (b *backend) initializeEnt(_ context.Context, _ *logical.InitializationRequest) error {
	return nil
}

func (b *backend) invalidateEnt(_ context.Context, _ string) {}

func (b *backend) periodicFuncEnt(_ context.Context, _ *logical.Request) error { return nil }

func (b *backend) cleanupEnt(_ context.Context) {}

func (b *backend) setupEnt() {}

func entEncodePrivateKey(_ string, p *keysutil.Policy, _ *keysutil.KeyEntry) (string, error) {
	return "", nil
}

func entEncodePublicKey(_ string, p *keysutil.Policy, _ *keysutil.KeyEntry) (string, error) {
	return "", nil
}
