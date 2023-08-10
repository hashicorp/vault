// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type SecretsSyncBackend struct {
	*framework.Backend
}

func NewSecretsSyncBackend(_ *Core, _ log.Logger) *SecretsSyncBackend {
	return &SecretsSyncBackend{
		&framework.Backend{
			Help:        stubBackendHelp,
			BackendType: logical.TypeLogical,
		},
	}
}

const stubBackendHelp = `
Unimplemented stub for the enterprise-only secrets sync feature.
`
