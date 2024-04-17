// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

// managedKeyRegistrySubPath is the storage prefix used by the registry.
// We need to define the constant even though managed keys is a Vault Enterprise
// feature in order to set up seal wrapping in the SystemBackend.
const managedKeyRegistrySubPath = "managed-key-registry/"

func (c *Core) setupManagedKeyRegistry() error {
	// Nothing to do, the registry is only used by enterprise features
	return nil
}

func (c *Core) ReloadManagedKeyRegistryConfig() {
	// Nothing to do, the registry is only used by enterprise features
}
