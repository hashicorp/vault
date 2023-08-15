// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TransformImportVersionCommand)(nil)
	_ cli.CommandAutocomplete = (*TransformImportVersionCommand)(nil)
)

type TransformImportVersionCommand struct {
	*BaseCommand
}

func (c *TransformImportVersionCommand) Synopsis() string {
	return "Import key material into a new key version in the Transform secrets engines."
}

func (c *TransformImportVersionCommand) Help() string {
	helpText := `
Usage: vault transform import-version PATH KEY [...]

  Using the Transform key wrapping system, imports new key material from
  the base64 encoded KEY (either directly on the CLI or via @path notation),
  into an existing tokenization transformation whose API path is PATH. 

  The remaining options after KEY (key=value style) are passed on to 
  Create/Update Tokenization Transformation API endpoint.

  For example:
  $ vault transform import-version transform/transformations/tokenization/application-form @path/to/new_version \        
       allowed_roles=legacy-system
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransformImportVersionCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *TransformImportVersionCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransformImportVersionCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransformImportVersionCommand) Run(args []string) int {
	return ImportKey(c.BaseCommand, "import_version", transformImportKeyPath, c.Flags(), args)
}
