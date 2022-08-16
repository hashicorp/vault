package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TransitImportVersionCommand)(nil)
	_ cli.CommandAutocomplete = (*TransitImportVersionCommand)(nil)
)

type TransitImportVersionCommand struct {
	*BaseCommand
}

func (c *TransitImportVersionCommand) Synopsis() string {
	return "Imports a new key version into a Transit key"
}

func (c *TransitImportVersionCommand) Help() string {
	helpText := `
Usage: vault transit import_version [options]

  

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *TransitImportVersionCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *TransitImportVersionCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitImportVersionCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitImportVersionCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Too few arguments (expected 3, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	return doImport(c.UI, "import_version", args, client)
}
