package pki

import (
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Factory creates a new backend implementing the logical.Backend interface
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

// Backend returns a new Backend framework struct
func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
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
			pathGenerateCA(&b),
			pathSignCA(&b),
			pathSetCA(&b),
			pathConfigCA(&b),
			pathConfigCRL(&b),
			pathIssue(&b),
			pathRotateCRL(&b),
			pathFetchCA(&b),
			pathFetchCRL(&b),
			pathFetchCRLViaCertPath(&b),
			pathFetchValid(&b),
			pathRevoke(&b),
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},
	}

	b.crlLifetime = time.Hour * 72
	b.revokeStorageLock = &sync.Mutex{}

	return b.Backend
}

type backend struct {
	*framework.Backend

	crlLifetime       time.Duration
	revokeStorageLock *sync.Mutex
}

const backendHelp = `
The PKI backend dynamically generates X509 server and client certificates.

After mounting this backend, configure the CA using the "pem_bundle" endpoint within
the "config/" path.
`
