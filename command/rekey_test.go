package command

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRekey(t *testing.T) {
	core, key, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Key: hex.EncodeToString(key),
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{"-address", addr}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.SealConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config.SecretShares != 5 {
		t.Fatal("should rekey")
	}
}

func TestRekey_arg(t *testing.T) {
	core, key, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{"-address", addr, hex.EncodeToString(key)}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.SealConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config.SecretShares != 5 {
		t.Fatal("should rekey")
	}
}

func TestRekey_init(t *testing.T) {
	core, key, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Key: hex.EncodeToString(key),
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{"-address", addr, "-init", "-key-threshold=10", "-key-shares=10"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.RekeyConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config.SecretShares != 10 {
		t.Fatal("should rekey")
	}
	if config.SecretThreshold != 10 {
		t.Fatal("should rekey")
	}
}

func TestRekey_cancel(t *testing.T) {
	core, key, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Key: hex.EncodeToString(key),
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{"-address", addr, "-init"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{"-address", addr, "-cancel"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.RekeyConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config != nil {
		t.Fatal("should not rekey")
	}
}

func TestRekey_status(t *testing.T) {
	core, key, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Key: hex.EncodeToString(key),
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{"-address", addr, "-init"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{"-address", addr, "-status"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !strings.Contains(string(ui.OutputWriter.Bytes()), "Started: true") {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}
