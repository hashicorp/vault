package rabbitmq

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/michaelklishin/rabbit-hole"
)

// Factory creates and configures Backends
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

// Backend creates a new Backend
func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathConfigLease(&b),
			pathRoleCreate(&b),
			pathRoles(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean: b.ResetClient,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend

	client *rabbithole.Client
	lock   sync.Mutex
}

// DB returns the database connection.
func (b *backend) Client(s logical.Storage) (*rabbithole.Client, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If we already have a client, we got it!
	if b.client != nil {
		return b.client, nil
	}

	// Otherwise, attempt to make connection
	entry, err := s.Get("config/connection")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil,
			fmt.Errorf("configure the client connection with config/connection first")
	}

	var connConfig connectionConfig
	if err := entry.DecodeJSON(&connConfig); err != nil {
		return nil, err
	}

	b.client, err = rabbithole.NewClient(connConfig.URI, connConfig.Username, connConfig.Password)
	if err != nil {
		return nil, err
	}

	return b.client, nil
}

// ResetClient forces a connection next time Client() is called.
func (b *backend) ResetClient() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.client = nil
}

// Lease returns the lease information
func (b *backend) Lease(s logical.Storage) (*configLease, error) {
	entry, err := s.Get(leasePatternLabel)
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
