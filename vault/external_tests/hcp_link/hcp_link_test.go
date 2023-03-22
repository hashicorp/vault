// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcp_link

import (
	"testing"

	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/vault"
)

func TestHCPLinkConnected(t *testing.T) {
	cluster := getTestCluster(t, 2)
	defer cluster.Cleanup()

	vaultHCPLink, _ := TestClusterWithHCPLinkEnabled(t, cluster, false, false)
	defer vaultHCPLink.Cleanup()

	for _, core := range cluster.Cores {
		checkLinkStatus(core.Client, scada.SessionStatusConnected, t)
	}
}

func TestHCPLinkNotConfigured(t *testing.T) {
	t.Parallel()
	cluster := getTestCluster(t, 2)
	defer cluster.Cleanup()

	cluster.Start()
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	for _, core := range cluster.Cores {
		checkLinkStatus(core.Client, "", t)
	}
}
