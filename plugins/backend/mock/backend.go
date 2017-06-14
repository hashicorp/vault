package mock

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// New returns a new backend as an interface. This func
// is only necessary for builtin backend plugins.
func New() (interface{}, error) {
	return Backend(), nil
}

// Factory returns a new backend a logical.Backend.
func Factory() (logical.Backend, error) {
	return Backend(), nil
}

// Backend returns a private embedded struct of framework.Backend.
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: "",
		Paths: []*framework.Path{
			pathTesting(&b),
			pathInternal(&b),
		},
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"special",
			},
		},
		Secrets:    []*framework.Secret{},
		Invalidate: b.invalidate,
	}
	b.internal = "bar"
	return &b
}

type backend struct {
	*framework.Backend

	// internal is used to test invalidate
	internal string
}

func (b *backend) invalidate(key string) {
	switch key {
	case "internal":
		b.internal = ""
	}
}
