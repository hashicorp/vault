// +build !enterprise

package vault

import (
	"context"
	"sync"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/logical"
)

func waitForReplicationState(context.Context, *Core, *logical.Request) (*sync.WaitGroup, error) {
	return nil, nil
}

func checkNeedsCG(context.Context, *Core, *logical.Request, *logical.Auth, error, []string) (error, *logical.Response, *logical.Auth, error) {
	return nil, nil, nil, nil
}

func checkErrControlGroupTokenNeedsCreated(err error) bool {
	return false
}

func shouldForward(c *Core, resp *logical.Response, err error) bool {
	return false
}

func syncCounters(c *Core) error {
	return nil
}

func syncBarrierEncryptionCounter(c *Core) error {
	return nil
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
