// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

var _ cli.Command = (*TokenCommand)(nil)

type TokenCommand struct {
	*base.BaseCommand
}

func (c *TokenCommand) Synopsis() string {
	return "Interact with tokens"
}

func (c *TokenCommand) Help() string {
	helpText := `
Usage: ` + c.BinaryName + ` token <subcommand> [options] [args]

  This command groups subcommands for interacting with tokens. Users can
  create, lookup, renew, and revoke tokens.

  Create a new token:

      $ ` + c.BinaryName + ` token create

  Revoke a token:

      $ ` + c.BinaryName + ` token revoke 96ddf4bc-d217-f3ba-f9bd-017055595017

  Renew a token:

      $ ` + c.BinaryName + ` token renew 96ddf4bc-d217-f3ba-f9bd-017055595017

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *TokenCommand) Run(args []string) int {
	return cli.RunResultHelp
}
