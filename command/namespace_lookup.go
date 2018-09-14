package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*NamespaceLookupCommand)(nil)
var _ cli.CommandAutocomplete = (*NamespaceLookupCommand)(nil)

type NamespaceLookupCommand struct {
	*BaseCommand
}

func (c *NamespaceLookupCommand) Synopsis() string {
	return "Look up an existing namespace"
}

func (c *NamespaceLookupCommand) Help() string {
	helpText := `
Usage: vault namespace create [options] PATH

  Create a child namespace. The namespace created will be relative to the
  namespace provided in either the VAULT_NAMESPACE environment variable or
  -namespace CLI flag.

  Get information about the namespace of the locally authenticated token:

      $ vault namespace lookup

  Get information about the namespace of a particular child token (e.g. ns1/ns2/):

      $ vault namespace lookup -namespace=ns1 ns2

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceLookupCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *NamespaceLookupCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *NamespaceLookupCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceLookupCommand) Run(args []string) int {
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

	secret, err := client.Logical().Read("sys/namespaces/" + namespacePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error looking up namespace: %s", err))
		return 2
	}
	if secret == nil {
		c.UI.Error("Namespace not found")
		return 2
	}

	return OutputSecret(c.UI, secret)
}
