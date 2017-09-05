package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault-enterprise/meta"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*RemountCommand)(nil)
var _ cli.CommandAutocomplete = (*RemountCommand)(nil)

// RemountCommand is a Command that remounts a mounted secret backend
// to a new endpoint.
type RemountCommand struct {
	*BaseCommand
}

func (c *RemountCommand) Synopsis() string {
	return "Remounts a secret backend to a new path"
}

func (c *RemountCommand) Help() string {
	helpText := `
Usage: vault remount [options] SOURCE DESTINATION

  Remounts an existing secret backend to a new path. Any leases from the old
  backend are revoked, but the data associated with the backend (such as
  configuration), is preserved.

  Move the existing mount at secret/ to generic/:

      $ vault remount secret/ generic/

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *RemountCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *RemountCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *RemountCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RemountCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch len(args) {
	case 0, 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case 2:
	default:
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
		c.UI.Error(fmt.Sprintf("Error remounting %s to %s: %s", source, destination, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Remounted %s to: %s", source, destination))
	return 0
}

func (c *RemountCommand) Synopsis() string {
	return "Remount a secret backend to a new path"
}

func (c *RemountCommand) Help() string {
	helpText := `
Usage: vault remount [options] from to

  Remount a mounted secret backend to a new path.

  This command remounts a secret backend that is already mounted to
  a new path. All the secrets from the old path will be revoked, but
  the data associated with the backend (such as configuration), will
  be preserved.

  Example: vault remount secret/ kv/

General Options:
` + meta.GeneralOptionsUsage()

	return strings.TrimSpace(helpText)
}
