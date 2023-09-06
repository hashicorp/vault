// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*PluginRuntimeCommand)(nil)

type PluginRuntimeCommand struct {
	*BaseCommand
}

func (c *PluginRuntimeCommand) Synopsis() string {
	return "Interact with Vault plugins and catalog"
}

func (c *PluginRuntimeCommand) Help() string {
	helpText := `
Usage: vault plugin runtime <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's plugin runtimes and the
  plugin runtime catalog. The plugin runtime catalog is divided into types. Currently,
  Vault only supports "container" plugin runtimes. A type must be specified on each call. Here 
  are a few examples of the plugin runtime commands.

  List all available plugin runtimes in the catalog of a particular type:

      $ vault plugin runtime list -type=container

  Register a new plugin runtime to the catalog as a particular type:

      $ vault plugin runtime register -type=container -oci_runtime=my-oci-runtime my-custom-plugin-runtime

  Get information about a plugin runtime in the catalog listed under a particular type:

      $ vault plugin runtime info -type=container my-custom-plugin

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PluginRuntimeCommand) Run(args []string) int {
	return cli.RunResultHelp
}
