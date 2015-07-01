package consul

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(*logical.BackendConfig) (logical.Backend, error) {
	return Backend(), nil
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
			},
		},

		Paths: []*framework.Path{
			pathConfigAccess(),
			pathRoles(),
			pathToken(&b),
		},

		Secrets: []*framework.Secret{
			secretToken(),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}
