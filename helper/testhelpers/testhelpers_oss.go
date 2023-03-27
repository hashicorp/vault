// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package testhelpers

import (
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

// WaitForActiveNodeAndStandbys does nothing more than wait for the active node
// on OSS. On enterprise it waits for perf standbys to be healthy too.
func WaitForActiveNodeAndStandbys(t testing.T, cluster *vault.TestCluster) {
	WaitForActiveNode(t, cluster)
	for _, core := range cluster.Cores {
		if standby, _ := core.Core.Standby(); standby {
			WaitForStandbyNode(t, core)
		}
	}
}

// WaitForActiveNodeAndSelectedStandbys is variation on WaitForActiveNodeAndStandbys,
// the difference is that this function can be given core indexes which it should not attempt to wait for.
func WaitForActiveNodeAndSelectedStandbys(t testing.T, cluster *vault.TestCluster, indexesToSkip ...int) {
	WaitForActiveNode(t, cluster)

	contains := func(elems []int, e int) bool {
		for _, v := range elems {
			if v == e {
				return true
			}
		}

		return false
	}
	for i, core := range cluster.Cores {
		if contains(indexesToSkip, i) {
			continue
		}

		if standby, _ := core.Core.Standby(); standby {
			WaitForStandbyNode(t, core)
		}
	}
}
