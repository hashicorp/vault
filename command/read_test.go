package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRead(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &ReadCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
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

func TestRead_notFound(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &ReadCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/nope",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestRead_field(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &ReadCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"-field", "value",
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

	// Run the read
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if output != "bar\n" {
		t.Fatalf("unexpectd output:\n%s", output)
	}
}

func TestRead_field_notFound(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &ReadCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"-field", "nope",
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

	// Run the read
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
