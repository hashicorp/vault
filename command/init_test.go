package command

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestInit(t *testing.T) {
	ui := new(cli.MockUi)
	c := &InitCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	init, err := core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if init {
		t.Fatal("should not be initialized")
	}

	args := []string{"-address", addr}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	init, err = core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !init {
		t.Fatal("should be initialized")
	}

	sealConf, err := core.SealConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected := &vault.SealConfig{
		SecretShares:    5,
		SecretThreshold: 3,
	}
	if !reflect.DeepEqual(expected, sealConf) {
		t.Fatalf("bad: %#v", sealConf)
	}
}

func TestInit_custom(t *testing.T) {
	ui := new(cli.MockUi)
	c := &InitCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	init, err := core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if init {
		t.Fatal("should not be initialized")
	}

	args := []string{
		"-address", addr,
		"-key-shares", "7",
		"-key-threshold", "3",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	init, err = core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !init {
		t.Fatal("should be initialized")
	}

	sealConf, err := core.SealConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected := &vault.SealConfig{
		SecretShares:    7,
		SecretThreshold: 3,
	}
	if !reflect.DeepEqual(expected, sealConf) {
		t.Fatalf("bad: %#v", sealConf)
	}
}
