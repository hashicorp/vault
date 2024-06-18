// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package command

import (
	_ "github.com/hashicorp/vault/helper/builtinplugins"
)

func extendAddonHandlers(*vaultHandlers) {
	// No-op
}
