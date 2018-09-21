// +build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
)

func waitForReplicationState(context.Context, *Core, *logical.Request) error { return nil }

func checkNeedsCG(context.Context, *Core, *logical.Request, *logical.Auth, error, []string) (error, *logical.Response, *logical.Auth, error) {
	return nil, nil, nil, nil
}

func possiblyForward(ctx context.Context, c *Core, req *logical.Request, resp *logical.Response, routeErr error) (*logical.Response, error) {
	return resp, routeErr
}

func getLeaseRegisterFunc(c *Core) (func(context.Context, *logical.Request, *logical.Response) (string, error), error) {
	return c.expiration.Register, nil
}

func getAuthRegisterFunc(c *Core) (RegisterAuthFunc, error) {
	return c.RegisterAuth, nil
}

func possiblyForwardAliasCreation(ctx context.Context, c *Core, inErr error, auth *logical.Auth, entity *identity.Entity) (*identity.Entity, error) {
	return entity, inErr
}
