package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*NamespaceDeleteCommand)(nil)
var _ cli.CommandAutocomplete = (*NamespaceDeleteCommand)(nil)

type NamespaceDeleteCommand struct {
	*BaseCommand
}

func (c *NamespaceDeleteCommand) Synopsis() string {
	return "Delete an existing namespace"
}

func (c *NamespaceDeleteCommand) Help() string {
	helpText := `
Usage: vault namespace delete [options] PATH

  Delete an existing namespace. The namespace deleted will be relative to the
  namespace provided in either the VAULT_NAMESPACE environment variable or
  -namespace CLI flag.

  Delete a namespace (e.g. ns1/):

      $ vault namespace delete ns1

  Delete a namespace namespace from a parent namespace (e.g. ns1/ns2/):

      $ vault namespace create -namespace=ns1 ns2

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceDeleteCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *NamespaceDeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *NamespaceDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceDeleteCommand) Run(args []string) int {
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

	secret, err := client.Logical().Delete("sys/namespaces/" + namespacePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting namespace: %s", err))
		return 2
	}

	if secret != nil {
		// Likely, we have warnings
		return OutputSecret(c.UI, secret)
	}

	if !strings.HasSuffix(namespacePath, "/") {
		namespacePath = namespacePath + "/"
	}

	c.UI.Output(fmt.Sprintf("Success! Namespace deleted at: %s", namespacePath))
	return 0
}
