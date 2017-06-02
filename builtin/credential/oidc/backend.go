package oidc

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
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathGroups(&b),
			pathUsersList(&b),
			pathGroupsList(&b),
			pathLogin(&b),
		}),

		AuthRenew: nil, // explicitly don't support renewal.
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The OpenID Connect provider allows Vault to issue Tokens for
holders of OpenID Connect identity tokens, which are self validating.

Only users that have an explicit mapping of username or group to a policy
will be granted Tokens.
`
