// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "github.com/hashicorp/vault/sdk/framework"

// tpmgroupPaths returns empty paths for OSS builds
func tpmgroupPaths(i *IdentityStore) []*framework.Path {
	return nil
}
