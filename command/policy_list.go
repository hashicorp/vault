package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PolicyListCommand)(nil)
var _ cli.CommandAutocomplete = (*PolicyListCommand)(nil)

type PolicyListCommand struct {
	*BaseCommand
}

func (c *PolicyListCommand) Synopsis() string {
	return "Lists the installed policies"
}

func (c *PolicyListCommand) Help() string {
	helpText := `
Usage: vault policy list [options]

  Lists the names of the policies that are installed on the Vault server.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *PolicyListCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PolicyListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PolicyListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) > 0:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	policies, err := client.Sys().ListPolicies()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing policies: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
		for _, p := range policies {
			c.UI.Output(p)
		}
		return 0
	default:
		return OutputData(c.UI, policies)
	}
}
