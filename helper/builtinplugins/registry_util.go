// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package builtinplugins

import "github.com/hashicorp/vault/sdk/helper/consts"

// IsBuiltinEntPlugin checks whether the plugin is an enterprise only builtin plugin
func (r *registry) IsBuiltinEntPlugin(name string, pluginType consts.PluginType) bool {
	return false
}
