package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*TransitCommand)(nil)

type TransitCommand struct {
	*BaseCommand
}

func (c *TransitCommand) Synopsis() string {
	return "Interact with Vault's Transit Secrets Engine"
}

func (c *TransitCommand) Help() string {
	helpText := `
Usage: vault transit <subcommand> [options] [args]

  This command has subcommands for interacting with Vault's Transit Secrets
  Engine. Here are some simple examples, and more detailed examples are
  available in the subcommands or the documentation.

  To import a key into the specified Transit or Transform mount:

  $ vault transit import transit/keys/newly-imported @path/to/key type=rsa-2048

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *TransitCommand) Run(args []string) int {
	return cli.RunResultHelp
}
