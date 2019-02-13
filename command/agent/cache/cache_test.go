package cache

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/command/agent/cache/cachememdb"
	"github.com/hashicorp/vault/logical"

	"github.com/go-test/deep"
	hclog "github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

const policyAdmin = `
path "*" {
	capabilities = ["sudo", "create", "read", "update", "delete", "list"]
}
`

// setupClusterAndAgent is a helper func used to set up a test cluster and
// caching agent. It returns a cleanup func that should be deferred immediately
// along with two clients, one for direct cluster communication and another to
// talk to the caching agent.
func setupClusterAndAgent(t *testing.T, coreConfig *vault.CoreConfig) (func(), *api.Client, *api.Client, *LeaseCache) {
	t.Helper()

	// Handle sane defaults
	if coreConfig == nil {
		coreConfig = &vault.CoreConfig{
			DisableMlock: true,
			DisableCache: true,
			Logger:       logging.NewVaultLogger(hclog.Trace),
			CredentialBackends: map[string]logical.Factory{
				"userpass": userpass.Factory,
			},
		}
	}

	if coreConfig.CredentialBackends == nil {
		coreConfig.CredentialBackends = map[string]logical.Factory{
			"userpass": userpass.Factory,
		}
	}

	// Init new test cluster
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)

	// clusterClient is the client that is used to talk directly to the cluster.
	clusterClient := cores[0].Client

	// Add an admin policy
	if err := clusterClient.Sys().PutPolicy("admin", policyAdmin); err != nil {
		t.Fatal(err)
	}

	// Set up the userpass auth backend and an admin user. Used for getting a token
	// for the agent later down in this func.
	clusterClient.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})

	_, err := clusterClient.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
		"policies": []string{"admin"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set up env vars for agent consumption
	origEnvVaultAddress := os.Getenv(api.EnvVaultAddress)
	os.Setenv(api.EnvVaultAddress, clusterClient.Address())

	origEnvVaultCACert := os.Getenv(api.EnvVaultCACert)
	os.Setenv(api.EnvVaultCACert, fmt.Sprintf("%s/ca_cert.pem", cluster.TempDir))

	cacheLogger := logging.NewVaultLogger(hclog.Trace).Named("cache")
	ctx := context.Background()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	// Create the API proxier
	apiProxy := NewAPIProxy(&APIProxyConfig{
		Logger: cacheLogger.Named("apiproxy"),
	})

	// Create the lease cache proxier and set its underlying proxier to
	// the API proxier.
	leaseCache, err := NewLeaseCache(&LeaseCacheConfig{
		BaseContext: ctx,
		Proxier:     apiProxy,
		Logger:      cacheLogger.Named("leasecache"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a muxer and add paths relevant for the lease cache layer
	mux := http.NewServeMux()
	mux.Handle("/v1/agent/cache-clear", leaseCache.HandleCacheClear(ctx))

	mux.Handle("/", Handler(ctx, cacheLogger, leaseCache, false, clusterClient))
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          cacheLogger.StandardLogger(nil),
	}
	go server.Serve(listener)

	// testClient is the client that is used to talk to the agent for proxying/caching behavior.
	testClient, err := clusterClient.Clone()
	if err != nil {
		t.Fatal(err)
	}

	if err := testClient.SetAddress("http://" + listener.Addr()); err != nil {
		t.Fatal(err)
	}

	// Login via userpass method to derive a managed token. Set that token as the
	// testClient's token
	resp, err := testClient.Logical().Write("auth/userpass/login/foo", map[string]interface{}{
		"password": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	testClient.SetToken(resp.Auth.ClientToken)

	cleanup := func() {
		cluster.Cleanup()
		os.Setenv(api.EnvVaultAddress, origEnvVaultAddress)
		os.Setenv(api.EnvVaultCACert, origEnvVaultCACert)
		listener.Close()
	}

	return cleanup, clusterClient, testClient, leaseCache
}

func TestCache_TokenRevocations(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, _, testClient, leaseCache := setupClusterAndAgent(t, coreConfig)
	defer cleanup()

	token1 := testClient.Token()

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
	lease11 := leaseResp.LeaseID

	resp, err := testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token2 := resp.Auth.ClientToken

	testClient.SetToken(token2)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease21 := leaseResp.LeaseID

	resp, err = testClient.Logical().Write("auth/token/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	token3 := resp.Auth.ClientToken

	testClient.SetToken(token3)

	leaseResp, err = testClient.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	lease31 := leaseResp.LeaseID

	// TODO: This test will be enhanced soon to use all the values here
	fmt.Printf("===== token1: %q\n", token1)
	fmt.Printf("===== lease11: %#v\n", lease11)
	fmt.Printf("===== token2: %q\n", token2)
	fmt.Printf("===== lease21: %#v\n", lease21)
	fmt.Printf("===== token3: %q\n", token3)
	fmt.Printf("===== lease31: %#v\n", lease31)

	/*
		testClient.SetToken(token1)
		err = testClient.Auth().Token().RevokeSelf("")
		if err != nil {
			t.Fatal(err)
		}
	*/

	indexes, err := leaseCache.db.GetAll(cachememdb.IndexNameID)
	if err != nil {
		t.Fatal(err)
	}
	if len(indexes) != 6 {
		t.Fatalf("bad: len(indexes); expected: 6, actual: %d", len(indexes))
	}

	leaseCache.baseCtxInfo.CancelFunc()

	time.Sleep(2 * time.Second)

	indexes, err = leaseCache.db.GetAll(cachememdb.IndexNameID)
	if err != nil {
		t.Fatal(err)
	}
	if len(indexes) != 0 {
		t.Fatalf("bad: len(indexes); expected: 0, actual: %d", len(indexes))
	}
}

func TestCache_NonCacheable(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}

	cleanup, _, testClient, _ := setupClusterAndAgent(t, coreConfig)
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
}

func TestCache_AuthResponse(t *testing.T) {
	cleanup, _, testClient, _ := setupClusterAndAgent(t, nil)
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

func TestCache_LeaseResponse(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}

	cleanup, client, testClient, _ := setupClusterAndAgent(t, coreConfig)
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
