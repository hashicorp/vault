// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package totp

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	ttlcache "github.com/jellydator/ttlcache/v3"
)

const operationPrefixTOTP = "totp"

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"key/",
			},
		},

		Paths: []*framework.Path{
			pathListKeys(&b),
			pathKeys(&b),
			pathCode(&b),
		},

		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	b.usedCodes = ttlcache.New[string, struct{}]()
	go b.usedCodes.Start()

	return &b
}

type backend struct {
	*framework.Backend

	usedCodes *ttlcache.Cache[string, struct{}]
}

const backendHelp = `
The TOTP backend dynamically generates time-based one-time use passwords.
`
