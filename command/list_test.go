package command

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestList(t *testing.T) {
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
		"-format", "json",
		"secret",
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

	data = map[string]interface{}{"value": "bar"}
	if _, err := client.Logical().Write("secret/foo/bar", data); err != nil {
		t.Fatalf("err: %s", err)
	}

	secret, err := client.Logical().List("secret/")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if secret == nil {
		t.Fatalf("err: No value found at secret/")
	}

	if secret.Data == nil {
		t.Fatalf("err: Data not found")
	}

	exp := map[string]interface{}{
		"keys": []interface{}{"foo", "foo/"},
	}

	if !reflect.DeepEqual(secret.Data, exp) {
		t.Fatalf("err: expected %#v, got %#v", exp, secret.Data)
	}
}
