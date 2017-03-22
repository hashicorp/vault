package command

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestGenerateShare_Cancel(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateShareCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	tempDir, pubFiles, err := getPubKeyFiles(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Init the attempt
	args := []string{
		"-address", addr,
		"-init",
		"-pgp-key", pubFiles[0],
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{"-address", addr, "-cancel"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if config != nil {
		t.Fatal("should not have a config for root generation")
	}
}

func TestGenerateShare_status(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateShareCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	tempDir, pubFiles, err := getPubKeyFiles(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Init the attempt
	args := []string{
		"-address", addr,
		"-init",
		"-pgp-key", pubFiles[0],
	}
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

func TestGenerateShare_PGP(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateShareCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	tempDir, pubFiles, err := getPubKeyFiles(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Init the attempt
	args := []string{
		"-address", addr,
		"-init",
		"-pgp-key", pubFiles[0],
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	for _, key := range keys {
		c = &GenerateShareCommand{
			Key: hex.EncodeToString(key),
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		// Provide the key
		args = []string{
			"-address", addr,
		}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	// Get new share
	beforeNAfter := strings.Split(ui.OutputWriter.String(), "Share: ")
	if len(beforeNAfter) != 2 {
		t.Fatalf("did not find share in %s", ui.OutputWriter.String())
	}
	encodedShare := strings.TrimSpace(beforeNAfter[1])

	ptBuf, err := pgpkeys.DecryptBytes(encodedShare, pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if ptBuf == nil {
		t.Fatal("returned plain text buffer is nil")
	}

	newShare := ptBuf.String()
	newShareBytes, err := base64.StdEncoding.DecodeString(newShare)

	keys[0] = newShareBytes

	// Seal the vault
	sc := &SealCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}
	args = []string{"-address", addr}
	if code := sc.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	sealed, err := core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !sealed {
		t.Fatal("should be sealed")
	}

	// Unseal using new share
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

	sealed, err = core.Sealed()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if sealed {
		t.Fatal("should not be sealed")
	}
}
