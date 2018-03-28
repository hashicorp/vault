package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*PolicyCommand)(nil)

// PolicyCommand is a Command that holds the audit commands
type PolicyCommand struct {
	*BaseCommand
}

func (c *PolicyCommand) Synopsis() string {
	return "Interact with policies"
}

func (c *PolicyCommand) Help() string {
	helpText := `
Usage: vault policy <subcommand> [options] [args]

  This command groups subcommands for interacting with policies.
  Users can write, read, and list policies in Vault.

  List all enabled policies:

      $ vault policy list

  Create a policy named "my-policy" from contents on local disk:

      $ vault policy write my-policy ./my-policy.hcl

  Delete the policy named my-policy:

      $ vault policy delete my-policy

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PolicyCommand) Run(args []string) int {
	return cli.RunResultHelp
}
