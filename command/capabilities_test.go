package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestCapabilities_Basic(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()
	ui := new(cli.MockUi)
	c := &CapabilitiesCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	var args []string

	args = []string{"-address", addr}
	if code := c.Run(args); code == 0 {
		t.Fatalf("expected failure due to no args")
	}

	args = []string{"-address", addr, "testpath"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{"-address", addr, token, "test"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{"-address", addr, "invalidtoken", "test"}
	if code := c.Run(args); code == 0 {
		t.Fatalf("expected failure due to invalid token")
	}
}
