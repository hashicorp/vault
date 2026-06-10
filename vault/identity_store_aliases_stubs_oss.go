// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

func (c *Core) validateSyntheticAliasAccessor(context.Context, string) (bool, error) {
	return false, nil
}

func (c *Core) generateSyntheticAliasAccessor(context.Context, string) (string, error) {
	return "", nil
}
