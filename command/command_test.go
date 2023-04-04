// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/benchhelpers"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/cli"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
)

var (
	defaultVaultLogger = log.NewNullLogger()

	defaultVaultCredentialBackends = map[string]logical.Factory{
		"userpass": credUserpass.Factory,
	}

	defaultVaultAuditBackends = map[string]audit.Factory{
		"file": auditFile.Factory,
	}

	defaultVaultLogicalBackends = map[string]logical.Factory{
		"generic-leased": vault.LeasedPassthroughBackendFactory,
		"pki":            pki.Factory,
		"ssh":            ssh.Factory,
		"transit":        transit.Factory,
		"kv":             kv.Factory,
	}
)

// assertNoTabs asserts the CLI help has no tab characters.
func assertNoTabs(tb testing.TB, c cli.Command) {
	tb.Helper()

	if strings.ContainsRune(c.Help(), '\t') {
		tb.Errorf("%#v help output contains tabs", c)
	}
}

// testVaultServer creates a test vault cluster and returns a configured API
// client and closer function.
func testVaultServer(tb testing.TB) (*api.Client, func()) {
	tb.Helper()

	client, _, closer := testVaultServerUnseal(tb)
	return client, closer
}

func testVaultServerWithKVVersion(tb testing.TB, kvVersion string) (*api.Client, func()) {
	tb.Helper()

	client, _, closer := testVaultServerUnsealWithKVVersionWithSeal(tb, kvVersion, nil)
	return client, closer
}

func testVaultServerAllBackends(tb testing.TB) (*api.Client, func()) {
	tb.Helper()

	client, _, closer := testVaultServerCoreConfig(tb, &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             defaultVaultLogger,
		CredentialBackends: credentialBackends,
		AuditBackends:      auditBackends,
		LogicalBackends:    logicalBackends,
		BuiltinRegistry:    builtinplugins.Registry,
	})
	return client, closer
}

// testVaultServerAutoUnseal creates a test vault cluster and sets it up with auto unseal
// the function returns a client, the recovery keys, and a closer function
func testVaultServerAutoUnseal(tb testing.TB) (*api.Client, []string, func()) {
	testSeal := seal.NewTestSeal(nil)
	autoSeal, err := vault.NewAutoSeal(testSeal)
	if err != nil {
		tb.Fatal("unable to create autoseal", err)
	}
	return testVaultServerUnsealWithKVVersionWithSeal(tb, "1", autoSeal)
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(tb testing.TB) (*api.Client, []string, func()) {
	return testVaultServerUnsealWithKVVersionWithSeal(tb, "1", nil)
}

func testVaultServerUnsealWithKVVersionWithSeal(tb testing.TB, kvVersion string, seal vault.Seal) (*api.Client, []string, func()) {
	tb.Helper()
	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Output:     log.DefaultOutput,
		Level:      log.Debug,
		JSONFormat: logging.ParseEnvLogFormat() == logging.JSONFormat,
	})

	return testVaultServerCoreConfigWithOpts(tb, &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             logger,
		CredentialBackends: defaultVaultCredentialBackends,
		AuditBackends:      defaultVaultAuditBackends,
		LogicalBackends:    defaultVaultLogicalBackends,
		BuiltinRegistry:    builtinplugins.Registry,
		Seal:               seal,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
		KVVersion:   kvVersion,
	})
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function
// configured with the given plugin directory.
func testVaultServerPluginDir(tb testing.TB, pluginDir string) (*api.Client, []string, func()) {
	tb.Helper()

	return testVaultServerCoreConfig(tb, &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             defaultVaultLogger,
		CredentialBackends: defaultVaultCredentialBackends,
		AuditBackends:      defaultVaultAuditBackends,
		LogicalBackends:    defaultVaultLogicalBackends,
		PluginDirectory:    pluginDir,
		BuiltinRegistry:    builtinplugins.Registry,
	})
}

func testVaultServerCoreConfig(tb testing.TB, coreConfig *vault.CoreConfig) (*api.Client, []string, func()) {
	return testVaultServerCoreConfigWithOpts(tb, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1, // Default is 3, but we don't need that many
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper. If the seal config supports recovery keys, then
// recovery keys are returned. Otherwise, unseal keys are returned
func testVaultServerCoreConfigWithOpts(tb testing.TB, coreConfig *vault.CoreConfig, opts *vault.TestClusterOptions) (*api.Client, []string, func()) {
	tb.Helper()

	cluster := vault.NewTestCluster(benchhelpers.TBtoT(tb), coreConfig, opts)
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(benchhelpers.TBtoT(tb), core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	var keys [][]byte
	if coreConfig.Seal != nil && coreConfig.Seal.RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}

	return client, encodeKeys(keys), cluster.Cleanup
}

// Convert the unseal keys to base64 encoded, since these are how the user
// will get them.
func encodeKeys(rawKeys [][]byte) []string {
	keys := make([]string, len(rawKeys))
	for i := range rawKeys {
		keys[i] = base64.StdEncoding.EncodeToString(rawKeys[i])
	}
	return keys
}

// testVaultServerUninit creates an uninitialized server.
func testVaultServerUninit(tb testing.TB) (*api.Client, func()) {
	tb.Helper()

	inm, err := inmem.NewInmem(nil, defaultVaultLogger)
	if err != nil {
		tb.Fatal(err)
	}

	core, err := vault.NewCore(&vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             defaultVaultLogger,
		Physical:           inm,
		CredentialBackends: defaultVaultCredentialBackends,
		AuditBackends:      defaultVaultAuditBackends,
		LogicalBackends:    defaultVaultLogicalBackends,
		BuiltinRegistry:    builtinplugins.Registry,
	})
	if err != nil {
		tb.Fatal(err)
	}

	ln, addr := vaulthttp.TestServer(tb, core)

	client, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		tb.Fatal(err)
	}

	closer := func() {
		core.Shutdown()
		ln.Close()
	}

	return client, closer
}

// testVaultServerBad creates an http server that returns a 500 on each request
// to simulate failures.
func testVaultServerBad(tb testing.TB) (*api.Client, func()) {
	tb.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatal(err)
	}

	server := &http.Server{
		Addr: "127.0.0.1:0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			tb.Fatal(err)
		}
	}()

	client, err := api.NewClient(&api.Config{
		Address: "http://" + listener.Addr().String(),
	})
	if err != nil {
		tb.Fatal(err)
	}

	return client, func() {
		ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		server.Shutdown(ctx)
	}
}

// testTokenAndAccessor creates a new authentication token capable of being renewed with
// the default policy attached. It returns the token and it's accessor.
func testTokenAndAccessor(tb testing.TB, client *api.Client) (string, string) {
	tb.Helper()

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"default"},
		TTL:      "30m",
	})
	if err != nil {
		tb.Fatal(err)
	}
	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		tb.Fatalf("missing auth data: %#v", secret)
	}
	return secret.Auth.ClientToken, secret.Auth.Accessor
}

func testClient(tb testing.TB, addr string, token string) *api.Client {
	tb.Helper()
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		tb.Fatal(err)
	}
	client.SetToken(token)

	return client
}
