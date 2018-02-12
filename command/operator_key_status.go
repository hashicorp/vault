package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorKeyStatusCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorKeyStatusCommand)(nil)

type OperatorKeyStatusCommand struct {
	*BaseCommand
}

func (c *OperatorKeyStatusCommand) Synopsis() string {
	return "Provides information about the active encryption key"
}

func (c *OperatorKeyStatusCommand) Help() string {
	helpText := `
Usage: vault operator key-status [options]

  Provides information about the active encryption key. Specifically,
  the current key term and the key installation time.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorKeyStatusCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *OperatorKeyStatusCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *OperatorKeyStatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorKeyStatusCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	status, err := client.Sys().KeyStatus()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading key status: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
		c.UI.Output(printKeyStatus(status))
		return 0
	default:
		return OutputData(c.UI, status)
	}
}
