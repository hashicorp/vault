package command

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestWrite(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/foo",
		"value=bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data["value"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestWrite_arbitrary(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	stdinR, stdinW := io.Pipe()
	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},

		testStdin: stdinR,
	}

	go func() {
		stdinW.Write([]byte(`{"foo":"bar"}`))
		stdinW.Close()
	}()

	args := []string{
		"-address", addr,
		"secret/foo",
		"-",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestWrite_escaped(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/foo",
		"value=\\@bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data["value"] != "@bar" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestWrite_file(t *testing.T) {
	tf, err := ioutil.TempFile("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Write([]byte(`{"foo":"bar"}`))
	tf.Close()
	defer os.Remove(tf.Name())

	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/foo",
		"@" + tf.Name(),
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestWrite_fileValue(t *testing.T) {
	tf, err := ioutil.TempFile("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Write([]byte("foo"))
	tf.Close()
	defer os.Remove(tf.Name())

	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"secret/foo",
		"value=@" + tf.Name(),
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.Data["value"] != "foo" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestWrite_Output(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"auth/token/create",
		"display_name=foo",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "Key") {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}

func TestWrite_force(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"-force",
		"sys/rotate",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
