// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
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
	agentConfig "github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/useragent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	BasicHclConfig = `
log_file = "TMPDIR/juan.log"
log_level="warn"
log_rotate_max_files=2
log_rotate_bytes=1048576
vault {
	address = "http://127.0.0.1:8200"
	retry {
		num_retries = 5
	}
}

listener "tcp" {
	address = "127.0.0.1:8100"
	tls_disable = false
	tls_cert_file = "TMPDIR/reload_cert.pem"
  	tls_key_file  = "TMPDIR/reload_key.pem"
}`
	BasicHclConfig2 = `
log_file = "TMPDIR/juan.log"
log_level="debug"
log_rotate_max_files=-1
log_rotate_bytes=1048576
vault {
	address = "http://127.0.0.1:8200"
	retry {
		num_retries = 5
	}
}

listener "tcp" {
	address = "127.0.0.1:8100"
	tls_disable = false
	tls_cert_file = "TMPDIR/reload_cert.pem"
  	tls_key_file  = "TMPDIR/reload_key.pem"
}`
)

func testAgentCommand(tb testing.TB, logger hclog.Logger) (*cli.MockUi, *AgentCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AgentCommand{
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

func TestAgent_ExitAfterAuth(t *testing.T) {
	t.Run("via_config", func(t *testing.T) {
		testAgentExitAfterAuth(t, false)
	})

	t.Run("via_flag", func(t *testing.T) {
		testAgentExitAfterAuth(t, true)
	})
}

func testAgentExitAfterAuth(t *testing.T, viaFlag bool) {
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
		ui, cmd := testAgentCommand(t, logger)
		cmd.client = client

		args := []string{"-config", conf}
		if viaFlag {
			args = append(args, "-exit-after-auth")
		}

		code := cmd.Run(args)
		if code != 0 {
			t.Errorf("expected %d to be %d", code, 0)
			t.Logf("output from agent:\n%s", ui.OutputWriter.String())
			t.Logf("error from agent:\n%s", ui.ErrorWriter.String())
		}
		close(doneCh)
	}()

	select {
	case <-doneCh:
		break
	case <-time.After(1 * time.Minute):
		t.Fatal("timeout reached while waiting for agent to exit")
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

func TestAgent_RequireRequestHeader(t *testing.T) {
	// newApiClient creates an *api.Client.
	newApiClient := func(addr string, includeVaultRequestHeader bool) *api.Client {
		conf := api.DefaultConfig()
		conf.Address = addr
		cli, err := api.NewClient(conf)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		h := cli.Headers()
		val, ok := h[consts.RequestHeaderName]
		if !ok || !reflect.DeepEqual(val, []string{"true"}) {
			t.Fatalf("invalid %s header", consts.RequestHeaderName)
		}
		if !includeVaultRequestHeader {
			delete(h, consts.RequestHeaderName)
			cli.SetHeaders(h)
		}

		return cli
	}

	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------

	// Start a vault server
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
		},
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()
	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Enable the approle auth method
	roleIDPath, secretIDPath := setupAppRole(t, serverClient)

	// Create a config file
	config := `
auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
        }
    }
}

cache {
    use_auto_auth_token = true
}

listener "tcp" {
    address = "%s"
    tls_disable = true
}
listener "tcp" {
    address = "%s"
    tls_disable = true
    require_request_header = false
}
listener "tcp" {
    address = "%s"
    tls_disable = true
    require_request_header = true
}
`
	listenAddr1 := generateListenerAddress(t)
	listenAddr2 := generateListenerAddress(t)
	listenAddr3 := generateListenerAddress(t)
	config = fmt.Sprintf(
		config,
		roleIDPath,
		secretIDPath,
		listenAddr1,
		listenAddr2,
		listenAddr3,
	)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	var output string
	var code int
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code = cmd.Run([]string{"-config", configPath})
		if code != 0 {
			output = ui.ErrorWriter.String() + ui.OutputWriter.String()
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}

	// defer agent shutdown
	defer func() {
		cmd.ShutdownCh <- struct{}{}
		wg.Wait()
		if code != 0 {
			t.Fatalf("got a non-zero exit status: %d, stdout/stderr: %s", code, output)
		}
	}()

	//----------------------------------------------------
	// Perform the tests
	//----------------------------------------------------

	// Test against a listener configuration that omits
	// 'require_request_header', with the header missing from the request.
	agentClient := newApiClient("http://"+listenAddr1, false)
	req := agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)

	// Test against a listener configuration that sets 'require_request_header'
	// to 'false', with the header missing from the request.
	agentClient = newApiClient("http://"+listenAddr2, false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with the header missing from the request.
	agentClient = newApiClient("http://"+listenAddr3, false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	resp, err := agentClient.RawRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if resp.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected status code %d, not %d", http.StatusPreconditionFailed, resp.StatusCode)
	}

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with an invalid header present in the request.
	agentClient = newApiClient("http://"+listenAddr3, false)
	h := agentClient.Headers()
	h[consts.RequestHeaderName] = []string{"bogus"}
	agentClient.SetHeaders(h)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	resp, err = agentClient.RawRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if resp.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected status code %d, not %d", http.StatusPreconditionFailed, resp.StatusCode)
	}

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with the proper header present in the request.
	agentClient = newApiClient("http://"+listenAddr3, true)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)
}

// TestAgent_RequireAutoAuthWithForce ensures that the client exits with a
// non-zero code if configured to force the use of an auto-auth token without
// configuring the auto_auth block
func TestAgent_RequireAutoAuthWithForce(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	// Create a config file
	config := fmt.Sprintf(`
cache {
    use_auto_auth_token = "force"
}

