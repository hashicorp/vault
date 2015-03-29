package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRead(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()
	http.TestServerAuth(t, addr, token)

	ui := new(cli.MockUi)
	c := &ReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"-address", addr,
		"sys/mounts",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
