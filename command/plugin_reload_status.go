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
	_ cli.Command             = (*PluginReloadCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginReloadCommand)(nil)
)

type PluginReloadStatusCommand struct {
	*BaseCommand
}

func (c *PluginReloadStatusCommand) Synopsis() string {
	return "Get the status of an active or recently completed global plugin reload"
}

func (c *PluginReloadStatusCommand) Help() string {
	helpText := `
Usage: vault plugin reload-status RELOAD_ID

  Retrieves the status of a recent cluster plugin reload.  The reload id must be provided.

	  $ vault plugin reload-status d60a3e83-a598-4f3a-879d-0ddd95f11d4e

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginReloadStatusCommand) Flags() *FlagSets {
	return c.FlagSet(FlagSetHTTP)
}

func (c *PluginReloadStatusCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PluginReloadStatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginReloadStatusCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	reloadId := strings.TrimSpace(args[0])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	r, err := client.Sys().ReloadPluginStatus(&api.ReloadPluginStatusInput{
		ReloadID: reloadId,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving plugin reload status: %s", err))
		return 2
	}
	out := []string{"Time | Participant | Success | Message "}
	for i, s := range r.Results {
		out = append(out, fmt.Sprintf("%s | %s | %t | %s ",
			s.Timestamp.Format("15:04:05"),
			i,
			s.Error == "",
			s.Error))
	}
	c.UI.Output(TableOutput(out, nil))
	return 0
}
