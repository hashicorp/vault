// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/sdk/logical"

	"nhooyr.io/websocket"

	"go.uber.org/atomic"

	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"

	vaulthttp "github.com/hashicorp/vault/http"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

// Avoiding a circular dependency in the test.
type mockSink struct {
	token *atomic.String
}

func (m *mockSink) Token() string {
	return m.token.Load()
}

func (m *mockSink) WriteToken(token string) error {
	m.token.Store(token)
	return nil
}

func newMockSink(t *testing.T) sink.Sink {
	t.Helper()

	return &mockSink{
		token: atomic.NewString(""),
	}
}

// testNewStaticSecretCacheUpdater returns a new StaticSecretCacheUpdater
// for use in tests.
func testNewStaticSecretCacheUpdater(t *testing.T, client *api.Client) *StaticSecretCacheUpdater {
	t.Helper()

	lc := testNewLeaseCache(t, []*SendResponse{})
	tokenSink := newMockSink(t)
	tokenSink.WriteToken(client.Token())

	updater, err := NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
		TokenSink:  tokenSink,
	})
	if err != nil {
		t.Fatal(err)
	}
	return updater
}

// TestNewStaticSecretCacheUpdater tests the NewStaticSecretCacheUpdater method,
// to ensure it errors out when appropriate.
func TestNewStaticSecretCacheUpdater(t *testing.T) {
	t.Parallel()

	lc := testNewLeaseCache(t, []*SendResponse{})
	config := api.DefaultConfig()
	logger := logging.NewVaultLogger(hclog.Trace).Named("cache.updater")
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	tokenSink := newMockSink(t)

	// Expect an error if any of the arguments are nil:
	updater, err := NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     nil,
		LeaseCache: lc,
		Logger:     logger,
		TokenSink:  tokenSink,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: nil,
		Logger:     logger,
		TokenSink:  tokenSink,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     nil,
		TokenSink:  tokenSink,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
		TokenSink:  nil,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	// Don't expect an error if the arguments are as expected
	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
		TokenSink:  tokenSink,
	})
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, updater)
}

// TestOpenWebSocketConnection tests that the openWebSocketConnection function
// works as expected. This uses a TLS enabled (wss) WebSocket connection.
func TestOpenWebSocketConnection(t *testing.T) {
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)
	updater.tokenSink.WriteToken(client.Token())

	conn, err := updater.openWebSocketConnection(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, conn)
}

// TestOpenWebSocketConnectionReceivesEvents tests that the openWebSocketConnection function
// works as expected with KVV1, and then the connection can be used to receive an event.
// This acts as more of an event system sanity check than a test of the updater
// logic. It's still important coverage, though.
func TestOpenWebSocketConnectionReceivesEvents(t *testing.T) {
	t.Parallel()
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

	t.Cleanup(func() {
		conn.Close(websocket.StatusNormalClosure, "")
	})

	makeData := func(i int) map[string]interface{} {
		return map[string]interface{}{
			"foo": fmt.Sprintf("bar%d", i),
		}
	}
	// Put a secret, which should trigger an event
	err = client.KVv1("secret").Put(context.Background(), "foo", makeData(100))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		// Do a fresh PUT just to refresh the secret and send a new message
		err = client.KVv1("secret").Put(context.Background(), "foo", makeData(i))
		if err != nil {
			t.Fatal(err)
		}

		// This method blocks until it gets a secret, so this test
		// will only pass if we're receiving events correctly.
		_, message, err := conn.Read(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(message))
	}
}

