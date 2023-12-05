// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*TransformCommand)(nil)

type TransformCommand struct {
	*BaseCommand
}

func (c *TransformCommand) Synopsis() string {
	return "Interact with Vault's Transform Secrets Engine"
}

func (c *TransformCommand) Help() string {
	helpText := `
Usage: vault transform <subcommand> [options] [args]

  This command has subcommands for interacting with Vault's Transform Secrets
  Engine. Here are some simple examples, and more detailed examples are
  available in the subcommands or the documentation.

  To import a key into a new FPE transformation:

  $ vault transform import transform/transformations/fpe/new-transformation @path/to/key \
      template=identifier \
	  allowed_roles=physical-access 

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *TransformCommand) Run(args []string) int {
	return cli.RunResultHelp
}
