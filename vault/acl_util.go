// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *Core) performEntPolicyChecks(ctx context.Context, acl *ACL, te *logical.TokenEntry, req *logical.Request, inEntity *identity.Entity, opts *PolicyCheckOpts, ret *AuthResults) {
	ret.Allowed = true
}
