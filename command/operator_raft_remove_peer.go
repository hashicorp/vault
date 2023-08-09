// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftRemovePeerCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftRemovePeerCommand)(nil)
)

type OperatorRaftRemovePeerCommand struct {
	*BaseCommand

	flagDRToken string
}

func (c *OperatorRaftRemovePeerCommand) Synopsis() string {
	return "Removes a node from the Raft cluster"
}

func (c *OperatorRaftRemovePeerCommand) Help() string {
	helpText := `
Usage: vault operator raft remove-peer <server_id>

  Removes a node from the Raft cluster.

	  $ vault operator raft remove-peer node1

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftRemovePeerCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "dr-token",
		Target:     &c.flagDRToken,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage:      "DR operation token used to authorize this request (if a DR secondary node).",
	})

	return set
}

func (c *OperatorRaftRemovePeerCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftRemovePeerCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftRemovePeerCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	serverID := ""

	args = f.Args()
	switch len(args) {
	case 1:
		serverID = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 1, got %d)", len(args)))
		return 1
	}

	if len(serverID) == 0 {
		c.UI.Error("Server id is required")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	_, err = client.Logical().Write("sys/storage/raft/remove-peer", map[string]interface{}{
		"server_id":          serverID,
		"dr_operation_token": c.flagDRToken,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error removing the peer from raft cluster: %s", err))
		return 2
	}

	c.UI.Output("Peer removed successfully!")

	return 0
}
