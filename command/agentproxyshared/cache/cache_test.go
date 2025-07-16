// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/mock"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func tokenRevocationValidation(t *testing.T, sampleSpace map[string]string, expected map[string]string, leaseCache *LeaseCache) {
	t.Helper()
	for val, valType := range sampleSpace {
		index, err := leaseCache.db.Get(valType, val)
		if err != nil && err != cachememdb.ErrCacheItemNotFound {
			t.Fatal(err)
		}
		if expected[val] == "" && index != nil {
			t.Fatalf("failed to evict index from the cache: type: %q, value: %q", valType, val)
		}
		if expected[val] != "" && index == nil {
			t.Fatalf("evicted an undesired index from cache: type: %q, value: %q", valType, val)
		}
	}
}

func TestCache_AutoAuthTokenStripping(t *testing.T) {
	response1 := `{"data": {"id": "testid", "accessor": "testaccessor", "request": "lookup-self"}}`
	response2 := `{"data": {"id": "testid", "accessor": "testaccessor", "request": "lookup"}}`
	response3 := `{"auth": {"client_token": "testid", "accessor": "testaccessor"}}`
	response4 := `{"auth": {"client_token": "testid", "accessor": "testaccessor"}}`
	responses := []*SendResponse{
		newTestSendResponse(http.StatusOK, response1),
		newTestSendResponse(http.StatusOK, response2),
		newTestSendResponse(http.StatusOK, response3),
		newTestSendResponse(http.StatusOK, response4),
	}

	leaseCache := testNewLeaseCache(t, responses)

	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	cacheLogger := logging.NewVaultLogger(hclog.Trace).Named("cache")
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	ctx := namespace.RootContext(nil)

	// Create a muxer and add paths relevant for the lease cache layer
	mux := http.NewServeMux()
	mux.Handle(consts.AgentPathCacheClear, leaseCache.HandleCacheClear(ctx))

	mux.Handle("/", ProxyHandler(ctx, cacheLogger, leaseCache, mock.NewSink("testid"), false, true, nil, nil))
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          cacheLogger.StandardLogger(nil),
	}
	go server.Serve(listener)

	testClient, err := client.Clone()
	if err != nil {
		t.Fatal(err)
	}

	if err := testClient.SetAddress("http://" + listener.Addr().String()); err != nil {
		t.Fatal(err)
	}

	// Empty the token in the client. Auto-auth token should be put to use.
	testClient.SetToken("")
	secret, err := testClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["id"] != nil || secret.Data["accessor"] != nil || secret.Data["request"].(string) != "lookup-self" {
		t.Fatalf("failed to strip off auto-auth token on lookup-self")
	}

	secret, err = testClient.Auth().Token().Lookup("")
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["id"] != nil || secret.Data["accessor"] != nil || secret.Data["request"].(string) != "lookup" {
		t.Fatalf("failed to strip off auto-auth token on lookup")
	}

	secret, err = testClient.Auth().Token().RenewSelf(1)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth == nil {
		secretJson, _ := json.Marshal(secret)
		t.Fatalf("Expected secret to have Auth but was %s", secretJson)
	}
	if secret.Auth.ClientToken != "" || secret.Auth.Accessor != "" {
		t.Fatalf("failed to strip off auto-auth token on renew-self")
	}

	secret, err = testClient.Auth().Token().Renew("testid", 1)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth == nil {
		secretJson, _ := json.Marshal(secret)
		t.Fatalf("Expected secret to have Auth but was %s", secretJson)
	}
	if secret.Auth.ClientToken != "" || secret.Auth.Accessor != "" {
		t.Fatalf("failed to strip off auto-auth token on renew")
	}
}

func TestCache_AutoAuthClientTokenProxyStripping(t *testing.T) {
	leaseCache := &mockTokenVerifierProxier{}
	dummyToken := "DUMMY"
	realToken := "testid"

	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	cacheLogger := logging.NewVaultLogger(hclog.Trace).Named("cache")
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	ctx := namespace.RootContext(nil)

	// Create a muxer and add paths relevant for the lease cache layer
	mux := http.NewServeMux()

	mux.Handle("/", ProxyHandler(ctx, cacheLogger, leaseCache, mock.NewSink(realToken), true, true, nil, nil))
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          cacheLogger.StandardLogger(nil),
	}
	go server.Serve(listener)

	testClient, err := client.Clone()
	if err != nil {
		t.Fatal(err)
	}

	if err := testClient.SetAddress("http://" + listener.Addr().String()); err != nil {
		t.Fatal(err)
	}

	// Empty the token in the client. Auto-auth token should be put to use.
	testClient.SetToken(dummyToken)
	_, err = testClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if leaseCache.currentToken != realToken {
		t.Fatalf("failed to use real token from auto-auth")
	}
}

