// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agent"
	proxyConfig "github.com/hashicorp/vault/command/proxy/config"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/useragent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testProxyCommand(tb testing.TB, logger hclog.Logger) (*cli.MockUi, *ProxyCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &ProxyCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		ShutdownCh: MakeShutdownCh(),
		SighupCh:   MakeSighupCh(),
		logger:     logger,
		startedCh:  make(chan struct{}, 5),
		reloadedCh: make(chan struct{}, 5),
	}
}

// TestProxy_ExitAfterAuth tests the exit_after_auth flag, provided both
// as config and via -exit-after-auth.
func TestProxy_ExitAfterAuth(t *testing.T) {
	t.Run("via_config", func(t *testing.T) {
		testProxyExitAfterAuth(t, false)
	})

	t.Run("via_flag", func(t *testing.T) {
		testProxyExitAfterAuth(t, true)
	})
}

func testProxyExitAfterAuth(t *testing.T, viaFlag bool) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"jwt": vaultjwt.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	// Setup Vault
	err := client.Sys().EnableAuthWithOptions("jwt", &api.EnableAuthOptions{
		Type: "jwt",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/config", map[string]interface{}{
		"bound_issuer":           "https://team-vault.auth0.com/",
		"jwt_validation_pubkeys": agent.TestECDSAPubKey,
		"jwt_supported_algs":     "ES256",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/role/test", map[string]interface{}{
		"role_type":       "jwt",
		"bound_subject":   "r3qXcK2bix9eFECzsU3Sbmh0K16fatW6@clients",
		"bound_audiences": "https://vault.plugin.auth.jwt.test",
		"user_claim":      "https://vault/user",
		"groups_claim":    "https://vault/groups",
		"policies":        "test",
		"period":          "3s",
	})
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	inf, err := os.CreateTemp(dir, "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	// We remove these files in this test since we don't need the files, we just need
	// a non-conflicting file name for the config.
	os.Remove(in)
	t.Logf("input: %s", in)

	sink1f, err := os.CreateTemp(dir, "sink1.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink1 := sink1f.Name()
	sink1f.Close()
	os.Remove(sink1)
	t.Logf("sink1: %s", sink1)

	sink2f, err := os.CreateTemp(dir, "sink2.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink2 := sink2f.Name()
	sink2f.Close()
	os.Remove(sink2)
	t.Logf("sink2: %s", sink2)

	conff, err := os.CreateTemp(dir, "conf.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	conf := conff.Name()
	conff.Close()
	os.Remove(conf)
	t.Logf("config: %s", conf)

	jwtToken, _ := agent.GetTestJWT(t)
	if err := os.WriteFile(in, []byte(jwtToken), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test jwt", "path", in)
	}

	exitAfterAuthTemplText := "exit_after_auth = true"
	if viaFlag {
		exitAfterAuthTemplText = ""
	}

	config := `
%s

auto_auth {
        method {
                type = "jwt"
                config = {
                        role = "test"
                        path = "%s"
                }
        }

        sink {
                type = "file"
                config = {
                        path = "%s"
                }
        }

        sink "file" {
                config = {
                        path = "%s"
                }
        }
}
`

	config = fmt.Sprintf(config, exitAfterAuthTemplText, in, sink1, sink2)
	if err := os.WriteFile(conf, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test config", "path", conf)
	}

	doneCh := make(chan struct{})
	go func() {
		ui, cmd := testProxyCommand(t, logger)
		cmd.client = client

		args := []string{"-config", conf}
		if viaFlag {
			args = append(args, "-exit-after-auth")
		}

		code := cmd.Run(args)
		if code != 0 {
			t.Errorf("expected %d to be %d", code, 0)
			t.Logf("output from proxy:\n%s", ui.OutputWriter.String())
			t.Logf("error from proxy:\n%s", ui.ErrorWriter.String())
		}
		close(doneCh)
	}()

	select {
	case <-doneCh:
		break
	case <-time.After(1 * time.Minute):
		t.Fatal("timeout reached while waiting for proxy to exit")
	}

	sink1Bytes, err := os.ReadFile(sink1)
	if err != nil {
		t.Fatal(err)
	}
	if len(sink1Bytes) == 0 {
		t.Fatal("got no output from sink 1")
	}

	sink2Bytes, err := os.ReadFile(sink2)
	if err != nil {
		t.Fatal(err)
	}
	if len(sink2Bytes) == 0 {
		t.Fatal("got no output from sink 2")
	}

	if string(sink1Bytes) != string(sink2Bytes) {
		t.Fatal("sink 1/2 values don't match")
	}
}

// TestProxy_AutoAuth_UserAgent tests that the User-Agent sent
// to Vault by Vault Proxy is correct when performing Auto-Auth.
// Uses the custom handler userAgentHandler (defined above) so
// that Vault validates the User-Agent on requests sent by Proxy.
func TestProxy_AutoAuth_UserAgent(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	var h userAgentHandler
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}, &vault.TestClusterOptions{
		NumCores: 1,
		HandlerFunc: vaulthttp.HandlerFunc(
			func(properties *vault.HandlerProperties) http.Handler {
				h.props = properties
				h.userAgentToCheckFor = useragent.ProxyAutoAuthString()
				h.requestMethodToCheck = "PUT"
				h.pathToCheck = "auth/approle/login"
				h.t = t
				return &h
			}),
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// Create a named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role")
	req.BodyBytes = []byte(`{
	  "secret_id_num_uses": "10",
	  "secret_id_ttl": "1m",
	  "token_max_ttl": "1m",
	  "token_num_uses": "10",
	  "token_ttl": "1m",
	  "policies": "default"
	}`)
	request(t, serverClient, req, 204)

	// Fetch the RoleID of the named role
	req = serverClient.NewRequest("GET", "/v1/auth/approle/role/test-role/role-id")
	body := request(t, serverClient, req, 200)
	data := body["data"].(map[string]interface{})
	roleID := data["role_id"].(string)

	// Get a SecretID issued against the named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role/secret-id")
	body = request(t, serverClient, req, 200)
	data = body["data"].(map[string]interface{})
	secretID := data["secret_id"].(string)

	// Write the RoleID and SecretID to temp files
	roleIDPath := makeTempFile(t, "role_id.txt", roleID+"\n")
	secretIDPath := makeTempFile(t, "secret_id.txt", secretID+"\n")
	defer os.Remove(roleIDPath)
	defer os.Remove(secretIDPath)

	sinkf, err := os.CreateTemp("", "sink.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink := sinkf.Name()
	sinkf.Close()
	os.Remove(sink)

	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
        }
    }

	sink "file" {
		config = {
			path = "%s"
		}
	}
}`, roleIDPath, secretIDPath, sink)

	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
api_proxy {
  use_auto_auth_token = true
}
%s
%s
`, serverClient.Address(), listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start proxy
	_, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	// Validate that the auto-auth token has been correctly attained
	// and works for LookupSelf
	conf := api.DefaultConfig()
	conf.Address = "http://" + listenAddr
	proxyClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	proxyClient.SetToken("")
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the token to be sent to syncs and be available to be used
	time.Sleep(5 * time.Second)

	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	body = request(t, proxyClient, req, 200)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_APIProxyWithoutCache_UserAgent tests that the User-Agent sent
// to Vault by Vault Proxy is correct using the API proxy without
// the cache configured. Uses the custom handler
// userAgentHandler struct defined in this test package, so that Vault validates the
// User-Agent on requests sent by Proxy.
func TestProxy_APIProxyWithoutCache_UserAgent(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	userAgentForProxiedClient := "proxied-client"
	var h userAgentHandler
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		NumCores: 1,
		HandlerFunc: vaulthttp.HandlerFunc(
			func(properties *vault.HandlerProperties) http.Handler {
				h.props = properties
				h.userAgentToCheckFor = useragent.ProxyStringWithProxiedUserAgent(userAgentForProxiedClient)
				h.pathToCheck = "/v1/auth/token/lookup-self"
				h.requestMethodToCheck = "GET"
				h.t = t
				return &h
			}),
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
`, serverClient.Address(), listenConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the proxy
	_, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.AddHeader("User-Agent", userAgentForProxiedClient)
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	_, err = proxyClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_APIProxyWithCache_UserAgent tests that the User-Agent sent
// to Vault by Vault Proxy is correct using the API proxy with
// the cache configured.  Uses the custom handler
// userAgentHandler struct defined in this test package, so that Vault validates the
// User-Agent on requests sent by Proxy.
func TestProxy_APIProxyWithCache_UserAgent(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	userAgentForProxiedClient := "proxied-client"
	var h userAgentHandler
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		NumCores: 1,
		HandlerFunc: vaulthttp.HandlerFunc(
			func(properties *vault.HandlerProperties) http.Handler {
				h.props = properties
				h.userAgentToCheckFor = useragent.ProxyStringWithProxiedUserAgent(userAgentForProxiedClient)
				h.pathToCheck = "/v1/auth/token/lookup-self"
				h.requestMethodToCheck = "GET"
				h.t = t
				return &h
			}),
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	cacheConfig := `
cache {
}`

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
`, serverClient.Address(), listenConfig, cacheConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the proxy
	_, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.AddHeader("User-Agent", userAgentForProxiedClient)
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	_, err = proxyClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_DynamicSecret tests that the cache successfully caches a dynamic secret
// going through the Proxy, and that a subsequent request will be served from the cache.
func TestProxy_Cache_DynamicSecret(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	cacheConfig := `
cache {
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
`, serverClient.Address(), cacheConfig, listenConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	_, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	renewable := true
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies:  []string{"default"},
		TTL:       "30m",
		Renewable: &renewable,
	}

	// This was the simplest test I could find to trigger the caching behaviour,
	// i.e. the most concise I could make the test that I can tell
	// creating an orphan token returns Auth, is renewable, and isn't a token
	// that's managed elsewhere (since it's an orphan)
	secret, err := proxyClient.Auth().Token().CreateOrphan(tokenCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Auth == nil {
		t.Fatalf("secret not as expected: %v", secret)
	}

	token := secret.Auth.ClientToken

	secret, err = proxyClient.Auth().Token().CreateOrphan(tokenCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Auth == nil {
		t.Fatalf("secret not as expected: %v", secret)
	}

	token2 := secret.Auth.ClientToken

	if token != token2 {
		t.Fatalf("token create response not cached when it should have been, as tokens differ")
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_DisableDynamicSecretCaching tests that the cache will not cache a dynamic secret
// if disabled in the options.
func TestProxy_Cache_DisableDynamicSecretCaching(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth for static secret caching.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	cacheConfig := `
cache {
	disable_caching_dynamic_secrets = true
	cache_static_secrets = true // We need to cache at least one kind of secret
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	_, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	renewable := true
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies:  []string{"default"},
		TTL:       "30m",
		Renewable: &renewable,
	}

	// This was the simplest test I could find to trigger the caching behaviour,
	// i.e. the most concise I could make the test that I can tell
	// creating an orphan token returns Auth, is renewable, and isn't a token
	// that's managed elsewhere (since it's an orphan)
	secret, err := proxyClient.Auth().Token().CreateOrphan(tokenCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Auth == nil {
		t.Fatalf("secret not as expected: %v", secret)
	}

	token := secret.Auth.ClientToken

	secret, err = proxyClient.Auth().Token().CreateOrphan(tokenCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Auth == nil {
		t.Fatalf("secret not as expected: %v", secret)
	}

	token2 := secret.Auth.ClientToken

	if token == token2 {
		t.Fatalf("token create response was cached, as the tokens differ")
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_StaticSecret Tests that the cache successfully caches a static secret
// going through the Proxy,
func TestProxy_Cache_StaticSecret(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth so that the event system can run.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	cacheConfig := `
cache {
	cache_static_secrets = true
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s
log_level = "trace"
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
		t.Errorf("stdout: %s", ui.OutputWriter.String())
		t.Errorf("stderr: %s", ui.ErrorWriter.String())
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	// Create kvv1 secret
	err = serverClient.KVv1("secret").Put(context.Background(), "my-secret", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// We use raw requests so we can check the headers for cache hit/miss.
	// We expect the first to miss, and the second to hit.
	req := proxyClient.NewRequest(http.MethodGet, "/v1/secret/my-secret")
	resp1, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue := resp1.Header.Get("X-Cache")
	require.Equal(t, "MISS", cacheValue)

	req = proxyClient.NewRequest(http.MethodGet, "/v1/secret/my-secret")
	resp2, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue = resp2.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	// Lastly, we check to make sure the actual data we received is
	// as we expect. We must use ParseSecret due to the raw requests.
	secret1, err := api.ParseSecret(resp1.Body)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, secretData, secret1.Data)

	secret2, err := api.ParseSecret(resp2.Body)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, secret1.Data, secret2.Data)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_EventSystemUpdatesCacheKVV1 Tests that the cache successfully caches a static secret
// going through the Proxy, and then the cache gets updated on a POST to the KVV1 secret due to an
// event.
func TestProxy_Cache_EventSystemUpdatesCacheKVV1(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth so that the event system can run.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	cacheConfig := `
cache {
	cache_static_secrets = true
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s
log_level = "trace"
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
		t.Errorf("stdout: %s", ui.OutputWriter.String())
		t.Errorf("stderr: %s", ui.ErrorWriter.String())
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	secretData2 := map[string]interface{}{
		"bar": "baz",
	}

	// Wait for the event system to successfully connect.
	// The test would pass without this time.Sleep, due to the call to updater.preEventStreamUpdate
	// but this Sleep ensures we test both paths (both updating from an event, and from
	// the pre event update).
	// As a result, we shouldn't remove this sleep, since it ensures we have greater coverage.
	time.Sleep(5 * time.Second)

	// Mount the KVV2 engine
	err = serverClient.Sys().Mount("secret-v1", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create kvv1 secret
	err = serverClient.KVv1("secret-v1").Put(context.Background(), "my-secret", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// We use raw requests so we can check the headers for cache hit/miss.
	req := proxyClient.NewRequest(http.MethodGet, "/v1/secret-v1/my-secret")
	resp1, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue := resp1.Header.Get("X-Cache")
	require.Equal(t, "MISS", cacheValue)

	// Update the secret using the proxy client
	err = proxyClient.KVv1("secret-v1").Put(context.Background(), "my-secret", secretData2)
	if err != nil {
		t.Fatal(err)
	}

	// Give some time for the event to actually get sent and the cache to be updated.
	// This is longer than it needs to be to account for unnatural slowness/avoiding
	// flakiness.
	time.Sleep(5 * time.Second)

	// We expect this to be a cache hit, with the new value
	resp2, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue = resp2.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	// Lastly, we check to make sure the actual data we received is
	// as we expect. We must use ParseSecret due to the raw requests.
	secret1, err := api.ParseSecret(resp1.Body)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, secretData, secret1.Data)

	secret2, err := api.ParseSecret(resp2.Body)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, secretData2, secret2.Data)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_EventSystemUpdatesCacheKVV2 Tests that the cache successfully caches a static secret
// going through the Proxy for a KVV2 secret, and then the cache gets updated on a POST to the secret due to an
// event.
func TestProxy_Cache_EventSystemUpdatesCacheKVV2(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth so that the event system can run.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	cacheConfig := `
cache {
	cache_static_secrets = true
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s
log_level = "trace"
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
		t.Errorf("stdout: %s", ui.OutputWriter.String())
		t.Errorf("stderr: %s", ui.ErrorWriter.String())
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	secretData2 := map[string]interface{}{
		"bar": "baz",
	}

	// Wait for the event system to successfully connect.
	// This is longer than it needs to be to account for unnatural slowness/avoiding
	// flakiness.
	time.Sleep(5 * time.Second)

	// Mount the KVV2 engine
	err = serverClient.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create kvv2 secret
	_, err = serverClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// We use raw requests so we can check the headers for cache hit/miss.
	req := proxyClient.NewRequest(http.MethodGet, "/v1/secret-v2/data/my-secret")
	resp1, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue := resp1.Header.Get("X-Cache")
	require.Equal(t, "MISS", cacheValue)

	// Update the secret using the proxy client
	_, err = proxyClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData2)
	if err != nil {
		t.Fatal(err)
	}

	// Give some time for the event to actually get sent and the cache to be updated.
	// This is longer than it needs to be to account for unnatural slowness/avoiding
	// flakiness.
	time.Sleep(5 * time.Second)

	// We expect this to be a cache hit, with the new value
	resp2, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue = resp2.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	// Lastly, we check to make sure the actual data we received is
	// as we expect. We must use ParseSecret due to the raw requests.
	secret1, err := api.ParseSecret(resp1.Body)
	if err != nil {
		t.Fatal(err)
	}
	data, ok := secret1.Data["data"]
	require.True(t, ok)
	require.Equal(t, secretData, data)

	secret2, err := api.ParseSecret(resp2.Body)
	if err != nil {
		t.Fatal(err)
	}
	data2, ok := secret2.Data["data"]
	require.True(t, ok)
	// We expect that the cached value got updated by the event system.
	require.Equal(t, secretData2, data2)

	// Lastly, ensure that a client without a token fails to access the secret.
	proxyClient.SetToken("")
	req = proxyClient.NewRequest(http.MethodGet, "/v1/secret-v2/data/my-secret")
	_, err = proxyClient.RawRequest(req)
	require.NotNil(t, err)

	_, err = proxyClient.KVv2("secret-v2").Get(context.Background(), "my-secret")
	require.NotNil(t, err)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_EventSystemUpdatesCacheUseAutoAuthToken Tests that the cache successfully caches a static secret
// going through the Proxy for a KVV2 secret, and that the cache works as expected with the
// use_auto_auth_token=force option.
func TestProxy_Cache_EventSystemUpdatesCacheUseAutoAuthToken(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth so that the event system can run.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	cacheConfig := `
cache {
	cache_static_secrets = true
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s

api_proxy {
  use_auto_auth_token = "force"
}

log_level = "trace"
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
		t.Errorf("stdout: %s", ui.OutputWriter.String())
		t.Errorf("stderr: %s", ui.ErrorWriter.String())
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	secretData2 := map[string]interface{}{
		"bar": "baz",
	}

	// Wait for the event system to successfully connect.
	// This is longer than it needs to be to account for unnatural slowness/avoiding
	// flakiness.
	time.Sleep(5 * time.Second)

	// Mount the KVV2 engine
	err = serverClient.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create kvv2 secret
	_, err = serverClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData)
	if err != nil {
		t.Fatal(err)
	}

	// We use raw requests so we can check the headers for cache hit/miss.
	req := proxyClient.NewRequest(http.MethodGet, "/v1/secret-v2/data/my-secret")
	resp1, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue := resp1.Header.Get("X-Cache")
	require.Equal(t, "MISS", cacheValue)

	// Update the secret using the proxy client
	_, err = proxyClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData2)
	if err != nil {
		t.Fatal(err)
	}

	// Give some time for the event to actually get sent and the cache to be updated.
	// This is longer than it needs to be to account for unnatural slowness/avoiding
	// flakiness.
	time.Sleep(5 * time.Second)

	// We expect this to be a cache hit, with the new value
	resp2, err := proxyClient.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	cacheValue = resp2.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	// Lastly, we check to make sure the actual data we received is
	// as we expect. We must use ParseSecret due to the raw requests.
	secret1, err := api.ParseSecret(resp1.Body)
	require.Nil(t, err)
	data, ok := secret1.Data["data"]
	require.True(t, ok)
	require.Equal(t, secretData, data)

	secret2, err := api.ParseSecret(resp2.Body)
	require.Nil(t, err)
	data2, ok := secret2.Data["data"]
	require.True(t, ok)
	// We expect that the cached value got updated by the event system.
	require.Equal(t, secretData2, data2)

	// Lastly, ensure that a client without a token succeeds
	// at accessing the secret, due to the use_auto_auth_token = "force"
	// option.
	proxyClient.SetToken("")
	req = proxyClient.NewRequest(http.MethodGet, "/v1/secret-v2/data/my-secret")
	resp3, err := proxyClient.RawRequest(req)
	require.Nil(t, err)
	cacheValue = resp3.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	secret3, err := api.ParseSecret(resp3.Body)
	require.Nil(t, err)
	data3, ok := secret3.Data["data"]
	require.True(t, ok)
	// We expect that the cached value got updated by the event system.
	require.Equal(t, secretData2, data3)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_EventSystemPreEventStreamUpdateWorks Tests that the pre-event stream update works
// (i.e. the method preEventStreamUpdate). This test is similar to TestProxy_Cache_EventSystemUpdatesCacheKVV2,
// but with the key difference of not waiting the five seconds for the event subsystem of Proxy to get running
// before it updates the secret, meaning that the event system should update it from the pre-event stream
// update as opposed to receiving the event.
func TestProxy_Cache_EventSystemPreEventStreamUpdateWorks(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	})

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth so that the event system can run.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	cacheConfig := `
cache {
	cache_static_secrets = true
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s
log_level = "trace"
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
		t.Errorf("stdout: %s", ui.OutputWriter.String())
		t.Errorf("stderr: %s", ui.ErrorWriter.String())
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken(serverClient.Token())
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	require.NoError(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	secretData2 := map[string]interface{}{
		"bar": "baz",
	}

	// Mount the KVV2 engine
	err = serverClient.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// Create kvv2 secret
	_, err = serverClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData)
	require.NoError(t, err)

	// We use raw requests so we can check the headers for cache hit/miss.
	req := proxyClient.NewRequest(http.MethodGet, "/v1/secret-v2/data/my-secret")
	resp1, err := proxyClient.RawRequest(req)
	require.NoError(t, err)

	cacheValue := resp1.Header.Get("X-Cache")
	require.Equal(t, "MISS", cacheValue)

	// Update the secret using the proxy client
	_, err = proxyClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData2)
	require.NoError(t, err)

	// Give some time for the event system to run and update the secret as part of the
	// pre-event stream run. Likely, this will be the period in which the event subsystem
	// of proxy actually starts up, so it will have missed the event, but this should still
	// result in an updated cache.
	time.Sleep(5 * time.Second)

	// We expect this to be a cache hit, with the new value
	resp2, err := proxyClient.RawRequest(req)
	require.NoError(t, err)

	cacheValue = resp2.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	// Lastly, we check to make sure the actual data we received is
	// as we expect. We must use ParseSecret due to the raw requests.
	secret1, err := api.ParseSecret(resp1.Body)
	require.NoError(t, err)
	data, ok := secret1.Data["data"]
	require.True(t, ok)
	require.Equal(t, secretData, data)

	secret2, err := api.ParseSecret(resp2.Body)
	require.NoError(t, err)
	data2, ok := secret2.Data["data"]
	require.True(t, ok)
	// We expect that the cached value got updated by the event system.
	require.Equal(t, secretData2, data2)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_StaticSecretPermissionsLost Tests that the cache successfully caches a static secret
// going through the Proxy for a KVV2 secret, and then the calling client loses permissions to the secret,
// so it can no longer access the cache.
func TestProxy_Cache_StaticSecretPermissionsLost(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())
	defer os.Remove(tokenFileName)
	// We need auto-auth so that the event system can run.
	// For ease, we use the token file path with the root token.
	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
}`, tokenFileName)

	// We make the token capability refresh interval one second, for ease of testing
	cacheConfig := `
cache {
	cache_static_secrets = true
	static_secret_token_capability_refresh_interval = "1s"
}
`
	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}
%s
%s
%s
log_level = "trace"
`, serverClient.Address(), cacheConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
		t.Errorf("stdout: %s", ui.OutputWriter.String())
		t.Errorf("stderr: %s", ui.ErrorWriter.String())
	}

	proxyClient, err := api.NewClient(api.DefaultConfig())
	require.Nil(t, err)
	proxyClient.SetMaxRetries(0)
	err = proxyClient.SetAddress("http://" + listenAddr)
	require.Nil(t, err)

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	// Mount the KVV2 engine
	err = serverClient.Sys().Mount("secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.Nil(t, err)

	err = serverClient.Sys().PutPolicy("kv-policy", `
   path "secret-v2/*" {
     capabilities = ["update", "read"]
   }`)
	require.Nil(t, err)

	// Setup a token that we can later revoke:
	renewable := true
	// Set the token's policies to 'default' and nothing else
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies:  []string{"default", "kv-policy"},
		TTL:       "2s",
		Renewable: &renewable,
	}

	secret, err := serverClient.Auth().Token().CreateOrphan(tokenCreateRequest)
	require.Nil(t, err)
	token := secret.Auth.ClientToken
	proxyClient.SetToken(token)

	// Create kvv2 secret
	_, err = serverClient.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData)
	require.Nil(t, err)

	// We use raw requests so we can check the headers for cache hit/miss.
	req := proxyClient.NewRequest(http.MethodGet, "/v1/secret-v2/data/my-secret")
	resp1, err := proxyClient.RawRequest(req)
	require.Nil(t, err)

	cacheValue := resp1.Header.Get("X-Cache")
	require.Equal(t, "MISS", cacheValue)

	// We expect this to be a cache hit, with the new value
	resp2, err := proxyClient.RawRequest(req)
	require.Nil(t, err)

	cacheValue = resp2.Header.Get("X-Cache")
	require.Equal(t, "HIT", cacheValue)

	// Lastly, we check to make sure the actual data we received is
	// as we expect. We must use ParseSecret due to the raw requests.
	secret1, err := api.ParseSecret(resp1.Body)
	if err != nil {
		t.Fatal(err)
	}
	data, ok := secret1.Data["data"]
	require.True(t, ok)
	require.Equal(t, secretData, data)

	secret2, err := api.ParseSecret(resp2.Body)
	if err != nil {
		t.Fatal(err)
	}
	data2, ok := secret2.Data["data"]
	require.True(t, ok)
	// We expect that the cached value got updated by the event system.
	require.Equal(t, secretData, data2)

	// Wait for the token to expire, and for the permissions to be revoked
	// The TTL on the token was 2s, and the capability refresh is every 1s,
	// so this should give us more than enough time!
	time.Sleep(5 * time.Second)
	kvSecret, err := proxyClient.KVv2("secret-v2").Get(context.Background(), "my-secret")
	if err == nil {
		t.Fatalf("expected error, but none found, secret:%v, err:%v", kvSecret, err)
	}
	// Make sure it's a permission denied error
	if !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("expected error on GET to secret after token revocation, secret:%v, err:%v", kvSecret, err)
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_ApiProxy_Retry Tests the retry functionalities of Vault Proxy's API Proxy
func TestProxy_ApiProxy_Retry(t *testing.T) {
	//----------------------------------------------------
	// Start the server and proxy
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	var h handler
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
			LogicalBackends: map[string]logical.Factory{
				"kv": logicalKv.Factory,
			},
		},
		&vault.TestClusterOptions{
			NumCores: 1,
			HandlerFunc: vaulthttp.HandlerFunc(func(properties *vault.HandlerProperties) http.Handler {
				h.props = properties
				h.t = t
				return &h
			}),
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	_, err := serverClient.Logical().Write("secret/foo", map[string]interface{}{
		"bar": "baz",
	})
	if err != nil {
		t.Fatal(err)
	}

	intRef := func(i int) *int {
		return &i
	}
	// start test cases here
	testCases := map[string]struct {
		retries     *int
		expectError bool
	}{
		"none": {
			retries:     intRef(-1),
			expectError: true,
		},
		"one": {
			retries:     intRef(1),
			expectError: true,
		},
		"two": {
			retries:     intRef(2),
			expectError: false,
		},
		"missing": {
			retries:     nil,
			expectError: false,
		},
		"default": {
			retries:     intRef(0),
			expectError: false,
		},
	}

	for tcname, tc := range testCases {
		t.Run(tcname, func(t *testing.T) {
			h.failCount = 2

			cacheConfig := `
cache {
}
`
			listenAddr := generateListenerAddress(t)
			listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

			var retryConf string
			if tc.retries != nil {
				retryConf = fmt.Sprintf("retry { num_retries = %d }", *tc.retries)
			}

			config := fmt.Sprintf(`
vault {
  address = "%s"
  %s
  tls_skip_verify = true
}
%s
%s
`, serverClient.Address(), retryConf, cacheConfig, listenConfig)
			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			_, cmd := testProxyCommand(t, logger)
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				cmd.Run([]string{"-config", configPath})
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			client, err := api.NewClient(api.DefaultConfig())
			if err != nil {
				t.Fatal(err)
			}
			client.SetToken(serverClient.Token())
			client.SetMaxRetries(0)
			err = client.SetAddress("http://" + listenAddr)
			if err != nil {
				t.Fatal(err)
			}
			secret, err := client.Logical().Read("secret/foo")
			switch {
			case (err != nil || secret == nil) && tc.expectError:
			case (err == nil || secret != nil) && !tc.expectError:
			default:
				t.Fatalf("%s expectError=%v error=%v secret=%v", tcname, tc.expectError, err, secret)
			}
			if secret != nil && secret.Data["foo"] != nil {
				val := secret.Data["foo"].(map[string]interface{})
				if !reflect.DeepEqual(val, map[string]interface{}{"bar": "baz"}) {
					t.Fatalf("expected key 'foo' to yield bar=baz, got: %v", val)
				}
			}
			time.Sleep(time.Second)

			close(cmd.ShutdownCh)
			wg.Wait()
		})
	}
}

// TestProxy_Metrics tests that metrics are being properly reported.
func TestProxy_Metrics(t *testing.T) {
	// Start a vault server
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, nil,
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()
	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Create a config file
	listenAddr := generateListenerAddress(t)
	config := fmt.Sprintf(`
cache {}

listener "tcp" {
    address = "%s"
    tls_disable = true
}
`, listenAddr)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	ui, cmd := testProxyCommand(t, logger)
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code := cmd.Run([]string{"-config", configPath})
		if code != 0 {
			t.Errorf("non-zero return code when running proxy: %d", code)
			t.Logf("STDOUT from proxy:\n%s", ui.OutputWriter.String())
			t.Logf("STDERR from proxy:\n%s", ui.ErrorWriter.String())
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	// defer proxy shutdown
	defer func() {
		cmd.ShutdownCh <- struct{}{}
		wg.Wait()
	}()

	conf := api.DefaultConfig()
	conf.Address = "http://" + listenAddr
	proxyClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	req := proxyClient.NewRequest("GET", "/proxy/v1/metrics")
	body := request(t, proxyClient, req, 200)
	keys := []string{}
	for k := range body {
		keys = append(keys, k)
	}
	require.ElementsMatch(t, keys, []string{
		"Counters",
		"Samples",
		"Timestamp",
		"Gauges",
		"Points",
	})
}

// TestProxy_QuitAPI Tests the /proxy/v1/quit API that can be enabled for the proxy.
func TestProxy_QuitAPI(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	err := os.Unsetenv(api.EnvVaultAddress)
	if err != nil {
		t.Fatal(err)
	}

	listenAddr := generateListenerAddress(t)
	listenAddr2 := generateListenerAddress(t)
	config := fmt.Sprintf(`
vault {
  address = "%s"
  tls_skip_verify = true
}

listener "tcp" {
	address = "%s"
	tls_disable = true
}

listener "tcp" {
	address = "%s"
	tls_disable = true
	proxy_api {
		enable_quit = true
	}
}

cache {}
`, serverClient.Address(), listenAddr, listenAddr2)

	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	_, cmd := testProxyCommand(t, nil)
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd.Run([]string{"-config", configPath})
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(serverClient.Token())
	client.SetMaxRetries(0)
	err = client.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	// First try on listener 1 where the API should be disabled.
	resp, err := client.RawRequest(client.NewRequest(http.MethodPost, "/proxy/v1/quit"))
	if err == nil {
		t.Fatalf("expected error")
	}
	if resp != nil && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected %d but got: %d", http.StatusNotFound, resp.StatusCode)
	}

	// Now try on listener 2 where the quit API should be enabled.
	err = client.SetAddress("http://" + listenAddr2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RawRequest(client.NewRequest(http.MethodPost, "/proxy/v1/quit"))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	select {
	case <-cmd.ShutdownCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	wg.Wait()
}

// TestProxy_LogFile_CliOverridesConfig tests that the CLI values
// override the config for log files
func TestProxy_LogFile_CliOverridesConfig(t *testing.T) {
	// Create basic config
	configFile := populateTempFile(t, "proxy-config.hcl", BasicHclConfig)
	cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
	if err != nil {
		t.Fatal("Cannot load config to test update/merge", err)
	}

	// Sanity check that the config value is the current value
	assert.Equal(t, "TMPDIR/juan.log", cfg.LogFile)

	// Initialize the command and parse any flags
	cmd := &ProxyCommand{BaseCommand: &BaseCommand{}}
	f := cmd.Flags()
	// Simulate the flag being specified
	err = f.Parse([]string{"-log-file=/foo/bar/test.log"})
	if err != nil {
		t.Fatal(err)
	}

	// Update the config based on the inputs.
	cmd.applyConfigOverrides(f, cfg)

	assert.NotEqual(t, "TMPDIR/juan.log", cfg.LogFile)
	assert.NotEqual(t, "/squiggle/logs.txt", cfg.LogFile)
	assert.Equal(t, "/foo/bar/test.log", cfg.LogFile)
}

// TestProxy_LogFile_Config tests log file config when loaded from config
func TestProxy_LogFile_Config(t *testing.T) {
	configFile := populateTempFile(t, "proxy-config.hcl", BasicHclConfig)

	cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
	if err != nil {
		t.Fatal("Cannot load config to test update/merge", err)
	}

	// Sanity check that the config value is the current value
	assert.Equal(t, "TMPDIR/juan.log", cfg.LogFile, "sanity check on log config failed")
	assert.Equal(t, 2, cfg.LogRotateMaxFiles)
	assert.Equal(t, 1048576, cfg.LogRotateBytes)

	// Parse the cli flags (but we pass in an empty slice)
	cmd := &ProxyCommand{BaseCommand: &BaseCommand{}}
	f := cmd.Flags()
	err = f.Parse([]string{})
	if err != nil {
		t.Fatal(err)
	}

	// Should change nothing...
	cmd.applyConfigOverrides(f, cfg)

	assert.Equal(t, "TMPDIR/juan.log", cfg.LogFile, "actual config check")
	assert.Equal(t, 2, cfg.LogRotateMaxFiles)
	assert.Equal(t, 1048576, cfg.LogRotateBytes)
}

// TestProxy_Config_NewLogger_Default Tests defaults for log level and
// specifically cmd.newLogger()
func TestProxy_Config_NewLogger_Default(t *testing.T) {
	cmd := &ProxyCommand{BaseCommand: &BaseCommand{}}
	cmd.config = proxyConfig.NewConfig()
	logger, err := cmd.newLogger()

	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.Equal(t, hclog.Info.String(), logger.GetLevel().String())
}

// TestProxy_Config_ReloadLogLevel Tests reloading updates the log
// level as expected.
func TestProxy_Config_ReloadLogLevel(t *testing.T) {
	cmd := &ProxyCommand{BaseCommand: &BaseCommand{}}
	var err error
	tempDir := t.TempDir()

	// Load an initial config
	hcl := strings.ReplaceAll(BasicHclConfig, "TMPDIR", tempDir)
	configFile := populateTempFile(t, "proxy-config.hcl", hcl)
	cmd.config, err = proxyConfig.LoadConfigFile(configFile.Name())
	if err != nil {
		t.Fatal("Cannot load config to test update/merge", err)
	}

	// Tweak the loaded config to make sure we can put log files into a temp dir
	// and systemd log attempts work fine, this would usually happen during Run.
	cmd.logWriter = os.Stdout
	cmd.logger, err = cmd.newLogger()
	if err != nil {
		t.Fatal("logger required for systemd log messages", err)
	}

	// Sanity check
	assert.Equal(t, "warn", cmd.config.LogLevel)

	// Load a new config
	hcl = strings.ReplaceAll(BasicHclConfig2, "TMPDIR", tempDir)
	configFile = populateTempFile(t, "proxy-config.hcl", hcl)
	err = cmd.reloadConfig([]string{configFile.Name()})
	assert.NoError(t, err)
	assert.Equal(t, "debug", cmd.config.LogLevel)
}

// TestProxy_Config_ReloadTls Tests that the TLS certs for the listener are
// correctly reloaded.
func TestProxy_Config_ReloadTls(t *testing.T) {
	var wg sync.WaitGroup
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("unable to get current working directory")
	}
	workingDir := filepath.Join(wd, "/proxy/test-fixtures/reload")
	fooCert := "reload_foo.pem"
	fooKey := "reload_foo.key"

	barCert := "reload_bar.pem"
	barKey := "reload_bar.key"

	reloadCert := "reload_cert.pem"
	reloadKey := "reload_key.pem"
	caPem := "reload_ca.pem"

	tempDir := t.TempDir()

	// Set up initial 'foo' certs
	inBytes, err := os.ReadFile(filepath.Join(workingDir, fooCert))
	if err != nil {
		t.Fatal("unable to read cert required for test", fooCert, err)
	}
	err = os.WriteFile(filepath.Join(tempDir, reloadCert), inBytes, 0o777)
	if err != nil {
		t.Fatal("unable to write temp cert required for test", reloadCert, err)
	}

	inBytes, err = os.ReadFile(filepath.Join(workingDir, fooKey))
	if err != nil {
		t.Fatal("unable to read cert key required for test", fooKey, err)
	}
	err = os.WriteFile(filepath.Join(tempDir, reloadKey), inBytes, 0o777)
	if err != nil {
		t.Fatal("unable to write temp cert key required for test", reloadKey, err)
	}

	inBytes, err = os.ReadFile(filepath.Join(workingDir, caPem))
	if err != nil {
		t.Fatal("unable to read CA pem required for test", caPem, err)
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(inBytes)
	if !ok {
		t.Fatal("not ok when appending CA cert")
	}

	replacedHcl := strings.ReplaceAll(BasicHclConfig, "TMPDIR", tempDir)
	configFile := populateTempFile(t, "proxy-config.hcl", replacedHcl)

	// Set up Proxy
	logger := logging.NewVaultLogger(hclog.Trace)
	ui, cmd := testProxyCommand(t, logger)

	var output string
	var code int
	wg.Add(1)
	args := []string{"-config", configFile.Name()}
	go func() {
		if code = cmd.Run(args); code != 0 {
			output = ui.ErrorWriter.String() + ui.OutputWriter.String()
		}
		wg.Done()
	}()

	testCertificateName := func(cn string) error {
		conn, err := tls.Dial("tcp", "127.0.0.1:8100", &tls.Config{
			RootCAs: certPool,
		})
		if err != nil {
			return err
		}
		defer conn.Close()
		if err = conn.Handshake(); err != nil {
			return err
		}
		servName := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
		if servName != cn {
			return fmt.Errorf("expected %s, got %s", cn, servName)
		}
		return nil
	}

	// Start
	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}

	if err := testCertificateName("foo.example.com"); err != nil {
		t.Fatalf("certificate name didn't check out: %s", err)
	}

	// Swap out certs
	inBytes, err = os.ReadFile(filepath.Join(workingDir, barCert))
	if err != nil {
		t.Fatal("unable to read cert required for test", barCert, err)
	}
	err = os.WriteFile(filepath.Join(tempDir, reloadCert), inBytes, 0o777)
	if err != nil {
		t.Fatal("unable to write temp cert required for test", reloadCert, err)
	}

	inBytes, err = os.ReadFile(filepath.Join(workingDir, barKey))
	if err != nil {
		t.Fatal("unable to read cert key required for test", barKey, err)
	}
	err = os.WriteFile(filepath.Join(tempDir, reloadKey), inBytes, 0o777)
	if err != nil {
		t.Fatal("unable to write temp cert key required for test", reloadKey, err)
	}

	// Reload
	cmd.SighupCh <- struct{}{}
	select {
	case <-cmd.reloadedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}

	if err := testCertificateName("bar.example.com"); err != nil {
		t.Fatalf("certificate name didn't check out: %s", err)
	}

	// Shut down
	cmd.ShutdownCh <- struct{}{}
	wg.Wait()

	if code != 0 {
		t.Fatalf("got a non-zero exit status: %d, stdout/stderr: %s", code, output)
	}
}
