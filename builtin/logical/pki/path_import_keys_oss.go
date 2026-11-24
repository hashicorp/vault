// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// unwrapAndImportKey returns an error on CE indicating secure import is only supported on Vault Enterprise.
func (b *backend) unwrapAndImportKey(ctx context.Context, storage logical.Storage, wrappedKeyB64, exportKeyHMAC string) (string, error) {
	return "", errutil.UserError{Err: "secure key import with wrapped_key is only supported on Vault Enterprise"}
}
