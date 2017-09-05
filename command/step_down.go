package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*StepDownCommand)(nil)
var _ cli.CommandAutocomplete = (*StepDownCommand)(nil)

// StepDownCommand is a Command that tells the Vault server to give up its
// leadership
type StepDownCommand struct {
	*BaseCommand
}

func (c *StepDownCommand) Synopsis() string {
	return "Forces Vault to resign active duty"
}

func (c *StepDownCommand) Help() string {
	helpText := `
Usage: vault step-down [options]

  Forces the Vault server at the given address to step down from active duty.
  While the affected node will have a delay before attempting to acquire the
  leader lock again, if no other Vault nodes acquire the lock beforehand, it
  is possible for the same node to re-acquire the lock and become active
  again.

  Force Vault to step down as the leader:

      $ vault step-down

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *StepDownCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *StepDownCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *StepDownCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *StepDownCommand) Run(args []string) int {
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

	if err := client.Sys().StepDown(); err != nil {
		c.UI.Error(fmt.Sprintf("Error stepping down: %s", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Stepped down: %s", client.Address()))
	return 0
}
