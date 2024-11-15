// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"

	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/client"
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/util"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	backend := newBackend(util.NewSecretsClient(conf.Logger), conf.System)
	if err := backend.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return backend, nil
}

func newBackend(client secretsClient, passwordGenerator passwordGenerator) *backend {
	adBackend := &backend{
		client:         client,
		roleCache:      cache.New(roleCacheExpiration, roleCacheCleanup),
		credCache:      cache.New(credCacheExpiration, credCacheCleanup),
		rotateRootLock: new(int32),
		checkOutHandler: &checkOutHandler{
			client:            client,
			passwordGenerator: passwordGenerator,
		},
		checkOutLocks: locksutil.CreateLocks(),
	}
	adBackend.Backend = &framework.Backend{
		Help: backendHelp,
		Paths: []*framework.Path{
			adBackend.pathConfig(),
			adBackend.pathRoles(),
			adBackend.pathListRoles(),
			adBackend.pathCreds(),
			adBackend.pathRotateRootCredentials(),
			adBackend.pathRotateCredentials(),

			// The following paths are for AD credential checkout.
			adBackend.pathSetCheckIn(),
			adBackend.pathSetManageCheckIn(),
			adBackend.pathSetCheckOut(),
			adBackend.pathSetStatus(),
			adBackend.pathSets(),
			adBackend.pathListSets(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				configPath,
				credPrefix,
			},
		},
		Invalidate:  adBackend.Invalidate,
		BackendType: logical.TypeLogical,
		Secrets: []*framework.Secret{
			adBackend.secretAccessKeys(),
		},
		WALRollback:       adBackend.walRollback,
		WALRollbackMinAge: 1 * time.Minute,
	}
	return adBackend
}

type backend struct {
	*framework.Backend

	client secretsClient

	roleCache      *cache.Cache
	credCache      *cache.Cache
	credLock       sync.Mutex
	rotateRootLock *int32

	checkOutHandler *checkOutHandler
	// checkOutLocks are used for avoiding races
	// when working with sets through the check-out system.
	checkOutLocks []*locksutil.LockEntry
}

func (b *backend) Invalidate(ctx context.Context, key string) {
	b.invalidateRole(ctx, key)
	b.invalidateCred(ctx, key)
}

// Wraps the *util.SecretsClient in an interface to support testing.
type secretsClient interface {
	Get(conf *client.ADConf, serviceAccountName string) (*client.Entry, error)
	GetPasswordLastSet(conf *client.ADConf, serviceAccountName string) (time.Time, error)
	UpdatePassword(conf *client.ADConf, serviceAccountName string, newPassword string) error
	UpdateRootPassword(conf *client.ADConf, bindDN string, newPassword string) error
}

const backendHelp = `
The Active Directory (AD) secrets engine rotates AD passwords dynamically,
and is designed for a high-load environment where many instances may be accessing
a shared password simultaneously. With a simple set up and a simple creds API,
it doesn't require instances to be manually registered in advance to gain access. 
As long as access has been granted to the creds path via a method like 
AppRole, they're available.

Passwords are lazily rotated based on preset TTLs and can have a length configured to meet 
your needs.
`
