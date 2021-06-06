package azuresecrets

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type azureSecretBackend struct {
	*framework.Backend

	getProvider func(*clientSettings) (AzureProvider, error)
	client      *client
	settings    *clientSettings
	lock        sync.RWMutex

	// Creating/deleting passwords against a single Application is a PATCH
	// operation that must be locked per Application Object ID.
	appLocks []*locksutil.LockEntry
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func backend() *azureSecretBackend {
	var b = azureSecretBackend{}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			pathsRole(&b),
			[]*framework.Path{
				pathConfig(&b),
				pathServicePrincipal(&b),
			},
		),
		Secrets: []*framework.Secret{
			secretServicePrincipal(&b),
			secretStaticServicePrincipal(&b),
		},
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,

		WALRollback: b.walRollback,

		// Role assignment can take up to a few minutes, so ensure we don't try
		// to roll back during creation.
		WALRollbackMinAge: 10 * time.Minute,
	}

	b.getProvider = newAzureProvider
	b.appLocks = locksutil.CreateLocks()

	return &b
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
	unlockFunc := b.lock.RUnlock
	defer func() { unlockFunc() }()

	if b.client.Valid() {
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

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

	p, err := b.getProvider(b.settings)
	if err != nil {
		return nil, err
	}

	passwords := passwords{
		policyGenerator: b.System(),
		policyName:      config.PasswordPolicy,
	}

	c := &client{
		provider:   p,
		settings:   b.settings,
		expiration: time.Now().Add(clientLifetime),
		passwords:  passwords,
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
