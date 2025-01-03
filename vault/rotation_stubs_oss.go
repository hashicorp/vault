// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1package vault

//go:build !enterprise

package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type RotationManager struct{}

func (c *Core) startRotation() error {
	return nil
}

func (c *Core) stopRotation() error {
	return nil
}

func (c *Core) RegisterRotationJob(_ context.Context, _ string, _ *logical.RotationJob) (string, error) {
	return "", automatedrotationutil.ErrRotationManagerUnsupported
}
