// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*DeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*DeleteCommand)(nil)
)

type DeleteCommand struct {
	*BaseCommand

	testStdin io.Reader // for tests
}

func (c *DeleteCommand) Synopsis() string {
	return "Delete secrets and configuration"
}

func (c *DeleteCommand) Help() string {
	helpText := `
Usage: vault delete [options] PATH

  Deletes secrets and configuration from Vault at the given path. The behavior
  of "delete" is delegated to the backend corresponding to the given path.

  Remove data in the status secret backend:

      $ vault delete secret/my-secret

  Uninstall an encryption key in the transit backend:

      $ vault delete transit/keys/my-key

  Delete an IAM role:

      $ vault delete aws/roles/ops

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret backend in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *DeleteCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *DeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *DeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *DeleteCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected at least 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	path := sanitizePath(args[0])

	data, err := parseArgsDataStringLists(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse string list data: %s", err))
		return 1
	}

	secret, err := client.Logical().DeleteWithData(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}

	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data deleted (if it existed) at: %s", path))
		}
		return 0
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
