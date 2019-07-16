package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*OperatorRaftCommand)(nil)

type OperatorRaftCommand struct {
	*BaseCommand
}

func (c *OperatorRaftCommand) Synopsis() string {
	return "Interact with Vault's raft storage backend"
}

func (c *OperatorRaftCommand) Help() string {
	helpText := `
Usage: vault operator raft <subcommand> [options] [args]

  This command groups subcommands for operators interacting with the Vault raft storage backend. Most
  users will not need to interact with these commands. Here are a few examples
  of the raft operator commands:

  Joins a node to the raft cluster:

      $ vault operator raft join https://127.0.0.1:8200

  Returns the raft cluster configuration:

      $ vault operator raft configuration

  Removes a node from the raft cluster:

      $ vault operator raft remove-peer

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftCommand) Run(args []string) int {
	return cli.RunResultHelp
}
