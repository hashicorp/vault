// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

// testNewStaticSecretCacheUpdater returns a new StaticSecretCacheUpdater
// for use in tests.
func testNewStaticSecretCacheUpdater(t *testing.T, client *api.Client) *StaticSecretCacheUpdater {
	t.Helper()

	lc := testNewLeaseCache(t, []*SendResponse{})

	updater, err := NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return updater
}

// TestOpenWebSocketConnection tests that the openWebSocketConnection function
// works as expected.
func TestOpenWebSocketConnection(t *testing.T) {
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, nil)
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, conn)
}
