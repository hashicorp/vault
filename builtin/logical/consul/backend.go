package consul

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Paths: []*framework.Path{
			pathConfigAccess(),
			pathRoles(),
			pathToken(&b),
		},

		Secrets: []*framework.Secret{
			secretToken(&b),
		},
	}

	return &b
}

type backend struct {
	*framework.Backend
}