listener "tcp" {
    address = "%s"
    tls_disable = true
}
`, generateListenerAddress(t))

	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	code := cmd.Run([]string{"-config", configPath})
	if code == 0 {
		t.Errorf("expected error code, but got 0: %d", code)
		t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
		t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
	}
}

// TestAgent_Template_UserAgent Validates that the User-Agent sent to Vault
// as part of Templating requests is correct. Uses the custom handler
// userAgentHandler struct defined in this test package, so that Vault validates the
// User-Agent on requests sent by Agent.
func TestAgent_Template_UserAgent(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	var h userAgentHandler
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
			HandlerFunc: vaulthttp.HandlerFunc(
				func(properties *vault.HandlerProperties) http.Handler {
					h.props = properties
					h.userAgentToCheckFor = useragent.AgentTemplatingString()
					h.pathToCheck = "/v1/secret/data"
					h.requestMethodToCheck = "GET"
					h.t = t
					return &h
				}),
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, serverClient.Address())

	roleIDPath, secretIDPath := setupAppRoleAndKVMounts(t, serverClient)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := os.MkdirTemp("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)
	// create temp dir for this test run
	tmpDir, err := os.MkdirTemp(tmpDirRoot, "TestAgent_Template_UserAgent")
	if err != nil {
		t.Fatal(err)
	}

	// make some template files
	var templatePaths []string
	fileName := filepath.Join(tmpDir, "render_0.tmpl")
	if err := os.WriteFile(fileName, []byte(templateContents(0)), 0o600); err != nil {
		t.Fatal(err)
	}
	templatePaths = append(templatePaths, fileName)

	// build up the template config to be added to the Agent config.hcl file
	var templateConfigStrings []string
	for i, t := range templatePaths {
		index := fmt.Sprintf("render_%d.json", i)
		s := fmt.Sprintf(templateConfigString, t, tmpDir, index)
		templateConfigStrings = append(templateConfigStrings, s)
	}

	// Create a config file
	config := `
vault {
  address = "%s"
	tls_skip_verify = true
}

auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
            remove_secret_id_file_after_reading = false
        }
    }
}

%s
`

	// flatten the template configs
	templateConfig := strings.Join(templateConfigStrings, " ")

	config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, templateConfig)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code := cmd.Run([]string{"-config", configPath})
		if code != 0 {
			t.Errorf("non-zero return code when running agent: %d", code)
			t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
			t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	// We need to shut down the Agent command
	defer func() {
		cmd.ShutdownCh <- struct{}{}
		wg.Wait()
	}()

	verify := func(suffix string) {
		t.Helper()
		// We need to poll for a bit to give Agent time to render the
		// templates. Without this, the test will attempt to read
		// the temp dir before Agent has had time to render and will
		// likely fail the test
		tick := time.Tick(1 * time.Second)
		timeout := time.After(10 * time.Second)
		var err error
		for {
			select {
			case <-timeout:
				t.Fatalf("timed out waiting for templates to render, last error: %v", err)
			case <-tick:
			}
			// Check for files rendered in the directory and break
			// early for shutdown if we do have all the files
			// rendered

			//----------------------------------------------------
			// Perform the tests
			//----------------------------------------------------

			if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != len(templatePaths) {
				err = fmt.Errorf("expected (%d) templates, got (%d)", len(templatePaths), numFiles)
				continue
			}

			for i := range templatePaths {
				fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.json", i))
				var c []byte
				c, err = os.ReadFile(fileName)
				if err != nil {
					continue
				}
				if string(c) != templateRendered(i)+suffix {
					err = fmt.Errorf("expected=%q, got=%q", templateRendered(i)+suffix, string(c))
					continue
				}
			}
			return
		}
	}

	verify("")

	fileName = filepath.Join(tmpDir, "render_0.tmpl")
	if err := os.WriteFile(fileName, []byte(templateContents(0)+"{}"), 0o600); err != nil {
		t.Fatal(err)
	}

	verify("{}")
}

// TestAgent_Template tests rendering templates
func TestAgent_Template_Basic(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
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
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, serverClient.Address())

	roleIDPath, secretIDPath := setupAppRoleAndKVMounts(t, serverClient)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := os.MkdirTemp("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// start test cases here
	testCases := map[string]struct {
		templateCount int
		exitAfterAuth bool
	}{
		"one": {
			templateCount: 1,
		},
		"one_with_exit": {
			templateCount: 1,
			exitAfterAuth: true,
		},
		"many": {
			templateCount: 15,
		},
		"many_with_exit": {
			templateCount: 13,
			exitAfterAuth: true,
		},
	}

	for tcname, tc := range testCases {
		t.Run(tcname, func(t *testing.T) {
			// create temp dir for this test run
			tmpDir, err := os.MkdirTemp(tmpDirRoot, tcname)
			if err != nil {
				t.Fatal(err)
			}

			// make some template files
			var templatePaths []string
			for i := 0; i < tc.templateCount; i++ {
				fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.tmpl", i))
				if err := os.WriteFile(fileName, []byte(templateContents(i)), 0o600); err != nil {
					t.Fatal(err)
				}
				templatePaths = append(templatePaths, fileName)
			}

			// build up the template config to be added to the Agent config.hcl file
			var templateConfigStrings []string
			for i, t := range templatePaths {
				index := fmt.Sprintf("render_%d.json", i)
				s := fmt.Sprintf(templateConfigString, t, tmpDir, index)
				templateConfigStrings = append(templateConfigStrings, s)
			}

			// Create a config file
			config := `
vault {
  address = "%s"
	tls_skip_verify = true
}

auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
            remove_secret_id_file_after_reading = false
        }
    }
}

%s

