// +build !race

package command

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/meta"
	"github.com/mitchellh/cli"
)

// The following tests have a go-metrics/exp manager race condition
func TestServer_CommonHA(t *testing.T) {
	ui := new(cli.MockUi)
	c := &ServerCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}

	tmpfile.WriteString(basehcl + consulhcl)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	args := []string{"-config", tmpfile.Name(), "-verify-only", "true"}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !strings.Contains(ui.OutputWriter.String(), "(HA available)") {
		t.Fatalf("did not find HA available: %s", ui.OutputWriter.String())
	}
}

func TestServer_GoodSeparateHA(t *testing.T) {
	ui := new(cli.MockUi)
	c := &ServerCommand{
		Meta: meta.Meta{
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

	if !strings.Contains(ui.OutputWriter.String(), "HA Backend:") {
		t.Fatalf("did not find HA Backend: %s", ui.OutputWriter.String())
	}
}

func TestServer_BadSeparateHA(t *testing.T) {
	ui := new(cli.MockUi)
	c := &ServerCommand{
		Meta: meta.Meta{
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
