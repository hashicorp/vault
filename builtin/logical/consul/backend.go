// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const operationPrefixConsul = "consul"

// ReportedVersion is used to report a specific version to Vault.
var ReportedVersion = ""

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
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/access",
			},
		},

		Paths: []*framework.Path{
			pathConfigAccess(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathToken(&b),
		},

		Secrets: []*framework.Secret{
			secretToken(&b),
		},
		BackendType:    logical.TypeLogical,
		RunningVersion: ReportedVersion,
	}

	return &b
}

type backend struct {
	*framework.Backend
}
