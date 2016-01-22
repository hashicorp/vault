package command

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
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
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
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

	pgpKeys := []string{}
	for _, pubFile := range pubFiles {
		pub, err := pgpkeys.ReadPGPFile(pubFile)
		if err != nil {
			t.Fatalf("bad: %v", err)
		}
		pgpKeys = append(pgpKeys, pub)
	}

	expected := &vault.SealConfig{
		SecretShares:    4,
		SecretThreshold: 2,
		PGPKeys:         pgpKeys,
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

	parseDecryptAndTestUnsealKeys(t, ui.OutputWriter.String(), rootToken, false, nil, core)
}

func TestInit_PGP_Idempotentcy(t *testing.T) {
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
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
		"-key-threshold", "2",
		"-idempotent",
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

	// Try again with same args, should be success
	args = []string{
		"-address", addr,
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
		"-key-threshold", "2",
		"-idempotent",
	}

	ui.OutputWriter.Reset()

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if strings.TrimSpace(ui.OutputWriter.String()) != "Vault is already initialized" {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}

	// These should all fail as they change parameters
	args = []string{
		"-address", addr,
		"-key-shares", "3",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
		"-key-threshold", "2",
		"-idempotent",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{
		"-address", addr,
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
		"-key-threshold", "1",
		"-idempotent",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{
		"-address", addr,
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[3] + "," + pubFiles[3],
		"-key-threshold", "2",
		"-idempotent",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
