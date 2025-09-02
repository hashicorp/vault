// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) doTidyCMPV2NonceStore(_ context.Context, _ logical.Storage) error {
	return nil
}
