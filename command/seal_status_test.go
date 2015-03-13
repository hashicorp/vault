package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestSealStatus(t *testing.T) {
	ui := new(cli.MockUi)
	c := &SealStatusCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	keys := vault.TestCoreInit(t, core)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	args := []string{"-address", addr}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	for _, k := range keys {
		if _, err := core.Unseal(k); err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
