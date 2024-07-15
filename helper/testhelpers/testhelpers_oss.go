// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package testhelpers

import (
	"testing"

	"github.com/hashicorp/vault/vault"
)

// WaitForActiveNodeAndStandbys does nothing more than wait for the active node
// on OSS. On enterprise it waits for perf standbys to be healthy too.
func WaitForActiveNodeAndStandbys(t testing.TB, cluster *vault.TestCluster) {
	WaitForActiveNode(t, cluster)
	for _, core := range cluster.Cores {
		if standby, _ := core.Core.Standby(); standby {
			WaitForStandbyNode(t, core)
		}
	}
}
