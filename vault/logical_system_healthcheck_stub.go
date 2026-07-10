// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type HealthCheckBackend struct {
	*framework.Backend
}

func NewHealthCheckBackend(_ *Core, _ log.Logger) *HealthCheckBackend {
	return &HealthCheckBackend{
		&framework.Backend{
			Help:        healthCheckStubBackendHelp,
			BackendType: logical.TypeLogical,
		},
	}
}

const healthCheckStubBackendHelp = `
Unimplemented stub for the enterprise-only health check feature.
`
