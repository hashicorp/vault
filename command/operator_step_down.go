package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorStepDownCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorStepDownCommand)(nil)

type OperatorStepDownCommand struct {
	*BaseCommand
}

func (c *OperatorStepDownCommand) Synopsis() string {
	return "Forces Vault to resign active duty"
}

func (c *OperatorStepDownCommand) Help() string {
	helpText := `
Usage: vault operator step-down [options]

  Forces the Vault server at the given address to step down from active duty.
  While the affected node will have a delay before attempting to acquire the
  leader lock again, if no other Vault nodes acquire the lock beforehand, it
  is possible for the same node to re-acquire the lock and become active
  again.

  Force Vault to step down as the leader:

      $ vault operator step-down

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorStepDownCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *OperatorStepDownCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *OperatorStepDownCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorStepDownCommand) Run(args []string) int {
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
