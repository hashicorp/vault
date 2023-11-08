// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginDeregisterCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginDeregisterCommand)(nil)
)

type PluginDeregisterCommand struct {
	*BaseCommand

	flagPluginVersion string
}

func (c *PluginDeregisterCommand) Synopsis() string {
	return "Deregister an existing plugin in the catalog"
}

func (c *PluginDeregisterCommand) Help() string {
	helpText := `
Usage: vault plugin deregister [options] TYPE NAME

  Deregister an existing plugin in the catalog. If the plugin does not exist,
  no action is taken (the command is idempotent). The TYPE argument
  takes "auth", "database", or "secret".

  Deregister the unversioned auth plugin named my-custom-plugin:

      $ vault plugin deregister auth my-custom-plugin

  Deregister the auth plugin named my-custom-plugin, version 1.0.0:

      $ vault plugin deregister -version=v1.0.0 auth my-custom-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginDeregisterCommand) Flags() *FlagSets {
	set := c.FlagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "version",
		Target:     &c.flagPluginVersion,
		Completion: complete.PredictAnything,
		Usage: "Semantic version of the plugin to deregister. If unset, " +
			"only an unversioned plugin may be deregistered.",
	})

	return set
}

func (c *PluginDeregisterCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins(api.PluginTypeUnknown)
}

func (c *PluginDeregisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginDeregisterCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var pluginNameRaw, pluginTypeRaw string
	args = f.Args()
	switch len(args) {
	case 0:
		c.UI.Error("Not enough arguments (expected 1, or 2, got 0)")
		return 1
	case 1:
		pluginTypeRaw = "unknown"
		pluginNameRaw = args[0]
	case 2:
		pluginTypeRaw = args[0]
		pluginNameRaw = args[1]
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, or 2, got %d)", len(args)))
		return 1
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
	if c.flagPluginVersion != "" {
		_, err := semver.NewSemver(c.flagPluginVersion)
		if err != nil {
			c.UI.Error(fmt.Sprintf("version %q is not a valid semantic version: %v", c.flagPluginVersion, err))
			return 2
		}
	}

	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    pluginName,
		Type:    pluginType,
		Version: c.flagPluginVersion,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error deregistering plugin named %s: %s", pluginName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Deregistered plugin (if it was registered): %s", pluginName))
	return 0
}
