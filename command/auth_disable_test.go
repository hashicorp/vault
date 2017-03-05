package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestAuthDisable(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &AuthDisableCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"noop",
	}

	// Run the command once to setup the client, it will fail
	c.Run(args)

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := client.Sys().EnableAuth("noop", "noop", ""); err != nil {
		t.Fatalf("err: %s", err)
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, ok := mounts["noop"]; ok {
		t.Fatal("should not have noop mount")
	}
}

func TestAuthDisableWithOptions(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &AuthDisableCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"noop",
	}

	// Run the command once to setup the client, it will fail
	c.Run(args)

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := client.Sys().EnableAuthWithOptions("noop", &api.EnableAuthOptions{
		Type:        "noop",
		Description: "",
	}); err != nil {
		t.Fatalf("err: %#v", err)
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, ok := mounts["noop"]; ok {
		t.Fatal("should not have noop mount")
	}
}
