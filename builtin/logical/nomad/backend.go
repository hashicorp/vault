package nomad

import (
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
}

func (b *backend) client(s logical.Storage) (*api.Client, error) {
	conf, err := b.readConfigAccess(s)
	if err != nil {
		return nil, err
	}

	nomadConf := api.DefaultConfig()
	nomadConf.Address = conf.Address
	nomadConf.SecretID = conf.Token

	client, err := api.NewClient(nomadConf)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Lease returns the lease information
func (b *backend) LeaseConfig(s logical.Storage) (*configLease, error) {
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
