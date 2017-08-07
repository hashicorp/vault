package totp

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/totputil"
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

func NewBackend() *totputil.Backend {
	var b totputil.Backend
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

const backendHelp = `
The TOTP backend dynamically generates time-based one-time use passwords.
`

func pathCode(b *totputil.Backend) *framework.Path {
	return b.PathCode("")
}

func pathListKeys(b *totputil.Backend) *framework.Path {
	return b.PathListKeys("")
}

func pathKeys(b *totputil.Backend) *framework.Path {
	return b.PathKeys("")
}
