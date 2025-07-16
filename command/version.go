// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/version"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*VersionCommand)(nil)
	_ cli.CommandAutocomplete = (*VersionCommand)(nil)
)

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	*BaseCommand

	VersionInfo *version.VersionInfo
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the Vault CLI version"
}

func (c *VersionCommand) Help() string {
	helpText := `
Usage: vault version

  Prints the version of this Vault CLI. This does not print the target Vault
  server version.

  Print the version:

      $ vault version

  There are no arguments or flags to this command. Any additional arguments or
  flags are ignored.
`
	return strings.TrimSpace(helpText)
}

func (c *VersionCommand) Flags() *FlagSets {
	return nil
}

func (c *VersionCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *VersionCommand) AutocompleteFlags() complete.Flags {
	return nil
}

func (c *VersionCommand) Run(_ []string) int {
	out := c.VersionInfo.FullVersionNumber(true)
	if version.CgoEnabled {
		out += " (cgo)"
	}
	c.UI.Output(out)
	return 0
}
