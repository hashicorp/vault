// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package seal_binary

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/stretchr/testify/assert"
)

func pkcsWrapper(string, int) testcluster.VaultNodeSealConfig {
	return testcluster.VaultNodeSealConfig{}
}

func getRewrappedEntryCount(client *api.Client) (uint32, error) {
	return 0, nil
}

func verifyRewrappedEntryCount(t *assert.CollectT, client *api.Client, initialProcessedEntries uint32) uint32 {
	return 0
}
