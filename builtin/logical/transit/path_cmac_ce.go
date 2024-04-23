// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package transit

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathCMACVerify(_ context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	return logical.ErrorResponse(ErrCmacEntOnly.Error()), nil
}
