// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1package vault

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/logical"
)

var ErrRotationManagerUnsupported = errors.New("rotation manager capabilities not supported in Vault community edition")

type RotationManager struct{}

func (c *Core) startRotation() error {
	return nil
}

func (c *Core) stopRotation() error {
	return nil
}

func (c *Core) RegisterRotationJob(_ context.Context, _ string, _ *logical.RotationJob) (string, error) {
	return "", ErrRotationManagerUnsupported
}
