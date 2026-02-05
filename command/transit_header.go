// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/cli"
	envenc "github.com/hashicorp/vault-envelope-encryption-sdk"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TransitEnvHeaderCommand)(nil)
	_ cli.CommandAutocomplete = (*TransitEnvHeaderCommand)(nil)
)

type TransitEnvHeaderCommand struct {
	*BaseCommand
}

func (c *TransitEnvHeaderCommand) Synopsis() string {
	return "Display the header of an envelope encrypted file/stream"
}

func (c *TransitEnvHeaderCommand) Help() string {
	helpText := `
Usage: vault transit envelope header [filenames...]

  Displays the headers of one or more envelope encrypted files/streams.  Note that this does not validate the authenticity
of the header, as this requires at least partial decryption.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransitEnvHeaderCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetOutputFormat)
}

func (c *TransitEnvHeaderCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitEnvHeaderCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitEnvHeaderCommand) Run(args []string) int {
	return c.envelopeHeader(c.BaseCommand, c.Flags(), args)
}

// error codes: 1: user error, 2: internal computation error, 3: remote api call error
func (tc *TransitEnvHeaderCommand) envelopeHeader(c *BaseCommand, flags *FlagSets, args []string) int {
	// Parse and validate the arguments.
	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.UI.Error("Incorrect argument count (expected 1+ files/streams, got 0)")
		return 1
	}

	for _, file := range args {
		df, in, _, err := openInput(file, c)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
		if df != nil {
			defer df()
		}
		h, err := envenc.ReadHeader(in)
		if err != nil {
			c.UI.Error(err.Error())
			return 3
		}
		OutputData(c.UI, h.Map())
	}
	return 0
}
