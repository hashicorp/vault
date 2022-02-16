package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*SecretsMoveCommand)(nil)
	_ cli.CommandAutocomplete = (*SecretsMoveCommand)(nil)
)

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

  This command works within or across namespaces, both source and destination paths
  can be prefixed with a namespace heirarchy relative to the current namespace.

  WARNING! Moving an existing secrets engine will revoke any leases from the
  old engine.

  Move the existing secrets engine at secret/ to generic/:

      $ vault secrets move secret/ generic/

  Move the existing secrets engine at ns1/secret/ across namespaces to ns2/generic/, 
  where ns1 and ns2 are child namespaces of the current namespace:

      $ vault secrets move ns1/secret/ ns2/generic/

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

	remountResp, err := client.Sys().StartRemount(source, destination)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error moving secrets engine %s to %s: %s", source, destination, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Started moving secrets engine %s to %s, with migration ID %s", source, destination, remountResp.MigrationID))
	return 0
}
