package command

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestGenerateRoot_Cancel(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateRootCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	otpBytes, err := vault.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	args := []string{"-address", addr, "-init", "-otp", otp}
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

func TestGenerateRoot_status(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateRootCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	otpBytes, err := vault.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	args := []string{"-address", addr, "-init", "-otp", otp}
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

func TestGenerateRoot_OTP(t *testing.T) {
	core, ts, keys, _ := vault.TestCoreWithTokenStore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateRootCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	// Generate an OTP
	otpBytes, err := vault.GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	otp := base64.StdEncoding.EncodeToString(otpBytes)

	// Init the attempt
	args := []string{
		"-address", addr,
		"-init",
		"-otp", otp,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	config, err := core.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, key := range keys {
		ui = new(cli.MockUi)
		c = &GenerateRootCommand{
			Key: hex.EncodeToString(key),
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		c.Nonce = config.Nonce

		// Provide the key
		args = []string{
			"-address", addr,
		}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	beforeNAfter := strings.Split(ui.OutputWriter.String(), "Encoded root token: ")
	if len(beforeNAfter) != 2 {
		t.Fatalf("did not find encoded root token in %s", ui.OutputWriter.String())
	}
	encodedToken := strings.TrimSpace(beforeNAfter[1])

	decodedToken, err := xor.XORBase64(encodedToken, otp)
	if err != nil {
		t.Fatal(err)
	}

	token, err := uuid.FormatUUID(decodedToken)
	if err != nil {
		t.Fatal(err)
	}

	req := logical.TestRequest(t, logical.ReadOperation, "lookup-self")
	req.ClientToken = token

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("error running token lookup-self: %v", err)
	}
	if resp == nil {
		t.Fatalf("got nil resp with token lookup-self")
	}
	if resp.Data == nil {
		t.Fatalf("got nil resp.Data with token lookup-self")
	}

	if resp.Data["orphan"].(bool) != true ||
		resp.Data["ttl"].(int64) != 0 ||
		resp.Data["num_uses"].(int) != 0 ||
		resp.Data["meta"].(map[string]string) != nil ||
		len(resp.Data["policies"].([]string)) != 1 ||
		resp.Data["policies"].([]string)[0] != "root" {
		t.Fatalf("bad: %#v", resp.Data)
	}

	// Clear the output and run a decode to verify we get the same result
	ui.OutputWriter.Reset()
	args = []string{
		"-address", addr,
		"-decode", encodedToken,
		"-otp", otp,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
	beforeNAfter = strings.Split(ui.OutputWriter.String(), "Root token: ")
	if len(beforeNAfter) != 2 {
		t.Fatalf("did not find decoded root token in %s", ui.OutputWriter.String())
	}

	outToken := strings.TrimSpace(beforeNAfter[1])
	if outToken != token {
		t.Fatalf("tokens do not match:\n%s\n%s", token, outToken)
	}
}

func TestGenerateRoot_PGP(t *testing.T) {
	core, ts, keys, _ := vault.TestCoreWithTokenStore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &GenerateRootCommand{
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

	config, err := core.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, key := range keys {
		c = &GenerateRootCommand{
			Key: hex.EncodeToString(key),
			Meta: meta.Meta{
				Ui: ui,
			},
		}

		c.Nonce = config.Nonce

		// Provide the key
		args = []string{
			"-address", addr,
		}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		}
	}

	beforeNAfter := strings.Split(ui.OutputWriter.String(), "Encoded root token: ")
	if len(beforeNAfter) != 2 {
		t.Fatalf("did not find encoded root token in %s", ui.OutputWriter.String())
	}
	encodedToken := strings.TrimSpace(beforeNAfter[1])

	ptBuf, err := pgpkeys.DecryptBytes(encodedToken, pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if ptBuf == nil {
		t.Fatal("returned plain text buffer is nil")
	}

	token := ptBuf.String()

	req := logical.TestRequest(t, logical.ReadOperation, "lookup-self")
	req.ClientToken = token

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("error running token lookup-self: %v", err)
	}
	if resp == nil {
		t.Fatalf("got nil resp with token lookup-self")
	}
	if resp.Data == nil {
		t.Fatalf("got nil resp.Data with token lookup-self")
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
