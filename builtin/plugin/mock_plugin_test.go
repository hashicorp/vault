// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	_ "github.com/hashicorp/vault-plugin-mock"
)

// This file exists to force an import of vault-plugin-mock (which itself does nothing),
// for purposes of CI and GitHub actions testing between plugin repos and Vault.
