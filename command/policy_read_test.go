package command

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestPolicyRead(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	// Register a policy first and then read back the same
	ui := new(cli.MockUi)
	writeCmd := &PolicyWriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	writeArgs := []string{
		"-address", addr,
		"foo",
		"./test-fixtures/policy.hcl",
	}
	if code := writeCmd.Run(writeArgs); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	readCmd := &PolicyReadCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}
	readArgs := []string{
		"-address", addr,
		"foo",
	}
	if code := readCmd.Run(readArgs); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	policyBytes, _ := ioutil.ReadFile("./test-fixtures/policy.hcl")
	if !strings.Contains(output, string(policyBytes)) {
		t.Fatalf("bad: %#v", output)
	}
}
