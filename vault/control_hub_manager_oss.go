// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

type ControlHubManager struct{}

// NO-OP
func NewControlHubManager(c *Core) *ControlHubManager {
	return &ControlHubManager{}
}

func (c ControlHubManager) WriteData(ctx context.Context, key string, value []byte) error {
	return nil
}
