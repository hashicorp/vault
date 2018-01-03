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

func (b *backend) client(s logical.Storage) (*api.Client, error) {
	conf, err := b.readConfigAccess(s)
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