%s
`

			// conditionally set the exit_after_auth flag
			exitAfterAuth := ""
			if tc.exitAfterAuth {
				exitAfterAuth = "exit_after_auth = true"
			}

			// flatten the template configs
			templateConfig := strings.Join(templateConfigStrings, " ")

			config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, templateConfig, exitAfterAuth)
			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			ui, cmd := testAgentCommand(t, logger)
			cmd.client = serverClient
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				code := cmd.Run([]string{"-config", configPath})
				if code != 0 {
					t.Errorf("non-zero return code when running agent: %d", code)
					t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
					t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
				}
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			// if using exit_after_auth, then the command will have returned at the
			// end and no longer be running. If we are not using exit_after_auth, then
			// we need to shut down the command
			if !tc.exitAfterAuth {
				defer func() {
					cmd.ShutdownCh <- struct{}{}
					wg.Wait()
				}()
			}

			verify := func(suffix string) {
				t.Helper()
				// We need to poll for a bit to give Agent time to render the
				// templates. Without this, the test will attempt to read
				// the temp dir before Agent has had time to render and will
				// likely fail the test
				tick := time.Tick(1 * time.Second)
				timeout := time.After(10 * time.Second)
				var err error
				for {
					select {
					case <-timeout:
						t.Fatalf("timed out waiting for templates to render, last error: %v", err)
					case <-tick:
					}
					// Check for files rendered in the directory and break
					// early for shutdown if we do have all the files
					// rendered

					//----------------------------------------------------
					// Perform the tests
					//----------------------------------------------------

					if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != len(templatePaths) {
						err = fmt.Errorf("expected (%d) templates, got (%d)", len(templatePaths), numFiles)
						continue
					}

					for i := range templatePaths {
						fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.json", i))
						var c []byte
						c, err = os.ReadFile(fileName)
						if err != nil {
							continue
						}
						if string(c) != templateRendered(i)+suffix {
							err = fmt.Errorf("expected=%q, got=%q", templateRendered(i)+suffix, string(c))
							continue
						}
					}
					return
				}
			}

			verify("")

			for i := 0; i < tc.templateCount; i++ {
				fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.tmpl", i))
				if err := os.WriteFile(fileName, []byte(templateContents(i)+"{}"), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			verify("{}")
		})
	}
}

func setupAppRole(t *testing.T, serverClient *api.Client) (string, string) {
	t.Helper()
	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// Create a named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role")
	req.BodyBytes = []byte(`{
	  "token_ttl": "5m",
		"token_policies":"default,myapp-read",
		"policies":"default,myapp-read"
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
	t.Cleanup(func() {
		os.Remove(roleIDPath)
		os.Remove(secretIDPath)
	})

	return roleIDPath, secretIDPath
}

func setupAppRoleAndKVMounts(t *testing.T, serverClient *api.Client) (string, string) {
	roleIDPath, secretIDPath := setupAppRole(t, serverClient)

	// give test-role permissions to read the kv secret
	req := serverClient.NewRequest("PUT", "/v1/sys/policy/myapp-read")
	req.BodyBytes = []byte(`{
	  "policy": "path \"secret/*\" { capabilities = [\"read\", \"list\"] }"
	}`)
	request(t, serverClient, req, 204)

	// setup the kv secrets
	req = serverClient.NewRequest("POST", "/v1/sys/mounts/secret/tune")
	req.BodyBytes = []byte(`{
	"options": {"version": "2"}
	}`)
	request(t, serverClient, req, 200)

	// Secret: myapp
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "bar",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// Secret: myapp2
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp2")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// Secret: otherapp
	req = serverClient.NewRequest("POST", "/v1/secret/data/otherapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap",
			"cert": "something"
    }
	}`)
	request(t, serverClient, req, 200)

	return roleIDPath, secretIDPath
}

// TestAgent_Template_VaultClientFromEnv tests that Vault Agent can read in its
// required `vault` client details from environment variables instead of config.
func TestAgent_Template_VaultClientFromEnv(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
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
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	roleIDPath, secretIDPath := setupAppRoleAndKVMounts(t, serverClient)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := os.MkdirTemp("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	vaultAddr := "https://" + cluster.Cores[0].Listeners[0].Address.String()
	testCases := map[string]struct {
		env map[string]string
	}{
		"VAULT_ADDR and VAULT_CACERT": {
			env: map[string]string{
				api.EnvVaultAddress: vaultAddr,
				api.EnvVaultCACert:  cluster.CACertPEMFile,
			},
		},
		"VAULT_ADDR and VAULT_CACERT_BYTES": {
			env: map[string]string{
				api.EnvVaultAddress:     vaultAddr,
				api.EnvVaultCACertBytes: string(cluster.CACertPEM),
			},
		},
	}

	for tcname, tc := range testCases {
		t.Run(tcname, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			tmpDir := t.TempDir()

			// Make a template.
			templateFile := filepath.Join(tmpDir, "render.tmpl")
			if err := os.WriteFile(templateFile, []byte(templateContents(0)), 0o600); err != nil {
				t.Fatal(err)
			}

			// build up the template config to be added to the Agent config.hcl file
			targetFile := filepath.Join(tmpDir, "render.json")
			templateConfig := fmt.Sprintf(`
template {
    source      = "%s"
    destination = "%s"
}
			`, templateFile, targetFile)

			// Create a config file
			config := `
auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
            remove_secret_id_file_after_reading = false
        }
    }
}

