package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestAuth_argsWithMethod(t *testing.T) {
	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"-method=foo",
		"bar",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestAuth_tooManyArgs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &AuthCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{
		"foo",
		"bar",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
