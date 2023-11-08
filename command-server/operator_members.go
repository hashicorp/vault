// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command_server

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/command"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorMembersCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorMembersCommand)(nil)
)

type OperatorMembersCommand struct {
	*command.BaseCommand
}

func (c *OperatorMembersCommand) Synopsis() string {
	return "Returns the list of nodes in the cluster"
}

func (c *OperatorMembersCommand) Help() string {
	helpText := `
Usage: vault operator members

  Provides the details of all the nodes in the cluster.

	  $ vault operator members

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorMembersCommand) Flags() *command.FlagSets {
	set := c.FlagSet(command.FlagSetHTTP | command.FlagSetOutputFormat)

	return set
}

func (c *OperatorMembersCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorMembersCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorMembersCommand) Run(args []string) int {
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

	resp, err := client.Sys().HAStatus()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	switch command.Format(c.UI) {
	case "table":
		out := make([]string, 0)
		cols := []string{"Host Name", "API Address", "Cluster Address", "Active Node", "Version", "Upgrade Version", "Redundancy Zone", "Last Echo"}
		out = append(out, strings.Join(cols, " | "))
		for _, node := range resp.Nodes {
			cols := []string{node.Hostname, node.APIAddress, node.ClusterAddress, fmt.Sprintf("%t", node.ActiveNode), node.Version, node.UpgradeVersion, node.RedundancyZone}
			if node.LastEcho != nil {
				cols = append(cols, node.LastEcho.Format(time.RFC3339))
			} else {
				cols = append(cols, "")
			}
			out = append(out, strings.Join(cols, " | "))
		}
		c.UI.Output(command.TableOutput(out, nil))
		return 0
	default:
		return command.OutputData(c.UI, resp)
	}
}
