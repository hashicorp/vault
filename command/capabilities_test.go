package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestCapabilities_Args(t *testing.T) {
	core, key, _ := vault.TestCoreUnsealed(t)
	ln, _ := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &CapabilitiesCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{}
	if code := c.Run(args); code == 0 {
		t.Fatalf("expected failure due to no args")
	}

	args = []string{"test"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{string(key), "test"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
