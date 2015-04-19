package aws

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(map[string]string) (logical.Backend, error) {
	return Backend(), nil
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
				"policy/*",
			},
		},

		Paths: []*framework.Path{
			pathConfigRoot(),
			pathConfigLease(&b),
			pathPolicy(),
			pathUser(&b),
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
the "policy/" endpoints before any access keys can be generated.
`
