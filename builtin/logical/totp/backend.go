package totp

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	cache "github.com/patrickmn/go-cache"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := NewBackend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func NewBackend() *Backend {
	var b Backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathListKeys(&b),
			pathKeys(&b),
			pathCode(&b),
		},

		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	b.UsedCodes = cache.New(0, 30*time.Second)

	return &b
}

type Backend struct {
	*framework.Backend

	UsedCodes *cache.Cache
}

const backendHelp = `
The TOTP backend dynamically generates time-based one-time use passwords.
`
