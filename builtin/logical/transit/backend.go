package transit

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
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
			// Rotate/Config needs to come before Keys
			// as the handler is greedy
			pathConfig(),
			pathRotate(),
			pathRewrap(),
			pathKeys(),
			pathEncrypt(),
			pathDecrypt(),
			pathDatakey(),
		},

		Secrets: []*framework.Secret{},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}
