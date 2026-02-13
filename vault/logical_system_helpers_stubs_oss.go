// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "github.com/hashicorp/vault/sdk/logical"

func forwardCertCounts(c *Core, inc logical.CertCount) bool {
	return false
}
