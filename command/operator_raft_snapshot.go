// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*OperatorRaftSnapshotCommand)(nil)

type OperatorRaftSnapshotCommand struct {
	*BaseCommand
}

func (c *OperatorRaftSnapshotCommand) Synopsis() string {
	return "Restores and saves snapshots from the Raft cluster"
}

func (c *OperatorRaftSnapshotCommand) Help() string {
	helpText := `
Usage: vault operator raft snapshot <subcommand> [options] [args]

  This command groups subcommands for operators interacting with the snapshot
  functionality of the integrated Raft storage backend. Here are a few examples of
  the Raft snapshot operator commands:

  Installs the provided snapshot, returning the cluster to the state defined in it:

      $ vault operator raft snapshot restore raft.snap

  Saves a snapshot of the current state of the Raft cluster into a file:

      $ vault operator raft snapshot save raft.snap

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftSnapshotCommand) Run(args []string) int {
	return cli.RunResultHelp
}