func TestCache_ConcurrentRequests(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, _, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("kv/foo/%d_%d", i, rand.Int())
			_, err := testClient.Logical().Write(key, map[string]interface{}{
				"key": key,
			})
			if err != nil {
				t.Error(err.Error())
				return
			}
			secret, err := testClient.Logical().Read(key)
			if err != nil {
				t.Error(err.Error())
				return
			}
			if secret == nil || secret.Data["key"].(string) != key {
				t.Errorf("failed to read value for key: %q", key)
				return
			}
		}(i)
	}
	wg.Wait()
}

func TestCache_TokenRevocations_RevokeOrphan(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	sampleSpace := make(map[string]string)

	cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()
	sampleSpace[token1] = "token"

	// Mount the kv backend
	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read the secret and create a lease
	leaseResp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease1 := leaseResp.LeaseID
	sampleSpace[lease1] = "lease"

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken
	sampleSpace[token2] = "token"

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease2 := leaseResp.LeaseID
	sampleSpace[lease2] = "lease"

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken
	sampleSpace[token3] = "token"

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease3 := leaseResp.LeaseID
	sampleSpace[lease3] = "lease"

	expected := make(map[string]string)
	for k, v := range sampleSpace {
		expected[k] = v
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)

	// Revoke-orphan the intermediate token. This should result in its own
	// eviction and evictions of the revoked token's leases. All other things
	// including the child tokens and leases of the child tokens should be
	// untouched.
	testClient.SetToken(token2)
	err = testClient.Auth().Token().RevokeOrphan(token2)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	expected = map[string]string{
		token1: "token",
		lease1: "lease",
		token3: "token",
		lease3: "lease",
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)
}

func TestCache_TokenRevocations_LeafLevelToken(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	sampleSpace := make(map[string]string)

	cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()
	sampleSpace[token1] = "token"

	// Mount the kv backend
	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read the secret and create a lease
	leaseResp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease1 := leaseResp.LeaseID
	sampleSpace[lease1] = "lease"

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken
	sampleSpace[token2] = "token"

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease2 := leaseResp.LeaseID
	sampleSpace[lease2] = "lease"

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken
	sampleSpace[token3] = "token"

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease3 := leaseResp.LeaseID
	sampleSpace[lease3] = "lease"

	expected := make(map[string]string)
	for k, v := range sampleSpace {
		expected[k] = v
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)

	// Revoke the lef token. This should evict all the leases belonging to this
	// token, evict entries for all the child tokens and their respective
	// leases.
	testClient.SetToken(token3)
	err = testClient.Auth().Token().RevokeSelf("")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	expected = map[string]string{
		token1: "token",
		lease1: "lease",
		token2: "token",
		lease2: "lease",
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)
}

func TestCache_TokenRevocations_IntermediateLevelToken(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	sampleSpace := make(map[string]string)

	cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()
	sampleSpace[token1] = "token"

	// Mount the kv backend
	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read the secret and create a lease
	leaseResp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease1 := leaseResp.LeaseID
	sampleSpace[lease1] = "lease"

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken
	sampleSpace[token2] = "token"

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease2 := leaseResp.LeaseID
	sampleSpace[lease2] = "lease"

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken
	sampleSpace[token3] = "token"

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease3 := leaseResp.LeaseID
	sampleSpace[lease3] = "lease"

	expected := make(map[string]string)
	for k, v := range sampleSpace {
		expected[k] = v
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)

	// Revoke the second level token. This should evict all the leases
	// belonging to this token, evict entries for all the child tokens and
	// their respective leases.
	testClient.SetToken(token2)
	err = testClient.Auth().Token().RevokeSelf("")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	expected = map[string]string{
		token1: "token",
		lease1: "lease",
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)
}

