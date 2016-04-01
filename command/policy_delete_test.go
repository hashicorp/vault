package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestPolicyDelete(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &PolicyDeleteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"foo",
	}

	// Run once so the client is setup, ignore errors
	c.Run(args)

	// Get the client so we can write data
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := client.Sys().PutPolicy("foo", testPolicyDeleteRules); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test that the delete works
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Test the policy is gone
	rules, err := client.Sys().GetPolicy("foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if rules != "" {
		t.Fatalf("bad: %#v", rules)
	}
}

const testPolicyDeleteRules = `
path "sys" {
	policy = "deny"
}
`
