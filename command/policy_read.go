package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PolicyReadCommand)(nil)
var _ cli.CommandAutocomplete = (*PolicyReadCommand)(nil)

type PolicyReadCommand struct {
	*BaseCommand
}

func (c *PolicyReadCommand) Synopsis() string {
	return "Prints the contents of a policy"
}

func (c *PolicyReadCommand) Help() string {
	helpText := `
Usage: vault policy read [options] [NAME]

  Prints the contents and metadata of the Vault policy named NAME. If the policy
  does not exist, an error is returned.

  Read the policy named "my-policy":

      $ vault policy read my-policy

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyReadCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *PolicyReadCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPolicies()
}

func (c *PolicyReadCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PolicyReadCommand) Run(args []string) int {
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

	name := strings.ToLower(strings.TrimSpace(args[0]))
	rules, err := client.Sys().GetPolicy(name)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading policy named %s: %s", name, err))
		return 2
	}
	if rules == "" {
		c.UI.Error(fmt.Sprintf("No policy named: %s", name))
		return 2
	}

	switch Format(c.UI) {
	case "table":
		c.UI.Output(strings.TrimSpace(rules))
		return 0
	default:
		resp := map[string]string{
			"policy": rules,
		}
		return OutputData(c.UI, &resp)
	}
}
