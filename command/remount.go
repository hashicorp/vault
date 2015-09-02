package command

import (
	"fmt"
	"strings"
)

// RemountCommand is a Command that remounts a mounted secret backend
// to a new endpoint.
type RemountCommand struct {
	Meta
}

func (c *RemountCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("remount", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 2 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nRemount expects two arguments: the from and to path"))
		return 1
	}

	from := args[0]
	to := args[1]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().Remount(from, to); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Unmount error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully remounted from '%s' to '%s'!", from, to))

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
  the Vault data associated with the backend will be preserved (such
  as configuration data).

  Example: vault remount secret/ generic/

General Options:

  ` + generalOptionsUsage()

	return strings.TrimSpace(helpText)
}
