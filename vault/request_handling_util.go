// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"sync"
	"time"

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

func getLeaseRegisterFunc(c *Core) (func(context.Context, *logical.Request, *logical.Response, string) (string, error), error) {
	return c.expiration.Register, nil
}

func getAuthRegisterFunc(c *Core) (RegisterAuthFunc, error) {
	return c.RegisterAuth, nil
}

func getUserFailedLoginInfo(ctx context.Context, c *Core, userInfo FailedLoginUser) (*FailedLoginInfo, error) {
	return c.LocalGetUserFailedLoginInfo(ctx, userInfo), nil
}

func updateUserFailedLoginInfo(ctx context.Context, c *Core, userInfo FailedLoginUser, failedLoginInfo *FailedLoginInfo, deleteEntry bool) error {
	return c.LocalUpdateUserFailedLoginInfo(ctx, userInfo, failedLoginInfo, deleteEntry)
}

func possiblyForwardAliasCreation(ctx context.Context, c *Core, inErr error, auth *logical.Auth, entity *identity.Entity) (*identity.Entity, bool, error) {
	return entity, false, inErr
}

var errCreateEntityUnimplemented = "create entity unimplemented in the server"

func registerLocalAlias(_ context.Context, _ *Core, _ *logical.Alias) (*identity.Entity, error) {
	return nil, logical.ErrReadOnly
}

func possiblyForwardSaveCachedAuthResponse(ctx context.Context, c *Core, respAuth *MFACachedAuthResponse) error {
	err := c.SaveMFAResponseAuth(respAuth)
	if err != nil {
		return err
	}

	return nil
}

func forwardCreateTokenRegisterAuth(ctx context.Context, c *Core, te *logical.TokenEntry, roleName string, renewable bool, periodToUse, explicitMaxTTLToUse time.Duration) (*logical.TokenEntry, error) {
	return nil, nil
}
