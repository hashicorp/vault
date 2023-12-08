// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

func stopPartialSealRewrapping(c *Core) {
	// nothing to do
}
