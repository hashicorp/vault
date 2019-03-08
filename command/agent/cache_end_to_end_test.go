package agent

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agent/auth"
	agentapprole "github.com/hashicorp/vault/command/agent/auth/approle"
	"github.com/hashicorp/vault/command/agent/cache"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/command/agent/sink/inmem"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

const policyAutoAuthAppRole = `
path "/kv/*" {
	capabilities = ["sudo", "create", "read", "update", "delete", "list"]
}

path "/auth/token/create" {
	capabilities = ["create", "update"]
}
`

func TestCache_UsingAutoAuthToken(t *testing.T) {
	var err error
	logger := logging.NewVaultLogger(log.Trace)
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, client.Address())

	defer os.Setenv(api.EnvVaultCACert, os.Getenv(api.EnvVaultCACert))
	os.Setenv(api.EnvVaultCACert, fmt.Sprintf("%s/ca_cert.pem", cluster.TempDir))

	err = client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a secret in the backend
	_, err = client.Logical().Write("kv/foo", map[string]interface{}{
		"value": "bar",
		"ttl":   "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Add an kv-admin policy
	if err := client.Sys().PutPolicy("test-autoauth", policyAutoAuthAppRole); err != nil {
		t.Fatal(err)
	}

	// Enable approle
	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test1", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "3s",
		"token_max_ttl":  "10s",
		"policies":       []string{"test-autoauth"},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Write("auth/approle/role/test1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID1 := resp.Data["secret_id"].(string)

	resp, err = client.Logical().Read("auth/approle/role/test1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID1 := resp.Data["role_id"].(string)

	rolef, err := ioutil.TempFile("", "auth.role-id.test.")
	if err != nil {
		t.Fatal(err)
	}
	role := rolef.Name()
	rolef.Close() // WriteFile doesn't need it open
	defer os.Remove(role)
	t.Logf("input role_id_file_path: %s", role)

	secretf, err := ioutil.TempFile("", "auth.secret-id.test.")
	if err != nil {
		t.Fatal(err)
	}
	secret := secretf.Name()
	secretf.Close()
	defer os.Remove(secret)
	t.Logf("input secret_id_file_path: %s", secret)

	// We close these right away because we're just basically testing
	// permissions and finding a usable file name
	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	out := ouf.Name()
	ouf.Close()
	os.Remove(out)
	t.Logf("output: %s", out)

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	conf := map[string]interface{}{
		"role_id_file_path":                   role,
		"secret_id_file_path":                 secret,
		"remove_secret_id_file_after_reading": true,
	}

	cacheLogger := logging.NewVaultLogger(hclog.Trace).Named("cache")

	// Create the API proxier
	apiProxy, err := cache.NewAPIProxy(&cache.APIProxyConfig{
		Client: client,
		Logger: cacheLogger.Named("apiproxy"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create the lease cache proxier and set its underlying proxier to
	// the API proxier.
	leaseCache, err := cache.NewLeaseCache(&cache.LeaseCacheConfig{
		Client:      client,
		BaseContext: ctx,
		Proxier:     apiProxy,
		Logger:      cacheLogger.Named("leasecache"),
	})
	if err != nil {
		t.Fatal(err)
	}

	am, err := agentapprole.NewApproleAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.approle"),
		MountPath: "auth/approle",
		Config:    conf,
	})
	if err != nil {
		t.Fatal(err)
	}
	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	}
	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
	}()

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": out,
		},
	}
	fs, err := file.NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})

	inmemSinkConfig := &sink.SinkConfig{
		Logger: logger.Named("sink.inmem"),
	}

	inmemSink, err := inmem.New(inmemSinkConfig, leaseCache)
	if err != nil {
		t.Fatal(err)
	}
	inmemSinkConfig.Sink = inmemSink

	go ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config, inmemSinkConfig})
	defer func() {
		<-ss.DoneCh
	}()

	// This has to be after the other defers so it happens first
	defer cancelFunc()

	// Check that no sink file exists
	_, err = os.Lstat(out)
	if err == nil {
		t.Fatal("expected err")
	}
	if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	if err := ioutil.WriteFile(role, []byte(roleID1), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test role 1", "path", role)
	}

	if err := ioutil.WriteFile(secret, []byte(secretID1), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test secret 1", "path", secret)
	}

	getToken := func() string {
		timeout := time.Now().Add(10 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("did not find a written token after timeout")
			}
			val, err := ioutil.ReadFile(out)
			if err == nil {
				os.Remove(out)
				if len(val) == 0 {
					t.Fatal("written token was empty")
				}

				_, err = os.Stat(secret)
				if err == nil {
					t.Fatal("secret file exists but was supposed to be removed")
				}

				return string(val)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}

	t.Logf("auto-auth token: %q", getToken())

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	// Create a muxer and add paths relevant for the lease cache layer
	mux := http.NewServeMux()
	mux.Handle(consts.AgentPathCacheClear, leaseCache.HandleCacheClear(ctx))

	// Passing a non-nil inmemsink tells the agent to use the auto-auth token
	mux.Handle("/", cache.Handler(ctx, cacheLogger, leaseCache, inmemSink))
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          cacheLogger.StandardLogger(nil),
	}
	go server.Serve(listener)

	testClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	if err := testClient.SetAddress("http://" + listener.Addr().String()); err != nil {
		t.Fatal(err)
	}

	// Wait for listeners to come up
	time.Sleep(2 * time.Second)

	// This block tests that no token on the client is detected by the agent
	// and the auto-auth token is used
	{
		// Empty the token in the client to ensure that auto-auth token is used
		testClient.SetToken("")

		resp, err = testClient.Logical().Read("auth/token/lookup-self")
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatalf("failed to use the auto-auth token to perform lookup-self")
		}
	}

	// This block tests lease creation caching using the auto-auth token.
	{
		resp, err = testClient.Logical().Read("kv/foo")
		if err != nil {
			t.Fatal(err)
		}

		origReqID := resp.RequestID

		resp, err = testClient.Logical().Read("kv/foo")
		if err != nil {
			t.Fatal(err)
		}

		// Sleep for a bit to allow renewer logic to kick in
		time.Sleep(20 * time.Millisecond)

		cacheReqID := resp.RequestID

		if origReqID != cacheReqID {
			t.Fatalf("request ID  mismatch, expected second request to be a cached response: %s != %s", origReqID, cacheReqID)
		}
	}

	// This block tests auth token creation caching (child, non-orphan tokens)
	// using the auto-auth token.
	{
		resp, err = testClient.Logical().Write("auth/token/create", nil)
		if err != nil {
			t.Fatal(err)
		}
		origReqID := resp.RequestID

		// Sleep for a bit to allow renewer logic to kick in
		time.Sleep(20 * time.Millisecond)

		resp, err = testClient.Logical().Write("auth/token/create", nil)
		if err != nil {
			t.Fatal(err)
		}
		cacheReqID := resp.RequestID

		if origReqID != cacheReqID {
			t.Fatalf("request ID mismatch, expected second request to be a cached response: %s != %s", origReqID, cacheReqID)
		}
	}

	// This blocks tests that despite being allowed to use auto-auth token, the
	// token on the request will be prioritized.
	{
		// Empty the token in the client to ensure that auto-auth token is used
		testClient.SetToken(client.Token())

		resp, err = testClient.Logical().Read("auth/token/lookup-self")
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || resp.Data["id"] != client.Token() {
			t.Fatalf("failed to use the cluster client token to perform lookup-self")
		}
	}
}
