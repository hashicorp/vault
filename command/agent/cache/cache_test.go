package cache

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/command/agent/sink/mock"
	"github.com/hashicorp/vault/logical"

	"github.com/go-test/deep"
	hclog "github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/namespace"
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
func setupClusterAndAgent(ctx context.Context, t *testing.T, coreConfig *vault.CoreConfig) (func(), *api.Client, *api.Client, *LeaseCache) {
	t.Helper()

	if ctx == nil {
		ctx = context.Background()
	}

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

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	// Create the API proxier
	apiProxy, err := NewAPIProxy(&APIProxyConfig{
		Client: clusterClient,
		Logger: cacheLogger.Named("apiproxy"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create the lease cache proxier and set its underlying proxier to
	// the API proxier.
	leaseCache, err := NewLeaseCache(&LeaseCacheConfig{
		Client:      clusterClient,
		BaseContext: ctx,
		Proxier:     apiProxy,
		Logger:      cacheLogger.Named("leasecache"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a muxer and add paths relevant for the lease cache layer
	mux := http.NewServeMux()
	mux.Handle("/agent/v1/cache-clear", leaseCache.HandleCacheClear(ctx))

	mux.Handle("/", Handler(ctx, cacheLogger, leaseCache, nil))
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

	if err := testClient.SetAddress("http://" + listener.Addr().String()); err != nil {
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

func tokenRevocationValidation(t *testing.T, sampleSpace map[string]string, expected map[string]string, leaseCache *LeaseCache) {
	t.Helper()
	for val, valType := range sampleSpace {
		index, err := leaseCache.db.Get(valType, val)
		if err != nil {
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
	responses := []*SendResponse{
		&SendResponse{
			Response: &api.Response{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader(response1)),
				},
			},
			ResponseBody: []byte(response1),
		},
		&SendResponse{
			Response: &api.Response{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader(response2)),
				},
			},
			ResponseBody: []byte(response2),
		},
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
	mux.Handle("/v1/agent/cache-clear", leaseCache.HandleCacheClear(ctx))

	mux.Handle("/", Handler(ctx, cacheLogger, leaseCache, mock.NewSink("testid")))
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
}

func TestCache_ConcurrentRequests(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
				t.Fatal(err)
			}
			secret, err := testClient.Logical().Read(key)
			if err != nil {
				t.Fatal(err)
			}
			if secret == nil || secret.Data["key"].(string) != key {
				t.Fatal(fmt.Sprintf("failed to read value for key: %q", key))
			}
		}(i)

	}
	wg.Wait()
}

func TestCache_TokenRevocations_RevokeOrphan(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
}

func TestCache_AuthResponse(t *testing.T) {
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

func TestCache_LeaseResponse(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
