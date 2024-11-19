// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"testing"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entIdentityStoreDeterminismTestSetup(t *testing.T, ctx context.Context, c *Core, upme, localme *MountEntry) {
	// no op
}

func entIdentityStoreDeterminismAssert(t *testing.T, i int, loadedIDs, lastIDs []string) {
	// no op
}
