// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// systemBackendPaths returns the full set of system backend paths.
// On OSS there are no namespace backends, so root always gets all paths.
func systemBackendPaths(b *SystemBackend, _ bool, _ *logical.BackendConfig) []*framework.Path {
	return operatorSystemBackendPaths(b)
}

func entWrappedPluginsCRUDPath(b *SystemBackend) []*framework.Path {
	return []*framework.Path{b.pluginsCatalogCRUDPath()}
}

func entWrappedAuthPath(b *SystemBackend) []*framework.Path {
	return b.authPaths()
}

func entWrappedMountsPath(b *SystemBackend) []*framework.Path {
	return b.mountsPaths()
}
