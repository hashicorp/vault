// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !testonly

package vault

import "github.com/hashicorp/vault/sdk/framework"

// entityTestonlyPaths is a stub for non-testonly builds.
func entityTestonlyPaths(i *IdentityStore) []*framework.Path {
	return nil
}
