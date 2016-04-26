package userpass

import (
	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathUsers(&b),
			pathUsersList(&b),
			pathUserPolicies(&b),
			pathUserPassword(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The "userpass" credential provider allows authentication using
a combination of a username and password. No additional factors
are supported.

The username/password combination is configured using the "users/"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
