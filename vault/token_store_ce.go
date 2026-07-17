// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/namespace"
)

func getOAuthJwtId(_ string) string {
	return ""
}

func normalizeOAuthJwtToId(token string) string {
	return token
}

func (ts *TokenStore) handleTidyEnterpriseTokens(_ context.Context, _ *namespace.Namespace, _ *multierror.Error) error {
	return nil
}
