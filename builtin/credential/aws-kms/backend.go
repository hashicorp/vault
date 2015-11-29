package awsKms

import (
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

type backend struct {
	*framework.Backend
	MapKey *framework.PolicyMap
}

func Backend(conf *logical.BackendConfig) (*framework.Backend, error) {

	var b backend
	b.MapKey = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "key-id",
			Schema: map[string]*framework.FieldSchema{
				"display_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "A name to map to this ARN for logs.",
				},

				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Policies for the ARN.",
				},
			},
		},
		DefaultKey: "default",
	}

	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			}, Root: []string{
				"config/*",
			},
		},

		Paths: framework.PathAppend([]*framework.Path{
			pathLogin(&b),
			pathConfigUser(),
		}, b.MapKey.Paths(),
		),
	}

	return b.Backend, nil
}

const backendHelp = `
TODO
`
