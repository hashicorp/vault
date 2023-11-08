// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command_server

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/command"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftAutopilotStateCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftAutopilotStateCommand)(nil)
)

type OperatorRaftAutopilotStateCommand struct {
	*command.BaseCommand
	flagDRToken string
}

func (c *OperatorRaftAutopilotStateCommand) Synopsis() string {
	return "Displays the state of the raft cluster under integrated storage as seen by autopilot"
}

func (c *OperatorRaftAutopilotStateCommand) Help() string {
	helpText := `
Usage: vault operator raft autopilot state

  Displays the state of the raft cluster under integrated storage as seen by autopilot.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftAutopilotStateCommand) Flags() *command.FlagSets {
	set := c.FlagSet(command.FlagSetHTTP | command.FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&command.StringVar{
		Name:       "dr-token",
		Target:     &c.flagDRToken,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage:      "DR operation token used to authorize this request (if a DR secondary node).",
	})

	// The output of the state endpoint contains nested values and is not fit for
	// the default "table" display format. Override the default display format to
	// "pretty", both in the flag and in the UI.
	set.MainSet.VisitAll(func(fl *flag.Flag) {
		if fl.Name == "format" {
			fl.DefValue = "pretty"
		}
	})
	ui, ok := c.UI.(*command.VaultUI)
	if ok && ui.Format == "table" {
		ui.Format = "pretty"
	}
	return set
}

func (c *OperatorRaftAutopilotStateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftAutopilotStateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftAutopilotStateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch len(args) {
	case 0:
	default:
		c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var state *api.AutopilotState
	switch {
	case c.flagDRToken != "":
		state, err = client.Sys().RaftAutopilotStateWithDRToken(c.flagDRToken)
	default:
		state, err = client.Sys().RaftAutopilotState()
	}

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error checking autopilot state: %s", err))
		return 2
	}

	if state == nil {
		return 0
	}

	return command.OutputData(c.UI, state)
}
