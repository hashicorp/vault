// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginReloadCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginReloadCommand)(nil)
)

type PluginReloadCommand struct {
	*BaseCommand
	plugin     string
	mounts     []string
	scope      string
	pluginType string
}

func (c *PluginReloadCommand) Synopsis() string {
	return "Reload mounted plugin backend"
}

func (c *PluginReloadCommand) Help() string {
	helpText := `
Usage: vault plugin reload [options]

  Reloads mounted plugins. Either the plugin name or the desired plugin 
  mount(s) must be provided, but not both. In case the plugin name is provided,
  all of its corresponding mounted paths that use the plugin backend will be reloaded.

  If run with a Vault namespace other than the root namespace, only plugins
  running in the same namespace will be reloaded.

  Reload the secret plugin named "my-custom-plugin" on the current node:

      $ vault plugin reload -type=secret -plugin=my-custom-plugin

  Reload the secret plugin named "my-custom-plugin" across all nodes and replicated clusters:

      $ vault plugin reload -type=secret -plugin=my-custom-plugin -scope=global

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginReloadCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "plugin",
		Target:     &c.plugin,
		Completion: complete.PredictAnything,
		Usage:      "The name of the plugin to reload, as registered in the plugin catalog.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:       "mounts",
		Target:     &c.mounts,
		Completion: complete.PredictAnything,
		Usage:      "Array or comma-separated string mount paths of the plugin backends to reload.",
	})

	f.StringVar(&StringVar{
		Name:       "scope",
		Target:     &c.scope,
		Completion: complete.PredictAnything,
		Usage:      "The scope of the reload, omitted for local, 'global', for replicated reloads.",
	})

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.pluginType,
		Completion: complete.PredictAnything,
		Usage: "The type of plugin to reload, one of auth, secret, or database. Mutually " +
			"exclusive with -mounts. If not provided, all plugins with a matching name will be reloaded.",
	})

	return set
}

func (c *PluginReloadCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PluginReloadCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginReloadCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	positionalArgs := len(f.Args())
	switch {
	case positionalArgs != 0:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", positionalArgs))
		return 1
	case c.plugin == "" && len(c.mounts) == 0:
		c.UI.Error("No plugins specified, must specify exactly one of -plugin or -mounts")
		return 1
	case c.plugin != "" && len(c.mounts) > 0:
		c.UI.Error("Must specify exactly one of -plugin or -mounts")
		return 1
	case c.scope != "" && c.scope != "global":
		c.UI.Error(fmt.Sprintf("Invalid reload scope: %s", c.scope))
		return 1
	case len(c.mounts) > 0 && c.pluginType != "":
		c.UI.Error("Cannot specify -type with -mounts")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var reloadID string
	if client.Namespace() == "" {
		pluginType := api.PluginTypeUnknown
		pluginTypeStr := strings.TrimSpace(c.pluginType)
		if pluginTypeStr != "" {
			var err error
			pluginType, err = api.ParsePluginType(pluginTypeStr)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error parsing -type as a plugin type, must be unset or one of auth, secret, or database: %s", err))
				return 1
			}
		}

		reloadID, err = client.Sys().RootReloadPlugin(context.Background(), &api.RootReloadPluginInput{
			Plugin: c.plugin,
			Type:   pluginType,
			Scope:  c.scope,
		})
	} else {
		reloadID, err = client.Sys().ReloadPlugin(&api.ReloadPluginInput{
			Plugin: c.plugin,
			Mounts: c.mounts,
			Scope:  c.scope,
		})
	}

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reloading plugin/mounts: %s", err))
		return 2
	}

	if len(c.mounts) > 0 {
		if reloadID != "" {
			c.UI.Output(fmt.Sprintf("Success! Reloading mounts: %s, reload_id: %s", c.mounts, reloadID))
		} else {
			c.UI.Output(fmt.Sprintf("Success! Reloaded mounts: %s", c.mounts))
		}
	} else {
		if reloadID != "" {
			c.UI.Output(fmt.Sprintf("Success! Reloading plugin: %s, reload_id: %s", c.plugin, reloadID))
		} else {
			c.UI.Output(fmt.Sprintf("Success! Reloaded plugin: %s", c.plugin))
		}
	}

	return 0
}