func TestCache_TokenRevocations_TopLevelToken(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	sampleSpace := make(map[string]string)

	cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()
	sampleSpace[token1] = "token"

	// Mount the kv backend
	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read the secret and create a lease
	leaseResp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease1 := leaseResp.LeaseID
	sampleSpace[lease1] = "lease"

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken
	sampleSpace[token2] = "token"

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease2 := leaseResp.LeaseID
	sampleSpace[lease2] = "lease"

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken
	sampleSpace[token3] = "token"

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease3 := leaseResp.LeaseID
	sampleSpace[lease3] = "lease"

	expected := make(map[string]string)
	for k, v := range sampleSpace {
		expected[k] = v
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)

	// Revoke the top level token. This should evict all the leases belonging
	// to this token, evict entries for all the child tokens and their
	// respective leases.
	testClient.SetToken(token1)
	err = testClient.Auth().Token().RevokeSelf("")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	expected = make(map[string]string)
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)
}

func TestCache_TokenRevocations_Shutdown(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	sampleSpace := make(map[string]string)

	ctx, rootCancelFunc := context.WithCancel(namespace.RootContext(nil))
	cleanup, _, testClient, leaseCache := setupClusterAndAgent(ctx, t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()
	sampleSpace[token1] = "token"

	// Mount the kv backend
	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read the secret and create a lease
	leaseResp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease1 := leaseResp.LeaseID
	sampleSpace[lease1] = "lease"

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken
	sampleSpace[token2] = "token"

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease2 := leaseResp.LeaseID
	sampleSpace[lease2] = "lease"

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken
	sampleSpace[token3] = "token"

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease3 := leaseResp.LeaseID
	sampleSpace[lease3] = "lease"

	expected := make(map[string]string)
	for k, v := range sampleSpace {
		expected[k] = v
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)

	rootCancelFunc()
	time.Sleep(1 * time.Second)

	// Ensure that all the entries are now gone
	expected = make(map[string]string)
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)
}

func TestCache_TokenRevocations_BaseContextCancellation(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	sampleSpace := make(map[string]string)

	cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()
	sampleSpace[token1] = "token"

	// Mount the kv backend
	err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read the secret and create a lease
	leaseResp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease1 := leaseResp.LeaseID
	sampleSpace[lease1] = "lease"

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken
	sampleSpace[token2] = "token"

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease2 := leaseResp.LeaseID
	sampleSpace[lease2] = "lease"

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken
	sampleSpace[token3] = "token"

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease3 := leaseResp.LeaseID
	sampleSpace[lease3] = "lease"

	expected := make(map[string]string)
	for k, v := range sampleSpace {
		expected[k] = v
	}
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)

	// Cancel the base context of the lease cache. This should trigger
	// evictions of all the entries from the cache.
	leaseCache.baseCtxInfo.CancelFunc()
	time.Sleep(1 * time.Second)

	// Ensure that all the entries are now gone
	expected = make(map[string]string)
	tokenRevocationValidation(t, sampleSpace, expected, leaseCache)
}

func TestCache_NonCacheable(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}

	cleanup, _, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	// Query mounts first
	origMounts, err := testClient.Sys().ListMounts()
	if err != nil {
		t.Fatal(err)
	}

	// Mount a kv backend
	if err := testClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}

	// Query mounts again
	newMounts, err := testClient.Sys().ListMounts()
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(origMounts, newMounts); diff == nil {
		t.Logf("response #1: %#v", origMounts)
		t.Logf("response #2: %#v", newMounts)
		t.Fatal("expected requests to be not cached")
	}

	// Query a non-existing mount, expect an error from api.Response
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	r := testClient.NewRequest("GET", "/v1/kv-invalid")

	apiResp, err := testClient.RawRequestWithContext(ctx, r)
	if apiResp != nil {
		defer apiResp.Body.Close()
	}
	if apiResp.Error() == nil || (apiResp != nil && apiResp.StatusCode != 404) {
		t.Fatalf("expected an error response and a 404 from requesting an invalid path, got: %#v", apiResp)
	}
	if err == nil {
		t.Fatal("expected an error from requesting an invalid path")
	}
}

