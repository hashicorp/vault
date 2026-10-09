// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

var _ cli.Command = (*TPMCommand)(nil)

type TPMCommand struct {
	*base.BaseCommand
}

func (c *TPMCommand) Synopsis() string {
	return "Interact with Vault's TPM authentication"
}

func (c *TPMCommand) Help() string {
	helpText := `
Usage: vault tpm <subcommand> [options] [args]

  Interact with Vault's TPM authentication. The primary command is "attest",
  which opens the local TPM device, runs the workflow against a Vault TPM auth 
  method mount, and writes the issued certificate and key material to a state directory.

      $ vault tpm attest -role-name=my-role

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *TPMCommand) Run(args []string) int {
	return cli.RunResultHelp
}
