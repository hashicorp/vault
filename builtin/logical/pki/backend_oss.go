// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

type entBackend struct{}

func (b *backend) initializeEnt(_ *storageContext, _ *logical.InitializationRequest) error {
	return nil
}

func (b *backend) invalidateEnt(_ context.Context, _ string) {}

func (b *backend) periodicFuncEnt(_ *storageContext, _ *logical.Request) error {
	return nil
}

func (b *backend) cleanupEnt(_ *storageContext) {}

func (b *backend) SetupEnt() {}
