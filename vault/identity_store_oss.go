//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
)

func (c *Core) SendGroupUpdate(context.Context, *identity.Group) (bool, error) {
	return false, nil
}

func (c *Core) CreateEntity(ctx context.Context) (*identity.Entity, error) {
	return nil, nil
}