// TestOpenWebSocketConnectionReceivesEvents tests that the openWebSocketConnection function
// works as expected with KVV2, and then the connection can be used to receive an event.
// This acts as more of an event system sanity check than a test of the updater
// logic. It's still important coverage, though.
func TestOpenWebSocketConnectionReceivesEventsKVV2(t *testing.T) {
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.VersionedKVFactory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, conn)

	t.Cleanup(func() {
		conn.Close(websocket.StatusNormalClosure, "")
	})

	makeData := func(i int) map[string]interface{} {
		return map[string]interface{}{
			"foo": fmt.Sprintf("bar%d", i),
		}
	}

	err = client.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Put a secret, which should trigger an event
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", makeData(100))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		// Do a fresh PUT just to refresh the secret and send a new message
		_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", makeData(i))
		if err != nil {
			t.Fatal(err)
		}

		// This method blocks until it gets a secret, so this test
		// will only pass if we're receiving events correctly.
		_, _, err := conn.Read(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestOpenWebSocketConnectionTestServer tests that the openWebSocketConnection function
// works as expected using vaulthttp.TestServer. This server isn't TLS enabled, so tests
// the ws path (as opposed to the wss) path.
func TestOpenWebSocketConnectionTestServer(t *testing.T) {
	t.Parallel()
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

// Test_StreamStaticSecretEvents_UpdatesCacheWithNewSecrets tests that an event will
// properly update the corresponding secret in Proxy's cache. This is a little more end-to-end-y
// than TestUpdateStaticSecret, and essentially is testing a similar thing, though is
// ensuring that updateStaticSecret gets called by the event arriving
// (as part of streamStaticSecretEvents) instead of testing calling it explicitly.
func Test_StreamStaticSecretEvents_UpdatesCacheWithNewSecrets(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.VersionedKVFactory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)
	leaseCache := updater.leaseCache

	wg := &sync.WaitGroup{}
	runStreamStaticSecretEvents := func() {
		wg.Add(1)
		err := updater.streamStaticSecretEvents(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}
	go runStreamStaticSecretEvents()

	// TODO: Use the new function to make the ID from the path
	// First, create the secret in the cache that we expect to be updated:
	path := "secret-v2/data/foo"
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: path,
			},
		},
	}
	indexId := computeStaticSecretCacheIndex(req)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: path,
		LastRenewed: initialTime,
		ID:          indexId,
		// Valid token provided, so update should work.
		Tokens:   []string{client.Token()},
		Response: []byte{},
	}
	err := leaseCache.db.Set(index)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	err = client.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Put a secret, which should trigger an event
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the event to arrive. Events are usually much, much faster
	// than this, but we make it five seconds to protect against CI flakiness.
	time.Sleep(5 * time.Second)

	// Then, do a GET to see if the event got updated
	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, newIndex)
	require.NotEqual(t, []byte{}, newIndex.Response)
	require.Truef(t, initialTime.Before(newIndex.LastRenewed), "last updated time not updated on index")
	require.Equal(t, index.RequestPath, newIndex.RequestPath)
	require.Equal(t, index.Tokens, newIndex.Tokens)

	wg.Done()
}

// TestUpdateStaticSecret tests that updateStaticSecret works as expected, reaching out
// to Vault to get an updated secret when called.
func TestUpdateStaticSecret(t *testing.T) {
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)
	leaseCache := updater.leaseCache

	// TODO: avoid using req here make a new method
	// that takes path and returns index
	path := "secret/foo"
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: path,
			},
		},
	}
	indexId := computeStaticSecretCacheIndex(req)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: "secret/foo",
		LastRenewed: initialTime,
		ID:          indexId,
		// Valid token provided, so update should work.
		Tokens:   []string{client.Token()},
		Response: []byte{},
	}
	err := leaseCache.db.Set(index)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	// create the secret in Vault. n.b. the test cluster has already mounted the KVv1 backend at "secret"
	err = client.KVv1("secret").Put(context.Background(), "foo", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// attempt the update
	err = updater.updateStaticSecret(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}

	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, newIndex)
	require.Truef(t, initialTime.Before(newIndex.LastRenewed), "last updated time not updated on index")
	require.NotEqual(t, []byte{}, newIndex.Response)
	require.Equal(t, index.RequestPath, newIndex.RequestPath)
	require.Equal(t, index.Tokens, newIndex.Tokens)
}

// TestUpdateStaticSecret_EvictsIfInvalidTokens tests that updateStaticSecret will
// evict secrets from the cache if no valid tokens are left.
func TestUpdateStaticSecret_EvictsIfInvalidTokens(t *testing.T) {
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)
	leaseCache := updater.leaseCache

	// TODO: avoid using req here make a new method
	// that takes path and returns index
	path := "secret/foo"
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: path,
			},
		},
	}
	indexId := computeStaticSecretCacheIndex(req)
	renewTime := time.Now().UTC()

	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: "secret/foo",
		LastRenewed: renewTime,
		ID:          indexId,
		// Note: invalid Tokens value provided, so this secret cannot be updated, and must be evicted
		Tokens: []string{"invalid token"},
	}
	err := leaseCache.db.Set(index)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	// create the secret in Vault. n.b. the test cluster has already mounted the KVv1 backend at "secret"
	err = client.KVv1("secret").Put(context.Background(), "foo", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// attempt the update
	err = updater.updateStaticSecret(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}

	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	if err != nil {
		t.Fatal(err)
	}

	require.NotEqual(t, index, newIndex)
	require.Nil(t, newIndex)
}

// TestUpdateStaticSecret_HandlesNonCachedPaths tests that updateStaticSecret
// doesn't fail or error if we try and give it an update to a path that isn't cached.
func TestUpdateStaticSecret_HandlesNonCachedPaths(t *testing.T) {
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	path := "secret/foo"

	// attempt the update
	err := updater.updateStaticSecret(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	require.Nil(t, err)
}
