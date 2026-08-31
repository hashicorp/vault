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

// GetControlHubManager returns the node's ControlHubManager. It returns nil if
// the manager has not been initialized yet.
func (c *Core) GetControlHubManager() *ControlHubManager {
	c.controlHubManagerLock.RLock()
	defer c.controlHubManagerLock.RUnlock()
	return c.ControlHubManager
}

// SetControlHubManager sets the node's ControlHubManager. The manager must be
// built before calling this: NewControlHubManager reads from storage and can
// call back into the Core it is handed as a ControlHubNode, which would
// deadlock on the non-reentrant lock.
func (c *Core) SetControlHubManager(manager *ControlHubManager) {
	c.controlHubManagerLock.Lock()
	defer c.controlHubManagerLock.Unlock()
	c.ControlHubManager = manager
}

// NO-OP
func NewControlHubManager(c *Core) *ControlHubManager {
	return &ControlHubManager{}
}

func (c *ControlHubManager) WriteClusterCredentialsToStorage(ctx context.Context, key string, value []byte) error {
	return fmt.Errorf("control hub is only available in Vault Enterprise")
}
