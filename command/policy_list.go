package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*PolicyListCommand)(nil)
var _ cli.CommandAutocomplete = (*PolicyListCommand)(nil)

// PolicyListCommand is a Command that enables a new endpoint.
type PolicyListCommand struct {
	*BaseCommand
}

func (c *PolicyListCommand) Synopsis() string {
	return "Lists the installed policies"
}

func (c *PolicyListCommand) Help() string {
	helpText := `
Usage: vault policies [options] [NAME]

  Lists the policies that are installed on the Vault server. If the optional
  argument is given, this command returns the policy's contents.

  List all policies stored in Vault:

      $ vault policies

  Read the contents of the policy named "my-policy":

      $ vault policies my-policy

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PolicyListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPolicies()
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
	switch len(args) {
	case 0, 1:
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0-2, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	switch len(args) {
	case 0:
		policies, err := client.Sys().ListPolicies()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error listing policies: %s", err))
			return 2
		}
		for _, p := range policies {
			c.UI.Output(p)
		}
	case 1:
		name := strings.ToLower(strings.TrimSpace(args[0]))
		rules, err := client.Sys().GetPolicy(name)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading policy %s: %s", name, err))
			return 2
		}
		if rules == "" {
			c.UI.Error(fmt.Sprintf("Error reading policy: no policy named: %s", name))
			return 2
		}
		c.UI.Output(strings.TrimSpace(rules))
	}

	return 0
}
