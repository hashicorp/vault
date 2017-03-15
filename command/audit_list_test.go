package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestAuditList(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &AuditListCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
	}

	// Run once to get the client
	c.Run(args)

	// Get the client
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if err := client.Sys().EnableAuditWithOptions("foo", &api.EnableAuditOptions{
		Type:        "noop",
		Description: "noop",
		Options:     nil,
	}); err != nil {
		t.Fatalf("err: %#v", err)
	}

	// Run again
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
