// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/logical"
)

var healthCheckNotImplemented error = errors.New("health checks are not implemented in Vault CE")

type HealthCheckManager struct{}

func (c *Core) handleExecHealthCheck(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	return nil, healthCheckNotImplemented
}

func (c *Core) handleReadLastHealthCheck(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	return nil, healthCheckNotImplemented
}

// setupHealthCheckManager allows HealthCheckManager to be set up in
// the unseal strategy of an active or a perf standby node
func (c *Core) setupHealthCheckManager() error {
	return nil
}
