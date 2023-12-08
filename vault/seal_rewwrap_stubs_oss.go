// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func stopPartialSealRewrapping(c *Core) {
	// nothing to do
}
