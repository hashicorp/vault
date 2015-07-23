package jwt

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Factory creates a new backend implementing the logical.Backend interface
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

// Backend returns a new Backend framework struct
func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{

		Paths: []*framework.Path{
			pathRoles(&b),
			pathIssue(&b),
		},

	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}
