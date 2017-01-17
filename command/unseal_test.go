package command

import (
	"encoding/hex"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestUnseal(t *testing.T) {
	core := vault.TestCore(t)
	keys, _ := vault.TestCoreInit(t, core)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)

	for _, key := range keys {
		c := &UnsealCommand{
			Key: hex.EncodeToString(key),
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		args := []string{"-address", addr}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	sealed, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if sealed {
		t.Fatal("should not be sealed")
	}
}

func TestUnseal_arg(t *testing.T) {
	core := vault.TestCore(t)
	keys, _ := vault.TestCoreInit(t, core)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)

	for _, key := range keys {
		c := &UnsealCommand{
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		args := []string{"-address", addr, hex.EncodeToString(key)}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	sealed, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if sealed {
		t.Fatal("should not be sealed")
	}
}
