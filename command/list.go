package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*ListCommand)(nil)
var _ cli.CommandAutocomplete = (*ListCommand)(nil)

// ListCommand is a Command that lists data from the Vault.
type ListCommand struct {
	*BaseCommand
}

func (c *ListCommand) Synopsis() string {
	return "Lists data or secrets"
}

func (c *ListCommand) Help() string {
	helpText := `

Usage: vault list [options] PATH

  Lists data from Vault at the given path. This can be used to list keys in a,
  given backend.

  List values under the "my-app" folder:

      $ vault list secret/my-app/

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret backend in use. Not all backends support listing.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *ListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *ListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	path, kvs, err := extractPath(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if len(kvs) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().List(path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing %s: %s", path, err))
		return 2
	}
	if secret == nil || secret.Data == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}

	// If the secret is wrapped, return the wrapped response.
	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, c.flagFormat, secret)
	}

	if _, ok := extractListData(secret); !ok {
		c.UI.Error(fmt.Sprintf("No entries found at %s", path))
		return 2
	}

	return OutputList(c.UI, c.flagFormat, secret)
}
