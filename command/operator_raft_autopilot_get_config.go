package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorRaftAutopilotGetConfigCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorRaftAutopilotGetConfigCommand)(nil)

type OperatorRaftAutopilotGetConfigCommand struct {
	*BaseCommand
}

func (c *OperatorRaftAutopilotGetConfigCommand) Synopsis() string {
	return "Returns the configuration of the autopilot subsystem under integrated storage"
}

func (c *OperatorRaftAutopilotGetConfigCommand) Help() string {
	helpText := `
Usage: vault operator raft autopilot get-config

  Returns the configuration of the autopilot subsystem under integrated storage.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftAutopilotGetConfigCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *OperatorRaftAutopilotGetConfigCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftAutopilotGetConfigCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftAutopilotGetConfigCommand) Run(args []string) int {
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

	resp, err := client.Logical().Read("sys/storage/raft/autopilot/configuration")
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if Format(c.UI) != "table" {
		return OutputData(c.UI, resp)
	}

	entries := []string{"Key | Value"}
	entries = append(entries, fmt.Sprintf("%s | %t", "Cleanup Dead Servers", resp.Data["cleanup_dead_servers"]))
	entries = append(entries, fmt.Sprintf("%s | %s", "Last Contact Threshold", resp.Data["last_contact_threshold"]))
	entries = append(entries, fmt.Sprintf("%s | %s", "Server Stabilization Time", resp.Data["server_stabilization_time"]))

	minQuorum, _ := resp.Data["min_quorum"].(json.Number).Int64()
	entries = append(entries, fmt.Sprintf("%s | %d", "Min Quorum", minQuorum))

	maxTrailingLogs, _ := resp.Data["max_trailing_logs"].(json.Number).Int64()
	entries = append(entries, fmt.Sprintf("%s | %d", "Max Trailing Logs", maxTrailingLogs))

	return OutputData(c.UI, entries)
}
