// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/framework"
)

func (c *Core) SendGroupUpdate(context.Context, *identity.Group) error {
	return nil
}

func (c *Core) CreateEntity(ctx context.Context) (*identity.Entity, error) {
	return nil, nil
}

func identityStoreLoginMFAEntUnauthedPaths() []string {
	return []string{}
}

func identityStoreSCIMUnauthedPaths() []string {
	return []string{}
}

func mfaLoginEnterprisePaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{}
}
