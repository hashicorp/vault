// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginReloadCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginReloadCommand)(nil)
)

type PluginReloadCommand struct {
	*BaseCommand
	plugin string
	mounts []string
	scope  string
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

  Reload the plugin named "my-custom-plugin":

	  $ vault plugin reload -plugin=my-custom-plugin

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
		Usage:      "The scope of the reload, omitted for local, 'global', for replicated reloads",
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

	switch {
	case c.plugin == "" && len(c.mounts) == 0:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case c.plugin != "" && len(c.mounts) > 0:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	case c.scope != "" && c.scope != "global":
		c.UI.Error(fmt.Sprintf("Invalid reload scope: %s", c.scope))
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	rid, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: c.plugin,
		Mounts: c.mounts,
		Scope:  c.scope,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reloading plugin/mounts: %s", err))
		return 2
	}

	if len(c.mounts) > 0 {
		if rid != "" {
			c.UI.Output(fmt.Sprintf("Success! Reloading mounts: %s, reload_id: %s", c.mounts, rid))
		} else {
			c.UI.Output(fmt.Sprintf("Success! Reloaded mounts: %s", c.mounts))
		}
	} else {
		if rid != "" {
			c.UI.Output(fmt.Sprintf("Success! Reloading plugin: %s, reload_id: %s", c.plugin, rid))
		} else {
			c.UI.Output(fmt.Sprintf("Success! Reloaded plugin: %s", c.plugin))
		}
	}

	return 0
}
