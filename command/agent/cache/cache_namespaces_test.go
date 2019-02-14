package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-test/deep"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestCache_Namespaces(t *testing.T) {
	t.Run("send", testSendNamespaces)

	t.Run("full_path", func(t *testing.T) {
		t.Run("handle_cacheclear", func(t *testing.T) {
			testHandleCacheClearNamespaces(t, true)
		})

		t.Run("eviction_on_revocation", func(t *testing.T) {
			testEvictionOnRevocationNamespaces(t, true)
		})
	})

	t.Run("namespace_header", func(t *testing.T) {
		t.Run("handle_cacheclear", func(t *testing.T) {
			testHandleCacheClearNamespaces(t, false)
		})

		t.Run("eviction_on_revocation", func(t *testing.T) {
			testEvictionOnRevocationNamespaces(t, false)
		})
	})
}

func testSendNamespaces(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, clusterClient, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	// Create a namespace
	_, err := clusterClient.Logical().Write("sys/namespaces/ns1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mount the leased KV into ns1
	clusterClient.SetNamespace("ns1/")
	err = clusterClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}
	clusterClient.SetNamespace("")

	// Try request using full path
	{
		// Write some random value
		_, err = clusterClient.Logical().Write("/ns1/kv/foo", map[string]interface{}{
			"value": "test",
			"ttl":   "1h",
		})
		if err != nil {
			t.Fatal(err)
		}

		proxiedResp, err := testClient.Logical().Read("/ns1/kv/foo")
		if err != nil {
			t.Fatal(err)
		}

		cachedResp, err := testClient.Logical().Read("/ns1/kv/foo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := deep.Equal(proxiedResp, cachedResp); diff != nil {
			t.Fatal(diff)
		}
	}

	// Try request using the namespace header
	{
		// Write some random value
		_, err = clusterClient.Logical().Write("/ns1/kv/bar", map[string]interface{}{
			"value": "test",
			"ttl":   "1h",
		})
		if err != nil {
			t.Fatal(err)
		}

		testClient.SetNamespace("ns1/")
		proxiedResp, err := testClient.Logical().Read("/kv/bar")
		if err != nil {
			t.Fatal(err)
		}

		cachedResp, err := testClient.Logical().Read("/kv/bar")
		if err != nil {
			t.Fatal(err)
		}

		if diff := deep.Equal(proxiedResp, cachedResp); diff != nil {
			t.Fatal(diff)
		}
		testClient.SetNamespace("")
	}

	// Try the same request using different namespace input methods (header vs
	// full path), they should not be the same cache entry (i.e. should produce
	// different lease ID's).
	{
		_, err := clusterClient.Logical().Write("/ns1/kv/baz", map[string]interface{}{
			"value": "test",
			"ttl":   "1h",
		})
		if err != nil {
			t.Fatal(err)
		}

		proxiedResp, err := testClient.Logical().Read("/ns1/kv/baz")
		if err != nil {
			t.Fatal(err)
		}

		testClient.SetNamespace("ns1/")
		cachedResp, err := testClient.Logical().Read("/kv/baz")
		if err != nil {
			t.Fatal(err)
		}
		testClient.SetNamespace("")

		if diff := deep.Equal(proxiedResp, cachedResp); diff == nil {
			t.Logf("response #1: %#v", proxiedResp)
			t.Logf("response #2: %#v", cachedResp)
			t.Fatal("expected requests to be not cached")
		}
	}
}

func testHandleCacheClearNamespaces(t *testing.T, fullPath bool) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, clusterClient, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	// Create a namespace
	_, err := clusterClient.Logical().Write("sys/namespaces/ns1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mount the leased KV into ns1
	clusterClient.SetNamespace("ns1/")
	err = clusterClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}
	clusterClient.SetNamespace("")

	// Write some random value
	_, err = clusterClient.Logical().Write("/ns1/kv/foo", map[string]interface{}{
		"value": "test",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	requestPath := "/kv/foo"
	testClient.SetNamespace("ns1/")
	if fullPath {
		requestPath = "/ns1" + requestPath
		testClient.SetNamespace("")
	}

	// Request the secret
	firstResp, err := testClient.Logical().Read(requestPath)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)

	// Clear by request_path and namespace
	requestPathValue := "/v1" + requestPath
	data := &cacheClearRequest{
		Type:  "request_path",
		Value: requestPathValue,
	}

	r := testClient.NewRequest("PUT", "/v1/agent/cache-clear")
	if err := r.SetJSONBody(data); err != nil {
		t.Fatal(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	_, err = clusterClient.RawRequestWithContext(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)

	secondResp, err := testClient.Logical().Read(requestPath)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(firstResp, secondResp); diff == nil {
		t.Logf("response #1: %#v", firstResp)
		t.Logf("response #2: %#v", secondResp)
		t.Fatal("expected requests to be not cached")
	}
}

func testEvictionOnRevocationNamespaces(t *testing.T, fullPath bool) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, clusterClient, testClient, _ := setupClusterAndAgent(namespace.RootContext(nil), t, coreConfig)
	defer cleanup()

	// Create a namespace
	_, err := clusterClient.Logical().Write("sys/namespaces/ns1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mount the leased KV into ns1
	clusterClient.SetNamespace("ns1/")
	err = clusterClient.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}
	clusterClient.SetNamespace("")

	// Write some random value
	_, err = clusterClient.Logical().Write("/ns1/kv/foo", map[string]interface{}{
		"value": "test",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	requestPath := "/kv/foo"
	testClient.SetNamespace("ns1/")
	if fullPath {
		requestPath = "/ns1" + requestPath
		testClient.SetNamespace("")
	}

	// Request the secret
	firstResp, err := testClient.Logical().Read(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	leaseID := firstResp.LeaseID

	time.Sleep(200 * time.Millisecond)

	revocationPath := "/sys/leases/revoke"
	if fullPath {
		revocationPath = "/ns1/" + revocationPath
		testClient.SetNamespace("")
	}

	_, err = testClient.Logical().Write(revocationPath, map[string]interface{}{
		"lease_id": leaseID,
	})

	secondResp, err := testClient.Logical().Read(requestPath)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(firstResp, secondResp); diff == nil {
		t.Logf("response #1: %#v", firstResp)
		t.Logf("response #2: %#v", secondResp)
		t.Fatal("expected requests to be not cached")
	}
}
