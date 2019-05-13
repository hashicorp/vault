package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*SecretsMoveCommand)(nil)
var _ cli.CommandAutocomplete = (*SecretsMoveCommand)(nil)

type SecretsMoveCommand struct {
	*BaseCommand
}

func (c *SecretsMoveCommand) Synopsis() string {
	return "Move a secrets engine to a new path"
}

func (c *SecretsMoveCommand) Help() string {
	helpText := `
Usage: vault secrets move [options] SOURCE DESTINATION

  Moves an existing secrets engine to a new path. Any leases from the old
  secrets engine are revoked, but all configuration associated with the engine
  is preserved.

  This command only works within a namespace; it cannot be used to move engines
  to different namespaces.

  WARNING! Moving an existing secrets engine will revoke any leases from the
  old engine.

  Move the existing secrets engine at secret/ to generic/:

      $ vault secrets move secret/ generic/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SecretsMoveCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *SecretsMoveCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *SecretsMoveCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SecretsMoveCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}

	// Grab the source and destination
	source := ensureTrailingSlash(args[0])
	destination := ensureTrailingSlash(args[1])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().Remount(source, destination); err != nil {
		c.UI.Error(fmt.Sprintf("Error moving secrets engine %s to %s: %s", source, destination, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Moved secrets engine %s to: %s", source, destination))
	return 0
}
