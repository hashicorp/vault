// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package github

import (
	"context"
	"net/url"

	"github.com/google/go-github/github"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.TeamMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "teams",
		},
		DefaultKey: "default",
	}

	b.UserMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "users",
		},
		DefaultKey: "default",
	}

	allPaths := append(b.TeamMap.Paths(), b.UserMap.Paths()...)
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},

		Paths:       append([]*framework.Path{pathConfig(&b), pathLogin(&b)}, allPaths...),
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
	}

	return &b
}

type backend struct {
	*framework.Backend

	TeamMap *framework.PolicyMap

	UserMap *framework.PolicyMap
}

// Client returns the GitHub client to communicate to GitHub via the
// configured settings.
func (b *backend) Client(token string) (*github.Client, error) {
	tc := cleanhttp.DefaultClient()
	if token != "" {
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tc)
		tc = oauth2.NewClient(ctx, &tokenSource{Value: token})
	}

	client := github.NewClient(tc)
	emptyUrl, err := url.Parse("")
	if err != nil {
		return nil, err
	}
	client.UploadURL = emptyUrl

	return client, nil
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
