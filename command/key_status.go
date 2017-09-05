package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*KeyStatusCommand)(nil)
var _ cli.CommandAutocomplete = (*KeyStatusCommand)(nil)

// KeyStatusCommand is a Command that provides information about the key status
type KeyStatusCommand struct {
	*BaseCommand
}

func (c *KeyStatusCommand) Synopsis() string {
	return "Provides information about the active encryption key"
}

func (c *KeyStatusCommand) Help() string {
	helpText := `
Usage: vault key-status [options]

  Provides information about the active encryption key. Specifically,
  the current key term and the key installation time.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KeyStatusCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *KeyStatusCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KeyStatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KeyStatusCommand) Run(args []string) int {
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

	c.UI.Output(printKeyStatus(status))
	return 0
}
