// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

func (c *Core) setupOAuthTokenDenylist(ctx context.Context) error {
	return nil
}
