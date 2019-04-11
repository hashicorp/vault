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

func checkErrControlGroupTokenNeedsCreated(err error) bool {
	return false
}

func shouldForward(c *Core, routeErr error) bool {
	return false
}

func syncCounter(c *Core) {
}

func couldForward(c *Core) bool {
	return false
}

func forward(ctx context.Context, c *Core, req *logical.Request) (*logical.Response, error) {
	panic("forward called in OSS Vault")
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
