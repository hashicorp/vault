// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestUnsealedStatus verifies that the Vault cluster is unsealed (not sealed).
// This test only checks seal status, not general cluster health.
func TestUnsealedStatus(t *testing.T) {
	v := blackbox.New(t)

	// Verify the cluster is unsealed
	v.AssertUnsealedAny()

	t.Log("Successfully verified Vault cluster is unsealed")
}
