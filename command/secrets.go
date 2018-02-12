package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*SecretsCommand)(nil)

type SecretsCommand struct {
	*BaseCommand
}

func (c *SecretsCommand) Synopsis() string {
	return "Interact with secrets engines"
}

func (c *SecretsCommand) Help() string {
	helpText := `
Usage: vault secrets <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's secrets engines.
  Each secret engine behaves differently. Please see the documentation for
  more information.

  List all enabled secrets engines:

      $ vault secrets list

  Enable a new secrets engine:

      $ vault secrets enable database

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *SecretsCommand) Run(args []string) int {
	return cli.RunResultHelp
}
