package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestUnwrap(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &UnwrapCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"-field", "zip",
	}

	// Run once so the client is setup, ignore errors
	c.Run(args)

	// Get the client so we can write data
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	wrapLookupFunc := func(method, path string) string {
		if method == "GET" && path == "secret/foo" {
			return "60s"
		}
		if method == "LIST" && path == "secret" {
			return "60s"
		}
		return ""
	}
	client.SetWrappingLookupFunc(wrapLookupFunc)

	data := map[string]interface{}{"zip": "zap"}
	if _, err := client.Logical().Write("secret/foo", data); err != nil {
		t.Fatalf("err: %s", err)
	}

	outer, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if outer == nil {
		t.Fatal("outer response was nil")
	}
	if outer.WrapInfo == nil {
		t.Fatalf("outer wrapinfo was nil, response was %#v", *outer)
	}

	args = append(args, outer.WrapInfo.Token)

	// Run the read
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if output != "zap\n" {
		t.Fatalf("unexpectd output:\n%s", output)
	}

	// Now test with list handling, specifically that it will be called with
	// the list output formatter
	ui.OutputWriter.Reset()

	outer, err = client.Logical().List("secret")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if outer == nil {
		t.Fatal("outer response was nil")
	}
	if outer.WrapInfo == nil {
		t.Fatalf("outer wrapinfo was nil, response was %#v", *outer)
	}

	args = []string{
		"-address", addr,
		outer.WrapInfo.Token,
	}
	// Run the read
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	output = ui.OutputWriter.String()
	if strings.TrimSpace(output) != "Keys\n----\nfoo" {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
