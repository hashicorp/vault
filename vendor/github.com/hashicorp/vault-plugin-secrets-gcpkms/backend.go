// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	kmsapi "cloud.google.com/go/kms/apiv1"
)

const (
	userAgentPluginName = "secrets-gcpkms"

	// operationPrefixGoogleCloudKMS is used as a prefix for OpenAPI operation id's.
	operationPrefixGoogleCloudKMS = "google-cloud-kms"
)

var (
	// defaultClientLifetime is the amount of time to cache the KMS client. This
	// has to be less than 60 minutes or the oauth token will expire and
	// subsequent requests will fail. The reason we cache the client is because
	// the process for looking up credentials is not performant and the overhead
	// is too significant for a plugin that will receive this much traffic.
	defaultClientLifetime = 30 * time.Minute
)

type backend struct {
	*framework.Backend

	// keysCache holds a temporal copy of keys retrieved from KMS
	keysCache *cache.Cache

	// kmsClient is the actual client for connecting to KMS. It is cached on
	// the backend for efficiency.
	kmsClient           *kmsapi.KeyManagementClient
	kmsClientCreateTime time.Time
	kmsClientLifetime   time.Duration
	kmsClientLock       sync.RWMutex

	// pluginEnv contains Vault version information. It is used in user-agent headers.
	pluginEnv *logical.PluginEnvironment

	// ctx and ctxCancel are used to control overall plugin shutdown. These
	// contexts are given to any client libraries or requests that should be
	// terminated during plugin termination.
	ctx       context.Context
	ctxCancel context.CancelFunc
	ctxLock   sync.Mutex
}

// Factory returns a configured instance of the backend.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns a configured instance of the backend.
func Backend() *backend {
	var b backend

	b.kmsClientLifetime = defaultClientLifetime
	b.ctx, b.ctxCancel = context.WithCancel(context.Background())
	b.keysCache = cache.New(cache.DefaultExpiration, 60*time.Minute)

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help: "The GCP KMS secrets engine provides pass-through encryption and " +
			"decryption to Google Cloud KMS keys.",

		Paths: []*framework.Path{
			b.pathConfig(),

			b.pathKeys(),
			b.pathKeysCRUD(),
			b.pathKeysConfigCRUD(),
			b.pathKeysDeregister(),
			b.pathKeysRegister(),
			b.pathKeysRotate(),
			b.pathKeysTrim(),

			b.pathDecrypt(),
			b.pathEncrypt(),
			b.pathPubkey(),
			b.pathReencrypt(),
			b.pathSign(),
			b.pathVerify(),
		},

		InitializeFunc: b.initialize,
		Invalidate:     b.invalidate,
		Clean:          b.clean,
	}

	return &b
}

func (b *backend) initialize(ctx context.Context, _ *logical.InitializationRequest) error {
	pluginEnv, err := b.System().PluginEnv(ctx)
	if err != nil {
		return fmt.Errorf("failed to read plugin environment: %w", err)
	}
	b.pluginEnv = pluginEnv

	return nil
}

// clean cancels the shared contexts. This is called just before unmounting
// the plugin.
func (b *backend) clean(_ context.Context) {
	b.ctxLock.Lock()
	b.ctxCancel()
	b.ctxLock.Unlock()
}

// invalidate resets the plugin. This is called when a key is updated via
// replication.
func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.ResetClient()
	}
}

// ResetClient closes any connected clients.
func (b *backend) ResetClient() {
	b.kmsClientLock.Lock()
	b.resetClient()
	b.kmsClientLock.Unlock()
}

// resetClient rests the underlying client. The caller is responsible for
// acquiring and releasing locks. This method is not safe to call concurrently.
func (b *backend) resetClient() {
	if b.kmsClient != nil {
		b.kmsClient.Close()
		b.kmsClient = nil
	}

	b.kmsClientCreateTime = time.Unix(0, 0).UTC()
}

// KMSClient creates a new client for talking to the GCP KMS service.
func (b *backend) KMSClient(s logical.Storage) (*kmsapi.KeyManagementClient, func(), error) {
	// If the client already exists and is valid, return it
	b.kmsClientLock.RLock()
	if b.kmsClient != nil && time.Now().UTC().Sub(b.kmsClientCreateTime) < b.kmsClientLifetime {
		closer := func() { b.kmsClientLock.RUnlock() }
		return b.kmsClient, closer, nil
	}
	b.kmsClientLock.RUnlock()

	// Acquire a full lock. Since all invocations acquire a read lock and defer
	// the release of that lock, this will block until all clients are no longer
	// in use. At that point, we can acquire a globally exclusive lock to close
	// any connections and create a new client.
	b.kmsClientLock.Lock()

	b.Logger().Debug("creating new KMS client")

	// Attempt to close an existing client if we have one.
	b.resetClient()

	// Get the config
	config, err := b.Config(b.ctx, s)
	if err != nil {
		b.kmsClientLock.Unlock()
		return nil, nil, err
	}

	// If credentials were provided, use those. Otherwise fall back to the
	// default application credentials.
	var creds *google.Credentials
	if config.Credentials != "" {
		creds, err = google.CredentialsFromJSON(b.ctx, []byte(config.Credentials), config.Scopes...)
		if err != nil {
			b.kmsClientLock.Unlock()
			return nil, nil, errwrap.Wrapf("failed to parse credentials: {{err}}", err)
		}
	} else {
		creds, err = google.FindDefaultCredentials(b.ctx, config.Scopes...)
		if err != nil {
			b.kmsClientLock.Unlock()
			return nil, nil, errwrap.Wrapf("failed to get default token source: {{err}}", err)
		}
	}

	// Create and return the KMS client with a custom user agent.
	client, err := kmsapi.NewKeyManagementClient(b.ctx,
		option.WithCredentials(creds),
		option.WithScopes(config.Scopes...),
		option.WithUserAgent(useragent.PluginString(b.pluginEnv, userAgentPluginName)),
	)
	if err != nil {
		b.kmsClientLock.Unlock()
		return nil, nil, errwrap.Wrapf("failed to create KMS client: {{err}}", err)
	}

	// Cache the client
	b.kmsClient = client
	b.kmsClientCreateTime = time.Now().UTC()
	b.kmsClientLock.Unlock()

	b.kmsClientLock.RLock()
	closer := func() { b.kmsClientLock.RUnlock() }
	return client, closer, nil
}

// Config parses and returns the configuration data from the storage backend.
// Even when no user-defined data exists in storage, a Config is returned with
// the default values.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*Config, error) {
	c := DefaultConfig()

	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, errwrap.Wrapf("failed to get configuration from storage: {{err}}", err)
	}
	if entry == nil || len(entry.Value) == 0 {
		return c, nil
	}

	if err := entry.DecodeJSON(&c); err != nil {
		return nil, errwrap.Wrapf("failed to decode configuration: {{err}}", err)
	}
	return c, nil
}
