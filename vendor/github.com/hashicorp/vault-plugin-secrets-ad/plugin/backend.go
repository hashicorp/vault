package plugin

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/client"
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/patrickmn/go-cache"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	backend := newBackend(util.NewSecretsClient(conf.Logger))
	backend.Setup(ctx, conf)
	return backend, nil
}

func newBackend(client secretsClient) *backend {
	adBackend := &backend{
		client:         client,
		roleCache:      cache.New(roleCacheExpiration, roleCacheCleanup),
		credCache:      cache.New(credCacheExpiration, credCacheCleanup),
		rotateRootLock: new(int32),
	}
	adBackend.Backend = &framework.Backend{
		Help: backendHelp,
		Paths: []*framework.Path{
			adBackend.pathConfig(),
			adBackend.pathRoles(),
			adBackend.pathListRoles(),
			adBackend.pathCreds(),
			adBackend.pathRotateCredentials(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				configPath,
				credPrefix,
			},
		},
		Invalidate:  adBackend.Invalidate,
		BackendType: logical.TypeLogical,
	}
	return adBackend
}

type backend struct {
	logical.Backend

	client secretsClient

	roleCache      *cache.Cache
	credCache      *cache.Cache
	credLock       sync.Mutex
	rotateRootLock *int32
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
