// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

var _ cli.Command = (*NamespaceCommand)(nil)

type NamespaceCommand struct {
	*base.BaseCommand
}

func (c *NamespaceCommand) Synopsis() string {
	return "Interact with namespaces"
}

func (c *NamespaceCommand) Help() string {
	helpText := `
Usage: ` + c.BinaryName + ` namespace <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault namespaces.
  These subcommands operate in the context of the namespace that the
  currently logged in token belongs to.

  List enabled child namespaces:

      $ ` + c.BinaryName + ` namespace list

  Look up an existing namespace:

      $ ` + c.BinaryName + ` namespace lookup

  Create a new namespace:

      $ ` + c.BinaryName + ` namespace create

  Patch an existing namespace:

      $ ` + c.BinaryName + ` namespace patch

  Delete an existing namespace:

      $ ` + c.BinaryName + ` namespace delete

  Lock the API for an existing namespace:

      $ ` + c.BinaryName + ` namespace lock

  Unlock the API for an existing namespace:

      $ ` + c.BinaryName + ` namespace unlock

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *NamespaceCommand) Run(args []string) int {
	return cli.RunResultHelp
}
