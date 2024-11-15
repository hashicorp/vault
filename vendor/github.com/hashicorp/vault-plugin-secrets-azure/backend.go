// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresecrets

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	userAgentPluginName = "secrets-azure"

	// operationPrefixAzure is used as a prefix for OpenAPI operation id's.
	operationPrefixAzure = "azure"
)

type azureSecretBackend struct {
	*framework.Backend

	getProvider func(context.Context, hclog.Logger, logical.SystemView, *clientSettings) (AzureProvider, error)
	client      *client
	settings    *clientSettings
	lock        sync.RWMutex

	// Creating/deleting passwords against a single Application is a PATCH
	// operation that must be locked per Application Object ID.
	appLocks       []*locksutil.LockEntry
	updatePassword bool
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func backend() *azureSecretBackend {
	b := azureSecretBackend{
		updatePassword: true,
	}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			pathsRole(&b),
			[]*framework.Path{
				pathConfig(&b),
				pathServicePrincipal(&b),
				pathRotateRoot(&b),
			},
		),
		Secrets: []*framework.Secret{
			secretServicePrincipal(&b),
			secretStaticServicePrincipal(&b),
		},
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,

		// Role assignment can take up to a few minutes, so ensure we don't try
		// to roll back during creation.
		WALRollbackMinAge: 10 * time.Minute,

		WALRollback:  b.walRollback,
		PeriodicFunc: b.periodicFunc,
	}
	b.getProvider = newAzureProvider
	b.appLocks = locksutil.CreateLocks()

	return &b
}

func (b *azureSecretBackend) periodicFunc(ctx context.Context, sys *logical.Request) error {
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

		config, err := b.getConfig(ctx, sys.Storage)
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
		client, err := b.getClient(ctx, sys.Storage)
		if err != nil {
			return err
		}

		apps, err := client.provider.ListApplications(ctx, fmt.Sprintf("appId eq '%s'", config.ClientID))
		if err != nil {
			return err
		}

		if len(apps) == 0 {
			return fmt.Errorf("no application found")
		}
		if len(apps) > 1 {
			return fmt.Errorf("multiple applications found - double check your client_id")
		}

		app := apps[0]

		credsToDelete := []string{}
		for _, cred := range app.PasswordCredentials {
			if cred.KeyID != config.NewClientSecretKeyID {
				credsToDelete = append(credsToDelete, cred.KeyID)
			}
		}

		if len(credsToDelete) != 0 {
			b.Logger().Debug("periodic func", "rotate-root", "removing old passwords from Azure")
			err = removeApplicationPasswords(ctx, client.provider, app.AppObjectID, credsToDelete...)
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

		err = b.saveConfig(ctx, config, sys.Storage)
		if err != nil {
			return err
		}

		b.updatePassword = false
	}

	return nil
}

// reset clears the backend's cached client
// This is used when the configuration changes and a new client should be
// created with the updated settings.
func (b *azureSecretBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.settings = nil
	b.client = nil
}

func (b *azureSecretBackend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.reset()
	}
}

func (b *azureSecretBackend) getClient(ctx context.Context, s logical.Storage) (*client, error) {
	b.lock.RLock()

	if b.client.Valid() {
		b.lock.RUnlock()
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.client.Valid() {
		return b.client, nil
	}

	config, err := b.getConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	if b.settings == nil {
		if config == nil {
			config = new(azureConfig)
		}

		settings, err := b.getClientSettings(ctx, config)
		if err != nil {
			return nil, err
		}
		b.settings = settings
	}

	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	p, err := b.getProvider(ctx, b.Logger(), b.System(), b.settings)
	if err != nil {
		return nil, err
	}

	c := &client{
		provider:   p,
		settings:   b.settings,
		expiration: time.Now().Add(clientLifetime),
	}
	b.client = c

	return c, nil
}

const backendHelp = `
The Azure secrets backend dynamically generates Azure service
principals. The SP credentials have a configurable lease and
are automatically revoked at the end of the lease.

After mounting this backend, credentials to manage Azure resources
must be configured with the "config/" endpoints and policies must be
written using the "roles/" endpoints before any credentials can be
generated.
`
