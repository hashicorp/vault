package consul

import (
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
			pathConfigAccess(),
			pathListRoles(&b),
			pathRoles(),
			pathToken(&b),
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
