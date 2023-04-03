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

// WaitForNodesExcludingSelectedStandbys is variation on WaitForActiveNodeAndStandbys.
// It waits for the active node before waiting for standby nodes, however
// it will not wait for cores with indexes that match those specified as arguments.
// Whilst you could specify index 0 which is likely to be the leader node, the function
// checks for the leader first regardless of the indexes to skip, so it would be redundant to do so.
// The intention/use case for this function is to allow a cluster to start and become active with one
// or more nodes not joined, so that we can test scenarios where a node joins later.
// e.g. 4 nodes in the cluster, only 3 nodes in cluster 'active', 1 node can be joined later in tests.
func WaitForNodesExcludingSelectedStandbys(t testing.T, cluster *vault.TestCluster, indexesToSkip ...int) {
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
