// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"fmt"
	"os"
	"sync"
	syncatomic "sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
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
	require.NoError(t, err)
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
	require.NoError(t, err)
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
	require.NoError(t, err)
	require.NotNil(t, updater)
}

// TestOpenWebSocketConnection tests that the openWebSocketConnection function
// works as expected (fails on CE, succeeds on ent).
// This uses a TLS enabled (wss) WebSocket connection.
func TestOpenWebSocketConnection(t *testing.T) {
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)
	updater.tokenSink.WriteToken(client.Token())

	conn, err := updater.openWebSocketConnection(context.Background())
	if constants.IsEnterprise {
		require.NoError(t, err)
		require.NotNil(t, conn)
	} else {
		require.Nil(t, conn)
		require.Errorf(t, err, "ensure Vault is Enterprise version 1.16 or above")
	}
}

// TestOpenWebSocketConnection_BadPolicyToken tests attempting to open a websocket
// connection to the events system using a token that has incorrect policy access
// will not trigger auto auth
func TestOpenWebSocketConnection_BadPolicyToken(t *testing.T) {
	// We need a valid cluster for the connection to succeed.
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	eventPolicy := `path "sys/events/subscribe/*" {
		capabilities = ["deny"]
	}`
	client.Sys().PutPolicy("no_events_access", eventPolicy)

	// Create a new token with a bad policy
	token, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"no_events_access"},
	})
	require.NoError(t, err)

	// Set the client token to one with an invalid policy
	updater.tokenSink.WriteToken(token.Auth.ClientToken)
	client.SetToken(token.Auth.ClientToken)

	ctx, cancelFunc := context.WithCancel(context.Background())

	authInProgress := &syncatomic.Bool{}
	renewalChannel := make(chan error)
	errCh := make(chan error)
	go func() {
		errCh <- updater.Run(ctx, authInProgress, renewalChannel)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			require.NoError(t, err)
		}
	}()

	defer cancelFunc()

	// Verify that the token has been written to the sink before checking auto auth
	// is not re-triggered
	err = updater.streamStaticSecretEvents(ctx)
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())

	// Auto auth should not be retriggered
	timeout := time.After(2 * time.Second)
	select {
	case <-renewalChannel:
		t.Fatal("incorrectly triggered auto auth")
	case <-ctx.Done():
		t.Fatal("context was closed before auto auth could be re-triggered")
	case <-timeout:
	}
}

// TestOpenWebSocketConnection_AutoAuthSelfHeal tests attempting to open a websocket
// connection to the events system using an invalid token will re-trigger
// auto auth.
func TestOpenWebSocketConnection_AutoAuthSelfHeal(t *testing.T) {
	// We need a valid cluster for the connection to succeed.
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	// Revoke the token before it can be used to open a connection to the events system
	client.Auth().Token().RevokeOrphan(client.Token())
	updater.tokenSink.WriteToken(client.Token())
	time.Sleep(100 * time.Millisecond)

	ctx, cancelFunc := context.WithCancel(context.Background())

	authInProgress := &syncatomic.Bool{}
	renewalChannel := make(chan error)
	errCh := make(chan error)
	go func() {
		errCh <- updater.Run(ctx, authInProgress, renewalChannel)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			require.NoError(t, err)
		}
	}()

	defer cancelFunc()

	// Wait for static secret updater to begin
	timeout := time.After(10 * time.Second)

	select {
	case <-renewalChannel:
	case <-ctx.Done():
		t.Fatal("context was closed before auto auth could be re-triggered")
	case <-timeout:
		t.Fatal("timed out before auto auth could be re-triggered")
	}
	authInProgress.Store(false)

	// Verify that auto auth is re-triggered again because another auth is "not in progress"
	timeout = time.After(15 * time.Second)
	select {
	case <-renewalChannel:
	case <-ctx.Done():
		t.Fatal("context was closed before auto auth could be re-triggered")
	case <-timeout:
		t.Fatal("timed out before auto auth could be re-triggered")
	}
	authInProgress.Store(true)

	// Verify that auto auth is NOT re-triggered again because another auth is in progress
	timeout = time.After(2 * time.Second)
	select {
	case <-renewalChannel:
		t.Fatal("auto auth was incorrectly re-triggered")
	case <-ctx.Done():
		t.Fatal("context was closed before auto auth could be re-triggered")
	case <-timeout:
	}
}

