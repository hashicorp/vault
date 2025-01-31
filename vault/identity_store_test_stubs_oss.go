// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"math/rand"
	"testing"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entIdentityStoreDeterminismSecondaryTestSetup(t *testing.T, ctx context.Context, c *Core, me, localme *MountEntry, seed *rand.Rand) {
	// no op
}

func entIdentityStoreDeterminismSecondaryAssert(t *testing.T, i int, loadedIDs, lastIDs []string) {
	// no op
}

func entIdentityStoreDuplicateReportTestSetup(t *testing.T, ctx context.Context, c *Core, rootToken string, seed *rand.Rand) {
	// no op
}

func identityStoreDuplicateReportTestWantDuplicateCounts() (int, int, int, int) {
	// Note that the second count is for local aliases. CE Vault doesn't really
	// distinguish between local and non-local aliases because it doesn't have any
	// support for Performance Replication. But it's possible in code at least to
	// set the local flag on a mount or alias during creation so we might as well
	// test it behaves as expected in the CE code. It's maybe just about possible
	// that this could happen in real life too because of a downgrade from
	// Enterprise.
	return 1, 1, 1, 1
}
