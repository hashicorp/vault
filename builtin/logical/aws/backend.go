package aws

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
			},
		},

		Paths: []*framework.Path{
			pathConfigRoot(),
			pathConfigLease(&b),
			pathRoles(),
			pathUser(&b),
			pathSTS(&b),
		},

		Secrets: []*framework.Secret{
			secretAccessKeys(&b),
		},

		Rollback:       rollback,
		RollbackMinAge: 5 * time.Minute,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The AWS backend dynamically generates AWS access keys for a set of
IAM policies. The AWS access keys have a configurable lease set and
are automatically revoked at the end of the lease.

After mounting this backend, credentials to generate IAM keys must
be configured with the "root" path and policies must be written using
the "roles/" endpoints before any access keys can be generated.
`
