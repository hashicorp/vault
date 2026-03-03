// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func getEnterpriseTokenId(_ string) string {
	return ""
}

func (ts *TokenStore) handleTidyEnterpriseTokens(ctx context.Context, req *logical.Request, data *framework.FieldData) {
	return nil
}
