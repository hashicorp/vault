package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func testAgentCommand(tb testing.TB, logger hclog.Logger) (*cli.MockUi, *AgentCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AgentCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		ShutdownCh: MakeShutdownCh(),
		logger:     logger,
	}
}

/*
func TestAgent_Cache_UnixListener(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger.Named("core"),
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

	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, client.Address())

	defer os.Setenv(api.EnvVaultCACert, os.Getenv(api.EnvVaultCACert))
	os.Setenv(api.EnvVaultCACert, fmt.Sprintf("%s/ca_cert.pem", cluster.TempDir))

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

	inf, err := ioutil.TempFile("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	os.Remove(in)
	t.Logf("input: %s", in)

	sink1f, err := ioutil.TempFile("", "sink1.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink1 := sink1f.Name()
	sink1f.Close()
	os.Remove(sink1)
	t.Logf("sink1: %s", sink1)

	sink2f, err := ioutil.TempFile("", "sink2.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink2 := sink2f.Name()
	sink2f.Close()
	os.Remove(sink2)
	t.Logf("sink2: %s", sink2)

	conff, err := ioutil.TempFile("", "conf.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	conf := conff.Name()
	conff.Close()
	os.Remove(conf)
	t.Logf("config: %s", conf)

	jwtToken, _ := agent.GetTestJWT(t)
	if err := ioutil.WriteFile(in, []byte(jwtToken), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test jwt", "path", in)
	}

	socketff, err := ioutil.TempFile("", "cache.socket.")
	if err != nil {
		t.Fatal(err)
	}
	socketf := socketff.Name()
	socketff.Close()
	os.Remove(socketf)
	t.Logf("socketf: %s", socketf)

	config := `
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

cache {
	use_auto_auth_token = true

	listener "unix" {
		address = "%s"
		tls_disable = true
	}
}
`

	config = fmt.Sprintf(config, in, sink1, sink2, socketf)
	if err := ioutil.WriteFile(conf, []byte(config), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test config", "path", conf)
	}

	_, cmd := testAgentCommand(t, logger)
	cmd.client = client

	// Kill the command 5 seconds after it starts
	go func() {
		select {
		case <-cmd.ShutdownCh:
		case <-time.After(5 * time.Second):
			cmd.ShutdownCh <- struct{}{}
		}
	}()

	originalVaultAgentAddress := os.Getenv(api.EnvVaultAgentAddr)

	// Create a client that talks to the agent
	os.Setenv(api.EnvVaultAgentAddr, socketf)
	testClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv(api.EnvVaultAgentAddr, originalVaultAgentAddress)

	// Start the agent
	go cmd.Run([]string{"-config", conf})

	// Give some time for the auto-auth to complete
	time.Sleep(1 * time.Second)

	// Invoke lookup self through the agent
	secret, err := testClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data == nil || secret.Data["id"].(string) == "" {
		t.Fatalf("failed to perform lookup self through agent")
	}
}
*/

