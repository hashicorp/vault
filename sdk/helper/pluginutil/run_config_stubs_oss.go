// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package pluginutil

type EntPluginClientConfig struct{}

func (p *PluginClientConfig) EntUpdate(_ *PluginRunner) {
	// no-op
}
