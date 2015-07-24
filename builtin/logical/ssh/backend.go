package ssh

import (
	"strings"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	return b.Setup(conf)
}

func Backend(conf *logical.BackendConfig) (*framework.Backend, error) {
	salt, err := salt.NewSalt(conf.View, nil)
	if err != nil {
		return nil, err
	}

	var b backend
	b.salt = salt
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
				"keys/*",
			},
			Unauthenticated: []string{
				"verify",
			},
		},

		Paths: []*framework.Path{
			pathConfigLease(&b),
			pathKeys(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathLookup(&b),
			pathVerify(&b),
		},

		Secrets: []*framework.Secret{
			secretDynamicKey(&b),
			secretOTP(&b),
		},
	}
	return b.Backend, nil
}

type backend struct {
	*framework.Backend
	salt *salt.Salt
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
