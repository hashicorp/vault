package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestDelete(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &DeleteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/foo",
	}

	// Run once so the client is setup, ignore errors
	c.Run(args)

	// Get the client so we can write data
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	data := map[string]interface{}{"value": "bar"}
	if _, err := client.Logical().Write("secret/foo", data); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Run the delete
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}
}
