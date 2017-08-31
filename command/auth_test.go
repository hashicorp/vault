package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/logical"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestAuth_methods(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	testAuthInit(t)

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
	}

	args := []string{
		"-address", addr,
		"-methods",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "token") {
		t.Fatalf("bad: %#v", output)
	}
}

func TestAuth_token(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	testAuthInit(t)

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
	}

	args := []string{
		"-address", addr,
		token,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	helper, err := c.TokenHelper()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := helper.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if actual != token {
		t.Fatalf("bad: %s", actual)
	}
}

func TestAuth_wrapping(t *testing.T) {
	baseConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, baseConfig, &vault.TestClusterOptions{
		HandlerFunc:       http.Handler,
		BaseListenAddress: "127.0.0.1:8200",
	})
	cluster.Start()
	defer cluster.Cleanup()

	testAuthInit(t)

	client := cluster.Cores[0].Client
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
		"policies": "zip,zap",
	})
	if err != nil {
		t.Fatal(err)
	}

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
		Handlers: map[string]AuthHandler{
			"userpass": &credUserpass.CLIHandler{DefaultMount: "userpass"},
		},
	}

	args := []string{
		"-address",
		"https://127.0.0.1:8200",
		"-tls-skip-verify",
		"-method",
		"userpass",
		"username=foo",
		"password=bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Test again with wrapping
	ui = new(cli.MockUi)
	c = &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
		Handlers: map[string]AuthHandler{
			"userpass": &credUserpass.CLIHandler{DefaultMount: "userpass"},
		},
	}

	args = []string{
		"-address",
		"https://127.0.0.1:8200",
		"-tls-skip-verify",
		"-wrap-ttl",
		"5m",
		"-method",
		"userpass",
		"username=foo",
		"password=bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Test again with no-store
	ui = new(cli.MockUi)
	c = &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
		Handlers: map[string]AuthHandler{
			"userpass": &credUserpass.CLIHandler{DefaultMount: "userpass"},
		},
	}

	args = []string{
		"-address",
		"https://127.0.0.1:8200",
		"-tls-skip-verify",
		"-wrap-ttl",
		"5m",
		"-no-store",
		"-method",
		"userpass",
		"username=foo",
		"password=bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Test again with wrapping and token-only
	ui = new(cli.MockUi)
	c = &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
		Handlers: map[string]AuthHandler{
			"userpass": &credUserpass.CLIHandler{DefaultMount: "userpass"},
		},
	}

	args = []string{
		"-address",
		"https://127.0.0.1:8200",
		"-tls-skip-verify",
		"-wrap-ttl",
		"5m",
		"-token-only",
		"-method",
		"userpass",
		"username=foo",
		"password=bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
	token := strings.TrimSpace(ui.OutputWriter.String())
	if token == "" {
		t.Fatal("expected to find token in output")
	}
	secret, err := client.Logical().Unwrap(token)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.ClientToken == "" {
		t.Fatal("no client token found")
	}
}

func TestAuth_token_nostore(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	testAuthInit(t)

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
	}

	args := []string{
		"-address", addr,
		"-no-store",
		token,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	helper, err := c.TokenHelper()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := helper.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if actual != "" {
		t.Fatalf("bad: %s", actual)
	}
}

func TestAuth_stdin(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	testAuthInit(t)

	stdinR, stdinW := io.Pipe()
	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
		testStdin: stdinR,
	}

	go func() {
		stdinW.Write([]byte(token))
		stdinW.Close()
	}()

	args := []string{
		"-address", addr,
		"-",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestAuth_badToken(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	testAuthInit(t)

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
	}

	args := []string{
		"-address", addr,
		"not-a-valid-token",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestAuth_method(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	testAuthInit(t)

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Handlers: map[string]AuthHandler{
			"test": &testAuthHandler{},
		},
		Meta: meta.Meta{
			Ui:          ui,
			TokenHelper: DefaultTokenHelper,
		},
	}

	args := []string{
		"-address", addr,
		"-method=test",
		"foo=" + token,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	helper, err := c.TokenHelper()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := helper.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if actual != token {
		t.Fatalf("bad: %s", actual)
	}
}

func testAuthInit(t *testing.T) {
	td, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Set the HOME env var so we get that right
	os.Setenv("HOME", td)

	// Write a .vault config to use our custom token helper
	config := fmt.Sprintf(
		"token_helper = \"\"\n")
	ioutil.WriteFile(filepath.Join(td, ".vault"), []byte(config), 0644)
}

type testAuthHandler struct{}

func (h *testAuthHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: m["foo"],
		},
	}, nil
}

func (h *testAuthHandler) Help() string { return "" }
