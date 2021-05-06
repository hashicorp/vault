package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*LicenseCommand)(nil)

type LicenseCommand struct {
	*BaseCommand
}

func (c *LicenseCommand) Synopsis() string {
	return "Interact with licenses"
}

func (c *LicenseCommand) Help() string {
	helpText := `
Usage: vault license <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault licenses.

  Get the current Vault license:

      $ vault license get

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *LicenseCommand) Run(args []string) int {
	return cli.RunResultHelp
}
