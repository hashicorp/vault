package ssh

import (
	"log"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(map[string]string) (logical.Backend, error) {
	return Backend(), nil
}

func Backend() *framework.Backend {
	log.Printf("Vishal: ssh.Backend\n")
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
			secretSshKey(&b),
		},
	}
	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The ssh backend enables secure connections to remote hosts.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
