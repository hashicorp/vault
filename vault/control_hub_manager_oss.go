// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
)

const (
	ControlHubClusterIdsPrefix = "cluster/"
)

type ControlHubManager struct{}

// NO-OP
func NewControlHubManager(c *Core) *ControlHubManager {
	return &ControlHubManager{}
}

func (c *ControlHubManager) WriteClusterCredentialsToStorage(ctx context.Context, key string, value []byte) error {
	return fmt.Errorf("control hub is only available in Vault Enterprise")
}
