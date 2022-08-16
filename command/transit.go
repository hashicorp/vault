package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*TransitCommand)(nil)

type TransitCommand struct {
	*BaseCommand
}

func (c *TransitCommand) Synopsis() string {
	return "Perform transit secrets engine specific tasks"
}

func (c *TransitCommand) Help() string {
	helpText := `
Usage: vault transit <subcommand> [options] [args]

  This command hosts assistance functions for interacting with the Transit
  secrets engine.`

	return strings.TrimSpace(helpText)
}

func (c *TransitCommand) Run(args []string) int {
	return cli.RunResultHelp
}
