// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
)

var _ cli.Command = (*NamespaceCommand)(nil)

type NamespaceCommand struct {
	*BaseCommand
}

func (c *NamespaceCommand) Synopsis() string {
	return "Interact with namespaces"
}

func (c *NamespaceCommand) Help() string {
	helpText := `
Usage: vault namespace <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault namespaces.
  These subcommands operate in the context of the namespace that the
  currently logged in token belongs to.

  List enabled child namespaces:

      $ vault namespace list

  Look up an existing namespace:

      $ vault namespace lookup

  Create a new namespace:

      $ vault namespace create

  Patch an existing namespace:

      $ vault namespace patch

  Delete an existing namespace:

      $ vault namespace delete

  Lock the API for an existing namespace:

      $ vault namespace lock

  Unlock the API for an existing namespace:

      $ vault namespace unlock

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *NamespaceCommand) Run(args []string) int {
	return cli.RunResultHelp
}
