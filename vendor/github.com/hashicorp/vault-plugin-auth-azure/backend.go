// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azureauth

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	userAgentPluginName = "auth-azure"

	// operationPrefixAzure is used as a prefix for OpenAPI operation id's.
	operationPrefixAzure = "azure"
)

// Factory is used by framework
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

type azureAuthBackend struct {
	*framework.Backend

	l sync.RWMutex

	provider provider

	updatePassword bool
	// resourceAPIVersionCache is a mapping of ResourceType to APIVersion
	// so that we don't query supported API versions on each call to login for
	// a given resource type
	resourceAPIVersionCache map[string]string
	cacheLock               sync.RWMutex
}

func backend() *azureAuthBackend {
	b := azureAuthBackend{
		updatePassword: true,
	}

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Invalidate:  b.invalidate,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			Unauthenticated: []string{
				"login",
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathLogin(&b),
				pathConfig(&b),
				pathRotateRoot(&b),
			},
			pathsRole(&b),
		),
		// Root rotation can take up to a few minutes, so ensure we do not
		// roll back a root credential rotation that is currently in flight
		WALRollbackMinAge: 3 * time.Minute,
		WALRollback:       b.walRollback,
		// periodicFunc to clean up old credentials
		PeriodicFunc: b.periodicFunc,
	}

	b.resourceAPIVersionCache = make(map[string]string)

	return &b
}

// The periodicFunc is responsible for eventually swapping out the root credential for rotation
// operations. Due to Azure's eventual consistency model, the new credential will not be
// available immediately, and hence we check periodically and delete the old credential
// only once the new credential is at least a minute old
func (b *azureAuthBackend) periodicFunc(ctx context.Context, req *logical.Request) error {
	// Root rotation through the periodic func writes to storage. Only run this on the
	// active instance in the primary cluster or local mounts. The periodic func doesn't
	// run on perf standbys or DR secondaries, but we still protect against this here.
	replicationState := b.System().ReplicationState()
	if (b.System().LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby) {

		b.Logger().Debug("starting periodic func")
		if !b.updatePassword {
			b.Logger().Debug("periodic func", "rotate-root", "no rotate-root update")
			return nil
		}

		config, err := b.config(ctx, req.Storage)
		if err != nil {
			return err
		}

		// Config can be nil if deleted or when the engine is enabled
		// but not yet configured.
		if config == nil {
			return nil
		}

		// Password should be at least a minute old before we process it
		if config.NewClientSecret == "" || (time.Since(config.NewClientSecretCreated) < time.Minute) {
			return nil
		}

		b.Logger().Debug("periodic func", "rotate-root", "new password detected, swapping in storage")
		provider, err := b.getProvider(ctx, config)
		if err != nil {
			return err
		}

		client, err := provider.MSGraphClient()
		if err != nil {
			return err
		}

		app, err := client.GetApplication(ctx, config.ClientID)
		if err != nil {
			return err
		}

		var credsToDelete []*uuid.UUID
		for _, cred := range app.GetPasswordCredentials() {
			if cred.GetKeyId().String() != config.NewClientSecretKeyID {
				credsToDelete = append(credsToDelete, cred.GetKeyId())
			}
		}

		if len(credsToDelete) != 0 {
			b.Logger().Debug("periodic func", "rotate-root", "removing old passwords from Azure")
			err = removeApplicationPasswords(ctx, client, *app.GetId(), credsToDelete...)
			if err != nil {
				return err
			}
		}

		b.Logger().Debug("periodic func", "rotate-root", "updating config with new password")
		config.ClientSecret = config.NewClientSecret
		config.ClientSecretKeyID = config.NewClientSecretKeyID
		config.RootPasswordExpirationDate = config.NewClientSecretExpirationDate
		config.NewClientSecret = ""
		config.NewClientSecretKeyID = ""
		config.NewClientSecretCreated = time.Time{}

		err = b.saveConfig(ctx, config, req.Storage)
		if err != nil {
			return err
		}

		b.updatePassword = false
	}

	return nil
}

func (b *azureAuthBackend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.reset()
	}
}

func (b *azureAuthBackend) getProvider(ctx context.Context, config *azureConfig) (provider, error) {
	b.l.RLock()
	unlockFunc := b.l.RUnlock
	defer func() { unlockFunc() }()

	if b.provider != nil {
		return b.provider, nil
	}

	// Upgrade lock
	b.l.RUnlock()
	b.l.Lock()
	unlockFunc = b.l.Unlock

	if b.provider != nil {
		return b.provider, nil
	}

	provider, err := b.newAzureProvider(ctx, config)
	if err != nil {
		return nil, err
	}

	b.provider = provider
	return b.provider, nil
}

func (b *azureAuthBackend) reset() {
	b.l.Lock()
	defer b.l.Unlock()

	b.provider = nil
}

const backendHelp = `
The Azure backend plugin allows authentication for Azure .
`
