package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginListCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginListCommand)(nil)

type PluginListCommand struct {
	*BaseCommand
}

func (c *PluginListCommand) Synopsis() string {
	return "Lists available plugins"
}

func (c *PluginListCommand) Help() string {
	helpText := `
Usage: vault plugin list [options]

  Lists available plugins registered in the catalog. This does not list whether
  plugins are in use, but rather just their availability.

  List all available plugins in the catalog:

      $ vault plugin list

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *PluginListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PluginListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing available plugins: %s", err))
		return 2
	}

	pluginNames := resp.Names
	sort.Strings(pluginNames)

	switch Format(c.UI) {
	case "table":
		list := append([]string{"Plugins"}, pluginNames...)
		c.UI.Output(tableOutput(list, nil))
		return 0
	default:
		return OutputData(c.UI, pluginNames)
	}
}