%s
`

			config = fmt.Sprintf(config, roleIDPath, secretIDPath, templateConfig)
			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			ui, cmd := testAgentCommand(t, logger)
			cmd.client = serverClient
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				code := cmd.Run([]string{"-config", configPath})
				if code != 0 {
					t.Errorf("non-zero return code when running agent: %d", code)
					t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
					t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
				}
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			defer func() {
				cmd.ShutdownCh <- struct{}{}
				wg.Wait()
			}()

			// We need to poll for a bit to give Agent time to render the
			// templates. Without this this, the test will attempt to read
			// the temp dir before Agent has had time to render and will
			// likely fail the test
			tick := time.Tick(1 * time.Second)
			timeout := time.After(10 * time.Second)
			for {
				select {
				case <-timeout:
					t.Fatalf("timed out waiting for templates to render, last error: %v", err)
				case <-tick:
				}

				contents, err := os.ReadFile(targetFile)
				if err != nil {
					// If the file simply doesn't exist, continue waiting for
					// the template rendering to complete.
					if os.IsNotExist(err) {
						continue
					}
					t.Fatal(err)
				}

				if string(contents) != templateRendered(0) {
					t.Fatalf("expected=%q, got=%q", templateRendered(0), string(contents))
				}

				// Success! Break out of the retry loop.
				break
			}
		})
	}
}

func testListFiles(t *testing.T, dir, extension string) int {
	t.Helper()

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for _, f := range files {
		if filepath.Ext(f.Name()) == extension {
			count++
		}
	}

	return count
}

// TestAgent_Template_ExitCounter tests that Vault Agent correctly renders all
// templates before exiting when the configuration uses exit_after_auth. This is
// similar to TestAgent_Template_Basic, but differs by using a consistent number
// of secrets from multiple sources, where as the basic test could possibly
// generate a random number of secrets, but all using the same source. This test
// reproduces https://github.com/hashicorp/vault/issues/7883
func TestAgent_Template_ExitCounter(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
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
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, serverClient.Address())

	roleIDPath, secretIDPath := setupAppRoleAndKVMounts(t, serverClient)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := os.MkdirTemp("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// create temp dir for this test run
	tmpDir, err := os.MkdirTemp(tmpDirRoot, "agent-test")
	if err != nil {
		t.Fatal(err)
	}

	// Create a config file
	config := `
vault {
  address = "%s"
  tls_skip_verify = true
}

auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
            remove_secret_id_file_after_reading = false
        }
    }
}

template {
    contents = "{{ with secret \"secret/myapp\" }}{{ range $k, $v := .Data.data }}{{ $v }}{{ end }}{{ end }}"
    destination = "%s/render-pass.txt"
}

template {
    contents = "{{ with secret \"secret/myapp2\" }}{{ .Data.data.username}}{{ end }}"
    destination = "%s/render-user.txt"
}

