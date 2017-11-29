package rabbitmq

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/michaelklishin/rabbit-hole"
)

// Factory creates and configures the backend
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(conf); err != nil {
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
func (b *backend) Client(s logical.Storage) (*rabbithole.Client, error) {
	b.lock.RLock()

	// If we already have a client, return it
	if b.client != nil {
		b.lock.RUnlock()
		return b.client, nil
	}

	b.lock.RUnlock()

	// Otherwise, attempt to make connection
	entry, err := s.Get("config/connection")
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

	// If the client was creted during the lock switch, return it
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
func (b *backend) resetClient() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.client = nil
}

func (b *backend) invalidate(key string) {
	switch key {
	case "config/connection":
		b.resetClient()
	}
}

// Lease returns the lease information
func (b *backend) Lease(s logical.Storage) (*configLease, error) {
	entry, err := s.Get("config/lease")
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
