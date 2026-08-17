// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
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

func (ts *TokenStore) revokeCommonJWT(ctx context.Context, req *logical.Request, id string) (*logical.Response, error) {
	return logical.ErrorResponse("cannot revoke JWTs"), nil
}

func (ts *TokenStore) revokeJWT(ctx context.Context, req *logical.Request, jwtToken string) (*logical.Response, error) {
	return logical.ErrorResponse("cannot revoke JWTs"), nil
}

func (c *Core) normalizeJwtForLookup(ctx context.Context, token string) (string, error) {
	return token, nil
}