template {
    contents = <<EOF
{{ with secret "secret/otherapp"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
{{ .Data.data.cert }}
}
{{ end }}
EOF
    destination = "%s/render-other.txt"
}

exit_after_auth = true
`

	config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, tmpDir, tmpDir, tmpDir)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code := cmd.Run([]string{"-config", configPath})
		if code != 0 {
			t.Errorf("non-zero return code when running agent: %d", code)
			t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
			t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	wg.Wait()

	//----------------------------------------------------
	// Perform the tests
	//----------------------------------------------------

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected (%d) templates, got (%d)", 3, len(files))
	}
}

// a slice of template options
var templates = []string{
	`{{- with secret "secret/otherapp"}}{"secret": "other",
{{- if .Data.data.username}}"username":"{{ .Data.data.username}}",{{- end }}
{{- if .Data.data.password }}"password":"{{ .Data.data.password }}"{{- end }}}
{{- end }}`,
	`{{- with secret "secret/myapp"}}{"secret": "myapp",
{{- if .Data.data.username}}"username":"{{ .Data.data.username}}",{{- end }}
{{- if .Data.data.password }}"password":"{{ .Data.data.password }}"{{- end }}}
{{- end }}`,
	`{{- with secret "secret/myapp"}}{"secret": "myapp",
{{- if .Data.data.password }}"password":"{{ .Data.data.password }}"{{- end }}}
{{- end }}`,
}

var rendered = []string{
	`{"secret": "other","username":"barstuff","password":"zap"}`,
	`{"secret": "myapp","username":"bar","password":"zap"}`,
	`{"secret": "myapp","password":"zap"}`,
}

// templateContents returns a template from the above templates slice. Each
// invocation with incrementing seed will return "the next" template, and loop.
// This ensures as we use multiple templates that we have a increasing number of
// sources before we reuse a template.
func templateContents(seed int) string {
	index := seed % len(templates)
	return templates[index]
}

func templateRendered(seed int) string {
	index := seed % len(templates)
	return rendered[index]
}

var templateConfigString = `
template {
  source      = "%s"
  destination = "%s/%s"
}
`

// request issues HTTP requests.
func request(t *testing.T, client *api.Client, req *api.Request, expectedStatusCode int) map[string]interface{} {
	t.Helper()
	resp, err := client.RawRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("expected status code %d, not %d", expectedStatusCode, resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(bytes) == 0 {
		return nil
	}

	var body map[string]interface{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return body
}

// makeTempFile creates a temp file and populates it.
func makeTempFile(t *testing.T, name, contents string) string {
	t.Helper()
	f, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.WriteString(contents)
	f.Close()
	return path
}

func populateTempFile(t *testing.T, name, contents string) *os.File {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.WriteString(contents)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	return file
}

// handler makes 500 errors happen for reads on /v1/secret.
// Definitely not thread-safe, do not use t.Parallel with this.
type handler struct {
	props     *vault.HandlerProperties
	failCount int
	t         *testing.T
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/v1/secret") {
		if h.failCount > 0 {
			h.failCount--
			h.t.Logf("%s failing GET request on %s, failures left: %d", time.Now(), req.URL.Path, h.failCount)
			resp.WriteHeader(500)
			return
		}
		h.t.Logf("passing GET request on %s", req.URL.Path)
	}
	vaulthttp.Handler.Handler(h.props).ServeHTTP(resp, req)
}

// userAgentHandler makes it easy to test the User-Agent header received
// by Vault
type userAgentHandler struct {
	props                *vault.HandlerProperties
	failCount            int
	userAgentToCheckFor  string
	pathToCheck          string
	requestMethodToCheck string
	t                    *testing.T
}

func (h *userAgentHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == h.requestMethodToCheck && strings.Contains(req.RequestURI, h.pathToCheck) {
		userAgent := req.UserAgent()
		if !(userAgent == h.userAgentToCheckFor) {
			h.t.Fatalf("User-Agent string not as expected. Expected to find %s, got %s", h.userAgentToCheckFor, userAgent)
		}
	}
	vaulthttp.Handler.Handler(h.props).ServeHTTP(w, req)
}

// TestAgent_Template_Retry verifies that the template server retries requests
// based on retry configuration.
func TestAgent_Template_Retry(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
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
			HandlerFunc: vaulthttp.HandlerFunc(
				func(properties *vault.HandlerProperties) http.Handler {
					h.props = properties
					h.t = t
					return &h
				}),
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	methodConf, cleanup := prepAgentApproleKV(t, serverClient)
	defer cleanup()

	err := serverClient.Sys().TuneMount("secret", api.MountConfigInput{
		Options: map[string]string{
			"version": "2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = serverClient.Logical().Write("secret/data/otherapp", map[string]interface{}{
		"data": map[string]interface{}{
			"username": "barstuff",
			"password": "zap",
			"cert":     "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := os.MkdirTemp("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

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
			// We fail the first 6 times.  The consul-template code creates
			// a Vault client with MaxRetries=2, so for every consul-template
			// retry configured, it will in practice make up to 3 requests.
			// Thus if consul-template is configured with "one" retry, it will
			// fail given our failCount, but if configured with "two" retries,
			// they will consume our 6th failure, and on the "third (from its
			// perspective) attempt, it will succeed.
			h.failCount = 6

			// create temp dir for this test run
			tmpDir, err := os.MkdirTemp(tmpDirRoot, tcname)
			if err != nil {
				t.Fatal(err)
			}

			// make some template files
			templatePath := filepath.Join(tmpDir, "render_0.tmpl")
			if err := os.WriteFile(templatePath, []byte(templateContents(0)), 0o600); err != nil {
				t.Fatal(err)
			}
			templateConfig := fmt.Sprintf(templateConfigString, templatePath, tmpDir, "render_0.json")

			var retryConf string
			if tc.retries != nil {
				retryConf = fmt.Sprintf("retry { num_retries = %d }", *tc.retries)
			}

			config := fmt.Sprintf(`
%s
vault {
  address = "%s"
  %s
  tls_skip_verify = true
}
%s
template_config {
  exit_on_retry_failure = true
}
`, methodConf, serverClient.Address(), retryConf, templateConfig)

			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			_, cmd := testAgentCommand(t, logger)
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			var code int
			go func() {
				code = cmd.Run([]string{"-config", configPath})
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			verify := func() error {
				t.Helper()
				// We need to poll for a bit to give Agent time to render the
				// templates. Without this this, the test will attempt to read
				// the temp dir before Agent has had time to render and will
				// likely fail the test
				tick := time.Tick(1 * time.Second)
				timeout := time.After(15 * time.Second)
				var err error
				for {
					select {
					case <-timeout:
						return fmt.Errorf("timed out waiting for templates to render, last error: %v", err)
					case <-tick:
					}
					// Check for files rendered in the directory and break
					// early for shutdown if we do have all the files
					// rendered

					//----------------------------------------------------
					// Perform the tests
					//----------------------------------------------------

					if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != 1 {
						err = fmt.Errorf("expected 1 template, got (%d)", numFiles)
						continue
					}

					fileName := filepath.Join(tmpDir, "render_0.json")
					var c []byte
					c, err = os.ReadFile(fileName)
					if err != nil {
						continue
					}
					if string(c) != templateRendered(0) {
						err = fmt.Errorf("expected=%q, got=%q", templateRendered(0), string(c))
						continue
					}
					return nil
				}
			}

			err = verify()
			close(cmd.ShutdownCh)
			wg.Wait()

			switch {
			case (code != 0 || err != nil) && tc.expectError:
			case code == 0 && err == nil && !tc.expectError:
			default:
				t.Fatalf("%s expectError=%v error=%v code=%d", tcname, tc.expectError, err, code)
			}
		})
	}
}

// prepAgentApproleKV configures a Vault instance for approle authentication,
// such that the resulting token will have global permissions across /kv
// and /secret mounts.  Returns the auto_auth config stanza to setup an Agent
// to connect using approle.
func prepAgentApproleKV(t *testing.T, client *api.Client) (string, func()) {
	t.Helper()

	policyAutoAuthAppRole := `
path "/kv/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
path "/secret/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
`
	// Add an kv-admin policy
	if err := client.Sys().PutPolicy("test-autoauth", policyAutoAuthAppRole); err != nil {
		t.Fatal(err)
	}

	// Enable approle
	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test1", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "1h",
		"token_max_ttl":  "2h",
		"policies":       []string{"test-autoauth"},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Write("auth/approle/role/test1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID := resp.Data["secret_id"].(string)
	secretIDFile := makeTempFile(t, "secret_id.txt", secretID+"\n")

	resp, err = client.Logical().Read("auth/approle/role/test1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID := resp.Data["role_id"].(string)
	roleIDFile := makeTempFile(t, "role_id.txt", roleID+"\n")

	config := fmt.Sprintf(`
auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
            remove_secret_id_file_after_reading = false
        }
    }
}
`, roleIDFile, secretIDFile)

	cleanup := func() {
		_ = os.Remove(roleIDFile)
		_ = os.Remove(secretIDFile)
	}
	return config, cleanup
}

// TestAgent_AutoAuth_UserAgent tests that the User-Agent sent
// to Vault by Vault Agent is correct when performing Auto-Auth.
// Uses the custom handler userAgentHandler (defined above) so
// that Vault validates the User-Agent on requests sent by Agent.
func TestAgent_AutoAuth_UserAgent(t *testing.T) {
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
				h.userAgentToCheckFor = useragent.AgentAutoAuthString()
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
	roleIDPath, secretIDPath := setupAppRole(t, serverClient)

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

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	// Start the agent
	_, cmd := testAgentCommand(t, logger)
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

	req := agentClient.NewRequest("GET", "/v1/auth/token/lookup-self")
	request(t, agentClient, req, 200)

	close(cmd.ShutdownCh)
	wg.Wait()
}

// TestAgent_APIProxyWithoutCache_UserAgent tests that the User-Agent sent
// to Vault by Vault Agent is correct using the API proxy without
// the cache configured. Uses the custom handler
// userAgentHandler struct defined in this test package, so that Vault validates the
// User-Agent on requests sent by Agent.
func TestAgent_APIProxyWithoutCache_UserAgent(t *testing.T) {
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

	// Unset the environment variable so that agent picks up the right test
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
	_, cmd := testAgentCommand(t, logger)
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

// TestAgent_APIProxyWithCache_UserAgent tests that the User-Agent sent
// to Vault by Vault Agent is correct using the API proxy with
// the cache configured.  Uses the custom handler
// userAgentHandler struct defined in this test package, so that Vault validates the
// User-Agent on requests sent by Agent.
func TestAgent_APIProxyWithCache_UserAgent(t *testing.T) {
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

	// Unset the environment variable so that agent picks up the right test
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
	_, cmd := testAgentCommand(t, logger)
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

func TestAgent_Cache_DynamicSecret(t *testing.T) {
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

	// Start the agent
	_, cmd := testAgentCommand(t, logger)
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

func TestAgent_ApiProxy_Retry(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
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

	// Unset the environment variable so that agent picks up the right test
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

			// Start the agent
			_, cmd := testAgentCommand(t, logger)
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

func TestAgent_TemplateConfig_ExitOnRetryFailure(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
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
			NumCores:    1,
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	autoAuthConfig, cleanup := prepAgentApproleKV(t, serverClient)
	defer cleanup()

	err := serverClient.Sys().TuneMount("secret", api.MountConfigInput{
		Options: map[string]string{
			"version": "2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = serverClient.Logical().Write("secret/data/otherapp", map[string]interface{}{
		"data": map[string]interface{}{
			"username": "barstuff",
			"password": "zap",
			"cert":     "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := os.MkdirTemp("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// Note that missing key is different from a non-existent secret. A missing
	// key (2xx response with missing keys in the response map) can still yield
	// a successful render unless error_on_missing_key is specified, whereas a
	// missing secret (4xx response) always results in an error.
	missingKeyTemplateContent := `{{- with secret "secret/otherapp"}}{"secret": "other",
{{- if .Data.data.foo}}"foo":"{{ .Data.data.foo}}"{{- end }}}
{{- end }}`
	missingKeyTemplateRender := `{"secret": "other",}`

	badTemplateContent := `{{- with secret "secret/non-existent"}}{"secret": "other",
{{- if .Data.data.foo}}"foo":"{{ .Data.data.foo}}"{{- end }}}
{{- end }}`

	testCases := map[string]struct {
		exitOnRetryFailure        *bool
		templateContents          string
		expectTemplateRender      string
		templateErrorOnMissingKey bool
		expectError               bool
		expectExitFromError       bool
	}{
		"true, no template error": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          templateContents(0),
			expectTemplateRender:      templateRendered(0),
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"true, with non-existent secret": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          badTemplateContent,
			expectTemplateRender:      "",
			templateErrorOnMissingKey: false,
			expectError:               true,
			expectExitFromError:       true,
		},
		"true, with missing key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      missingKeyTemplateRender,
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"true, with missing key, with error_on_missing_key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      "",
			templateErrorOnMissingKey: true,
			expectError:               true,
			expectExitFromError:       true,
		},
		"false, no template error": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          templateContents(0),
			expectTemplateRender:      templateRendered(0),
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"false, with non-existent secret": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          badTemplateContent,
			expectTemplateRender:      "",
			templateErrorOnMissingKey: false,
			expectError:               true,
			expectExitFromError:       false,
		},
		"false, with missing key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      missingKeyTemplateRender,
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"false, with missing key, with error_on_missing_key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      missingKeyTemplateRender,
			templateErrorOnMissingKey: true,
			expectError:               true,
			expectExitFromError:       false,
		},
		"missing": {
			exitOnRetryFailure:        nil,
			templateContents:          templateContents(0),
			expectTemplateRender:      templateRendered(0),
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
	}

	for tcName, tc := range testCases {
		t.Run(tcName, func(t *testing.T) {
			// create temp dir for this test run
			tmpDir, err := os.MkdirTemp(tmpDirRoot, tcName)
			if err != nil {
				t.Fatal(err)
			}

			listenAddr := generateListenerAddress(t)
			listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

			var exitOnRetryFailure string
			if tc.exitOnRetryFailure != nil {
				exitOnRetryFailure = fmt.Sprintf("exit_on_retry_failure = %t", *tc.exitOnRetryFailure)
			}
			templateConfig := fmt.Sprintf(`
