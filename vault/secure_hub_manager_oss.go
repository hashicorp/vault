// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
)

const (
	SecureHubClusterIdsPrefix = "cluster/"
)

type SecureHubManager struct{}

// GetSecureHubManager returns the node's SecureHubManager. It returns nil if
// the manager has not been initialized yet.
func (c *Core) GetSecureHubManager() *SecureHubManager {
	c.secureHubManagerLock.RLock()
	defer c.secureHubManagerLock.RUnlock()
	return c.SecureHubManager
}

// SetSecureHubManager sets the node's SecureHubManager. The manager must be
// built before calling this: NewSecureHubManager reads from storage and can
// call back into the Core it is handed as a SecureHubNode, which would
// deadlock on the non-reentrant lock.
func (c *Core) SetSecureHubManager(manager *SecureHubManager) {
	c.secureHubManagerLock.Lock()
	defer c.secureHubManagerLock.Unlock()
	c.SecureHubManager = manager
}

// NO-OP
func NewSecureHubManager(c *Core) *SecureHubManager {
	return &SecureHubManager{}
}

func (c *SecureHubManager) WriteClusterCredentialsToStorage(ctx context.Context, key string, value []byte) error {
	return fmt.Errorf("secure hub is only available in Vault Enterprise")
}
