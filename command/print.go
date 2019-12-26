package command

import (
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PrintCommand)(nil)
var _ cli.CommandAutocomplete = (*PrintCommand)(nil)

type PrintCommand struct {
	*BaseCommand
}

func (c *PrintCommand) Synopsis() string {
	return "Prints runtime configurations"
}

func (c *PrintCommand) Help() string {
	helpText := `
Usage: vault print <subcommand>

	This command groups subcommands for interacting with Vault's runtime values.

Subcommands:
	token    Token currently in use
`
	return strings.TrimSpace(helpText)
}

func (c *PrintCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PrintCommand) AutocompleteFlags() complete.Flags {
	return nil
}

func (c *PrintCommand) Run(args []string) int {
	return cli.RunResultHelp
}