// TestOpenWebSocketConnectionReceivesEventsDefaultMount tests that the openWebSocketConnection function
// works as expected with the default KVV1 mount, and then the connection can be used to receive an event.
// This acts as more of an event system sanity check than a test of the updater
// logic. It's still important coverage, though.
// It also adds a client timeout of 1 second and checks that the connection does not timeout as this is a
// streaming request.
func TestOpenWebSocketConnectionReceivesEventsDefaultMount(t *testing.T) {
	if !constants.IsEnterprise {
		t.Skip("test can only run on enterprise due to requiring the event notification system")
	}
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	oldClientTimeout := os.Getenv("VAULT_CLIENT_TIMEOUT")
	os.Setenv("VAULT_CLIENT_TIMEOUT", "1")
	defer os.Setenv("VAULT_CLIENT_TIMEOUT", oldClientTimeout)

	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	require.NoError(t, err)
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
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		// Do a fresh PUT just to refresh the secret and send a new message
		err = client.KVv1("secret").Put(context.Background(), "foo", makeData(i))
		require.NoError(t, err)

		// This method blocks until it gets a secret, so this test
		// will only pass if we're receiving events correctly.
		// It will fail here if the connection times out.
		_, _, err = conn.Read(context.Background())
		require.NoError(t, err)
	}
}

// TestOpenWebSocketConnectionReceivesEventsKVV1 tests that the openWebSocketConnection function
// works as expected with KVV1, and then the connection can be used to receive an event.
// This acts as more of an event system sanity check than a test of the updater
// logic. It's still important coverage, though.
func TestOpenWebSocketConnectionReceivesEventsKVV1(t *testing.T) {
	if !constants.IsEnterprise {
		t.Skip("test can only run on enterprise due to requiring the event notification system")
	}
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	require.NoError(t, err)
	require.NotNil(t, conn)

	t.Cleanup(func() {
		conn.Close(websocket.StatusNormalClosure, "")
	})

	err = client.Sys().Mount("secret-v1", &api.MountInput{
		Type: "kv",
	})
	require.NoError(t, err)

	makeData := func(i int) map[string]interface{} {
		return map[string]interface{}{
			"foo": fmt.Sprintf("bar%d", i),
		}
	}
	// Put a secret, which should trigger an event
	err = client.KVv1("secret-v1").Put(context.Background(), "foo", makeData(100))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		// Do a fresh PUT just to refresh the secret and send a new message
		err = client.KVv1("secret-v1").Put(context.Background(), "foo", makeData(i))
		require.NoError(t, err)

		// This method blocks until it gets a secret, so this test
		// will only pass if we're receiving events correctly.
		_, _, err := conn.Read(context.Background())
		require.NoError(t, err)
	}
}

// TestOpenWebSocketConnectionReceivesEventsKVV2 tests that the openWebSocketConnection function
// works as expected with KVV2, and then the connection can be used to receive an event.
// This acts as more of an event system sanity check than a test of the updater
// logic. It's still important coverage, though.
func TestOpenWebSocketConnectionReceivesEventsKVV2(t *testing.T) {
	if !constants.IsEnterprise {
		t.Skip("test can only run on enterprise due to requiring the event notification system")
	}
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
	require.NoError(t, err)
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
	require.NoError(t, err)

	// Put a secret, which should trigger an event
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", makeData(100))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		// Do a fresh PUT just to refresh the secret and send a new message
		_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", makeData(i))
		require.NoError(t, err)

		// This method blocks until it gets a secret, so this test
		// will only pass if we're receiving events correctly.
		_, _, err := conn.Read(context.Background())
		require.NoError(t, err)
	}
}