func TestCache_Caching_AuthResponse(t *testing.T) {
	cleanup, _, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
	defer cleanup()

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token := resp.Auth.ClientToken
	testClient.SetToken(token)

	authTokeCreateReq := func(t *testing.T, policies map[string]interface{}) *api.Secret {
		resp, err := testClient.Logical().Write("auth/token/create", policies)
		if err != nil {
			t.Fatal(err)
		}
		if resp.Auth == nil || resp.Auth.ClientToken == "" {
			t.Fatalf("expected a valid client token in the response, got = %#v", resp)
		}

		return resp
	}

	// Test on auth response by creating a child token
	{
		proxiedResp := authTokeCreateReq(t, map[string]interface{}{
			"policies": "default",
		})

		cachedResp := authTokeCreateReq(t, map[string]interface{}{
			"policies": "default",
		})

		if diff := deep.Equal(proxiedResp.Auth.ClientToken, cachedResp.Auth.ClientToken); diff != nil {
			t.Fatal(diff)
		}
	}

	// Test on *non-renewable* auth response by creating a child root token
	{
		proxiedResp := authTokeCreateReq(t, nil)

		cachedResp := authTokeCreateReq(t, nil)

		if diff := deep.Equal(proxiedResp.Auth.ClientToken, cachedResp.Auth.ClientToken); diff != nil {
			t.Fatal(diff)
		}
	}
}

