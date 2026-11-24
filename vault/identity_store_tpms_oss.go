// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
)

// tpmPaths returns empty paths for OSS builds
func tpmPaths(i *IdentityStore) []*framework.Path {
	return nil
}
