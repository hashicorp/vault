package nomad

import (
	"fmt"
	"sync"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Paths: []*framework.Path{
			pathConfigAccess(),
			pathConfigLease(&b),
			pathListRoles(&b),
			pathRoles(),
			pathToken(&b),
		},

		Secrets: []*framework.Secret{
			secretToken(&b),
		},
		BackendType: logical.TypeLogical,
		Clean:       b.resetClient,
	}

	return &b
}

type backend struct {
	*framework.Backend

	client *api.Client
	lock   sync.RWMutex
}

func (b *backend) Client(s logical.Storage) (*api.Client, error) {

	b.lock.RLock()

	// If we already have a client, return it
	if b.client != nil {
		b.lock.RUnlock()
		return b.client, nil
	}

	b.lock.RUnlock()

	conf, intErr := readConfigAccess(s)
	if intErr != nil {
		return nil, intErr
	}
	if conf == nil {
		return nil, fmt.Errorf("no error received but no configuration found")
	}

	nomadConf := api.DefaultConfig()
	nomadConf.Address = conf.Address
	nomadConf.SecretID = conf.Token

	b.lock.Lock()
	defer b.lock.Unlock()

	// If the client was creted during the lock switch, return it
	if b.client != nil {
		return b.client, nil
	}
	var err error
	b.client, err = api.NewClient(nomadConf)
	if err != nil {
		return nil, err
	}
	return b.client, nil
}

func (b *backend) resetClient() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.client = nil
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
