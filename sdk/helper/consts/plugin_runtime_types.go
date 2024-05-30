// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consts

// NOTE: this file has been copied to
// https://github.com/hashicorp/vault/blob/main/api/plugin_runtime_types.go
// Any changes made should be made to both files at the same time.

import "fmt"

var PluginRuntimeTypes = _PluginRuntimeTypeValues

//go:generate enumer -type=PluginRuntimeType -trimprefix=PluginRuntimeType -transform=snake
type PluginRuntimeType uint32

// This is a list of PluginRuntimeTypes used by Vault.
const (
	DefaultContainerPluginOCIRuntime = "runsc"

	PluginRuntimeTypeUnsupported PluginRuntimeType = iota
	PluginRuntimeTypeContainer
)

// ParsePluginRuntimeType is a wrapper around PluginRuntimeTypeString kept for backwards compatibility.
func ParsePluginRuntimeType(PluginRuntimeType string) (PluginRuntimeType, error) {
	t, err := PluginRuntimeTypeString(PluginRuntimeType)
	if err != nil {
		return PluginRuntimeTypeUnsupported, fmt.Errorf("%q is not a supported plugin runtime type", PluginRuntimeType)
	}
	return t, nil
}