// TestOpenWebSocketConnectionTestServer tests that the openWebSocketConnection function
// works as expected using vaulthttp.TestServer. This server isn't TLS enabled, so tests
// the ws path (as opposed to the wss) path.
func TestOpenWebSocketConnectionTestServer(t *testing.T) {
	if !constants.IsEnterprise {
		t.Skip("test can only run on enterprise due to requiring the event notification system")
	}
	t.Parallel()
	// We need a valid cluster for the connection to succeed.
	core := vault.TestCoreWithConfig(t, &vault.CoreConfig{})
	ln, addr := vaulthttp.TestServer(t, core)
	defer ln.Close()

	keys, rootToken := vault.TestCoreInit(t, core)
	for _, key := range keys {
		_, err := core.Unseal(key)
		require.NoError(t, err)
	}

	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	require.NoError(t, err)
	client.SetToken(rootToken)
	updater := testNewStaticSecretCacheUpdater(t, client)

	conn, err := updater.openWebSocketConnection(context.Background())
	require.NoError(t, err)
	require.NotNil(t, conn)
}

// Test_StreamStaticSecretEvents_UpdatesCacheWithNewSecrets tests that an event will
// properly update the corresponding secret in Proxy's cache. This is a little more end-to-end-y
// than TestUpdateStaticSecret, and essentially is testing a similar thing, though is
// ensuring that updateStaticSecret gets called by the event arriving
// (as part of streamStaticSecretEvents) instead of testing calling it explicitly.
func Test_StreamStaticSecretEvents_UpdatesCacheWithNewSecrets(t *testing.T) {
	if !constants.IsEnterprise {
		t.Skip("test can only run on enterprise due to requiring the event notification system")
	}
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
		require.NoError(t, err)
	}
	go runStreamStaticSecretEvents()

	// First, create the secret in the cache that we expect to be updated:
	path := "secret-v2/data/foo"
	indexId := hashStaticSecretIndex(path)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: path,
		Versions:    map[int][]byte{},
		LastRenewed: initialTime,
		ID:          indexId,
		// Valid token provided, so update should work.
		Tokens:   map[string]struct{}{client.Token(): {}},
		Response: []byte{},
	}
	err := leaseCache.db.Set(index)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	err = client.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// Wait for the event stream to be fully up and running. Should be faster than this in reality, but
	// we make it five seconds to protect against CI flakiness.
	time.Sleep(5 * time.Second)

	// Put a secret, which should trigger an event
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", secretData)
	require.NoError(t, err)

	// Wait for the event to arrive. Events are usually much, much faster
	// than this, but we make it five seconds to protect against CI flakiness.
	time.Sleep(5 * time.Second)

	// Then, do a GET to see if the index got updated by the event
	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	require.NoError(t, err)
	require.NotNil(t, newIndex)
	require.NotEqual(t, []byte{}, newIndex.Response)
	require.Truef(t, initialTime.Before(newIndex.LastRenewed), "last updated time not updated on index")
	require.Equal(t, index.RequestPath, newIndex.RequestPath)
	require.Equal(t, index.Tokens, newIndex.Tokens)

	// Assert that the corresponding version got updated too
	require.Len(t, newIndex.Versions, 1)
	require.NotNil(t, newIndex.Versions)
	require.NotNil(t, newIndex.Versions[1])
	require.Equal(t, newIndex.Versions[1], newIndex.Response)

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

	path := "secret/foo"
	indexId := hashStaticSecretIndex(path)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: "secret/foo",
		LastRenewed: initialTime,
		ID:          indexId,
		Versions:    map[int][]byte{},
		// Valid token provided, so update should work.
		Tokens:   map[string]struct{}{client.Token(): {}},
		Response: []byte{},
	}
	err := leaseCache.db.Set(index)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	// create the secret in Vault. n.b. the test cluster has already mounted the KVv1 backend at "secret"
	err = client.KVv1("secret").Put(context.Background(), "foo", secretData)
	require.NoError(t, err)

	// attempt the update
	err = updater.updateStaticSecret(context.Background(), path)
	require.NoError(t, err)

	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	require.NoError(t, err)
	require.NotNil(t, newIndex)
	require.Truef(t, initialTime.Before(newIndex.LastRenewed), "last updated time not updated on index")
	require.NotEqual(t, []byte{}, newIndex.Response)
	require.Equal(t, index.RequestPath, newIndex.RequestPath)
	require.Equal(t, index.Tokens, newIndex.Tokens)
	require.Len(t, newIndex.Versions, 0)
}

