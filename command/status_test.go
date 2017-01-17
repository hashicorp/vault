package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestStatus(t *testing.T) {
	ui := new(cli.MockUi)
	c := &StatusCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	keys, _ := vault.TestCoreInit(t, core)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	args := []string{"-address", addr}
	if code := c.Run(args); code != 2 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	for _, key := range keys {
		if _, err := core.Unseal(key); err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
