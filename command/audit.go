// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*AuditCommand)(nil)

type AuditCommand struct {
	*BaseCommand
}

func (c *AuditCommand) Synopsis() string {
	return "Interact with audit devices"
}

func (c *AuditCommand) Help() string {
	helpText := `
Usage: vault audit <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's audit devices.
  Users can list, enable, and disable audit devices.

  *NOTE*: Once an audit device has been enabled, failure to audit could prevent
  Vault from servicing future requests. It is highly recommended that you enable
  multiple audit devices.

  List all enabled audit devices:

      $ vault audit list

  Enable a new audit device "file";

       $ vault audit enable file file_path=/var/log/audit.log

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *AuditCommand) Run(args []string) int {
	return cli.RunResultHelp
}
