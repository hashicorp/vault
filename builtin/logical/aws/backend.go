package aws

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config/root",
			},
		},

		Paths: []*framework.Path{
			pathConfigRoot(),
			pathConfigLease(&b),
			pathRoles(),
			pathListRoles(&b),
			pathUser(&b),
			pathSTS(&b),
		},

		Secrets: []*framework.Secret{
			secretAccessKeys(&b),
		},

		WALRollback:       walRollback,
		WALRollbackMinAge: 5 * time.Minute,
		BackendType:       logical.TypeLogical,
	}

	return &b
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
