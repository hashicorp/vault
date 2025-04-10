// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

// extendedSystemView (Vault Community edition) is a logical.SystemView
// that is extended with logical.ExtendedSystemView and SudoPrivilege
type extendedSystemView interface {
	logical.SystemView
	logical.ExtendedSystemView
	// SudoPrivilege won't work over the plugin system so we keep it here
	// instead of in sdk/logical to avoid exposing to plugins
	SudoPrivilege(context.Context, string, string) bool
}
