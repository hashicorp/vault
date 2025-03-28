// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly && !enterprise

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/stretchr/testify/require"
)

// TestActivityLog_setupClientIDsUsageInfo_CE verifies that upon startup, the client IDs are not loaded in CE
func TestActivityLog_setupClientIDsUsageInfo_CE(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	core.setupClientIDsUsageInfo(context.Background())

	// wait for clientIDs to be loaded into memory
	verifyClientsLoadedInMemory := func() {
		corehelpers.RetryUntil(t, 60*time.Second, func() error {
			if a.GetClientIDsUsageInfoLoaded() {
				return fmt.Errorf("loaded clientIDs to memory")
			}
			return nil
		})
	}
	verifyClientsLoadedInMemory()

	require.Len(t, a.GetClientIDsUsageInfo(), 0)
}
