package aws

import (
	"github.com/hashicorp/vault/logical/framework"
)

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
