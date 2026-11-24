// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"testing"
)

// TestGetParentNamespaceIDCE verifies that the default CE parent namespace
// resolver remains a no-op.
func TestGetParentNamespaceIDCE(t *testing.T) {
	if got := getParentNamespaceID(nil, "ns1/ns2/"); got != "" {
		t.Fatalf("expected empty parent namespace id in CE, got %q", got)
	}
}
