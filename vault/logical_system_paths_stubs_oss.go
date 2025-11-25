// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "github.com/hashicorp/vault/sdk/framework"

func entWrappedPluginsCRUDPath(b *SystemBackend) []*framework.Path {
	return []*framework.Path{b.pluginsCatalogCRUDPath()}
}

func entWrappedAuthPath(b *SystemBackend) []*framework.Path {
	return b.authPaths()
}

func entWrappedMountsPath(b *SystemBackend) []*framework.Path {
	return b.mountsPaths()
}
