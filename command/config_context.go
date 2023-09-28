// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*ConfigCommand)(nil)

type ConfigCommand struct {
	*BaseCommand
}

func (c *ConfigCommand) Synopsis() string {
	return "Facilitates access to multiple clusters by using configuration files"
}

func (c *ConfigCommand) Help() string {
	helpText := `
Usage: vault config <subcommand> [options] [args]

  Set a context:

      $ vault config set-context vault_1 --addr=http://127.0.0.1:8200 --token=hvs. --namespace=ns1

  Get a context:

      $ vault config get-context vault_1

  Get the current context:

      $ vault config current-context

  Delete a context:

      $ vault config delete-context vault_1

  Rename a context

      $ vault config rename-context vault_1 vault_1_new_name

  Use a context

      $ vault config use-context vault_2
`

	return strings.TrimSpace(helpText)
}

func (c *ConfigCommand) Run(args []string) int {
	return cli.RunResultHelp
}
