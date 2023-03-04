package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*CloudCommand)(nil)

type CloudCommand struct {
	*BaseCommand
}

func (c *CloudCommand) Synopsis() string {
	return "Interact with HCP Vault clusters"
}

func (c *CloudCommand) Help() string {
	helpText := `
Usage: vault cloud <subcommand> [options] [args]
`

	return strings.TrimSpace(helpText)
}

func (c *CloudCommand) Run(args []string) int {
	return cli.RunResultHelp
}
