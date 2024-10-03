// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/useragent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const policyAdmin = `
path "*" {
	capabilities = ["sudo", "create", "read", "update", "delete", "list"]
}
`

func TestAPIProxy(t *testing.T) {
	cleanup, client, _, _ := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client:                  client,
		Logger:                  logging.NewVaultLogger(hclog.Trace),
		UserAgentStringFunction: useragent.ProxyStringWithProxiedUserAgent,
		UserAgentString:         useragent.ProxyAPIProxyString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}
}

func TestAPIProxyNoCache(t *testing.T) {
	cleanup, client, _, _ := setupClusterAndAgentNoCache(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client:                  client,
		Logger:                  logging.NewVaultLogger(hclog.Trace),
		UserAgentStringFunction: useragent.ProxyStringWithProxiedUserAgent,
		UserAgentString:         useragent.ProxyAPIProxyString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}
}

func TestAPIProxy_queryParams(t *testing.T) {
	// Set up an agent that points to a standby node for this particular test
	// since it needs to proxy a /sys/health?standbyok=true request to a standby
	cleanup, client, _, _ := setupClusterAndAgentOnStandby(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client:                  client,
		Logger:                  logging.NewVaultLogger(hclog.Trace),
		UserAgentStringFunction: useragent.ProxyStringWithProxiedUserAgent,
		UserAgentString:         useragent.ProxyAPIProxyString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	// Add a query parameter for testing
	q := req.URL.Query()
	q.Add("standbyok", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || !result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}

	if resp.Response.StatusCode != http.StatusOK {
		t.Fatalf("exptected standby to return 200, got: %v", resp.Response.StatusCode)
	}
}

// setupClusterAndAgent is a helper func used to set up a test cluster and
// caching agent against the active node. It returns a cleanup func that should
// be deferred immediately along with two clients, one for direct cluster
// communication and another to talk to the caching agent.
func setupClusterAndAgent(ctx context.Context, t *testing.T, coreConfig *vault.CoreConfig) (func(), *api.Client, *api.Client, *LeaseCache) {
	return setupClusterAndAgentCommon(ctx, t, coreConfig, false, true)
}

// setupClusterAndAgentNoCache is a helper func used to set up a test cluster and
// proxying agent against the active node. It returns a cleanup func that should
// be deferred immediately along with two clients, one for direct cluster
// communication and another to talk to the caching agent.
func setupClusterAndAgentNoCache(ctx context.Context, t *testing.T, coreConfig *vault.CoreConfig) (func(), *api.Client, *api.Client, *LeaseCache) {
	return setupClusterAndAgentCommon(ctx, t, coreConfig, false, false)
}

// setupClusterAndAgentOnStandby is a helper func used to set up a test cluster
// and caching agent against a standby node. It returns a cleanup func that
// should be deferred immediately along with two clients, one for direct cluster
// communication and another to talk to the caching agent.
func setupClusterAndAgentOnStandby(ctx context.Context, t *testing.T, coreConfig *vault.CoreConfig) (func(), *api.Client, *api.Client, *LeaseCache) {
	return setupClusterAndAgentCommon(ctx, t, coreConfig, true, true)
}

func setupClusterAndAgentCommon(ctx context.Context, t *testing.T, coreConfig *vault.CoreConfig, onStandby bool, useCache bool) (func(), *api.Client, *api.Client, *LeaseCache) {
	t.Helper()

	if ctx == nil {
		ctx = context.Background()
	}

	if coreConfig == nil {
		coreConfig = &vault.CoreConfig{}
	}
	// Always set up the userpass backend since we use that to generate an admin
	// token for the client that will make proxied requests to through the agent.
	if coreConfig.CredentialBackends == nil || coreConfig.CredentialBackends["userpass"] == nil {
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

	activeClient := cores[0].Client
	standbyClient := cores[1].Client

	// clienToUse is the client for the agent to point to.
	clienToUse := activeClient
	if onStandby {
		clienToUse = standbyClient
	}

	// Add an admin policy
	if err := activeClient.Sys().PutPolicy("admin", policyAdmin); err != nil {
		t.Fatal(err)
	}

	// Set up the userpass auth backend and an admin user. Used for getting a token
	// for the agent later down in this func.
	err := activeClient.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = activeClient.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
		"policies": []string{"admin"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set up env vars for agent consumption
	origEnvVaultAddress := os.Getenv(api.EnvVaultAddress)
	os.Setenv(api.EnvVaultAddress, clienToUse.Address())

	origEnvVaultCACert := os.Getenv(api.EnvVaultCACert)
	os.Setenv(api.EnvVaultCACert, fmt.Sprintf("%s/ca_cert.pem", cluster.TempDir))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	apiProxyLogger := cluster.Logger.Named("apiproxy")

	// Create the API proxier
	apiProxy, err := NewAPIProxy(&APIProxyConfig{
		Client:                  clienToUse,
		Logger:                  apiProxyLogger,
		UserAgentStringFunction: useragent.ProxyStringWithProxiedUserAgent,
		UserAgentString:         useragent.ProxyAPIProxyString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a muxer and add paths relevant for the lease cache layer and API proxy layer
	mux := http.NewServeMux()

	var leaseCache *LeaseCache
	if useCache {
		cacheLogger := cluster.Logger.Named("cache")

		// Create the lease cache proxier and set its underlying proxier to
		// the API proxier.
		leaseCache, err = NewLeaseCache(&LeaseCacheConfig{
			Client:              clienToUse,
			BaseContext:         ctx,
			Proxier:             apiProxy,
			Logger:              cacheLogger.Named("leasecache"),
			CacheDynamicSecrets: true,
			UserAgentToUse:      "test",
		})
		if err != nil {
			t.Fatal(err)
		}

		mux.Handle("/agent/v1/cache-clear", leaseCache.HandleCacheClear(ctx))

		mux.Handle("/", ProxyHandler(ctx, cacheLogger, leaseCache, nil, false, false, nil, nil))
	} else {
		mux.Handle("/", ProxyHandler(ctx, apiProxyLogger, apiProxy, nil, false, false, nil, nil))
	}

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          apiProxyLogger.StandardLogger(nil),
	}
	go server.Serve(listener)

	// testClient is the client that is used to talk to the agent for proxying/caching behavior.
	testClient, err := activeClient.Clone()
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
		// We wait for a tiny bit for things such as agent renewal to exit properly
		time.Sleep(50 * time.Millisecond)

		cluster.Cleanup()
		os.Setenv(api.EnvVaultAddress, origEnvVaultAddress)
		os.Setenv(api.EnvVaultCACert, origEnvVaultCACert)
		listener.Close()
	}

	return cleanup, clienToUse, testClient, leaseCache
}
