package google

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const BackendName = "google"

//Factory for google backend
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	_, err := b.Setup(conf)
	if err != nil {
		return b, err
	}
	return b, nil
}

const googleBackendHelp = `
The Google credential provider allows you to authenticate with Google. You must own a registered google application.
`+ writeConfigPathHelp + `
Then, proceed to generate a personal access token by browsing to a google url.
` + readCodeUrlPathHelp + `

    Example: vault auth -method=` + BackendName + ` ` + googleAuthCodeParameterName + `=<code>

    the user's google domain will be matched against the domain you configured for the backend, e.g. example.com (or empty string for none)

Key/Value Pairs:

    mount=` + BackendName + `      The mountpoint for the Google credential provider.
                      Defaults to "` + BackendName + `"

    ` + googleAuthCodeParameterName + `=<code>     The Google access code for authentication.
`

const usersToPoliciesMapPath = "users"

//Backend for google
func Backend() *backend {
	var b backend
	b.Map = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: usersToPoliciesMapPath,
		},
	}
	b.Backend = &framework.Backend{
		Help: googleBackendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				loginPath,
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

	return &b
}

type backend struct {
	*framework.Backend

	Map *framework.PolicyMap
}


