package ssh

import (
	"strings"

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
			Root: []string{"config/*"},
		},

		Paths: []*framework.Path{
			pathConfigLease(&b),
			pathKeys(&b),
			pathRoles(&b),
			pathRoleCreate(&b),
			pathLookup(&b),
		},

		Secrets: []*framework.Secret{
			secretSSHKey(&b),
		},
	}
	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The SSH backend dynamically generates SSH private keys for 
remote hosts.The generated key has a configurable lease set
and are automatically revoked at the end of the lease.

After mounting this backend, configure the lease using the
'config/lease' endpoint. The shared SSH key belonging to any
infrastructure should be registered with the 'roles/' endpoint
before dynamic keys for remote hosts can be generated.
`
