// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *SystemBackend) verifyDROperationToken(f framework.OperationFunc, lock bool) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return f(ctx, req, d)
	}
}
