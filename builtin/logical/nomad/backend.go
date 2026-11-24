// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package nomad

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// nomadMaxConnsPerHost is the number of concurrent TCP connections the HTTP transport
// will open per Nomad host. 10 was chosen to avoid exceeding Nomad's default
// http_max_conns_per_client (100) under high concurrency — see VAULT-20923.
const nomadMaxConnsPerHost = 10

const operationPrefixNomad = "nomad"

// Factory returns a Nomad backend that satisfies the logical.Backend interface
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns the configured Nomad backend
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/access",
			},
		},

		Paths: []*framework.Path{
			pathConfigAccess(&b),
			pathConfigLease(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretToken(&b),
		},
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	// cachedClient is a shared Nomad API client reused across requests to
	// avoid opening a new TCP connection per credential operation. Access
	// must be guarded by clientLock.
	cachedClient *api.Client
	clientLock   sync.RWMutex
}

// clientFromConfig builds a new Nomad API client from the given access config.
// It applies a MaxConnsPerHost limit on the HTTP transport so that Vault never
// exceeds Nomad's http_max_conns_per_client limit, and explicitly wires up TLS
// via api.ConfigureTLS so that CA cert / client cert / client key are honoured
// even though we supply our own http.Client.
func clientFromConfig(conf *accessConfig) (*api.Client, error) {
	nomadConf := api.DefaultConfig()
	if conf != nil {
		if conf.Address != "" {
			nomadConf.Address = conf.Address
		}
		if conf.Token != "" {
			nomadConf.SecretID = conf.Token
		}
		// Populate TLSConfig so ConfigureTLS below can apply it.
		if conf.CACert != "" {
			nomadConf.TLSConfig.CACertPEM = []byte(conf.CACert)
		}
		if conf.ClientCert != "" {
			nomadConf.TLSConfig.ClientCertPEM = []byte(conf.ClientCert)
		}
		if conf.ClientKey != "" {
			nomadConf.TLSConfig.ClientKeyPEM = []byte(conf.ClientKey)
		}
	}

	// Build a pooled HTTP transport with a connection cap.
	// NOTE: When HttpClient is set, the Nomad SDK ignores TLSConfig, so we
	// must call api.ConfigureTLS manually afterwards.
	httpClient := cleanhttp.DefaultPooledClient()
	if t, ok := httpClient.Transport.(*http.Transport); ok {
		t.MaxConnsPerHost = nomadMaxConnsPerHost
		// Ensure TLSClientConfig is non-nil so api.ConfigureTLS can populate
		// it safely. cleanhttp.DefaultPooledClient leaves it nil.
		if t.TLSClientConfig == nil {
			t.TLSClientConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}
	}

	// Apply TLS (CA cert, client cert/key) to our custom transport.
	// This is required because the Nomad SDK skips TLS setup when HttpClient
	// is provided — failing to do this would silently drop mTLS config.
	if err := api.ConfigureTLS(httpClient, nomadConf.TLSConfig); err != nil {
		return nil, err
	}

	nomadConf.HttpClient = httpClient
	return api.NewClient(nomadConf)
}

// client returns the cached Nomad API client, building and caching it on first
// use. It is safe for concurrent use.
func (b *backend) client(ctx context.Context, s logical.Storage) (*api.Client, error) {
	// Fast path: return cached client under read lock.
	b.clientLock.RLock()
	if b.cachedClient != nil {
		defer b.clientLock.RUnlock()
		return b.cachedClient, nil
	}
	b.clientLock.RUnlock()

	// Slow path: build and cache under write lock.
	b.clientLock.Lock()
	defer b.clientLock.Unlock()

	// Re-check after acquiring the write lock (double-checked locking).
	if b.cachedClient != nil {
		return b.cachedClient, nil
	}

	conf, err := b.readConfigAccess(ctx, s)
	if err != nil {
		return nil, err
	}

	c, err := clientFromConfig(conf)
	if err != nil {
		return nil, err
	}

	b.cachedClient = c
	return b.cachedClient, nil
}

// resetClient closes idle connections on the cached client and clears it,
// forcing a rebuild on the next request. Must be called whenever config/access
// is written or deleted so the new configuration is picked up.
func (b *backend) resetClient() {
	b.clientLock.Lock()
	defer b.clientLock.Unlock()
	if b.cachedClient != nil {
		b.cachedClient.Close()
		b.cachedClient = nil
	}
}

// isCachedClientNil reports whether the cached client is nil.
// It acquires the read lock so callers (including tests) are race-safe.
func (b *backend) isCachedClientNil() bool {
	b.clientLock.RLock()
	defer b.clientLock.RUnlock()
	return b.cachedClient == nil
}
