package command

import (
	"io"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestWrite(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()
	http.TestServerAuth(t, addr, token)

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/foo",
		"bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data[DefaultDataKey] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestWrite_arbitrary(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()
	http.TestServerAuth(t, addr, token)

	stdinR, stdinW := io.Pipe()
	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: Meta{
			Ui: ui,
		},

		testStdin: stdinR,
	}

	go func() {
		stdinW.Write([]byte(`{"foo":"bar"}`))
		stdinW.Close()
	}()

	args := []string{
		"-address", addr,
		"secret/foo",
		"-",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}
}
