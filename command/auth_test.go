package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	tokenDisk "github.com/hashicorp/vault/builtin/token/disk"
	"github.com/hashicorp/vault/command/token"
	"github.com/mitchellh/cli"
)

func TestAuth_token(t *testing.T) {
	testAuthInit(t)

	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"foo",
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

	if actual != "foo" {
		t.Fatalf("bad: %s", actual)
	}
}

func TestAuth_argsWithMethod(t *testing.T) {
	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"-method=foo",
		"bar",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestAuth_tooManyArgs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"foo",
		"bar",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
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
