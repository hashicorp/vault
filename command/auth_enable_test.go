package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestAuthEnable(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &AuthEnableCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"noop",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mount, ok := mounts["noop/"]
	if !ok {
		t.Fatal("should have noop mount")
	}
	if mount.Type != "noop" {
		t.Fatal("should be noop type")
	}
}
