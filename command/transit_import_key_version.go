// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TransitImportVersionCommand)(nil)
	_ cli.CommandAutocomplete = (*TransitImportVersionCommand)(nil)
)

type TransitImportVersionCommand struct {
	*BaseCommand
}

func (c *TransitImportVersionCommand) Synopsis() string {
	return "Import key material into a new key version in the Transit secrets engines."
}

func (c *TransitImportVersionCommand) Help() string {
	helpText := `
Usage: vault transit import-version PATH KEY [...]

  Using the Transit key wrapping system, imports key material from
  the base64 encoded KEY (either directly on the CLI or via @path notation),
  into a new key whose API path is PATH.  To import a new Transit
  key, use the import command instead.  The remaining options after KEY
  (key=value style) are passed on to the Transit create key endpoint.
  If your system or device natively supports the RSA AES key wrap mechanism
  (such as the PKCS#11 mechanism CKM_RSA_AES_KEY_WRAP), you should use it
  directly rather than this command.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransitImportVersionCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *TransitImportVersionCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitImportVersionCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitImportVersionCommand) Run(args []string) int {
	return ImportKey(c.BaseCommand, "import_version", transitImportKeyPath, c.Flags(), args)
}
