package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*NamespaceListCommand)(nil)
var _ cli.CommandAutocomplete = (*NamespaceListCommand)(nil)

type NamespaceListCommand struct {
	*BaseCommand
}

func (c *NamespaceListCommand) Synopsis() string {
	return "List child namespaces"
}

func (c *NamespaceListCommand) Help() string {
	helpText := `
Usage: vault namespaces list [options]

  Lists the enabled child namespaces.

  List all enabled child namespaces:

      $ vault namespaces list

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *NamespaceListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *NamespaceListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceListCommand) Run(args []string) int {
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

	namespaces, err := client.Sys().ListNamespaces()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing namespaces: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
		for _, ns := range namespaces.NamespacePaths {
			c.UI.Output(ns)
		}
		return 0
	default:
		return OutputData(c.UI, namespaces)
	}
}
