package github

import (
	"net/http"

	"github.com/google/go-github/github"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/oauth2"
)

func Factory(map[string]string) (logical.Backend, error) {
	return Backend(), nil
}

func Backend() *framework.Backend {
	var b backend
	b.Map = &framework.PolicyMap{
		PathMap:    &framework.PathMap{"teams"},
		DefaultKey: "default",
	}
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config",
			},

			Unauthenticated: []string{
				"login",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(),
			pathLogin(&b),
		}, b.Map.Paths()...),
	}

	return b.Backend
}

type backend struct {
	*framework.Backend

	Map *framework.PolicyMap
}

// Client returns the GitHub client to communicate to GitHub via the
// configured settings.
func (b *backend) Client(token string) (*github.Client, error) {
	var tc *http.Client
	if token != "" {
		tc = oauth2.NewClient(oauth2.NoContext, &tokenSource{Value: token})
	}

	return github.NewClient(tc), nil
}

// tokenSource is an oauth2.TokenSource implementation.
type tokenSource struct {
	Value string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: t.Value}, nil
}

const backendHelp = `
The GitHub credential provider allows authentication via GitHub.

Users provide a personal access token to log in, and the credential
provider verifies they're part of the correct organization and then
maps the user to a set of Vault policies according to the teams they're
part of.

After enabling the credential provider, use the "config" route to
configure it.
`
