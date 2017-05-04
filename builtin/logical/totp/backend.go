package totp

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend(conf).Setup(conf)
}

func Backend(conf *logical.BackendConfig) *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathListKeys(&b),
			pathKeys(&b),
			pathCode(&b),
		},

		Secrets: []*framework.Secret{},
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The TOTP backend dynamically generates time-based one-time use passwords.
`
