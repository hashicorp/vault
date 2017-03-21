package command

import (
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
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
	core, ts, keys, _ := vault.TestCoreWithTokenStore(t)
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

	share := ptBuf.String()

	req := logical.TestRequest(t, logical.ReadOperation, "lookup-self")
	req.ClientToken = share

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("error running share lookup-self: %v", err)
	}
	if resp == nil {
		t.Fatalf("got nil resp with share lookup-self")
	}
	if resp.Data == nil {
		t.Fatalf("got nil resp.Data with share lookup-self")
	}

	if resp.Data["orphan"].(bool) != true ||
		resp.Data["ttl"].(int64) != 0 ||
		resp.Data["num_uses"].(int) != 0 ||
		resp.Data["meta"].(map[string]string) != nil ||
		len(resp.Data["policies"].([]string)) != 1 ||
		resp.Data["policies"].([]string)[0] != "root" {
		t.Fatalf("bad: %#v", resp.Data)
	}
}