// TestUpdateStaticSecret_KVv2 tests that updateStaticSecret works as expected, reaching out
// to Vault to get an updated secret when called. It should also update the corresponding
// version of that secret in the cache index's Versions field.
func TestUpdateStaticSecret_KVv2(t *testing.T) {
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
	leaseCache := updater.leaseCache

	path := "secret-v2/data/foo"
	indexId := hashStaticSecretIndex(path)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: path,
		LastRenewed: initialTime,
		ID:          indexId,
		Versions:    map[int][]byte{},
		// Valid token provided, so update should work.
		Tokens:   map[string]struct{}{client.Token(): {}},
		Response: []byte{},
	}
	err := leaseCache.db.Set(index)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	err = client.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// create the secret in Vault
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", secretData)
	require.NoError(t, err)

	// attempt the update
	err = updater.updateStaticSecret(context.Background(), path)
	require.NoError(t, err)

	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	require.NoError(t, err)
	require.NotNil(t, newIndex)
	require.Truef(t, initialTime.Before(newIndex.LastRenewed), "last updated time not updated on index")
	require.NotEqual(t, []byte{}, newIndex.Response)
	require.Equal(t, index.RequestPath, newIndex.RequestPath)
	require.Equal(t, index.Tokens, newIndex.Tokens)

	// It should have also updated version 1 with the same version.
	require.Len(t, newIndex.Versions, 1)
	require.NotNil(t, newIndex.Versions[1])
	require.Equal(t, newIndex.Versions[1], newIndex.Response)
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

	path := "secret/foo"
	indexId := hashStaticSecretIndex(path)
	renewTime := time.Now().UTC()

	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: "secret/foo",
		LastRenewed: renewTime,
		ID:          indexId,
		// Note: invalid Tokens value provided, so this secret cannot be updated, and must be evicted
		Tokens: map[string]struct{}{"invalid token": {}},
	}
	err := leaseCache.db.Set(index)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	// create the secret in Vault. n.b. the test cluster has already mounted the KVv1 backend at "secret"
	err = client.KVv1("secret").Put(context.Background(), "foo", secretData)
	require.NoError(t, err)

	// attempt the update
	err = updater.updateStaticSecret(context.Background(), path)
	require.NoError(t, err)

	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	require.Equal(t, cachememdb.ErrCacheItemNotFound, err)
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

	// Attempt the update
	err := updater.updateStaticSecret(context.Background(), path)
	require.NoError(t, err)
	require.Nil(t, err)
}

// TestPreEventStreamUpdate tests that preEventStreamUpdate correctly
// updates old static secrets in the cache.
func TestPreEventStreamUpdate(t *testing.T) {
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

	// First, create the secret in the cache that we expect to be updated:
	path := "secret-v2/data/foo"
	indexId := hashStaticSecretIndex(path)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: path,
		LastRenewed: initialTime,
		ID:          indexId,
		Versions:    map[int][]byte{},
		// Valid token provided, so update should work.
		Tokens:   map[string]struct{}{client.Token(): {}},
		Response: []byte{},
		Type:     cacheboltdb.StaticSecretType,
	}
	err := leaseCache.db.Set(index)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	err = client.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// Put a secret (with different values to what's currently in the cache)
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", secretData)
	require.NoError(t, err)

	// perform the pre-event stream update:
	err = updater.preEventStreamUpdate(context.Background())
	require.Nil(t, err)

	// Then, do a GET to see if the event got updated
	newIndex, err := leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	require.Nil(t, err)
	require.NotNil(t, newIndex)
	require.NotEqual(t, []byte{}, newIndex.Response)
	require.Truef(t, initialTime.Before(newIndex.LastRenewed), "last updated time not updated on index")
	require.Equal(t, index.RequestPath, newIndex.RequestPath)
	require.Equal(t, index.Tokens, newIndex.Tokens)
	require.Equal(t, index.Versions, newIndex.Versions)
}

