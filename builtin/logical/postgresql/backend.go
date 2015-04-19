package postgresql

import (
	"strings"

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
			},
		},

		Paths: []*framework.Path{
			pathConfigConnection(),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The PostgreSQL backend dynamically generates database users.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
