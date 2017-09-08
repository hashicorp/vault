package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*TokenCommand)(nil)

type TokenCommand struct {
	*BaseCommand
}

func (c *TokenCommand) Synopsis() string {
	return "Interact with tokens"
}

func (c *TokenCommand) Help() string {
	helpText := `
Usage: vault token <subcommand> [options] [args]

  This command groups subcommands for interacting with tokens. Users can
  create, lookup, renew, and revoke tokens.

  Create a new token:

      $ vault token create

  Revoke a token:

      $ vault token revoke 96ddf4bc-d217-f3ba-f9bd-017055595017

  Renew a token:

      $ vault token renew 96ddf4bc-d217-f3ba-f9bd-017055595017

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *TokenCommand) Run(args []string) int {
	return cli.RunResultHelp
}
