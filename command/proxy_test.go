// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
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

	sinkFileName1 := makeTempFile(t, "sink-file", "")
	sinkFileName2 := makeTempFile(t, "sink-file", "")

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

	config = fmt.Sprintf(config, exitAfterAuthTemplText, in, sinkFileName1, sinkFileName2)
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

	sink1Bytes, err := os.ReadFile(sinkFileName1)
	if err != nil {
		t.Fatal(err)
	}
	if len(sink1Bytes) == 0 {
		t.Fatal("got no output from sink 1")
	}

	sink2Bytes, err := os.ReadFile(sinkFileName2)
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

// TestProxy_NoTriggerAutoAuth_BadPolicy tests that auto auth is not re-triggered
// if Proxy uses a token with incorrect policy access.
func TestProxy_NoTriggerAutoAuth_BadPolicy(t *testing.T) {
	proxyLogger := logging.NewVaultLogger(hclog.Trace)
	vaultLogger := logging.NewVaultLogger(hclog.Info)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: vaulthttp.Handler,
		Logger:      vaultLogger,
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Add a secret to the KV engine
	_, err := serverClient.Logical().Write("secret/foo", map[string]interface{}{"user": "something"})
	require.NoError(t, err)

	// Create kv read policy
	noKvAccess := `path "secret/*" {
capabilities = ["deny"]
}`
	err = serverClient.Sys().PutPolicy("noKvAccess", noKvAccess)
	require.NoError(t, err)

	// Create a token with that policy
	opts := &api.TokenCreateRequest{Policies: []string{"noKvAccess"}}
	tokenResp, err := serverClient.Auth().Token().Create(opts)
	require.NoError(t, err)
	firstToken := tokenResp.Auth.ClientToken

	// Create token file
	tokenFileName := makeTempFile(t, "token-file", firstToken)

	sinkFileName := makeTempFile(t, "sink-file", "")

	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
	sink "file" {
		config = {
			path = "%s"
		}
	}
}`, tokenFileName, sinkFileName)

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
	 use_auto_auth_token = "force"
	}
	%s
	%s
	`, serverClient.Address(), listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start proxy
	_, cmd := testProxyCommand(t, proxyLogger)
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
		t.Fatalf("timeout")
	}

	// Validate that the auto-auth token has been correctly attained
	// and works for LookupSelf
	conf := api.DefaultConfig()
	conf.Address = "http://" + listenAddr
	proxyClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}

	proxyClient.SetToken("")
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	// Wait for re-triggered auto auth to write new token to sink
	waitForFile := func(prevModTime time.Time) time.Time {
		ticker := time.Tick(100 * time.Millisecond)
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-ticker:
			case <-timeout:
				return prevModTime
			}
			modTime, err := os.Stat(sinkFileName)
			require.NoError(t, err)
			if modTime.ModTime().After(prevModTime) {
				return modTime.ModTime()
			}
		}
	}

	// Wait for the token to be sent to syncs and be available to be used
	initialModTime := waitForFile(time.Time{})
	req := proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	_ = request(t, proxyClient, req, 200)

	// Write a new token to the token file
	newTokenResp, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	secondToken := newTokenResp.Auth.ClientToken
	err = os.WriteFile(tokenFileName, []byte(secondToken), 0o600)
	require.NoError(t, err)

	// Make a request to a path that the token does not have access to
	req = proxyClient.NewRequest("GET", "/v1/secret/foo")
	_, err = proxyClient.RawRequest(req)
	require.Error(t, err)
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())
	require.NotContains(t, err.Error(), logical.ErrInvalidToken.Error())

	// Sleep for a bit to ensure that auto auth is not re-triggered
	newModTime := waitForFile(initialModTime)
	if newModTime.After(initialModTime) {
		t.Fatal("auto auth was incorrectly re-triggered")
	}

	// Read from the sink file and verify that the token has not changed
	newToken, err := os.ReadFile(sinkFileName)
	require.Equal(t, firstToken, string(newToken))

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_NoTriggerAutoAuth_ProxyTokenNotAutoAuth tests that auto auth is not re-triggered
// if Proxy uses a token that is not equal to the auto auth token
func TestProxy_NoTriggerAutoAuth_ProxyTokenNotAutoAuth(t *testing.T) {
	proxyLogger := logging.NewVaultLogger(hclog.Info)
	cluster := minimal.NewTestSoloCluster(t, nil)

	serverClient := cluster.Cores[0].Client

	// Create a token
	tokenResp, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	firstToken := tokenResp.Auth.ClientToken

	// Create token file
	tokenFileName := makeTempFile(t, "token-file", firstToken)

	sinkFileName := makeTempFile(t, "sink-file", "")

	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
	sink "file" {
		config = {
			path = "%s"
		}
	}
}`, tokenFileName, sinkFileName)

	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
	listener "tcp" {
	 address = "%s"
	 tls_disable = true
	}
	`, listenAddr)

	// Do not use the auto auth token if a token is provided with the proxy client
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

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start proxy
	_, cmd := testProxyCommand(t, proxyLogger)
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
		t.Fatalf("timeout")
	}

	// Validate that the auto-auth token has been correctly attained
	// and works for LookupSelf
	conf := api.DefaultConfig()
	conf.Address = "http://" + listenAddr
	proxyClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}

	proxyClient.SetToken(firstToken)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	// Wait for re-triggered auto auth to write new token to sink
	waitForFile := func(prevModTime time.Time) time.Time {
		ticker := time.Tick(100 * time.Millisecond)
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-ticker:
			case <-timeout:
				return prevModTime
			}
			modTime, err := os.Stat(sinkFileName)
			require.NoError(t, err)
			if modTime.ModTime().After(prevModTime) {
				return modTime.ModTime()
			}
		}
	}

	// Wait for the token is available to be used
	createTime := waitForFile(time.Time{})
	require.NoError(t, err)
	_, err = serverClient.Auth().Token().LookupSelf()
	require.NoError(t, err)

	// Revoke token
	err = serverClient.Auth().Token().RevokeOrphan(firstToken)
	require.NoError(t, err)

	// Write a new token to the token file
	newTokenResp, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	secondToken := newTokenResp.Auth.ClientToken
	err = os.WriteFile(tokenFileName, []byte(secondToken), 0o600)
	require.NoError(t, err)

	// Proxy uses revoked token to make request and should result in an error
	proxyClient.SetToken("random token")
	_, err = proxyClient.Auth().Token().LookupSelf()
	require.Error(t, err)

	// Wait to see if the sink file is modified
	newModTime := waitForFile(createTime)
	if newModTime.After(createTime) {
		t.Fatal("auto auth was incorrectly re-triggered")
	}

	// Read from the sink and verify that the token has not changed
	newToken, err := os.ReadFile(sinkFileName)
	require.Equal(t, firstToken, string(newToken))

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_ReTriggerAutoAuth_ForceAutoAuthToken tests that auto auth is re-triggered
// if Proxy always forcibly uses the auto auth token
func TestProxy_ReTriggerAutoAuth_ForceAutoAuthToken(t *testing.T) {
	proxyLogger := logging.NewVaultLogger(hclog.Trace)
	vaultLogger := logging.NewVaultLogger(hclog.Info)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: vaulthttp.Handler,
		Logger:      vaultLogger,
	})
	cluster.Start()
	defer cluster.Cleanup()

	serverClient := cluster.Cores[0].Client

	// Create a token
	tokenResp, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	firstToken := tokenResp.Auth.ClientToken

	// Create token file
	tokenFileName := makeTempFile(t, "token-file", firstToken)

	sinkFileName := makeTempFile(t, "sink-file", "")

	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }
	sink "file" {
		config = {
			path = "%s"
		}
	}
}`, tokenFileName, sinkFileName)

	listenAddr := generateListenerAddress(t)
	listenConfig := fmt.Sprintf(`
	listener "tcp" {
	 address = "%s"
	 tls_disable = true
	}
	`, listenAddr)

	// Do not use the auto auth token if a token is provided with the proxy client
	config := fmt.Sprintf(`
	vault {
	 address = "%s"
	 tls_skip_verify = true
	}
	api_proxy {
	 use_auto_auth_token = "force"
	}
	%s
	%s
	`, serverClient.Address(), listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start proxy
	_, cmd := testProxyCommand(t, proxyLogger)
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
		t.Fatalf("timeout")
	}

	// Validate that the auto-auth token has been correctly attained
	// and works for LookupSelf
	conf := api.DefaultConfig()
	conf.Address = "http://" + listenAddr
	proxyClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}

	proxyClient.SetToken(firstToken)
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	// Wait for re-triggered auto auth to write new token to sink
	waitForFile := func(prevModTime time.Time) time.Time {
		ticker := time.Tick(100 * time.Millisecond)
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-ticker:
			case <-timeout:
				return prevModTime
			}
			modTime, err := os.Stat(sinkFileName)
			require.NoError(t, err)
			if modTime.ModTime().After(prevModTime) {
				return modTime.ModTime()
			}
		}
	}

	// Wait for the token is available to be used
	createTime := waitForFile(time.Time{})
	require.NoError(t, err)
	req := proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	_, err = proxyClient.RawRequest(req)
	require.NoError(t, err)

	// Revoke token
	req = serverClient.NewRequest("PUT", "/v1/auth/token/revoke")
	req.BodyBytes = []byte(fmt.Sprintf(`{
	  "token": "%s"
	}`, firstToken))
	_ = request(t, serverClient, req, 204)

	// Create new token
	newTokenResp, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	secondToken := newTokenResp.Auth.ClientToken

	// Proxy uses the same token in the token file to make a request, which should result in error
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	_, err = proxyClient.RawRequest(req)
	require.Error(t, err)

	// Write a new token to the token file so that auto auth can write new token to sink
	err = os.WriteFile(tokenFileName, []byte(secondToken), 0o600)
	require.NoError(t, err)

	// Wait to see if that the sink file is modified
	waitForFile(createTime)

	// Read from the sink and verify that the sink contains the new token
	newToken, err := os.ReadFile(sinkFileName)
	require.Equal(t, secondToken, string(newToken))

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_ReTriggerAutoAuth_ProxyIsAutoAuthToken tests that auto auth is re-triggered
// the proxy client uses a token that is equal to the auto auth token
func TestProxy_ReTriggerAutoAuth_ProxyIsAutoAuthToken(t *testing.T) {
	proxyLogger := logging.NewVaultLogger(hclog.Trace)
	vaultLogger := logging.NewVaultLogger(hclog.Info)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: vaulthttp.Handler,
		Logger:      vaultLogger,
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
	  "token_max_ttl": "4m",
	  "token_num_uses": "10",
	  "token_ttl": "4m",
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

	sinkFileName := makeTempFile(t, "sink-file", "")

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
}`, roleIDPath, secretIDPath, sinkFileName)

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

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start proxy
	_, cmd := testProxyCommand(t, proxyLogger)
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
		t.Fatal(err)
	}

	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	// Wait for re-triggered auto auth to write new token to sink
	waitForFile := func(prevModTime time.Time) {
		ticker := time.Tick(100 * time.Millisecond)
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-ticker:
			case <-timeout:
				t.Fatal("timed out waiting for re-triggered auto auth to complete")
			}
			modTime, err := os.Stat(sinkFileName)
			require.NoError(t, err)
			if modTime.ModTime().After(prevModTime) {
				return
			}
		}
	}

	// Wait for the token to be sent to syncs and be available to be used
	waitForFile(time.Time{})
	oldToken, err := os.ReadFile(sinkFileName)
	require.NoError(t, err)
	prevModTime, err := os.Stat(sinkFileName)
	require.NoError(t, err)

	// Set proxy token
	proxyClient.SetToken(string(oldToken))

	// Make request using proxy client to test that token is valid
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	body = request(t, proxyClient, req, 200)

	// Revoke token
	req = serverClient.NewRequest("PUT", "/v1/auth/token/revoke")
	req.BodyBytes = []byte(fmt.Sprintf(`{
	  "token": "%s"
	}`, oldToken))
	body = request(t, serverClient, req, 204)

	// Proxy uses revoked token to make request and should result in an error
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	_, err = proxyClient.RawRequest(req)
	require.Error(t, err)

	// Wait for new token to be written and available to use
	waitForFile(prevModTime.ModTime())

	// Verify new token is not equal to the old token
	newToken, err := os.ReadFile(sinkFileName)
	require.NoError(t, err)
	require.NotEqual(t, string(newToken), string(oldToken))

	// Verify that proxy no longer fails when making a request with the new token
	proxyClient.SetToken(string(newToken))
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	body = request(t, proxyClient, req, 200)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_ReTriggerAutoAuth_RevokedToken tests that auto auth is re-triggered
// when Proxy uses a revoked auto auth token to make a request
func TestProxy_ReTriggerAutoAuth_RevokedToken(t *testing.T) {
	proxyLogger := logging.NewVaultLogger(hclog.Trace)
	vaultLogger := logging.NewVaultLogger(hclog.Info)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: vaulthttp.Handler,
		Logger:      vaultLogger,
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
	  "token_max_ttl": "4m",
	  "token_num_uses": "10",
	  "token_ttl": "4m",
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

	sinkFileName := makeTempFile(t, "sink-file", "")
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
}`, roleIDPath, secretIDPath, sinkFileName)

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
  use_auto_auth_token = "force"
}
%s
%s
`, serverClient.Address(), listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)

	// Unset the environment variable so that proxy picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start proxy
	_, cmd := testProxyCommand(t, proxyLogger)
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
		t.Fatal(err)
	}

	proxyClient.SetToken("")
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	// Wait for re-triggered auto auth to write new token to sink
	waitForFile := func(prevModTime time.Time) {
		ticker := time.Tick(100 * time.Millisecond)
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-ticker:
			case <-timeout:
				t.Fatal("timed out waiting for re-triggered auto auth to complete")
			}
			modTime, err := os.Stat(sinkFileName)
			require.NoError(t, err)
			if modTime.ModTime().After(prevModTime) {
				return
			}
		}
	}

	// Wait for the token to be sent to syncs and be available to be used
	waitForFile(time.Time{})
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	body = request(t, proxyClient, req, 200)

	oldToken, err := os.ReadFile(sinkFileName)
	require.NoError(t, err)
	prevModTime, err := os.Stat(sinkFileName)
	require.NoError(t, err)

	// Revoke token
	req = serverClient.NewRequest("PUT", "/v1/auth/token/revoke")
	req.BodyBytes = []byte(fmt.Sprintf(`{
	  "token": "%s"
	}`, oldToken))
	body = request(t, serverClient, req, 204)

	// Proxy uses revoked token to make request and should result in an error
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	_, err = proxyClient.RawRequest(req)
	require.Error(t, err)

	// Wait for new token to be written and available to use
	waitForFile(prevModTime.ModTime())

	// Verify new token is not equal to the old token
	newToken, err := os.ReadFile(sinkFileName)
	require.NoError(t, err)
	require.NotEqual(t, string(newToken), string(oldToken))

	// Verify that proxy no longer fails when making a request
	req = proxyClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	body = request(t, proxyClient, req, 200)

	close(cmd.ShutdownCh)
	wg.Wait()
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

	sinkFileName := makeTempFile(t, "sink-file", "")
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
}`, roleIDPath, secretIDPath, sinkFileName)

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

// TestProxy_NoAutoAuthTokenIfNotConfigured tests that Proxy will not use the auto-auth token
// unless configured to.
func TestProxy_NoAutoAuthTokenIfNotConfigured(t *testing.T) {
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

	// Create token file
	tokenFileName := makeTempFile(t, "token-file", serverClient.Token())

	sinkFileName := makeTempFile(t, "sink-file", "")

	autoAuthConfig := fmt.Sprintf(`
