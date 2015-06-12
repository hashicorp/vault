package pki

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Factory creates a new backend implementing the logical.Backend interface
func Factory(map[string]string) (logical.Backend, error) {
	return Backend(), nil
}

// Backend returns a new Backend framework struct
func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
				"revoked/*",
				"revoke/*",
				"crl/rotate",
			},
			Unauthenticated: []string{
				"cert/*",
				"ca/pem",
				"ca",
				"crl/pem",
				"crl",
			},
		},

		Paths: []*framework.Path{
			pathRoles(&b),
			pathConfigCA(&b),
			pathIssue(&b),
			pathRotateCRL(&b),
			pathFetchCA(&b),
			pathFetchCRL(&b),
			pathFetchCRLViaCertPath(&b),
			pathFetchValid(&b),
			pathFetchRevoked(&b),
			pathRevoke(&b),
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The PKI backend dynamically generates X509 server and client certificates.

After mounting this backend, configure the CA using the "ca_bundle" endpoint within
the "config/" path.
`
