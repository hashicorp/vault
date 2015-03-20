package aws

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(map[string]string) (logical.Backend, error) {
	return Backend(), nil
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		PathsRoot: []string{
			"root",
			"policy/*",
		},

		Paths: []*framework.Path{
			pathRoot(),
			pathPolicy(),
			pathUser(&b),
		},

		Secrets: []*framework.Secret{
			secretAccessKeys(),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}
