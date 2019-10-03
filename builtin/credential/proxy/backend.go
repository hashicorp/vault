package proxy

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	rolePrefixStoragePath string = "role/"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	b := backend{}
	b.Backend = &framework.Backend{
		Help:        backendHelp,
		BackendType: logical.TypeCredential,
		AuthRenew:   b.pathLoginRenew,
		Paths: []*framework.Path{
			pathLogin(&b),
			pathConfig(&b),
			pathRoleList(&b),
			pathRole(&b),
		},
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The "proxy" credential provider allows authentication using headers provided
by a trusted proxy server.  A client connects to vault via a trusted proxy
server, which performs its own authentication, adds a header with the
authenticated username, and then forwards the request to vault.  The
"proxy" credential provider can then issue a token based on the header
that was written by the proxy.

The "proxy" credential provider must only ever be used when all requests
from vault are first processed by a trusted proxy server.
`
