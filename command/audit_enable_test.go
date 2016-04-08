package command

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestAuditEnable(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &AuditEnableCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"noop",
		"foo=bar",
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Get the client
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %#v", err)
	}

	audits, err := client.Sys().ListAudit()
	if err != nil {
		t.Fatalf("err: %#v", err)
	}

	audit, ok := audits["noop/"]
	if !ok {
		t.Fatalf("err: %#v", audits)
	}

	expected := map[string]string{"foo": "bar"}
	if !reflect.DeepEqual(audit.Options, expected) {
		t.Fatalf("err: %#v", audit)
	}
}
