// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftListPeersCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftListPeersCommand)(nil)
)

type OperatorRaftListPeersCommand struct {
	*BaseCommand
	flagDRToken string
}

func (c *OperatorRaftListPeersCommand) Synopsis() string {
	return "Returns the Raft peer set"
}

func (c *OperatorRaftListPeersCommand) Help() string {
	helpText := `
Usage: vault operator raft list-peers

  Provides the details of all the peers in the Raft cluster.

	  $ vault operator raft list-peers

  Provides the details of all the peers in the Raft cluster of a DR secondary
  cluster. This command should be invoked on the DR secondary nodes.

      $ vault operator raft list-peers -dr-token <dr-operation-token>

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftListPeersCommand) Flags() *FlagSets {
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

func (c *OperatorRaftListPeersCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftListPeersCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftListPeersCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var secret *api.Secret
	switch {
	case c.flagDRToken != "":
		secret, err = client.Logical().Write("sys/storage/raft/configuration", map[string]interface{}{
			"dr_operation_token": c.flagDRToken,
		})
	default:
		secret, err = client.Logical().Read("sys/storage/raft/configuration")
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading the raft cluster configuration: %s", err))
		return 2
	}
	if secret == nil {
		c.UI.Error("No raft cluster configuration found")
		return 2
	}

	if Format(c.UI) != "table" {
		return OutputSecret(c.UI, secret)
	}

	config := secret.Data["config"].(map[string]interface{})

	servers := config["servers"].([]interface{})
	out := []string{"Node | Address | State | Voter"}
	for _, serverRaw := range servers {
		server := serverRaw.(map[string]interface{})
		state := "follower"
		if server["leader"].(bool) {
			state = "leader"
		}

		out = append(out, fmt.Sprintf("%s | %s | %s | %t", server["node_id"].(string), server["address"].(string), state, server["voter"].(bool)))
	}

	c.UI.Output(tableOutput(out, nil))
	return 0
}
