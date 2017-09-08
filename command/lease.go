package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*LeaseCommand)(nil)

type LeaseCommand struct {
	*BaseCommand
}

func (c *LeaseCommand) Synopsis() string {
	return "Interact with leases"
}

func (c *LeaseCommand) Help() string {
	helpText := `
Usage: vault lease <subcommand> [options] [args]

  This command groups subcommands for interacting with leases. Users can revoke
  or renew leases.

  Renew a lease:

      $ vault lease renew database/creds/readonly/2f6a614c...

  Revoke a lease:

      $ vault lease revoke database/creds/readonly/2f6a614c...
`

	return strings.TrimSpace(helpText)
}

func (c *LeaseCommand) Run(args []string) int {
	return cli.RunResultHelp
}
