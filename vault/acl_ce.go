// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

type (
	entAcl struct{}
)

func (a *ACL) performEnterpriseAclChecks(_ context.Context, _ *logical.Request, _ bool) (ret *ACLResults) {
	return nil
}
