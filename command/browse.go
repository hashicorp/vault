package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/hashicorp/vault/command/tui"
)

var _ cli.Command = (*BrowseCommand)(nil)

type BrowseCommand struct {
	*BaseCommand
}

func (c *BrowseCommand) Synopsis() string {
	return "Interactively browse Vault's Key-Value storage"
}

func (c *BrowseCommand) Help() string {
	helpText := `
Usage: vault browse [options] PATH

  This command opens a terminal UI for interacting with Vault's key-value
  store. You can browse through the store by using the arrow keys or vim
  keybindings. Quit browsing by pressing ESC or q.

  List values under the "my-app" folder of the generic secret engine:

      $ vault browse secret/my-app/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *BrowseCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *BrowseCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *BrowseCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *BrowseCommand) Run(args []string) int {
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

	path := ensureTrailingSlash(sanitizePath(args[0]))

	t, err := tui.New(client, path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not initialize terminal UI at %s: %s", path, err.Error()))
		return 2
	}

	err = t.Start()
	if err != nil {
		c.UI.Error(fmt.Sprintf("The Terminal UI encountered an error: %s", err.Error()))
		return 2
	}

	return 0
}
