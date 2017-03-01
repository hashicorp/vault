package httpBasic

import (
	"github.com/hashicorp/vault/helper/mfa"
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
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login",
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathUsersList(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew: b.pathLoginRenew,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The "http-basic" credential provider allows authentication against
a HTTP server, checking username/password and associating users
to set of policies.

Configuration of the server is done through the "config" and "users"
endpoints by a user with approriate access mandated by policy.
Authentication is then done by suppying the two fields for "login".

The backend optionally allows to grant a set of policies to any 
user that successfully authenticates against the HTTP server, 
without them being explicitly mapped in vault.
`
