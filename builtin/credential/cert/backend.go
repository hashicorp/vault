package cert

import (
	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	_, err := b.Setup(conf)
	if err != nil {
		return b, err
	}
	return b, b.populateCRLs(conf.StorageView)
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
			pathLogin(&b),
			pathListCerts(&b),
			pathCerts(&b),
			pathCRLs(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}

	b.crls = map[string]CRLInfo{}
	b.crlUpdateMutex = &sync.RWMutex{}

	return &b
}

type backend struct {
	*framework.Backend
	MapCertId *framework.PathMap

	crls           map[string]CRLInfo
	crlUpdateMutex *sync.RWMutex
}

const backendHelp = `
The "cert" credential provider allows authentication using
TLS client certificates. A client connects to Vault and uses
the "login" endpoint to generate a client token.

Trusted certificates are configured using the "certs/" endpoint
by a user with root access. A certificate authority can be trusted,
which permits all keys signed by it. Alternatively, self-signed
certificates can be trusted avoiding the need for a CA.
`
