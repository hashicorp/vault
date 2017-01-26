package chefnode

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
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
			pathLogin(&b),
			pathConfig(&b),
			pathEnvironments(&b),
			pathEnvironmentsList(&b),
			pathRoles(&b),
			pathRolesList(&b),
			pathTags(&b),
			pathTagsList(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}
	return &b
}

type backend struct {
	*framework.Backend
}

func parsePrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("Couldn't parse PEM data")
	}
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privkey, nil
}

const backendHelp = `
"chef-node" authentication backend
`
