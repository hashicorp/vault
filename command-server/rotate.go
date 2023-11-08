// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command_server

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/command"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRotateCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRotateCommand)(nil)
)

type OperatorRotateCommand struct {
	*command.BaseCommand
}

func (c *OperatorRotateCommand) Synopsis() string {
	return "Rotates the underlying encryption key"
}

func (c *OperatorRotateCommand) Help() string {
	helpText := `
Usage: vault operator rotate [options]

  Rotates the underlying encryption key which is used to secure data written
  to the storage backend. This installs a new key in the key ring. This new
  key is used to encrypted new data, while older keys in the ring are used to
  decrypt older data.

  This is an online operation and does not cause downtime. This command is run
  per-cluster (not per-server), since Vault servers in HA mode share the same
  storage backend.

  Rotate Vault's encryption key:

      $ vault operator rotate

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRotateCommand) Flags() *command.FlagSets {
	return c.FlagSet(command.FlagSetHTTP | command.FlagSetOutputFormat)
}

func (c *OperatorRotateCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *OperatorRotateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRotateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Rotate the key
	err = client.Sys().Rotate()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error rotating key: %s", err))
		return 2
	}

	// Print the key status
	status, err := client.Sys().KeyStatus()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading key status: %s", err))
		return 2
	}

	switch command.Format(c.UI) {
	case "table":
		c.UI.Output("Success! Rotated key")
		c.UI.Output("")
		c.UI.Output(command.PrintKeyStatus(status))
		return 0
	default:
		return command.OutputData(c.UI, status)
	}
}
