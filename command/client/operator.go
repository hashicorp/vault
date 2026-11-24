// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

var _ cli.Command = (*OperatorCommand)(nil)

type OperatorCommand struct {
	*base.BaseCommand

	// clientHelp limits the examples to the operator subcommands that
	// vault-client supports.
	clientHelp bool
}

func (c *OperatorCommand) Synopsis() string {
	return "Perform operator-specific tasks"
}

func (c *OperatorCommand) Help() string {
	if c.clientHelp {
		return c.clientHelpText()
	}

	helpText := `
Usage: ` + c.BinaryName + ` operator <subcommand> [options] [args]

  This command groups subcommands for operators interacting with Vault. Most
  users will not need to interact with these commands. Here are a few examples
  of the operator commands:

  Initialize a new Vault cluster:

      $ ` + c.BinaryName + ` operator init

  Force a Vault to resign leadership in a cluster:

      $ ` + c.BinaryName + ` operator step-down

  Rotate Vault's underlying encryption key:

      $ ` + c.BinaryName + ` operator rotate

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *OperatorCommand) clientHelpText() string {
	helpText := `
Usage: ` + c.BinaryName + ` operator <subcommand> [options] [args]

  This command groups subcommands for operators interacting with Vault. Most
  users will not need to interact with these commands. Here are examples of
  the supported operator commands:

  Show information about Vault's active encryption key:

      $ ` + c.BinaryName + ` operator key-status

  Rotate Vault's underlying encryption key:

      $ ` + c.BinaryName + ` operator rotate

  Provide an unseal key to unseal a Vault server:

      $ ` + c.BinaryName + ` operator unseal

  Other operator subcommands require the full Vault Server image.

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *OperatorCommand) Run(args []string) int {
	return cli.RunResultHelp
}
