package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginReloadCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginReloadCommand)(nil)

type PluginReloadCommand struct {
	*BaseCommand
	plugin string
	mounts []string
}

func (c *PluginReloadCommand) Synopsis() string {
	return "Reload mounted plugin backend"
}

func (c *PluginReloadCommand) Help() string {
	helpText := `
Usage: vault plugin reload [options]

  Reloads mounted plugin backends. Either the plugin name or the desired plugin 
  backend mounts (mounts) must be provided, but not both.

  Reload the plugin named my-custom-plugin:

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
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: c.plugin,
		Mounts: c.mounts,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error reloading plugin/mounts: %s", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Reloaded plugin/mounts: %s%s", c.plugin, c.mounts))
	return 0
}
