package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/posener/complete"
)

var _ cli.Command = (*PolicyFmtCommand)(nil)
var _ cli.CommandAutocomplete = (*PolicyFmtCommand)(nil)

type PolicyFmtCommand struct {
	*BaseCommand
}

func (c *PolicyFmtCommand) Synopsis() string {
	return "Formats a policy on disk"
}

func (c *PolicyFmtCommand) Help() string {
	helpText := `
Usage: vault policy fmt [options] PATH

  Formats a local policy file to the policy specification. This command will
  overwrite the file at the given PATH with the properly-formatted policy
  file contents.

  Format the local file "my-policy.hcl" as a policy file:

      $ vault policy fmt my-policy.hcl

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyFmtCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetNone)
}

func (c *PolicyFmtCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictFiles("*.hcl")
}

func (c *PolicyFmtCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PolicyFmtCommand) Run(args []string) int {
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

	// Get the filepath, accounting for ~ and stuff
	path, err := homedir.Expand(strings.TrimSpace(args[0]))
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to expand path: %s", err))
		return 1
	}

	// Read the entire contents into memory - it would be nice if we could use
	// a buffer, but hcl wants the full contents.
	b, err := ioutil.ReadFile(path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading source file: %s", err))
		return 1
	}

	// Actually parse the policy. We always use the root namespace here because
	// we don't want to modify the results.
	if _, err := vault.ParseACLPolicy(namespace.RootNamespace, string(b)); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Generate final contents
	result, err := printer.Format(b)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error printing result: %s", err))
		return 1
	}

	// Write them back out
	if err := ioutil.WriteFile(path, result, 0644); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing result: %s", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf("Success! Formatted policy: %s", path))
	return 0
}
