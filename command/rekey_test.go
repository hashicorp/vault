package command

import (
	"encoding/hex"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRekey(t *testing.T) {
	core, keys, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)

	for i, key := range keys {
		c := &RekeyCommand{
			Key:         hex.EncodeToString(key),
			RecoveryKey: false,
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		if i > 0 {
			conf, err := core.RekeyConfig(false)
			if err != nil {
				t.Fatal(err)
			}
			c.Nonce = conf.Nonce
		}

		args := []string{"-address", addr}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	config, err := core.SealAccess().BarrierConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config.SecretShares != 5 {
		t.Fatal("should rekey")
	}
}

func TestRekey_arg(t *testing.T) {
	core, keys, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)

	for i, key := range keys {
		c := &RekeyCommand{
			RecoveryKey: false,
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		if i > 0 {
			conf, err := core.RekeyConfig(false)
			if err != nil {
				t.Fatal(err)
			}
			c.Nonce = conf.Nonce
		}

		args := []string{"-address", addr, hex.EncodeToString(key)}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	config, err := core.SealAccess().BarrierConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config.SecretShares != 5 {
		t.Fatal("should rekey")
	}
}

func TestRekey_init(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)

	c := &RekeyCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	args := []string{
		"-address", addr,
		"-init",
		"-key-threshold", "10",
		"-key-shares", "10",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.RekeyConfig(false)
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
	core, keys, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Key: hex.EncodeToString(keys[0]),
		Meta: meta.Meta{
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

	config, err := core.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config != nil {
		t.Fatal("should not rekey")
	}
}

func TestRekey_status(t *testing.T) {
	core, keys, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Key: hex.EncodeToString(keys[0]),
		Meta: meta.Meta{
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

	if !strings.Contains(ui.OutputWriter.String(), "Started: true") {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}

func TestRekey_init_pgp(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	bc := &logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}
	sysBackend := vault.NewSystemBackend(core)
	err := sysBackend.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	ui := new(cli.MockUi)
	c := &RekeyCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	tempDir, pubFiles, err := getPubKeyFiles(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	args := []string{
		"-address", addr,
		"-init",
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
		"-key-threshold", "2",
		"-backup", "true",
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config.SecretShares != 4 {
		t.Fatal("should rekey")
	}
	if config.SecretThreshold != 2 {
		t.Fatal("should rekey")
	}

	for _, key := range keys {
		c = &RekeyCommand{
			Key: hex.EncodeToString(key),
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		c.Nonce = config.Nonce

		args = []string{
			"-address", addr,
		}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	type backupStruct struct {
		Keys    map[string][]string
		KeysB64 map[string][]string
	}
	backupVals := &backupStruct{}

	req := logical.TestRequest(t, logical.ReadOperation, "rekey/backup")
	resp, err := sysBackend.HandleRequest(req)
	if err != nil {
		t.Fatalf("error running backed-up unseal key fetch: %v", err)
	}
	if resp == nil {
		t.Fatalf("got nil resp with unseal key fetch")
	}
	if resp.Data["keys"] == nil {
		t.Fatalf("could not retrieve unseal keys from token")
	}
	if resp.Data["nonce"] != config.Nonce {
		t.Fatalf("nonce mismatch between rekey and backed-up keys")
	}

	backupVals.Keys = resp.Data["keys"].(map[string][]string)
	backupVals.KeysB64 = resp.Data["keys_base64"].(map[string][]string)

	// Now delete and try again; the values should be inaccessible
	req = logical.TestRequest(t, logical.DeleteOperation, "rekey/backup")
	resp, err = sysBackend.HandleRequest(req)
	if err != nil {
		t.Fatalf("error running backed-up unseal key delete: %v", err)
	}
	req = logical.TestRequest(t, logical.ReadOperation, "rekey/backup")
	resp, err = sysBackend.HandleRequest(req)
	if err != nil {
		t.Fatalf("error running backed-up unseal key fetch: %v", err)
	}
	if resp == nil {
		t.Fatalf("got nil resp with unseal key fetch")
	}
	if resp.Data["keys"] != nil {
		t.Fatalf("keys found when they should have been deleted")
	}

	// Sort, because it'll be tested with DeepEqual later
	for k, _ := range backupVals.Keys {
		sort.Strings(backupVals.Keys[k])
		sort.Strings(backupVals.KeysB64[k])
	}

	parseDecryptAndTestUnsealKeys(t, ui.OutputWriter.String(), token, true, backupVals.Keys, backupVals.KeysB64, core)
}
