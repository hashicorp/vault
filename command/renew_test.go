package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRenew(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RenewCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	// write a secret with a lease
	client := testClient(t, addr, token)
	_, err := client.Logical().Write("secret/foo", map[string]interface{}{
		"key":   "value",
		"lease": "1m",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// read the secret to get its lease ID
	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	args := []string{
		"-address", addr,
		secret.LeaseID,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
