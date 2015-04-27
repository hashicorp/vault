package transit

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
		PathsSpecial: &logical.Paths{
			Root: []string{
				"keys/*",
			},
		},

		Paths: []*framework.Path{
			pathKeys(),
			pathEncrypt(),
			pathDecrypt(),
		},

		Secrets: []*framework.Secret{},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}
