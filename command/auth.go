package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*AuthCommand)(nil)

type AuthCommand struct {
	*BaseCommand
}

func (c *AuthCommand) Synopsis() string {
	return "Interact with auth methods"
}

func (c *AuthCommand) Help() string {
	return strings.TrimSpace(`
Usage: vault auth <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's auth methods.
  Users can list, enable, disable, and get help for different auth methods.

  To authenticate to Vault as a user or machine, use the "vault login" command
  instead. This command is for interacting with the auth methods themselves, not
  authenticating to Vault.

  List all enabled auth methods:

      $ vault auth list

  Enable a new auth method "userpass";

      $ vault auth enable userpass

  Get detailed help information about how to authenticate to a particular auth
  method:

      $ vault auth help github

  Please see the individual subcommand help for detailed usage information.
`)
}

func (c *AuthCommand) Run(args []string) int {
	return cli.RunResultHelp
}
