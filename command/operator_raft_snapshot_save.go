package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftSnapshotSaveCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftSnapshotSaveCommand)(nil)
)

type OperatorRaftSnapshotSaveCommand struct {
	*BaseCommand
}

func (c *OperatorRaftSnapshotSaveCommand) Synopsis() string {
	return "Saves a snapshot of the current state of the Raft cluster into a file"
}

func (c *OperatorRaftSnapshotSaveCommand) Help() string {
	helpText := `
Usage: vault operator raft snapshot save <snapshot_file>

  Saves a snapshot of the current state of the Raft cluster into a file.

	  $ vault operator raft snapshot save raft.snap

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftSnapshotSaveCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *OperatorRaftSnapshotSaveCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftSnapshotSaveCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftSnapshotSaveCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	path := ""

	args = f.Args()
	switch len(args) {
	case 1:
		path = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 1, got %d)", len(args)))
		return 1
	}

	if len(path) == 0 {
		c.UI.Error("Output file name is required")
		return 1
	}

	w := &lazyOpenWriter{
		openFunc: func() (io.WriteCloser, error) {
			return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		},
	}

	client, err := c.Client()
	if err != nil {
		w.Close()
		c.UI.Error(err.Error())
		return 2
	}

	err = client.Sys().RaftSnapshot(w)
	if err != nil {
		w.Close()
		c.UI.Error(fmt.Sprintf("Error taking the snapshot: %s", err))
		return 2
	}

	err = w.Close()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error taking the snapshot: %s", err))
		return 2
	}
	return 0
}

type lazyOpenWriter struct {
	openFunc func() (io.WriteCloser, error)
	writer   io.WriteCloser
}

func (l *lazyOpenWriter) Write(p []byte) (n int, err error) {
	if l.writer == nil {
		var err error
		l.writer, err = l.openFunc()
		if err != nil {
			return 0, err
		}
	}
	return l.writer.Write(p)
}

func (l *lazyOpenWriter) Close() error {
	if l.writer != nil {
		return l.writer.Close()
	}
	return nil
}
