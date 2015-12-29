package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/cli"
)

var (
	basehcl = `
disable_mlock = true

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = "true"
}
`

	consulhcl = `
backend "consul" {
    prefix = "foo/"
	advertise_addr = "http://127.0.0.1:8200"
}
`
	haconsulhcl = `
ha_backend "consul" {
    prefix = "bar/"
	advertise_addr = "http://127.0.0.1:8200"
}
`

	badhaconsulhcl = `
ha_backend "file" {
    path = "/dev/null"
}
`
)

func TestServer_GoodSeparateHA(t *testing.T) {
	ui := new(cli.MockUi)
	c := &ServerCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}

	tmpfile.WriteString(basehcl + consulhcl + haconsulhcl)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	args := []string{"-config", tmpfile.Name(), "-verify-only", "true"}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestServer_BadSeparateHA(t *testing.T) {
	ui := new(cli.MockUi)
	c := &ServerCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}

	tmpfile.WriteString(basehcl + consulhcl + badhaconsulhcl)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	args := []string{"-config", tmpfile.Name()}

	if code := c.Run(args); code == 0 {
		t.Fatalf("bad: should have gotten an error on a bad HA config")
	}
}
