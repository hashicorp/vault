// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	return "Interact with Vault's PKI Secrets Engine"
}

func (c *PKICommand) Help() string {
	helpText := `
Usage: vault pki <subcommand> [options] [args]

  This command has subcommands for interacting with Vault's PKI Secrets
  Engine. Here are some simple examples, and more detailed examples are
  available in the subcommands or the documentation.

  Check the health of a PKI mount, to the best of this token's abilities:

      $ vault pki health-check pki

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PKICommand) Run(args []string) int {
	return cli.RunResultHelp
}
