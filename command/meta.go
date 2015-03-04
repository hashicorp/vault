package command

import (
	"bufio"
	"flag"
	"io"

	"github.com/mitchellh/cli"
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetServer  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetServer
)

// Meta contains the meta-options and functionality that nearly every
// Vault command inherits.
type Meta struct {
	Ui cli.Ui

	// These are set by the command line flags.
	flagAddress  string
	flagCACert   string
	flagCAPath   string
	flagInsecure bool
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	// FlagSetServer tells us to enable the settings for selecting
	// the server information.
	if fs&FlagSetServer != 0 {
		f.StringVar(&m.flagAddress, "address", "", "")
		f.StringVar(&m.flagCACert, "ca-cert", "", "")
		f.StringVar(&m.flagCAPath, "ca-path", "", "")
		f.BoolVar(&m.flagInsecure, "insecure", false, "")
	}

	// Create an io.Writer that writes to our Ui properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.Ui.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}
