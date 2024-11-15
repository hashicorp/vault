// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (*oidc.IDToken, error)
}
