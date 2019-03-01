package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
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

func TestExitAfterAuth(t *testing.T) {
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
