// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"testing"

	vaulthttp "github.com/hashicorp/vault/http"

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
// works as expected. This uses a TLS enabled (wss) WebSocket connection.
func TestOpenWebSocketConnection(t *testing.T) {
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, conn)
}

// TestOpenWebSocketConnectionTestServer tests that the openWebSocketConnection function
// works as expected using vaulthttp.TestServer. This server isn't TLS enabled, so tests
// the ws path (as opposed to the wss) path.
func TestOpenWebSocketConnectionTestServer(t *testing.T) {
	// We need a valid cluster for the connection to succeed.
	core := vault.TestCoreWithConfig(t, &vault.CoreConfig{})
	ln, addr := vaulthttp.TestServer(t, core)
	defer ln.Close()

	keys, rootToken := vault.TestCoreInit(t, core)
	for _, key := range keys {
		_, err := core.Unseal(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)
	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, conn)
}
