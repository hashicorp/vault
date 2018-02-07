package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PolicyDeleteCommand)(nil)
var _ cli.CommandAutocomplete = (*PolicyDeleteCommand)(nil)

type PolicyDeleteCommand struct {
	*BaseCommand
}

func (c *PolicyDeleteCommand) Synopsis() string {
	return "Deletes a policy by name"
}

func (c *PolicyDeleteCommand) Help() string {
	helpText := `
Usage: vault policy delete [options] NAME

  Deletes the policy named NAME in the Vault server. Once the policy is deleted,
  all tokens associated with the policy are affected immediately.

  Delete the policy named "my-policy":

      $ vault policy delete my-policy

  Note that it is not possible to delete the "default" or "root" policies.
  These are built-in policies.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyDeleteCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PolicyDeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPolicies()
}

func (c *PolicyDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PolicyDeleteCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	name := strings.TrimSpace(strings.ToLower(args[0]))
	if err := client.Sys().DeletePolicy(name); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", name, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Deleted policy: %s", name))
	return 0
}
