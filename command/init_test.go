package command

import (
	"os"
	"reflect"
	"regexp"
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

func TestInit_PGP(t *testing.T) {
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

	tempDir, pubFiles, err := getPubKeyFiles(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	args := []string{
		"-address", addr,
		"-key-shares", "2",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2],
		"-key-threshold", "2",
	}

	// This should fail, as key-shares does not match pgp-keys size
	if code := c.Run(args); code == 0 {
		t.Fatalf("bad (command should have failed): %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{
		"-address", addr,
		"-key-shares", "3",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2],
		"-key-threshold", "2",
	}

	ui.OutputWriter.Reset()

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
		SecretShares:    3,
		SecretThreshold: 2,
	}
	if !reflect.DeepEqual(expected, sealConf) {
		t.Fatalf("bad:\nexpected: %#v\ngot: %#v", expected, sealConf)
	}

	re, err := regexp.Compile("\\s+Initial Root Token:\\s+(.*)")
	if err != nil {
		t.Fatalf("Error compiling regex: %s", err)
	}
	matches := re.FindAllStringSubmatch(ui.OutputWriter.String(), -1)
	if len(matches) != 1 {
		t.Fatalf("Unexpected number of tokens found, got %d", len(matches))
	}

	rootToken := matches[0][1]

	parseDecryptAndTestUnsealKeys(t, ui.OutputWriter.String(), rootToken, core)
}