auto_auth {
    method {
		type = "token_file"
        config = {
            token_file_path = "%s"
        }
    }

	sink "file" {
		config = {
			path = "%s"
		}
	}
}`, tokenFileName, sinkFileName)

	apiProxyConfig := `
api_proxy {
	use_auto_auth_token = false
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
`, serverClient.Address(), apiProxyConfig, listenConfig, autoAuthConfig)
	configPath := makeTempFile(t, "config.hcl", config)

	// Start proxy
	ui, cmd := testProxyCommand(t, logger)
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

	proxyClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	proxyClient.SetToken("")
	err = proxyClient.SetAddress("http://" + listenAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the sink to be populated.
	// Realistically won't be this long, but keeping it long just in case, for CI.
	time.Sleep(10 * time.Second)

	secret, err := proxyClient.Auth().Token().CreateOrphan(&api.TokenCreateRequest{
		Policies: []string{"default"},
		TTL:      "30m",
	})
	if secret != nil || err == nil {
		t.Fatal("expected this to fail, since without a token you should not be able to make a token")
	}

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestProxy_ApiProxy_Retry Tests the retry functionalities of Vault Proxy's API Proxy
func TestProxy_ApiProxy_Retry(t *testing.T) {
	// ----------------------------------------------------
	// Start the server and proxy
	// ----------------------------------------------------
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

// TestProxy_EnvVar_Overrides tests that environment variables are properly
// parsed and override defaults.
func TestProxy_EnvVar_Overrides(t *testing.T) {
	configFile := populateTempFile(t, "proxy-config.hcl", BasicHclConfig)

	cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
	if err != nil {
		t.Fatal("Cannot load config to test update/merge", err)
	}

	assert.Equal(t, false, cfg.Vault.TLSSkipVerify)

	t.Setenv("VAULT_SKIP_VERIFY", "true")
	// Parse the cli flags (but we pass in an empty slice)
	cmd := &ProxyCommand{BaseCommand: &BaseCommand{}}
	f := cmd.Flags()
	err = f.Parse([]string{})
	if err != nil {
		t.Fatal(err)
	}

	cmd.applyConfigOverrides(f, cfg)
	assert.Equal(t, true, cfg.Vault.TLSSkipVerify)

	t.Setenv("VAULT_SKIP_VERIFY", "false")

	cmd.applyConfigOverrides(f, cfg)
	assert.Equal(t, false, cfg.Vault.TLSSkipVerify)
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
