// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *Core) validateEnterpriseTokenAndFetchEntity(ctx context.Context, tokenString string) (bool, map[string]interface{}, *identity.Entity, *identity.Entity, error) {
	return false, nil, nil, nil, errors.New("not implemented")
}

func (c *Core) createAndStoreEnterpriseTokenEntry(ctx context.Context, req *logical.Request, allClaims map[string]interface{}, entity *identity.Entity) error {
	return nil
}

func getEnterpriseTokenMetadata(_ map[string]interface{}) string {
	return ""
}

func (c *Core) performSecondaryEntityTokenChecks(_ context.Context, _ *ACL, _ *identity.Entity, _ map[string][]string) (*ACL, error) {
	return nil, errors.New("not implemented")
}
