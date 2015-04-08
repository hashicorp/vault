package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	tokenDisk "github.com/hashicorp/vault/builtin/token/disk"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/http"
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
		Meta: Meta{
			ClientToken: token,
			Ui:          ui,
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
		Meta: Meta{
			Ui: ui,
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
		Meta: Meta{
			Ui: ui,
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
		"token_helper = \"%s\"\n", token.TestProcessPath(t))
	ioutil.WriteFile(filepath.Join(td, ".vault"), []byte(config), 0644)
}

func TestHelperProcess(t *testing.T) {
	token.TestHelperProcessCLI(t, &tokenDisk.Command{})
}

type testAuthHandler struct{}

func (h *testAuthHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	return m["foo"], nil
}

func (h *testAuthHandler) Help() string { return "" }
