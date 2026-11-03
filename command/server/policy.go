// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

var _ cli.Command = (*PolicyCommand)(nil)

// PolicyCommand is a Command that holds the audit commands
type PolicyCommand struct {
	*base.BaseCommand
}

func (c *PolicyCommand) Synopsis() string {
	return "Interact with policies"
}

func (c *PolicyCommand) Help() string {
	helpText := `
Usage: ` + c.BinaryName + ` policy <subcommand> [options] [args]

  This command groups subcommands for interacting with policies.
  Users can write, read, and list policies in Vault.

  List all enabled policies:

      $ ` + c.BinaryName + ` policy list

  Create a policy named "my-policy" from contents on local disk:

      $ ` + c.BinaryName + ` policy write my-policy ./my-policy.hcl

  Delete the policy named my-policy:

      $ ` + c.BinaryName + ` policy delete my-policy

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PolicyCommand) Run(args []string) int {
	return cli.RunResultHelp
}
