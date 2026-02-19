// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// NewAgentRegistry

type AgentRegistry struct{}

func (ar *AgentRegistry) loadRegistrations(_ context.Context, _ bool) error {
	return nil
}

func NewAgentRegistry(core *Core, config *logical.BackendConfig, logger log.Logger) (*AgentRegistry, error) {
	return &AgentRegistry{}, nil
}