func TestCache_Caching_LeaseResponse(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, client, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	err := client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test proxy by issuing two different requests
	{
		// Write data to the lease-kv backend
		_, err := testClient.Logical().Write("kv/foo", map[string]interface{}{
			"value": "bar",
			"ttl":   "1h",
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = testClient.Logical().Write("kv/foobar", map[string]interface{}{
			"value": "bar",
			"ttl":   "1h",
		})
		if err != nil {
			t.Fatal(err)
		}

		firstResp, err := testClient.Logical().Read("kv/foo")
		if err != nil {
			t.Fatal(err)
		}

		secondResp, err := testClient.Logical().Read("kv/foobar")
		if err != nil {
			t.Fatal(err)
		}

		if diff := deep.Equal(firstResp, secondResp); diff == nil {
			t.Logf("response: %#v", firstResp)
			t.Fatal("expected proxied responses, got cached response on second request")
		}
	}

	// Test caching behavior by issue the same request twice
	{
		_, err := testClient.Logical().Write("kv/baz", map[string]interface{}{
			"value": "foo",
			"ttl":   "1h",
		})
		if err != nil {
			t.Fatal(err)
		}

		proxiedResp, err := testClient.Logical().Read("kv/baz")
		if err != nil {
			t.Fatal(err)
		}

		cachedResp, err := testClient.Logical().Read("kv/baz")
		if err != nil {
			t.Fatal(err)
		}

		if diff := deep.Equal(proxiedResp, cachedResp); diff != nil {
			t.Fatal(diff)
		}
	}
}

func TestCache_Caching_CacheClear(t *testing.T) {
	t.Run("request_path", func(t *testing.T) {
		testCachingCacheClearCommon(t, "request_path")
	})

	t.Run("lease", func(t *testing.T) {
		testCachingCacheClearCommon(t, "lease")
	})

	t.Run("token", func(t *testing.T) {
		testCachingCacheClearCommon(t, "token")
	})

	t.Run("token_accessor", func(t *testing.T) {
		testCachingCacheClearCommon(t, "token_accessor")
	})

	t.Run("all", func(t *testing.T) {
		testCachingCacheClearCommon(t, "all")
	})
}

func testCachingCacheClearCommon(t *testing.T, clearType string) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, client, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	err := client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write data to the lease-kv backend
	_, err = testClient.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Proxy this request, agent should cache the response
	resp, err := testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	gotLeaseID := resp.LeaseID

	// Verify the entry exists
	idx, err := leaseCache.db.Get(cachememdb.IndexNameLease, gotLeaseID)
	if err != nil {
		t.Fatal(err)
	}

	if idx == nil {
		t.Fatalf("expected cached entry, got: %v", idx)
	}

	data := map[string]interface{}{
		"type": clearType,
	}

	// We need to set the value here depending on what we're trying to test.
	// Some values are be static, but others are dynamically generated at runtime.
	switch clearType {
	case "request_path":
		data["value"] = "/v1/kv/foo"
	case "lease":
		data["value"] = resp.LeaseID
	case "token":
		data["value"] = testClient.Token()
	case "token_accessor":
		lookupResp, err := client.Auth().Token().Lookup(testClient.Token())
		if err != nil {
			t.Fatal(err)
		}
		data["value"] = lookupResp.Data["accessor"]
	case "all":
	default:
		t.Fatalf("invalid type provided: %v", clearType)
	}

	r := testClient.NewRequest("PUT", consts.AgentPathCacheClear)
	if err := r.SetJSONBody(data); err != nil {
		t.Fatal(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	apiResp, err := testClient.RawRequestWithContext(ctx, r)
	if apiResp != nil {
		defer apiResp.Body.Close()
	}
	if apiResp != nil && apiResp.StatusCode == 404 {
		_, parseErr := api.ParseSecret(apiResp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
		default:
			t.Fatal(err)
		}
	}
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify the entry is cleared
	idx, err = leaseCache.db.Get(cachememdb.IndexNameLease, gotLeaseID)
	if err != cachememdb.ErrCacheItemNotFound {
		t.Fatal("expected entry to be nil, got", err)
	}
}

func TestCache_AuthTokenCreateOrphan(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("managed", func(t *testing.T) {
			cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
			defer cleanup()

			reqOpts := &api.TokenCreateRequest{
				Policies: []string{"default"},
				NoParent: true,
			}
			resp, err := testClient.Auth().Token().Create(reqOpts)
			if err != nil {
				t.Fatal(err)
			}
			token := resp.Auth.ClientToken

			idx, err := leaseCache.db.Get(cachememdb.IndexNameToken, token)
			if err != nil {
				t.Fatal(err)
			}
			if idx == nil {
				t.Fatalf("expected entry to be non-nil, got: %#v", idx)
			}
		})

		t.Run("non-managed", func(t *testing.T) {
			cleanup, clusterClient, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
			defer cleanup()

			reqOpts := &api.TokenCreateRequest{
				Policies: []string{"default"},
				NoParent: true,
			}

			// Use the test client but set the token to one that's not managed by agent
			testClient.SetToken(clusterClient.Token())

			resp, err := testClient.Auth().Token().Create(reqOpts)
			if err != nil {
				t.Fatal(err)
			}
			token := resp.Auth.ClientToken

			idx, err := leaseCache.db.Get(cachememdb.IndexNameToken, token)
			if err != nil {
				t.Fatal(err)
			}
			if idx == nil {
				t.Fatalf("expected entry to be non-nil, got: %#v", idx)
			}
		})
	})

	t.Run("create-orphan", func(t *testing.T) {
		t.Run("managed", func(t *testing.T) {
			cleanup, _, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
			defer cleanup()

			reqOpts := &api.TokenCreateRequest{
				Policies: []string{"default"},
			}
			resp, err := testClient.Auth().Token().CreateOrphan(reqOpts)
			if err != nil {
				t.Fatal(err)
			}
			token := resp.Auth.ClientToken

			idx, err := leaseCache.db.Get(cachememdb.IndexNameToken, token)
			if err != nil {
				t.Fatal(err)
			}
			if idx == nil {
				t.Fatalf("expected entry to be non-nil, got: %#v", idx)
			}
		})

		t.Run("non-managed", func(t *testing.T) {
			cleanup, clusterClient, testClient, leaseCache := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
			defer cleanup()

			reqOpts := &api.TokenCreateRequest{
				Policies: []string{"default"},
			}

			// Use the test client but set the token to one that's not managed by agent
			testClient.SetToken(clusterClient.Token())

			resp, err := testClient.Auth().Token().CreateOrphan(reqOpts)
			if err != nil {
				t.Fatal(err)
			}
			token := resp.Auth.ClientToken

			idx, err := leaseCache.db.Get(cachememdb.IndexNameToken, token)
			if err != nil {
				t.Fatal(err)
			}
			if idx == nil {
				t.Fatalf("expected entry to be non-nil, got: %#v", idx)
			}
		})
	})
}
