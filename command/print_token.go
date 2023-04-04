// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PrintTokenCommand)(nil)
	_ cli.CommandAutocomplete = (*PrintTokenCommand)(nil)
)

type PrintTokenCommand struct {
	*BaseCommand
}

func (c *PrintTokenCommand) Synopsis() string {
	return "Prints the vault token currently in use"
}

func (c *PrintTokenCommand) Help() string {
	helpText := `
Usage: vault print token

  Prints the value of the Vault token that will be used for commands, after
  taking into account the configured token-helper and the environment.

      $ vault print token

`
	return strings.TrimSpace(helpText)
}

func (c *PrintTokenCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PrintTokenCommand) AutocompleteFlags() complete.Flags {
	return nil
}

func (c *PrintTokenCommand) Run(args []string) int {
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	c.UI.Output(client.Token())
	return 0
}
