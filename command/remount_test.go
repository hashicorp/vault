package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRemount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RemountCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/", "kv",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	_, ok := mounts["secret/"]
	if ok {
		t.Fatal("should not have mount")
	}

	_, ok = mounts["kv/"]
	if !ok {
		t.Fatal("should have kv")
	}
}
