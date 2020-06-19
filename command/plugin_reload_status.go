package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginReloadCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginReloadCommand)(nil)

type PluginReloadStatusCommand struct {
	*BaseCommand
	reload_id string
}

func (c *PluginReloadStatusCommand) Synopsis() string {
	return "Get the status of an active or recently completed cluster plugin reload"
}

func (c *PluginReloadStatusCommand) Help() string {
	helpText := `
Usage: vault plugin reload-status [options]

  Retrieves the status of a recent cluster plugin reload.  The reload id must be provided.

	  $ vault plugin reload -reload-id=d60a3e83-a598-4f3a-879d-0ddd95f11d4e

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginReloadStatusCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "reload-id",
		Target:     &c.reload_id,
		Completion: complete.PredictAnything,
		Usage:      "The reload id of the recently started plugin reload.",
	})

	return set
}

func (c *PluginReloadStatusCommand) AutocompleteArgs() complete.Predictor {
	return nil
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

	if c.reload_id == "" {
		c.UI.Error("reload-id must be specified")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	r, err := client.Sys().ReloadPluginStatus(c.reload_id)

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving plugin reload status: %s", err))
		return 2
	}

	// This is almost certainly not right and very gross, but how do we output a typed struct?
	c.UI.Output("Time    \tParticipant                         \tSuccess\tMessage")
	c.UI.Output("--------\t------------------------------------\t-------\t-------")
	for i, s := range r["results"].(map[string]interface{}) {
		m := s.(map[string]interface{})
		ts, err := time.Parse(m["timestamp"].(string), time.RFC3339)
		if err != nil {
			return 3
		}
		c.UI.Output(fmt.Sprintf("%s\t%s\t%t\t%s", ts.Format("15:04:05"), i, m["success"].(bool), m["message"].(string)))
	}
	return 0
}
