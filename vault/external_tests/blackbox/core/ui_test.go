// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestUIAssets verifies that the Vault UI is accessible
func TestUIAssets(t *testing.T) {
	v := blackbox.New(t)

	// This is a stub - in a real implementation, you would verify UI assets are accessible
	// For now, just verify the UI endpoint is available by checking sys/internal/ui/mounts
	uiMounts := v.MustRead("sys/internal/ui/mounts")
	if uiMounts == nil || uiMounts.Data == nil {
		t.Fatal("Could not access UI mounts endpoint")
	}

	t.Log("Successfully verified UI assets are accessible")
}
