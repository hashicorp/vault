package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorRaftSnapshotRestoreCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorRaftSnapshotRestoreCommand)(nil)

type OperatorRaftSnapshotRestoreCommand struct {
	flagForce bool
	*BaseCommand
}

func (c *OperatorRaftSnapshotRestoreCommand) Synopsis() string {
	return "Installs the provided snapshot, returning the cluster to the state defined in it"
}

func (c *OperatorRaftSnapshotRestoreCommand) Help() string {
	helpText := `
Usage: vault operator raft snapshot restore <snapshot_file>

  Installs the provided snapshot, returning the cluster to the state defined in it.

	  $ vault operator raft snapshot restore raft.snap

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftSnapshotRestoreCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "force",
		Target:  &c.flagForce,
		Default: false,
		Usage:   "This bypasses checks ensuring the Autounseal or shamir keys are consistent with the snapshot data.",
	})

	return set
}

func (c *OperatorRaftSnapshotRestoreCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftSnapshotRestoreCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftSnapshotRestoreCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	snapFile := ""

	args = f.Args()
	switch len(args) {
	case 1:
		snapFile = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 1, got %d)", len(args)))
		return 1
	}

	if len(snapFile) == 0 {
		c.UI.Error("Snapshot file name is required")
		return 1
	}

	snapReader, err := os.Open(snapFile)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error opening policy file: %s", err))
		return 2
	}
	defer snapReader.Close()

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	err = client.Sys().RaftSnapshotRestore(snapReader, c.flagForce)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error installing the snapshot: %s", err))
		return 2
	}

	return 0
}
