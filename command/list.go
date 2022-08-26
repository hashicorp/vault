package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*ListCommand)(nil)
	_ cli.CommandAutocomplete = (*ListCommand)(nil)
)

type ListCommand struct {
	*BaseCommand
}

func (c *ListCommand) Synopsis() string {
	return "List data or secrets"
}

func (c *ListCommand) Help() string {
	helpText := `

Usage: vault list [options] PATH

  Lists data from Vault at the given path. This can be used to list keys in a,
  given secret engine.

  List values under the "my-app" folder of the generic secret engine:

      $ vault list secret/my-app/

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret engine in use. Not all engines support listing.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat | FlagSetOutputDetailed)
	return set
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

	// Append trailing slash
	path := args[0]
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	path = sanitizePath(path)
	secret, err := client.Logical().List(path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing %s: %s", path, err))
		return 2
	}

	// If the secret is wrapped, return the wrapped response.
	if secret != nil && secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, secret)
	}

	_, ok := extractListData(secret)
	if Format(c.UI) != "table" {
		if secret == nil || secret.Data == nil || !ok {
			OutputData(c.UI, map[string]interface{}{})
			return 2
		}
	}

	if secret == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}
	if secret.Data == nil {
		// If secret wasn't nil, we have warnings, so output them anyways. We
		// may also have non-keys info.
		return OutputSecret(c.UI, secret)
	}

	if !ok {
		c.UI.Error(fmt.Sprintf("No entries found at %s", path))
		return 2
	}

	return OutputList(c.UI, secret)
}
