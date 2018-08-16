package nomad

import (
	"context"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

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
}

func (b *backend) client(ctx context.Context, s logical.Storage) (*api.Client, error) {
	conf, err := b.readConfigAccess(ctx, s)
	if err != nil {
		return nil, err
	}

	nomadConf := api.DefaultConfig()
	if conf != nil {
		if conf.Address != "" {
			nomadConf.Address = conf.Address
		}
		if conf.Token != "" {
			nomadConf.SecretID = conf.Token
		}
	}

	client, err := api.NewClient(nomadConf)
	if err != nil {
		return nil, err
	}

	return client, nil
}
