package google

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

//Factory for google backend
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

//Backend for google
func Backend() *framework.Backend {
	var b backend
	b.Map = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "teams",
		},
		DefaultKey: "default",
	}
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
				codeURLPath,
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathLogin(&b),
			pathCodeURL(&b),
		}, b.Map.Paths()...),

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend

	Map *framework.PolicyMap
}

const backendHelp = `
The Google credential provider allows authentication via Google.

Users provide a personal access code to log in, and the credential
provider verifies they're part of the correct domain and then
maps the user to a set of Vault policies according to the teams they're
part of.

After enabling the credential provider, use the "config" route to
configure it.
`
