// +build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
)

func (c *Core) PersistTOTPKey(context.Context, string, string, string) error {
	return nil
}

func (c *Core) SendGroupUpdate(context.Context, *identity.Group) (bool, error) {
	return false, nil
}