template_config = {
	%s
}
`, exitOnRetryFailure)

			template := fmt.Sprintf(`
template {
	contents = <<EOF
%s
EOF
	destination = "%s/render_0.json"
	error_on_missing_key = %t
}
`, tc.templateContents, tmpDir, tc.templateErrorOnMissingKey)

			config := fmt.Sprintf(`
# auto-auth stanza
%s

vault {
	address = "%s"
	tls_skip_verify = true
	retry {
		num_retries = 3
	}
}

# listener stanza
%s

# template_config stanza
%s

# template stanza
%s
`, autoAuthConfig, serverClient.Address(), listenConfig, templateConfig, template)

			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			ui, cmd := testAgentCommand(t, logger)
			cmd.startedCh = make(chan struct{})

			// Channel to let verify() know to stop early if agent
			// has exited
			cmdRunDoneCh := make(chan struct{})
			var exitedEarly bool

			wg := &sync.WaitGroup{}
			wg.Add(1)
			var code int
			go func() {
				code = cmd.Run([]string{"-config", configPath})
				close(cmdRunDoneCh)
				wg.Done()
			}()

			verify := func() error {
				t.Helper()
				// We need to poll for a bit to give Agent time to render the
				// templates. Without this this, the test will attempt to read
				// the temp dir before Agent has had time to render and will
				// likely fail the test
				tick := time.Tick(1 * time.Second)
				timeout := time.After(15 * time.Second)
				var err error
				for {
					select {
					case <-cmdRunDoneCh:
						exitedEarly = true
						return nil
					case <-timeout:
						return fmt.Errorf("timed out waiting for templates to render, last error: %w", err)
					case <-tick:
					}
					// Check for files rendered in the directory and break
					// early for shutdown if we do have all the files
					// rendered

					//----------------------------------------------------
					// Perform the tests
					//----------------------------------------------------

					if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != 1 {
						err = fmt.Errorf("expected 1 template, got (%d)", numFiles)
						continue
					}

					fileName := filepath.Join(tmpDir, "render_0.json")
					var c []byte
					c, err = os.ReadFile(fileName)
					if err != nil {
						continue
					}
					if strings.TrimSpace(string(c)) != tc.expectTemplateRender {
						err = fmt.Errorf("expected=%q, got=%q", tc.expectTemplateRender, strings.TrimSpace(string(c)))
						continue
					}
					return nil
				}
			}

			err = verify()
			close(cmd.ShutdownCh)
			wg.Wait()

			switch {
			case (code != 0 || err != nil) && tc.expectError:
				if exitedEarly != tc.expectExitFromError {
					t.Fatalf("expected program exit due to error to be '%t', got '%t'", tc.expectExitFromError, exitedEarly)
				}
			case code == 0 && err == nil && !tc.expectError:
				if exitedEarly {
					t.Fatalf("did not expect program to exit before verify completes")
				}
			default:
				if code != 0 {
					t.Logf("output from agent:\n%s", ui.OutputWriter.String())
					t.Logf("error from agent:\n%s", ui.ErrorWriter.String())
				}
				t.Fatalf("expectError=%v error=%v code=%d", tc.expectError, err, code)
			}
		})
	}
}

func TestAgent_Metrics(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------

	// Start a vault server
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

	// Start the agent
	ui, cmd := testAgentCommand(t, logging.NewVaultLogger(hclog.Trace))
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	var output string
	var code int
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code = cmd.Run([]string{"-config", configPath})
		if code != 0 {
			output = ui.ErrorWriter.String() + ui.OutputWriter.String()
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	// defer agent shutdown
	defer func() {
		cmd.ShutdownCh <- struct{}{}
		wg.Wait()
		if code != 0 {
			t.Fatalf("got a non-zero exit status: %d, stdout/stderr: %s", code, output)
		}
	}()

	conf := api.DefaultConfig()
	conf.Address = "http://" + listenAddr
	agentClient, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	req := agentClient.NewRequest("GET", "/agent/v1/metrics")
	body := request(t, agentClient, req, 200)
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

func TestAgent_Quit(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	cluster := minimal.NewTestSoloCluster(t, nil)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
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
	agent_api {
		enable_quit = true
	}
}

cache {}
`, serverClient.Address(), listenAddr, listenAddr2)

	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	_, cmd := testAgentCommand(t, nil)
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
	resp, err := client.RawRequest(client.NewRequest(http.MethodPost, "/agent/v1/quit"))
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

	_, err = client.RawRequest(client.NewRequest(http.MethodPost, "/agent/v1/quit"))
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

