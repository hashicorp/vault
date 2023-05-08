// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agent"
	"github.com/hashicorp/vault/helper/useragent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
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
		Logger: logger,
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

	inf, err := os.CreateTemp("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	os.Remove(in)
	t.Logf("input: %s", in)

	sink1f, err := os.CreateTemp("", "sink1.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink1 := sink1f.Name()
	sink1f.Close()
	os.Remove(sink1)
	t.Logf("sink1: %s", sink1)

	sink2f, err := os.CreateTemp("", "sink2.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink2 := sink2f.Name()
	sink2f.Close()
	os.Remove(sink2)
	t.Logf("sink2: %s", sink2)

	conff, err := os.CreateTemp("", "conf.jwt.test.")
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
		Logger: logger,
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
	agentClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	agentClient.SetToken("")
	err = agentClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the token to be sent to syncs and be available to be used
	time.Sleep(5 * time.Second)

	req = agentClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	body = request(t, agentClient, req, 200)

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
				h.userAgentToCheckFor = useragent.AgentProxyStringWithProxiedUserAgent(userAgentForProxiedClient)
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

	// Start the agent
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

	// Start the agent
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

	agentClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	agentClient.AddHeader("User-Agent", userAgentForProxiedClient)
	agentClient.SetToken(serverClient.Token())
	agentClient.SetMaxRetries(0)
	err = agentClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	_, err = agentClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_Cache_DynamicSecret Tests that the cache successfully caches a dynamic secret
// going through the Proxy,
func TestProxy_Cache_DynamicSecret(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
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

	agentClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	agentClient.SetToken(serverClient.Token())
	agentClient.SetMaxRetries(0)
	err = agentClient.SetAddress("http://" + listenAddr)
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
	secret, err := agentClient.Auth().Token().CreateOrphan(tokenCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Auth == nil {
		t.Fatalf("secret not as expected: %v", secret)
	}

	token := secret.Auth.ClientToken

	secret, err = agentClient.Auth().Token().CreateOrphan(tokenCreateRequest)
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
