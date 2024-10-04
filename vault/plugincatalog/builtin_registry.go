// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import "github.com/hashicorp/vault/sdk/helper/consts"

// BuiltinRegistry is an interface that allows the "vault" package to use
// the registry of builtin plugins without getting an import cycle. It
// also allows for mocking the registry easily.
type BuiltinRegistry interface {
	Contains(name string, pluginType consts.PluginType) bool
	Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool)
	Keys(pluginType consts.PluginType) []string
	DeprecationStatus(name string, pluginType consts.PluginType) (consts.DeprecationStatus, bool)
	IsBuiltinEntPlugin(name string, pluginType consts.PluginType) bool
}
