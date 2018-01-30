package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*OperatorCommand)(nil)

type OperatorCommand struct {
	*BaseCommand
}

func (c *OperatorCommand) Synopsis() string {
	return "Perform operator-specific tasks"
}

func (c *OperatorCommand) Help() string {
	helpText := `
Usage: vault operator <subcommand> [options] [args]

  This command groups subcommands for operators interacting with Vault. Most
  users will not need to interact with these commands. Here are a few examples
  of the operator commands:

  Initialize a new Vault cluster:

      $ vault operator init

  Force a Vault to resign leadership in a cluster:

      $ vault operator step-down

  Rotate Vault's underlying encryption key:

      $ vault operator rotate

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *OperatorCommand) Run(args []string) int {
	return cli.RunResultHelp
}
