package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PathHelpCommand)(nil)
var _ cli.CommandAutocomplete = (*PathHelpCommand)(nil)

var pathHelpVaultSealedMessage = strings.TrimSpace(`
Error: Vault is sealed.

The "path-help" command requires the Vault to be unsealed so that the mount
points of the secret engines are known.
`)

type PathHelpCommand struct {
	*BaseCommand
}

func (c *PathHelpCommand) Synopsis() string {
	return "Retrieve API help for paths"
}

func (c *PathHelpCommand) Help() string {
	helpText := `
Usage: vault path-help [options] PATH

  Retrieves API help for paths. All endpoints in Vault provide built-in help
  in markdown format. This includes system paths, secret engines, and auth
  methods.

  Get help for the thing mounted at database/:

      $ vault path-help database/

  The response object will return additional paths to retrieve help:

      $ vault path-help database/roles/

  Each secret engine produces different help output.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PathHelpCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PathHelpCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything // TODO: programatic way to invoke help
}

func (c *PathHelpCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PathHelpCommand) Run(args []string) int {
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

	path := sanitizePath(args[0])

	help, err := client.Help(path)
	if err != nil {
		if strings.Contains(err.Error(), "Vault is sealed") {
			c.UI.Error(pathHelpVaultSealedMessage)
		} else {
			c.UI.Error(fmt.Sprintf("Error retrieving help: %s", err))
		}
		return 2
	}

	c.UI.Output(help.Help)
	return 0
}
