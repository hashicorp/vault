package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestTokenRevokeAccessor(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenRevokeCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
	}

	// Run it once for client
	c.Run(args)

	// Create a token
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := client.Auth().Token().Create(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Treat the argument as accessor
	args = append(args, "-accessor")
	if code := c.Run(args); code == 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Verify it worked with proper accessor
	args1 := append(args, resp.Auth.Accessor)
	if code := c.Run(args1); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Fail if mode is set to 'orphan' when accessor is set
	args2 := append(args, "-mode=\"orphan\"")
	if code := c.Run(args2); code == 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Fail if mode is set to 'path' when accessor is set
	args3 := append(args, "-mode=\"path\"")
	if code := c.Run(args3); code == 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestTokenRevoke(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenRevokeCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
	}

	// Run it once for client
	c.Run(args)

	// Create a token
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := client.Auth().Token().Create(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify it worked
	args = append(args, resp.Auth.ClientToken)
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
