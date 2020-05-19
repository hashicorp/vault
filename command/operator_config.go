package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*OperatorConfigCommand)(nil)

type OperatorConfigCommand struct {
	*BaseCommand
}

func (c *OperatorConfigCommand) Synopsis() string {
	return "Manage sensitive values in Vault's configuration files"
}

func (c *OperatorConfigCommand) Help() string {
	helpText := `
Usage: vault operator config <subcommand> [options] [args]

  This command groups subcommands for operators interacting with Vault's
  config files. Here are a few examples of config operator commands:

  Encrypt sensitive values in a config file:

      $ vault operator config encrypt config.hcl

  Decrypt sensitive values in a config file:

      $ vault operator config decrypt config.hcl

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *OperatorConfigCommand) Run(args []string) int {
	return cli.RunResultHelp
}
