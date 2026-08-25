// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
)

func (c *Core) GetSCIMClientByEntityID(ctx context.Context, entityID string) (*identity.ScimClient, error) {
	return nil, nil
}
