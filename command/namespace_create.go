package command

import (
	"fmt"
	"path"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*NamespaceCreateCommand)(nil)
var _ cli.CommandAutocomplete = (*NamespaceCreateCommand)(nil)

type NamespaceCreateCommand struct {
	*BaseCommand
}

func (c *NamespaceCreateCommand) Synopsis() string {
	return "Create a new namespace"
}

func (c *NamespaceCreateCommand) Help() string {
	helpText := `
Usage: vault namespace create [options] PATH

  Create a child namespace. The namespace created will be relative to the 
  namespace provided in either VAULT_NAMESPACE environemnt variable or
  -namespace CLI flag.

  Create a child namespace (e.g. ns1/):

      $ vault namespace create ns1

  Create a child namespace from a parent namespace (e.g. ns1/ns2/):

      $ vault namespace create -namespace=ns1 ns2

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceCreateCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *NamespaceCreateCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *NamespaceCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceCreateCommand) Run(args []string) int {
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

	namespacePath := strings.TrimSpace(args[0])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	err = client.Sys().CreateNamespace(namespacePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating namespace: %s", err))
		return 2
	}

	// Output full path
	fullPath := path.Join(c.flagNamespace, namespacePath) + "/"
	c.UI.Output(fmt.Sprintf("Success! Namespace created at: %s", fullPath))
	return 0
}
