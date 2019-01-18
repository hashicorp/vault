package rabbitmq

import (
	"context"
	"fmt"
	"strings"
	"sync"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

// Factory creates and configures the backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Creates a new backend with all the paths and secrets belonging to it
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/connection",
			},
		},

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathConfigLease(&b),
			pathListRoles(&b),
			pathCreds(&b),
			pathRoles(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean:       b.resetClient,
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	client *rabbithole.Client
	lock   sync.RWMutex
}

// DB returns the database connection.
func (b *backend) Client(ctx context.Context, s logical.Storage) (*rabbithole.Client, error) {
	b.lock.RLock()

	// If we already have a client, return it
	if b.client != nil {
		b.lock.RUnlock()
		return b.client, nil
	}

	b.lock.RUnlock()

	// Otherwise, attempt to make connection
	entry, err := s.Get(ctx, "config/connection")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("configure the client connection with config/connection first")
	}

	var connConfig connectionConfig
	if err := entry.DecodeJSON(&connConfig); err != nil {
		return nil, err
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	// If the client was created during the lock switch, return it
	if b.client != nil {
		return b.client, nil
	}

	b.client, err = rabbithole.NewClient(connConfig.URI, connConfig.Username, connConfig.Password)
	if err != nil {
		return nil, err
	}
	// Use a default pooled transport so there would be no leaked file descriptors
	b.client.SetTransport(cleanhttp.DefaultPooledTransport())

	return b.client, nil
}

// resetClient forces a connection next time Client() is called.
func (b *backend) resetClient(_ context.Context) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.client = nil
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config/connection":
		b.resetClient(ctx)
	}
}

// Lease returns the lease information
func (b *backend) Lease(ctx context.Context, s logical.Storage) (*configLease, error) {
	entry, err := s.Get(ctx, "config/lease")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const backendHelp = `
The RabbitMQ backend dynamically generates RabbitMQ users.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
