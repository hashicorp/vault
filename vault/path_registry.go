// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import "github.com/hashicorp/vault/sdk/physical"

// This is a separate file in vault package for now to keep it simple. In the
// future if we choose to break up vault package such that namespaces or mount
// table are moved, this could become a separate helper package but since it's
// so trivial and not clear if it should be used for other types of paths it
// seems like overkill to design a whole package API when we don't know the
// future requirements yet.

var registeredMountOrNamespaceTableKeys []string

// registerMountOrNamespaceTablePaths is meant to be called in init() to allow
// code within vault core to register "special" backend storage keys that
// should be considered part of the mount table and/or namespace metadata in a
// natural place alongside the rest of the code that matters for that subsystem.
// We need to know them centrally within NewCore so that we can register them
// with Backends that want to apply different limits to mount table entries and
// namespace config.
func registerMountOrNamespaceTablePaths(paths ...string) {
	registeredMountOrNamespaceTableKeys = append(registeredMountOrNamespaceTableKeys, paths...)
}

func applyMountAndNamespaceTableKeys(b physical.Backend) {
	if b, ok := b.(physical.MountTableLimitingBackend); ok {
		for _, path := range registeredMountOrNamespaceTableKeys {
			b.RegisterMountTablePath(path)
		}
	}
}
