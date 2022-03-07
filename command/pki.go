package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*PKICommand)(nil)

type PKICommand struct {
	*BaseCommand
}

func (c *PKICommand) Synopsis() string {
	return "Interact with PKI Secret Engines"
}

func (c *PKICommand) Help() string {
	helpText := `
Usage: vault pki <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's PKI Secrets
  Engine. Operators can manage PKI mounts and roles.

  To test role based issuance:

       $ vault pki role-test -mount=pki-int server-role example.com

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PKICommand) Run(args []string) int {
	return cli.RunResultHelp
}
