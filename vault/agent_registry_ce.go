// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type AgentRegistry struct {
	*framework.Backend
}

func (ar *AgentRegistry) loadRegistrations(_ context.Context, _ bool) error {
	return nil
}

func (ar *AgentRegistry) paths() []*framework.Path {
	return []*framework.Path{}
}

func (ar *AgentRegistry) Invalidate(_ context.Context, _ string) {}

func NewAgentRegistry(core *Core, config *logical.BackendConfig, logger log.Logger) (*AgentRegistry, error) {
	ar := &AgentRegistry{}

	ar.Backend = &framework.Backend{
		BackendType:    logical.TypeLogical,
		Paths:          ar.paths(),
		Invalidate:     ar.Invalidate,
		RunningVersion: versions.DefaultBuiltinVersion,
	}
	return ar, nil
}