// TestPreEventStreamUpdateErrorUpdating tests that preEventStreamUpdate correctly responds
// to errors on secret updates
func TestPreEventStreamUpdateErrorUpdating(t *testing.T) {
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

	// First, create the secret in the cache that we expect to be updated:
	path := "secret-v2/data/foo"
	indexId := hashStaticSecretIndex(path)
	initialTime := time.Now().UTC()
	// pre-populate the leaseCache with a secret to update
	index := &cachememdb.Index{
		Namespace:   "root/",
		RequestPath: path,
		LastRenewed: initialTime,
		ID:          indexId,
		// Valid token provided, so update should work.
		Tokens:   map[string]struct{}{client.Token(): {}},
		Response: []byte{},
		Type:     cacheboltdb.StaticSecretType,
	}
	err := leaseCache.db.Set(index)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	err = client.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// Put a secret (with different values to what's currently in the cache)
	_, err = client.KVv2("secret-v2").Put(context.Background(), "foo", secretData)
	require.NoError(t, err)

	// Seal Vault, so that the update will fail
	cluster.EnsureCoresSealed(t)

	// perform the pre-event stream update:
	err = updater.preEventStreamUpdate(context.Background())
	require.Nil(t, err)

	// Then, we expect the index to be evicted since the token failed to update
	_, err = leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	require.Equal(t, cachememdb.ErrCacheItemNotFound, err)
}

// TestCheckForDeleteOrDestroyEvent tests the behaviour of checkForDeleteOrDestroyEvent
// and assures it gives the right responses for different events.
func TestCheckForDeleteOrDestroyEvent(t *testing.T) {
	t.Parallel()

	expectedVersions := []int{1, 3, 5}
	jsonFormatExpectedVersions := "[1,3,5]"
	expectedPath := "secret-v2/data/my-secret"
	deletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":  "2",
					"deleted_versions": jsonFormatExpectedVersions,
					"modified":         true,
					"operation":        "delete",
					"path":             "secret-v2/delete/my-secret",
				},
			},
			"event_type": "kv-v2/delete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	undeletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"undeleted_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "undelete",
					"path":               "secret-v2/undelete/my-secret",
				},
			},
			"event_type": "kv-v2/undelete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	destroyedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"destroyed_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "destroy",
					"path":               "secret-v2/destroy/my-secret",
				},
			},
			"event_type": "kv-v2/destroy",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	actualVersions, actualPath := checkForDeleteOrDestroyEvent(deletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(undeletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(destroyedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)
}

// TestCheckForDeleteOrDestroyNamespacedEvent tests the behaviour of checkForDeleteOrDestroyEvent
// with namespaces in paths.
func TestCheckForDeleteOrDestroyNamespacedEvent(t *testing.T) {
	t.Parallel()

	expectedVersions := []int{1, 3, 5}
	jsonFormatExpectedVersions := "[1,3,5]"
	expectedPath := "ns/secret-v2/data/my-secret"
	deletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":  "2",
					"deleted_versions": jsonFormatExpectedVersions,
					"modified":         true,
					"operation":        "delete",
					"data_path":        "secret-v2/data/my-secret",
					"path":             "secret-v2/delete/my-secret",
				},
			},
			"namespace":  "ns/",
			"event_type": "kv-v2/delete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	undeletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"undeleted_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "undelete",
					"data_path":          "secret-v2/data/my-secret",
					"path":               "secret-v2/undelete/my-secret",
				},
			},
			"namespace":  "ns/",
			"event_type": "kv-v2/undelete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	destroyedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"destroyed_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "destroy",
					"data_path":          "secret-v2/data/my-secret",
					"path":               "secret-v2/destroy/my-secret",
				},
			},
			"namespace":  "ns/",
			"event_type": "kv-v2/destroy",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	actualVersions, actualPath := checkForDeleteOrDestroyEvent(deletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(undeletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(destroyedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)
}
