package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func assertNoTabs(tb testing.TB, c cli.Command) {
	if strings.ContainsRune(c.Help(), '\t') {
		tb.Errorf("%#v help output contains tabs", c)
	}
}