func TestAgent_ExitAfterAuth(t *testing.T) {
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

	inf, err := ioutil.TempFile("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	os.Remove(in)
	t.Logf("input: %s", in)

	sink1f, err := ioutil.TempFile("", "sink1.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink1 := sink1f.Name()
	sink1f.Close()
	os.Remove(sink1)
	t.Logf("sink1: %s", sink1)

	sink2f, err := ioutil.TempFile("", "sink2.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink2 := sink2f.Name()
	sink2f.Close()
	os.Remove(sink2)
	t.Logf("sink2: %s", sink2)

	conff, err := ioutil.TempFile("", "conf.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	conf := conff.Name()
	conff.Close()
	os.Remove(conf)
	t.Logf("config: %s", conf)

	jwtToken, _ := agent.GetTestJWT(t)
	if err := ioutil.WriteFile(in, []byte(jwtToken), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test jwt", "path", in)
	}

	config := `
exit_after_auth = true

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

	config = fmt.Sprintf(config, in, sink1, sink2)
	if err := ioutil.WriteFile(conf, []byte(config), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test config", "path", conf)
	}

	// If this hangs forever until the test times out, exit-after-auth isn't
	// working
	ui, cmd := testAgentCommand(t, logger)
	cmd.client = client

	code := cmd.Run([]string{"-config", conf})
	if code != 0 {
		t.Errorf("expected %d to be %d", code, 0)
		t.Logf("output from agent:\n%s", ui.OutputWriter.String())
		t.Logf("error from agent:\n%s", ui.ErrorWriter.String())
	}

	sink1Bytes, err := ioutil.ReadFile(sink1)
	if err != nil {
		t.Fatal(err)
	}
	if len(sink1Bytes) == 0 {
		t.Fatal("got no output from sink 1")
	}

	sink2Bytes, err := ioutil.ReadFile(sink2)
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
			Logger: logger,
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
	  "token_ttl": "1m"
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

	// Get a temp file path we can use for the sink
	sinkPath := makeTempFile(t, "sink.txt", "")
	defer os.Remove(sinkPath)

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

    sink "file" {
        config = {
            path = "%s"
        }
    }
}

cache {
    use_auto_auth_token = true
}

listener "tcp" {
    address = "127.0.0.1:8101"
    tls_disable = true
}
listener "tcp" {
    address = "127.0.0.1:8102"
    tls_disable = true
    require_request_header = false
}
listener "tcp" {
    address = "127.0.0.1:8103"
    tls_disable = true
    require_request_header = true
}
`
	config = fmt.Sprintf(config, roleIDPath, secretIDPath, sinkPath)
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

	// defer agent shutdown
	defer func() {
		cmd.ShutdownCh <- struct{}{}
		wg.Wait()
	}()

	//----------------------------------------------------
	// Perform the tests
	//----------------------------------------------------

	// Test against a listener configuration that omits
	// 'require_request_header', with the header missing from the request.
	agentClient := newApiClient("http://127.0.0.1:8101", false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)

	// Test against a listener configuration that sets 'require_request_header'
	// to 'false', with the header missing from the request.
	agentClient = newApiClient("http://127.0.0.1:8102", false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with the header missing from the request.
	agentClient = newApiClient("http://127.0.0.1:8103", false)
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
	agentClient = newApiClient("http://127.0.0.1:8103", false)
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
	agentClient = newApiClient("http://127.0.0.1:8103", true)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)
}

// TestAgent_Template tests rendering templates
func TestAgent_Template_Basic(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			Logger: logger,
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

	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// give test-role permissions to read the kv secret
	req = serverClient.NewRequest("PUT", "/v1/sys/policy/myapp-read")
	req.BodyBytes = []byte(`{
	  "policy": "path \"secret/*\" { capabilities = [\"read\", \"list\"] }"
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
	defer os.Remove(roleIDPath)
	defer os.Remove(secretIDPath)

	// setup the kv secrets
	req = serverClient.NewRequest("POST", "/v1/sys/mounts/secret/tune")
	req.BodyBytes = []byte(`{
	"options": {"version": "2"}
	}`)
	request(t, serverClient, req, 200)

	// populate a secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "bar",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// populate another secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/otherapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap",
			"cert": "something"
    }
	}`)
	request(t, serverClient, req, 200)

	// Get a temp file path we can use for the sink
	sinkPath := makeTempFile(t, "sink.txt", "")
	defer os.Remove(sinkPath)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := ioutil.TempDir("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// start test cases here
	testCases := map[string]struct {
		templateCount int
		exitAfterAuth bool
	}{
		"zero": {},
		"zero-with-exit": {
			exitAfterAuth: true,
		},
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
			// make some template files
			var templatePaths []string
			for i := 0; i < tc.templateCount; i++ {
				path := makeTempFile(t, fmt.Sprintf("render_%d", i), templateContents())
				templatePaths = append(templatePaths, path)
			}

			// create temp dir for this test run
			tmpDir, err := ioutil.TempDir(tmpDirRoot, tcname)
			if err != nil {
				t.Fatal(err)
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

    sink "file" {
        config = {
            path = "%s"
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

			config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, sinkPath, templateConfig, exitAfterAuth)
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
				// We need to sleep to give Agent time to render the templates. Without this
				// sleep, the test will attempt to read the temp dir before Agent has had time
				// to render and will likely fail the test
				time.Sleep(5 * time.Second)
				cmd.ShutdownCh <- struct{}{}
			}
			wg.Wait()

			//----------------------------------------------------
			// Perform the tests
			//----------------------------------------------------

			files, err := ioutil.ReadDir(tmpDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(files) != len(templatePaths) {
				t.Fatalf("expected (%d) templates, got (%d)", len(templatePaths), len(files))
			}
		})
	}
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
			Logger: logger,
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

	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// give test-role permissions to read the kv secret
	req = serverClient.NewRequest("PUT", "/v1/sys/policy/myapp-read")
	req.BodyBytes = []byte(`{
	  "policy": "path \"secret/*\" { capabilities = [\"read\", \"list\"] }"
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
	defer os.Remove(roleIDPath)
	defer os.Remove(secretIDPath)

	// setup the kv secrets
	req = serverClient.NewRequest("POST", "/v1/sys/mounts/secret/tune")
	req.BodyBytes = []byte(`{
	"options": {"version": "2"}
	}`)
	request(t, serverClient, req, 200)

	// populate a secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "bar",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// populate another secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp2")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// populate another, another secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/otherapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap",
			"cert": "something"
    }
	}`)
	request(t, serverClient, req, 200)

	// Get a temp file path we can use for the sink
	sinkPath := makeTempFile(t, "sink.txt", "")
	defer os.Remove(sinkPath)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := ioutil.TempDir("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// create temp dir for this test run
	tmpDir, err := ioutil.TempDir(tmpDirRoot, "agent-test")
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

    sink "file" {
        config = {
            path = "%s"
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

	config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, sinkPath, tmpDir, tmpDir, tmpDir)
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

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected (%d) templates, got (%d)", 3, len(files))
	}
}

// a slice of template options
var templates = []string{
	`{{ with secret "secret/otherapp"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
{{ .Data.data.cert }}
}
{{ end }}`,
	`{{ with secret "secret/myapp"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
}
{{ end }}`,
	`{{ with secret "secret/myapp"}}
{
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
}
{{ end }}`,
}

// templateContents returns one of 2 templates at pseudo-random. We do this to
// create a dynamic, changing number of template stanzas with different sources,
// to simulate multiple stanzas with a varying number of source templates.
func templateContents() string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(templates))
	return templates[index]
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

	bytes, err := ioutil.ReadAll(resp.Body)
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
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.WriteString(contents)
	f.Close()
	return path
}
