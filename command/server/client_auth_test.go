// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// clientAuthNoHCPTokenHelper prevents commands from reading the user's HCP cache.
type clientAuthNoHCPTokenHelper struct{}

func (clientAuthNoHCPTokenHelper) GetHCPToken(string) (*hcpvlib.HCPToken, error) {
	return nil, nil
}

// TestClientCommands_LoginUserpass_HelperRoundTrip verifies that real userpass
// login stores its issued token and fresh commands cannot reuse setup credentials.
func TestClientCommands_LoginUserpass_HelperRoundTrip(t *testing.T) {
	// Normal helper loading constructs a client from process-wide configuration,
	// so this test and its cases must remain serial to use Setenv safely.
	for _, mount := range []string{"userpass", "employees"} {
		t.Run(mount, func(t *testing.T) {
			for _, key := range []string{
				api.EnvVaultToken, api.EnvVaultAddress, api.EnvVaultAgentAddr,
				api.EnvVaultNamespace, api.EnvVaultHeaders, api.EnvVaultCACert,
				api.EnvVaultCACertBytes, api.EnvVaultCAPath, api.EnvVaultClientCert,
				api.EnvVaultClientKey, api.EnvVaultTLSServerName, api.EnvVaultSkipVerify,
				api.EnvHTTPProxy, api.EnvVaultProxyAddr, api.EnvVaultSRVLookup,
				api.EnvRateLimit, api.EnvVaultDisableRedirects, api.EnvVaultWrapTTL,
				api.EnvVaultMFA, api.EnvVaultClientTimeout,
			} {
				t.Setenv(key, "")
			}
			t.Setenv(api.EnvVaultMaxRetries, "0")
			t.Setenv(base.EnvVaultFormat, "json")
			t.Setenv(base.EnvVaultDetailed, "false")
			t.Setenv(base.EnvVaultCLINoColor, "1")

			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			helper := token.NewTestingTokenHelper()
			setupUI := cli.NewMockUi()
			commands := clientcmd.InitClientCommands(setupUI, setupUI, &base.RunOptions{
				Client: client, TokenHelper: helper, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
			})
			enable, err := commands["auth enable"]()
			require.NoError(t, err, "client auth enable factory should construct a command")
			require.Equal(t, 0, enable.Run([]string{"-path=" + mount, "userpass"}), "setup should enable a real userpass mount")
			write, err := commands["write"]()
			require.NoError(t, err, "client write factory should construct a command")
			require.Equal(t, 0, write.Run([]string{
				"auth/" + mount + "/users/test-user", "password=test-password", "token_policies=default", "token_ttl=1h",
			}), "setup should create a disposable user through the client command")
			require.NoError(t, helper.Store("stale-helper-token"), "seed only the isolated in-memory helper")

			args := []string{
				"login", "-method=userpass", "-format=json",
				"-address=" + client.Address(), "-ca-cert=" + cluster.CACertPEMFile,
			}
			if mount != "userpass" {
				args = append(args, "-path="+mount)
			}
			args = append(args, "username=test-user", "password=test-password")
			var stdout, stderr bytes.Buffer
			code := clientcmd.RunCustom(args, &base.RunOptions{
				BinaryName: "vault-client", Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
				TokenHelper: helper, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
			})
			require.Equal(t, 0, code, "registered userpass login should authenticate against Vault")
			require.Zero(t, stderr.Len(), "login should not report an unknown method or authentication error")
			var login api.Secret
			require.NoError(t, json.Unmarshal(stdout.Bytes(), &login), "login should print JSON token metadata")
			require.NotNil(t, login.Auth, "login output should include authentication metadata")
			require.NotEmpty(t, login.Auth.Accessor, "issued token should have an accessor")
			require.Equal(t, 3600, login.Auth.LeaseDuration, "issued token should have the user's configured TTL")
			require.True(t, login.Auth.Renewable, "userpass should issue a renewable token")
			require.Contains(t, login.Auth.Policies, "default", "issued token should have the user's policy")
			require.Equal(t, "test-user", login.Auth.Metadata["username"], "metadata should identify the authenticated user")
			storedToken, err := helper.Get()
			require.NoError(t, err, "read the configured helper after login")
			require.True(t, storedToken != "" && storedToken == login.Auth.ClientToken, "helper should store exactly the token printed by login")
			require.True(t, storedToken != client.Token() && storedToken != "stale-helper-token", "login must replace stale credentials, not persist the setup token")

			var lookupOut, lookupErr bytes.Buffer
			code = clientcmd.RunCustom([]string{
				"token", "lookup", "-format=json", "-address=" + client.Address(), "-ca-cert=" + cluster.CACertPEMFile,
			}, &base.RunOptions{
				BinaryName: "vault-client", Stdin: strings.NewReader(""), Stdout: &lookupOut, Stderr: &lookupErr,
				TokenHelper: helper, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
			})
			require.Equal(t, 0, code, "a fresh command with no injected client or explicit token must load the helper")
			require.Zero(t, lookupErr.Len(), "helper-authenticated lookup should not report an error")
			var lookup struct {
				Data struct {
					ID   string `json:"id"`
					Path string `json:"path"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal(lookupOut.Bytes(), &lookup), "lookup should print token metadata")
			require.True(t, lookup.Data.ID == storedToken, "lookup must identify the userpass token, not a setup token")
			require.Equal(t, "auth/"+mount+"/login/test-user", lookup.Data.Path, "lookup should identify the actual login path")

			require.NoError(t, helper.Erase(), "erase only the isolated helper")
			var erasedOut, erasedErr bytes.Buffer
			code = clientcmd.RunCustom([]string{
				"token", "lookup", "-format=json", "-address=" + client.Address(), "-ca-cert=" + cluster.CACertPEMFile,
			}, &base.RunOptions{
				BinaryName: "vault-client", Stdin: strings.NewReader(""), Stdout: &erasedOut, Stderr: &erasedErr,
				TokenHelper: helper, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
			})
			require.Equal(t, 2, code, "an erased helper must not be masked by a cached client or setup token")
			require.Zero(t, erasedOut.Len(), "unauthenticated lookup must not print token metadata")
			require.Contains(t, erasedErr.String(), "Error looking up token", "erasing the helper should cause an authentication failure")
		})
	}
}

// TestClientCommands_LoginUserpass_HelperPreservation verifies that the registered
// handler preserves an existing helper token on bad credentials or explicit no-store.
func TestClientCommands_LoginUserpass_HelperPreservation(t *testing.T) {
	// Cluster construction reads process-wide API settings before client injection,
	// so this test and its cases must remain serial to use Setenv safely.
	for _, tc := range []struct {
		name     string
		password string
		noStore  bool
		code     int
	}{
		{name: "wrong_password", password: "wrong-password", code: 2},
		{name: "no_store", password: "test-password", noStore: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, key := range []string{
				api.EnvVaultToken, api.EnvVaultAddress, api.EnvVaultAgentAddr,
				api.EnvVaultNamespace, api.EnvVaultHeaders, api.EnvVaultCACert,
				api.EnvVaultCACertBytes, api.EnvVaultCAPath, api.EnvVaultClientCert,
				api.EnvVaultClientKey, api.EnvVaultTLSServerName, api.EnvVaultSkipVerify,
				api.EnvHTTPProxy, api.EnvVaultProxyAddr, api.EnvVaultSRVLookup,
				api.EnvRateLimit, api.EnvVaultDisableRedirects, api.EnvVaultWrapTTL,
				api.EnvVaultMFA, api.EnvVaultClientTimeout,
			} {
				t.Setenv(key, "")
			}
			t.Setenv(api.EnvVaultMaxRetries, "0")

			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			require.NoError(t, client.Sys().EnableAuthWithOptionsWithContext(t.Context(), "userpass", &api.EnableAuthOptions{
				Type: "userpass",
			}), "enable a real userpass mount")
			_, err := client.Logical().WriteWithContext(t.Context(), "auth/userpass/users/test-user", map[string]any{
				"password": "test-password", "token_policies": "default",
			})
			require.NoError(t, err, "create a disposable user")
			loginClient, err := client.Clone()
			require.NoError(t, err, "use a separate client so login cannot modify the setup client")
			helper := token.NewTestingTokenHelper()
			require.NoError(t, helper.Store("existing-helper-token"), "seed the isolated helper")
			args := []string{"login", "-method=userpass", "-format=json"}
			if tc.noStore {
				args = append(args, "-no-store")
			}
			args = append(args, "username=test-user", "password="+tc.password)
			var stdout, stderr bytes.Buffer
			code := clientcmd.RunCustom(args, &base.RunOptions{
				Client: loginClient, BinaryName: "vault-client", Stdin: strings.NewReader(""),
				Stdout: &stdout, Stderr: &stderr, TokenHelper: helper, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
			})
			require.Equal(t, tc.code, code, "userpass should preserve the failure/no-store exit contract")
			storedToken, err := helper.Get()
			require.NoError(t, err, "read the helper after login")
			require.True(t, storedToken == "existing-helper-token", "failed login and no-store must preserve the helper")
			require.NotContains(t, stderr.String(), "Unknown auth method", "userpass must reach its registered handler")
			if tc.noStore {
				var login api.Secret
				require.NoError(t, json.Unmarshal(stdout.Bytes(), &login), "no-store should still print token metadata")
				require.NotNil(t, login.Auth, "successful no-store login should include authentication metadata")
				require.True(t, login.Auth.ClientToken != "" && login.Auth.ClientToken != storedToken && login.Auth.ClientToken != client.Token(), "no-store should issue a new user token without storing it")
				require.Contains(t, stderr.String(), "not stored in token helper", "no-store should explain how to use the returned token")
			} else {
				require.Zero(t, stdout.Len(), "bad credentials must not print successful token metadata")
				require.Contains(t, stderr.String(), "Error authenticating", "Vault should reject the wrong password")
			}
		})
	}
}

// TestClientCommands_AuthHelp_MountedMethods verifies userpass and AWS type/path
// help through the client factories with real mounts, without cloud authentication.
func TestClientCommands_AuthHelp_MountedMethods(t *testing.T) {
	// Cluster construction reads process-wide API settings before client injection,
	// so this test and its cases must remain serial to use Setenv safely.
	for _, method := range []string{"userpass", "aws"} {
		t.Run(method, func(t *testing.T) {
			for _, key := range []string{
				api.EnvVaultToken, api.EnvVaultAddress, api.EnvVaultAgentAddr,
				api.EnvVaultNamespace, api.EnvVaultHeaders, api.EnvVaultCACert,
				api.EnvVaultCACertBytes, api.EnvVaultCAPath, api.EnvVaultClientCert,
				api.EnvVaultClientKey, api.EnvVaultTLSServerName, api.EnvVaultSkipVerify,
				api.EnvHTTPProxy, api.EnvVaultProxyAddr, api.EnvVaultSRVLookup,
				api.EnvRateLimit, api.EnvVaultDisableRedirects, api.EnvVaultWrapTTL,
				api.EnvVaultMFA, api.EnvVaultClientTimeout,
			} {
				t.Setenv(key, "")
			}
			t.Setenv(api.EnvVaultMaxRetries, "0")

			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			require.NoError(t, client.Sys().EnableAuthWithOptionsWithContext(t.Context(), method, &api.EnableAuthOptions{
				Type: method,
			}), "enable the real auth backend without configuring cloud credentials")
			for _, argument := range []string{method, method + "/"} {
				ui := cli.NewMockUi()
				commands := clientcmd.InitClientCommands(ui, ui, &base.RunOptions{
					Client: client, TokenHelper: token.NewTestingTokenHelper(), HCPTokenHelper: clientAuthNoHCPTokenHelper{},
				})
				help, err := commands["auth help"]()
				require.NoError(t, err, "client auth help factory should construct a command")
				require.Equal(t, 0, help.Run([]string{argument}), "type and mounted-path help should resolve %s", argument)
				require.Empty(t, ui.ErrorWriter.String(), "method help must not report the historical missing-method error")
				require.Contains(t, ui.OutputWriter.String(), "login -method="+method, "help should identify the actual auth method")
			}
		})
	}
}