func TestAgent_LogFile_CliOverridesConfig(t *testing.T) {
	// Create basic config
	configFile := populateTempFile(t, "agent-config.hcl", BasicHclConfig)
	cfg, err := agentConfig.LoadConfigFile(configFile.Name())
	if err != nil {
		t.Fatal("Cannot load config to test update/merge", err)
	}

	// Sanity check that the config value is the current value
	assert.Equal(t, "TMPDIR/juan.log", cfg.LogFile)

	// Initialize the command and parse any flags
	cmd := &AgentCommand{BaseCommand: &BaseCommand{}}
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

func TestAgent_LogFile_Config(t *testing.T) {
	configFile := populateTempFile(t, "agent-config.hcl", BasicHclConfig)

	cfg, err := agentConfig.LoadConfigFile(configFile.Name())
	if err != nil {
		t.Fatal("Cannot load config to test update/merge", err)
	}

	// Sanity check that the config value is the current value
	assert.Equal(t, "TMPDIR/juan.log", cfg.LogFile, "sanity check on log config failed")
	assert.Equal(t, 2, cfg.LogRotateMaxFiles)
	assert.Equal(t, 1048576, cfg.LogRotateBytes)

	// Parse the cli flags (but we pass in an empty slice)
	cmd := &AgentCommand{BaseCommand: &BaseCommand{}}
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

func TestAgent_Config_NewLogger_Default(t *testing.T) {
	cmd := &AgentCommand{BaseCommand: &BaseCommand{}}
	cmd.config = agentConfig.NewConfig()
	logger, err := cmd.newLogger()

	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.Equal(t, hclog.Info.String(), logger.GetLevel().String())
}

func TestAgent_Config_ReloadLogLevel(t *testing.T) {
	cmd := &AgentCommand{BaseCommand: &BaseCommand{}}
	var err error
	tempDir := t.TempDir()

	// Load an initial config
	hcl := strings.ReplaceAll(BasicHclConfig, "TMPDIR", tempDir)
	configFile := populateTempFile(t, "agent-config.hcl", hcl)
	cmd.config, err = agentConfig.LoadConfigFile(configFile.Name())
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
	configFile = populateTempFile(t, "agent-config.hcl", hcl)
	err = cmd.reloadConfig([]string{configFile.Name()})
	assert.NoError(t, err)
	assert.Equal(t, "debug", cmd.config.LogLevel)
}

func TestAgent_Config_ReloadTls(t *testing.T) {
	var wg sync.WaitGroup
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("unable to get current working directory")
	}
	workingDir := filepath.Join(wd, "/agent/test-fixtures/reload")
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
	configFile := populateTempFile(t, "agent-config.hcl", replacedHcl)

	// Set up Agent/cmd
	logger := logging.NewVaultLogger(hclog.Trace)
	ui, cmd := testAgentCommand(t, logger)

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

// TestAgent_NonTLSListener_SIGHUP tests giving a SIGHUP signal to a listener
// without a TLS configuration. Prior to fixing GitHub issue #19480, this
// would cause a panic.
func TestAgent_NonTLSListener_SIGHUP(t *testing.T) {
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
	ui, cmd := testAgentCommand(t, logger)

	cmd.startedCh = make(chan struct{})

	var output string
	var code int
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if code = cmd.Run([]string{"-config", configPath}); code != 0 {
			output = ui.ErrorWriter.String() + ui.OutputWriter.String()
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	// Reload
	cmd.SighupCh <- struct{}{}
	select {
	case <-cmd.reloadedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}

	close(cmd.ShutdownCh)
	wg.Wait()

	if code != 0 {
		t.Fatalf("got a non-zero exit status: %d, stdout/stderr: %s", code, output)
	}
}

// TestAgent_Logging_ConsulTemplate attempts to ensure two things about Vault Agent logs:
// 1. When -log-format command line arg is set to JSON, it is honored as the output format
// for messages generated from within the consul-template library.
// 2. When -log-file command line arg is supplied, a file receives all log messages
// generated by the consul-template library (they don't just go to stdout/stderr).
// Should prevent a regression of: https://github.com/hashicorp/vault/issues/21109
func TestAgent_Logging_ConsulTemplate(t *testing.T) {
	const (
		runnerLogMessage = "(runner) creating new runner (dry: false, once: false)"
	)

	// Configure a Vault server so Agent can successfully communicate and render its templates
	cluster := minimal.NewTestSoloCluster(t, nil)
	apiClient := cluster.Cores[0].Client
	t.Setenv(api.EnvVaultAddress, apiClient.Address())
	tempDir := t.TempDir()
	roleIDPath, secretIDPath := setupAppRoleAndKVMounts(t, apiClient)

	// Create relevant configs for Vault Agent (config, template config)
	templateSrc := filepath.Join(tempDir, "render_1.tmpl")
	err := os.WriteFile(templateSrc, []byte(templateContents(1)), 0o600)
	require.NoError(t, err)
	templateConfig := fmt.Sprintf(templateConfigString, templateSrc, tempDir, "render_1.json")

	config := `
vault {
  address = "%s"
	tls_skip_verify = true
}

auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
            remove_secret_id_file_after_reading = false
        }
    }
}

%s
`
	config = fmt.Sprintf(config, apiClient.Address(), roleIDPath, secretIDPath, templateConfig)
	configFileName := filepath.Join(tempDir, "config.hcl")
	err = os.WriteFile(configFileName, []byte(config), 0o600)
	require.NoError(t, err)
	_, cmd := testAgentCommand(t, nil)
	logFilePath := filepath.Join(tempDir, "agent")

	// Start Vault Agent
	go func() {
		code := cmd.Run([]string{"-config", configFileName, "-log-format", "json", "-log-file", logFilePath, "-log-level", "trace"})
		require.Equalf(t, 0, code, "Vault Agent returned a non-zero exit code")
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout starting agent")
	}

	// Give Vault Agent some time to render our template.
	time.Sleep(3 * time.Second)

	// This flag will be used to capture whether we saw a consul-template log
	// message in the log file (the presence of the log file is also part of the test)
	found := false

	// Vault Agent file logs will match agent-{timestamp}.log based on the
	// cmd line argument we supplied, e.g. agent-1701258869573205000.log
	m, err := filepath.Glob(logFilePath + "*")
	require.NoError(t, err)
	require.Truef(t, len(m) > 0, "no files were found")

	for _, p := range m {
		f, err := os.Open(p)
		require.NoError(t, err)

		fs := bufio.NewScanner(f)
		fs.Split(bufio.ScanLines)

		for fs.Scan() {
			s := fs.Text()
			entry := make(map[string]string)
			err := json.Unmarshal([]byte(s), &entry)
			require.NoError(t, err)
			v, ok := entry["@message"]
			if !ok {
				continue
			}
			if v == runnerLogMessage {
				found = true
				break
			}
		}
	}

	require.Truef(t, found, "unable to find consul-template partial message in logs", runnerLogMessage)
}

// Get a randomly assigned port and then free it again before returning it.
// There is still a race when trying to use it, but should work better
// than a static port.
func generateListenerAddress(t *testing.T) string {
	t.Helper()

	ln1, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listenAddr := ln1.Addr().String()
	ln1.Close()
	return listenAddr
}
