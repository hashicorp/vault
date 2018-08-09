package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*PluginCommand)(nil)

type PluginCommand struct {
	*BaseCommand
}

func (c *PluginCommand) Synopsis() string {
	return "Interact with Vault plugins and catalog"
}

func (c *PluginCommand) Help() string {
	helpText := `
Usage: vault plugin <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's plugins and the
  plugin catalog. Here are a few examples of the plugin commands:

  List all available plugins in the catalog:

      $ vault plugin list

  Register a new plugin to the catalog:

      $ vault plugin register -sha256=d3f0a8b... my-custom-plugin

  Get information about a plugin in the catalog:

      $ vault plugin info my-custom-plugin

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PluginCommand) Run(args []string) int {
	return cli.RunResultHelp
}
