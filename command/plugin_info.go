// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginInfoCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginInfoCommand)(nil)
)

type PluginInfoCommand struct {
	*BaseCommand

	flagVersion string
}

func (c *PluginInfoCommand) Synopsis() string {
	return "Read information about a plugin in the catalog"
}

func (c *PluginInfoCommand) Help() string {
	helpText := `
Usage: vault plugin info [options] TYPE NAME

  Displays information about a plugin in the catalog with the given name. If
  the plugin does not exist, an error is returned. The argument of type
  takes "auth", "database", or "secret".

  Get info about a plugin:

      $ vault plugin info database mysql-database-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginInfoCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "version",
		Target:     &c.flagVersion,
		Completion: complete.PredictAnything,
		Usage:      "Semantic version of the plugin. Optional.",
	})

	return set
}

func (c *PluginInfoCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins(api.PluginTypeUnknown)
}

func (c *PluginInfoCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginInfoCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var pluginNameRaw, pluginTypeRaw string
	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1 or 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1 or 2, got %d)", len(args)))
		return 1

	// These cases should come after invalid cases have been checked
	case len(args) == 1:
		pluginTypeRaw = "unknown"
		pluginNameRaw = args[0]
	case len(args) == 2:
		pluginTypeRaw = args[0]
		pluginNameRaw = args[1]
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	pluginType, err := api.ParsePluginType(strings.TrimSpace(pluginTypeRaw))
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	pluginName := strings.TrimSpace(pluginNameRaw)

	resp, err := client.Sys().GetPlugin(&api.GetPluginInput{
		Name:    pluginName,
		Type:    pluginType,
		Version: c.flagVersion,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading plugin named %s: %s", pluginName, err))
		return 2
	}

	if resp == nil {
		c.UI.Error(fmt.Sprintf("No value found for plugin %q", pluginName))
		return 2
	}

	data := map[string]interface{}{
		"args":               resp.Args,
		"builtin":            resp.Builtin,
		"command":            resp.Command,
		"oci_image":          resp.OCIImage,
		"runtime":            resp.Runtime,
		"name":               resp.Name,
		"sha256":             resp.SHA256,
		"deprecation_status": resp.DeprecationStatus,
		"version":            resp.Version,
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, data, c.flagField)
	}
	return OutputData(c.UI, data)
}
