// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import hcpvlib "github.com/hashicorp/vault-hcp-lib"

// DefaultHCPTokenHelper returns the HCP token helper that is configured for Vault.
// This helper should only be used for non-server CLI commands.
func DefaultHCPTokenHelper() hcpvlib.HCPTokenHelper {
	return &hcpvlib.InternalHCPTokenHelper{}
}
