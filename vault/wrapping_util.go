// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

func forwardWrapRequest(context.Context, *Core, *logical.Request, *logical.Response, *logical.Auth) (*logical.Response, error) {
	return nil, nil
}
