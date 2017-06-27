package totp

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// New returns a new backend as an interface. This func
// is only necessary for builtin backend plugins.
func New() (interface{}, error) {
	return Backend(), nil
}

// Factory returns a new backend a logical.Backend.
func Factory() (logical.Backend, error) {
	return Backend(), nil
}

// Backend returns a private embedded struct of framework.Backend.
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathListKeys(&b),
			pathKeys(&b),
			pathCode(&b),
		},

		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The TOTP backend dynamically generates time-based one-time use passwords.
`
