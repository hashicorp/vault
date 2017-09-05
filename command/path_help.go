package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*PathHelpCommand)(nil)
var _ cli.CommandAutocomplete = (*PathHelpCommand)(nil)

var pathHelpVaultSealedMessage = strings.TrimSpace(`
Error: Vault is sealed.

The "path-help" command requires the Vault to be unsealed so that the mount
points of the secret backends are known.
`)

// PathHelpCommand is a Command that lists the mounts.
type PathHelpCommand struct {
	*BaseCommand
}

func (c *PathHelpCommand) Synopsis() string {
	return "Retrieves API help for paths"
}

func (c *PathHelpCommand) Help() string {
	helpText := `
Usage: vault path-help [options] path

  Retrieves API help for paths. All endpoints in Vault provide built-in help
  in markdown format. This includes system paths, secret paths, and credential
  providers.

  A backend must be mounted before help is available:

      $ vault mount database
      $ vault path-help database/

  The response object will return additional paths to retrieve help:

      $ vault path-help database/roles/

  Each backend produces different help output. For additional information,
  please view the online documentation.

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
